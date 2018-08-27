package silvia

import "sync"

type (
	AdjustRingItem struct {
		Event *AdjustEvent
		Error error
	}

	AdjustRing struct {
		sync.RWMutex
		Ring  []*AdjustRingItem
		Total int
		Size  int
	}

	SnowplowRingItem struct {
		Event *SnowplowEvent
		Error error
	}

	SnowplowRing struct {
		sync.RWMutex
		Ring  []*SnowplowRingItem
		Total int
		Size  int
	}
)

func (ring *AdjustRing) Add(event *AdjustEvent, err error) {
	ringItem := &AdjustRingItem{
		Event: event,
		Error: err,
	}

	ring.Lock()
	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size {
		ring.Ring = ring.Ring[1:]
	}
	ring.Unlock()
}

func (ring *SnowplowRing) Add(event *SnowplowEvent, err error) {
	ringItem := &SnowplowRingItem{
		Event: event,
		Error: err,
	}

	ring.Lock()
	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size {
		ring.Ring = ring.Ring[1:]
	}
	ring.Unlock()
}

func (ring *AdjustRing) Display() AdjustRing {
	ring.RLock()
	defer ring.RUnlock()
	return *ring
}

func (ring *SnowplowRing) Display() SnowplowRing {
	ring.RLock()
	defer ring.RUnlock()
	return *ring
}
