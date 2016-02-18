package silvia

type(
	AdjustRingItem struct {
		Event *AdjustEvent
		Error error
	}

	AdjustRing struct {
		Ring  []*AdjustRingItem
		Total int
		Size  int
	}

	SnowplowRingItem struct {
		Event *SnowplowEvent
		Error error
	}

	SnowplowRing struct {
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

	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size { ring.Ring = ring.Ring[1:] }
}

func (ring *SnowplowRing) Add(event *SnowplowEvent, err error) {
	ringItem := &SnowplowRingItem{
		Event: event,
		Error: err,
	}

	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size { ring.Ring = ring.Ring[1:] }
}
