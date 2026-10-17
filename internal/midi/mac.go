//go:build darwin && cgo

package midi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include "native.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

type Mac struct{}

func (Mac) Devices() ([]Device, error) {
	var devices [128]C.cw_device
	n := int(C.cw_devices(&devices[0], C.int(len(devices))))
	if n < 0 || n > len(devices) {
		return nil, fmt.Errorf("CoreMIDI discovery failed (%d)", n)
	}
	var result []Device
	for i := 0; i < n; i++ {
		result = append(result, Device{int32(devices[i].id), C.GoString(&devices[i].name[0])})
	}
	return result, nil
}
func (Mac) Connect(id int32) (Port, error) {
	var status C.int
	native := C.cw_open(C.int32_t(id), &status)
	if native == nil {
		return nil, fmt.Errorf("MIDI connection failed (%d); select an online Vail input with a matching output", int(status))
	}
	return &macPort{native: native}, nil
}
func HostTime() time.Duration { return time.Duration(C.cw_now()) }

type macPort struct {
	mu      sync.Mutex
	native  *C.cw_midi
	decoder Decoder
}

func (p *macPort) Poll() ([]Event, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.native == nil {
		return nil, fmt.Errorf("MIDI connection closed")
	}
	var result []Event
	// Bound each poll so a busy device cannot starve keyboard processing.
	for i := 0; i < 1024; i++ {
		var packet C.cw_packet
		n := int(C.cw_poll(p.native, &packet))
		if n < 0 {
			return nil, fmt.Errorf("MIDI disconnected or input overflowed; reconnect and retry")
		}
		if n == 0 {
			break
		}
		bytes := C.GoBytes(unsafe.Pointer(&packet.bytes[0]), C.int(packet.length))
		result = append(result, p.decoder.Feed(time.Duration(packet.nanos), bytes)...)
	}
	return result, nil
}
func (p *macPort) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.native == nil {
		return nil
	}
	status := C.cw_close(p.native)
	p.native = nil
	if status != 0 {
		return fmt.Errorf("could not restore keyboard mode (%d); reconnect adapter", int(status))
	}
	return nil
}
