package onearrrev

import (
	"io"

	"github.com/egonelbre/exp/niterator/shape"
)

type Index struct {
	Track  uint32
	Shape  uint32
	Stride uint32
}

type Iterator struct {
	*shape.AP

	Track     []Index
	NextIndex uint32
	Done      bool
}

func NewIterator(ap *shape.AP) *Iterator {
	track := make([]Index, len(ap.Shape))
	for i := range ap.Shape {
		track[3-i].Track = uint32(ap.Shape[i])
		track[3-i].Shape = uint32(ap.Shape[i])
		track[3-i].Stride = uint32(ap.Stride[i])
	}

	return &Iterator{
		AP:    ap,
		Track: track,
	}
}

func (it *Iterator) IsDone() bool {
	return it.Done
}

func (it *Iterator) Next() (int, error) {
	if it.Done {
		return 0, io.EOF
	}

	next := it.NextIndex
	result := next
	track := it.Track
	for i := range track {
		x := &track[i]
		x.Track--
		if x.Track > 0 {
			next += x.Stride
			it.NextIndex = next
			return int(result), nil
		}
		x.Track = x.Shape
		next -= (x.Shape - 1) * x.Stride
	}

	it.Done = true
	it.NextIndex = next
	return int(result), nil
}
