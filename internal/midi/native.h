#ifndef CW_MIDI_H
#define CW_MIDI_H
#include <stdint.h>
typedef struct cw_midi cw_midi;
typedef struct { int32_t id; char name[256]; } cw_device;
typedef struct { uint64_t nanos; uint16_t length; uint8_t bytes[256]; } cw_packet;
int cw_devices(cw_device *devices, int capacity);
cw_midi *cw_open(int32_t id, int *status);
int cw_poll(cw_midi *m, cw_packet *packet);
int cw_close(cw_midi *m);
uint64_t cw_now(void);
#endif
