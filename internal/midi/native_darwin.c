//go:build darwin && cgo

#include "native.h"
#include <CoreMIDI/CoreMIDI.h>
#include <CoreFoundation/CoreFoundation.h>
#include <mach/mach_time.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

#define QUEUE_SIZE 1024
struct cw_midi {
 MIDIClientRef client;
 MIDIPortRef input, output;
 MIDIEndpointRef source, destination;
 MIDIUniqueID source_id, destination_id;
 pthread_mutex_t mutex;
 cw_packet queue[QUEUE_SIZE];
 int head, tail, overflow, mode_changed;
 mach_timebase_info_data_t timebase;
};
static uint64_t nanos(uint64_t ticks, mach_timebase_info_data_t tb) {
 return (uint64_t)(((__uint128_t)ticks * tb.numer) / tb.denom);
}
uint64_t cw_now(void) {
 mach_timebase_info_data_t tb; mach_timebase_info(&tb);
 return nanos(mach_absolute_time(), tb);
}
static int available(MIDIEndpointRef endpoint) {
 SInt32 offline=0;
 if (MIDIObjectGetIntegerProperty(endpoint,kMIDIPropertyOffline,&offline)==noErr && offline) return 0;
 return endpoint!=0;
}
int cw_devices(cw_device *devices, int capacity) {
 MIDIClientRef client=0;
 OSStatus status=MIDIClientCreate(CFSTR("CW device discovery"),NULL,NULL,&client);
 if(status) return status<0?status:-status;
 int count=0;
 for(ItemCount i=0;i<MIDIGetNumberOfSources() && count<capacity;i++) {
  MIDIEndpointRef source=MIDIGetSource(i);
  if(!available(source)) continue;
  SInt32 id=0;
  if(MIDIObjectGetIntegerProperty(source,kMIDIPropertyUniqueID,&id)) continue;
  devices[count].id=id;
  CFStringRef name=NULL;
  if(MIDIObjectGetStringProperty(source,kMIDIPropertyDisplayName,&name)==noErr && name) {
   CFStringGetCString(name,devices[count].name,sizeof(devices[count].name),kCFStringEncodingUTF8);
   CFRelease(name);
  }
  count++;
 }
 MIDIClientDispose(client);
 return count;
}
static int same_endpoint(MIDIUniqueID id,MIDIEndpointRef expected) {
 MIDIObjectRef found=0; MIDIObjectType type;
 return MIDIObjectFindByUniqueID(id,&found,&type)==noErr && found==expected && available(found);
}
static OSStatus mode(cw_midi *m,Byte value) {
 MIDIPacketList list;
 MIDIPacket *packet=MIDIPacketListInit(&list);
 Byte bytes[3]={0xb0,0,value};
 MIDIPacketListAdd(&list,sizeof(list),packet,0,3,bytes);
 return MIDISend(m->output,m->destination,&list);
}
static void receive(const MIDIPacketList *packets,void *context,void *connection) {
 (void)connection;
 cw_midi *m=context;
 const MIDIPacket *packet=&packets->packet[0];
 pthread_mutex_lock(&m->mutex);
 for(UInt32 i=0;i<packets->numPackets;i++) {
  int next=(m->tail+1)%QUEUE_SIZE;
  if(next==m->head || packet->length>256) {m->overflow=1;break;}
  cw_packet *out=&m->queue[m->tail];
  out->nanos=nanos(packet->timeStamp?packet->timeStamp:mach_absolute_time(),m->timebase);
  out->length=packet->length;
  memcpy(out->bytes,packet->data,packet->length);
  m->tail=next;
  packet=MIDIPacketNext(packet);
 }
 pthread_mutex_unlock(&m->mutex);
}
int cw_close(cw_midi *m) {
 if(!m) return 0;
 int status=0;
 if(m->mode_changed) {
  if(same_endpoint(m->destination_id,m->destination)) status=mode(m,127);
  else status=-5;
 }
 if(m->input) { MIDIPortDisconnectSource(m->input,m->source);MIDIPortDispose(m->input); }
 if(m->output) MIDIPortDispose(m->output);
 if(m->client) MIDIClientDispose(m->client);
 pthread_mutex_destroy(&m->mutex);
 free(m);
 return status;
}
cw_midi *cw_open(int32_t id,int *status) {
 cw_midi *m=calloc(1,sizeof(cw_midi));
 if(!m) {*status=-1;return NULL;}
 pthread_mutex_init(&m->mutex,NULL);mach_timebase_info(&m->timebase);
 *status=MIDIClientCreate(CFSTR("CW practice input"),NULL,NULL,&m->client);
 if(*status) goto fail;
 MIDIObjectRef object=0; MIDIObjectType type;
 *status=MIDIObjectFindByUniqueID(id,&object,&type);
 if(*status || type!=kMIDIObjectType_Source) {*status=-2;goto fail;}
 m->source=object;m->source_id=id;
 MIDIEntityRef entity=0;
 MIDIEndpointGetEntity(m->source,&entity);
 if(!entity) {*status=-3;goto fail;}
 for(ItemCount i=0;i<MIDIGetNumberOfDestinations();i++) {
  MIDIEndpointRef target=MIDIGetDestination(i);MIDIEntityRef other=0;
  MIDIEndpointGetEntity(target,&other);
  if(other==entity && available(target)) {
   if(m->destination) {*status=-4;goto fail;}
   m->destination=target;
  }
 }
 if(!m->destination || !available(m->source)) {*status=-5;goto fail;}
 *status=MIDIObjectGetIntegerProperty(m->destination,kMIDIPropertyUniqueID,&m->destination_id);
 if(*status) goto fail;
 *status=MIDIInputPortCreate(m->client,CFSTR("CW marks"),receive,m,&m->input);
 if(*status) goto fail;
 *status=MIDIOutputPortCreate(m->client,CFSTR("Vail mode"),&m->output);
 if(*status) goto fail;
 *status=MIDIPortConnectSource(m->input,m->source,NULL);
 if(*status) goto fail;
 m->mode_changed=1;
 *status=mode(m,0);
 if(*status) goto fail;
 return m;
fail:
 cw_close(m);return NULL;
}
int cw_poll(cw_midi *m,cw_packet *packet) {
 if(!same_endpoint(m->source_id,m->source)||!same_endpoint(m->destination_id,m->destination)) return -1;
 pthread_mutex_lock(&m->mutex);
 int result=0;
 if(m->overflow) result=-2;
 else if(m->head!=m->tail) {*packet=m->queue[m->head];m->head=(m->head+1)%QUEUE_SIZE;result=1;}
 pthread_mutex_unlock(&m->mutex);
 return result;
}
