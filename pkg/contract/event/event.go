package event

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
	tmcmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	EventKey     = "event"
	EventDataKey = EventKey + ".data"
	EventNameKey = EventKey + ".name"

	AddressKey           = "address"
	ContractKey          = "contract"
	ContractAddressKey   = ContractKey + "." + AddressKey
	ContractEventDataKey = ContractKey + "." + EventDataKey
	ContractEventNameKey = ContractKey + "." + EventNameKey
)

type Event struct {
	address common.Address
	entries []*Entry
}

func NewEvent(address common.Address, entries []*Entry) *Event {
	return &Event{
		address: address,
		entries: entries,
	}
}

func (es Event) Address() common.Address {
	return es.address
}

func (es Event) Entries() []*Entry {
	return es.entries
}

type Entry struct {
	Name  []byte
	Value []byte
}

func (e Entry) String() string {
	return fmt.Sprintf("%v{0x%X}", string(e.Name), e.Value)
}

func (e Entry) Validate() error {
	if len(e.Name) == 0 {
		return fmt.Errorf("length of Name must be greater than 0")
	}
	if len(e.Name) > 32 {
		return fmt.Errorf("length of Name must be not greater than 32")
	}
	if len(e.Value) == 0 {
		return fmt.Errorf("length of Value must be greater than 0")
	}
	if len(e.Value) > 1024 {
		return fmt.Errorf("length of Value must be not greater than 1024")
	}
	return nil
}

func (e Entry) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(len(e.Name)))
	buf.Write(e.Name)
	buf.Write(e.Value)
	return []byte(hex.EncodeToString(buf.Bytes()))
}

func MakeTMEvents(evs []*Event) (types.Events, error) {
	var events types.Events
	for _, ev := range evs {
		event, err := MakeTMEvent(ev.Address(), ev.Entries())
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}
	return events, nil
}

func MakeTMEvent(contractAddr common.Address, es []*Entry) (*types.Event, error) {
	pairs, err := eventsToPairs(es)
	if err != nil {
		return nil, err
	}
	pairs = append(pairs, tmcmn.KVPair{
		Key:   []byte(AddressKey),
		Value: []byte(contractAddr.Hex()),
	})
	e := types.Event{Type: ContractKey}
	e.Attributes = pairs
	return &e, nil
}

func eventsToPairs(es []*Entry) (tmcmn.KVPairs, error) {
	var pairs tmcmn.KVPairs
	for _, e := range es {
		if err := e.Validate(); err != nil {
			return nil, err
		}
		pairs = append(pairs, tmcmn.KVPair{Key: []byte(EventNameKey), Value: e.Name})
		pairs = append(pairs, tmcmn.KVPair{Key: []byte(EventDataKey), Value: e.Bytes()})
	}
	return pairs, nil
}

func ParseEntry(hexStrBytes []byte) (*Entry, error) {
	b, err := hex.DecodeString(string(hexStrBytes))
	if err != nil {
		return nil, err
	}

	size := b[0]
	name := b[1 : 1+size]
	value := b[1+size : len(b)]

	return &Entry{Name: name, Value: value}, nil
}

func MakeEntry(name, value string) (*Entry, error) {
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
	return &Entry{Name: []byte(name), Value: v}, nil
}

func MakeEntryBytes(name, value string) ([]byte, error) {
	e, err := MakeEntry(name, value)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

// MakeEventSearchQuery returns a query for searching events on transaction
// contractAddr: target contract address
// eventName: event name
// eventValue: value corresponding to event name. NOTE: if value has a prefix "0x", value will be decoded into []byte
func MakeEventSearchQuery(contractAddr common.Address, eventName, eventValue string) (string, error) {
	var parts []string
	parts = append(parts, fmt.Sprintf("%v='%v'", ContractAddressKey, contractAddr.Hex()))
	parts = append(parts, fmt.Sprintf("%v='%v'", ContractEventNameKey, eventName))
	if len(eventValue) > 0 {
		ev, err := MakeEntryBytes(eventName, eventValue)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("%v='%v'", ContractEventDataKey, string(ev)))
	}
	return strings.Join(parts, " AND "), nil
}

// GetContractEventsFromResultTx returns events that matches a given contract address from ResultTx.
func GetContractEventsFromResultTx(contractAddr common.Address, result *ctypes.ResultTx) ([]types.Event, error) {
	if !result.TxResult.IsOK() {
		return nil, errors.New("result has an error")
	}

	var events []types.Event
L:
	for _, ev := range result.TxResult.GetEvents() {
		for _, attr := range ev.GetAttributes() {
			if bytes.Equal([]byte(AddressKey), attr.GetKey()) {
				if bytes.Equal([]byte(contractAddr.Hex()), attr.GetValue()) {
					events = append(events, types.Event(ev))
				}
				continue L
			}
		}
	}
	return events, nil
}

// FilterContractEvents returns events that includes a given event name and value. (value is optional)
func FilterContractEvents(events []types.Event, eventName, eventValue string) ([]types.Event, error) {
	var value []byte
	var checkValue bool
	if len(eventValue) == 0 {
		checkValue = false
	} else {
		checkValue = true
		v, err := MakeEntryBytes(eventName, eventValue)
		if err != nil {
			return nil, err
		}
		value = v
	}

	var rets []types.Event
	for _, ev := range events {
		for _, attr := range ev.Attributes {
			if !checkValue {
				if bytes.Equal([]byte(EventNameKey), attr.GetKey()) {
					if bytes.Equal([]byte(eventName), attr.GetValue()) {
						rets = append(rets, ev)
					}
				}
			} else {
				if bytes.Equal([]byte(EventDataKey), attr.GetKey()) {
					if bytes.Equal(value, attr.GetValue()) {
						rets = append(rets, ev)
					}
				}
			}
		}
	}
	return rets, nil
}

func GetEntryFromEvent(ev types.Event) ([]*Entry, error) {
	var es []*Entry
	for _, attr := range ev.Attributes {
		if !bytes.Equal([]byte(EventDataKey), attr.GetKey()) {
			continue
		}
		e, err := ParseEntry(attr.GetValue())
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}
