package bind

import (
	"errors"
	"fmt"
	"github.com/bluele/hypermint/pkg/contract/event"
	"reflect"
)

var ErrTypeMismatch = errors.New("type mismatch")
var ErrUnknownType = errors.New("unknown type")
var ErrNotFound = errors.New("event not found")

type Event interface {
	Decode(bs []byte) error
}

type EventInfo struct {
	ID string
	EventCreator func() Event
}

type EventDecoder struct {
	events map[string]EventInfo
}

func NewEventDecoder() *EventDecoder {
	return &EventDecoder{
		events: make(map[string]EventInfo),
	}
}

func (d *EventDecoder) Register(event EventInfo) {
	d.events[event.ID] = event
}

func (d *EventDecoder) Decode(v interface{}, entry *event.Entry) error {
	if reflect.Ptr != reflect.ValueOf(v).Kind() {
		return fmt.Errorf("abi: Decode(non-pointer %T)", v)
	}
	ei, ok := d.events[string(entry.Name)]
	if !ok {
		return ErrUnknownType
	}
	e := ei.EventCreator()
	if reflect.Ptr != reflect.ValueOf(e).Kind() {
		return fmt.Errorf("abi: Decode(non-pointer %T)", e)
	}
	if reflect.TypeOf(v) != reflect.TypeOf(e) {
		return ErrTypeMismatch
	}
	err := e.Decode(entry.Value)
	if err != nil {
		return err
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(e).Elem())
	return nil
}

func (d *EventDecoder) FindFirst(v interface{}, entries []*event.Entry) error {
	for _, e := range entries {
		if err := d.Decode(v, e); err != nil && err != ErrTypeMismatch {
			return err
		} else if err == nil {
			return nil
		}
	}
	return ErrNotFound
}
