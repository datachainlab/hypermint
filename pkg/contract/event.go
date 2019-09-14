package contract

import (
	"bytes"
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
	tmcmn "github.com/tendermint/tendermint/libs/common"
)

type Event struct {
	Name  []byte
	Value []byte
}

func (ev Event) String() string {
	return fmt.Sprintf("%v{0x%X}", string(ev.Name), ev.Value)
}

func MakeTMEvents(txAddr common.Address, evs []*Event) (types.Events, error) {
	pairs, err := eventsToPairs(evs)
	if err != nil {
		return nil, err
	}
	pairs = append(pairs, tmcmn.KVPair{
		Key:   []byte("address"),
		Value: []byte(txAddr.Hex()),
	})
	e := types.Event{Type: "contract"}
	e.Attributes = pairs
	return types.Events{e}, nil
}

func eventsToPairs(evs []*Event) (tmcmn.KVPairs, error) {
	var pairs tmcmn.KVPairs
	for _, ev := range evs {
		if err := validateEvent(ev); err != nil {
			return nil, err
		}
		key := []byte("event.name")
		pairs = append(pairs, tmcmn.KVPair{Key: key, Value: ev.Name})
		dataKey := []byte("event.data")
		pairs = append(pairs, tmcmn.KVPair{Key: dataKey, Value: makeEventData(ev)})
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
