package morse

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"
	"sync/atomic"
)

// PCM streams the bounded timeline without allocating audio for the entire line.
// Reads and seeks are serialized by the audio player.
type PCM struct {
	timeline  Timeline
	hz        float64
	position  int64
	published atomic.Int64
}

func NewPCM(timeline Timeline, hz float64) (*PCM, error) {
	if math.IsNaN(hz) || hz < 100 || hz > 2000 {
		return nil, fmt.Errorf("tone must be 100–2000 Hz")
	}
	return &PCM{timeline: timeline, hz: hz}, nil
}
func (p *PCM) Read(b []byte) (int, error) {
	size := p.timeline.Samples * 4
	if p.position >= size {
		return 0, io.EOF
	}
	n := int(min(int64(len(b)), size-p.position))
	for i := 0; i < n; {
		sampleIndex := p.position / 4
		j := sort.Search(len(p.timeline.Segments), func(j int) bool { return p.timeline.Segments[j].End > sampleIndex })
		segment := p.timeline.Segments[j]
		value := float32(0)
		if segment.Tone {
			offset := sampleIndex - segment.Start
			ramp := min(int64(SampleRate*3/1000), (segment.End-segment.Start)/2)
			envelope := min(1.0, float64(min(offset, segment.End-sampleIndex-1))/float64(max(1, ramp)))
			value = float32(0.20 * envelope * math.Sin(2*math.Pi*p.hz*float64(offset)/SampleRate))
		}
		var encoded [4]byte
		binary.LittleEndian.PutUint32(encoded[:], math.Float32bits(value))
		copied := copy(b[i:n], encoded[p.position%4:])
		i += copied
		p.position += int64(copied)
	}
	p.published.Store(p.position)
	return n, nil
}
func (p *PCM) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		offset += p.position
	case io.SeekEnd:
		offset += p.timeline.Samples * 4
	default:
		return p.position, fmt.Errorf("invalid seek origin")
	}
	if offset < 0 || offset > p.timeline.Samples*4 {
		return p.position, fmt.Errorf("seek outside PCM")
	}
	p.position = offset
	p.published.Store(offset)
	return offset, nil
}

// ReadSamples is safe to observe while the audio player reads PCM.
func (p *PCM) ReadSamples() int64 { return p.published.Load() / 4 }
