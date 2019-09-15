package contract

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

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

func MakeTMEvent(contractAddr common.Address, evs []*Event) (*types.Event, error) {
	pairs, err := eventsToPairs(evs)
	if err != nil {
		return nil, err
	}
	pairs = append(pairs, tmcmn.KVPair{
		Key:   []byte("address"),
		Value: []byte(contractAddr.Hex()),
	})
	e := types.Event{Type: "contract"}
	e.Attributes = pairs
	return &e, nil
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
		pairs = append(pairs, tmcmn.KVPair{Key: dataKey, Value: ev.Bytes()})
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

func (ev *Event) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(len(ev.Name)))
	buf.Write(ev.Name)
	buf.Write(ev.Value)
	return []byte(hex.EncodeToString(buf.Bytes()))
}

func ParseEventData(hexStrBytes []byte) (*Event, error) {
	b, err := hex.DecodeString(string(hexStrBytes))
	if err != nil {
		return nil, err
	}

	size := b[0]
	name := b[1 : 1+size]
	value := b[1+size : len(b)]

	return &Event{Name: name, Value: value}, nil
}

func MakeEvent(name, value string) (*Event, error) {
	var v []byte
	if strings.Contains(value, "0x") {
		h, err := hex.DecodeString(value[2:])
		if err != nil {
			return nil, err
		}
		v = h
	} else {
		v = []byte(value)
	}
	return &Event{Name: []byte(name), Value: v}, nil
}

func MakeEventBytes(name, value string) ([]byte, error) {
	ev, err := MakeEvent(name, value)
	if err != nil {
		return nil, err
	}
	return ev.Bytes(), nil
}

// MakeEventSearchQuery returns a query for searching events on transaction
// contractAddr: target contract address
// eventName: event name
// eventValue: value corresponding to event name. NOTE: if value has a prefix "0x", value will be decoded into []byte
func MakeEventSearchQuery(contractAddr string, eventName, eventValue string) (string, error) {
	var parts []string

	parts = append(parts, fmt.Sprintf("contract.address='%v'", contractAddr))
	parts = append(parts, fmt.Sprintf("contract.event.name='%v'", eventName))
	if len(eventValue) > 0 {
		ev, err := MakeEventBytes(eventName, eventValue)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("contract.event.data='%v'", string(ev)))
	}
	return strings.Join(parts, " AND "), nil
}
