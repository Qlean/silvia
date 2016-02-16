package main

type(
	AdjustRingItem struct {
		Body  string
		Event *AdjustEvent
		Error error
	}

	AdjustRing struct {
		Ring  []*AdjustRingItem
		Total int
		Size  int
	}

	SnowplowRingItem struct {
		Body  string
		Event *SnowplowEvent
		Error error
	}

	AdjustRing struct {
		Ring  []*AdjustRingItem
		Total int
		Size  int
	}
)

func (ring *AdjustRing) Add(event *AdjustEvent, err error) {
	ringItem := &AdjustRingItem{
		Body: event.Body,
		AdjustEvent: event,
		Error: err,
	}

	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size { ring.Ring = ring.Ring[1:] }
}

func (ring *SnowplowRing) Add(event *SnowplowEvent, err error) {
	ringItem := &SnowplowRingItem{
		Body: event.Body,
		AdjustEvent: event,
		Error: err,
	}

	ring.Ring = append(ring.Ring, ringItem)
	ring.Total++

	if len(ring.Ring) > ring.Size { ring.Ring = ring.Ring[1:] }
}
