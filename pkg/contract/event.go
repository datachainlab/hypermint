package contract

import (
	"bytes"
	"fmt"

	"github.com/tendermint/tendermint/libs/common"
)

type Event struct {
	Name  []byte
	Value []byte
}

func (ev Event) String() string {
	return fmt.Sprintf("%v{0x%X}", string(ev.Name), ev.Value)
}

func EventsToTags(evs []*Event) (common.KVPairs, error) {
	var pairs common.KVPairs
	for _, ev := range evs {
		if err := validateEvent(ev); err != nil {
			return nil, err
		}
		key := []byte("event.name")
		pairs = append(pairs, common.KVPair{Key: key, Value: ev.Name})
		dataKey := []byte("event.data")
		pairs = append(pairs, common.KVPair{Key: dataKey, Value: makeEventData(ev)})
	}
	return pairs, nil
}

func validateEvent(ev *Event) error {
	if len(ev.Name) == 0 {
		return fmt.Errorf("length of Name must be greater than 0")
	}
	if len(ev.Name) > 32 {
		return fmt.Errorf("length of Name must be not greater than 32")
	}
	if len(ev.Value) == 0 {
		return fmt.Errorf("length of Value must be greater than 0")
	}
	if len(ev.Value) > 1024 {
		return fmt.Errorf("length of Value must be not greater than 1024")
	}
	return nil
}

func makeEventData(ev *Event) []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(len(ev.Name)))
	buf.Write(ev.Name)
	buf.Write(ev.Value)
	return buf.Bytes()
}

func ParseEventData(b []byte) (*Event, error) {
	size := b[0]
	name := b[1 : 1+size]
	value := b[1+size : len(b)]

	return &Event{Name: name, Value: value}, nil
}
