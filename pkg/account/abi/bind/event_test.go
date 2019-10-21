package bind

import (
	"encoding/json"
	"github.com/bluele/hypermint/pkg/contract/event"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

type TransferRaw struct {
	From    Address `json:"from"`
	To      Address `json:"to"`
	TokenID U64     `json:"tokenID"`
}

type Transfer struct {
	From    common.Address
	To      common.Address
	TokenID uint64
}

func (_Transfer *Transfer) Decode(bs []byte) error {
	var raw TransferRaw
	if err := json.Unmarshal(bs, &raw); err != nil {
		return err
	}
	if err := DeepCopy(_Transfer, &raw); err != nil {
		return err
	}
	return nil
}

var TransferInfo = EventInfo {
	ID: "Transfer",
	EventCreator: func() Event {
		return &Transfer{}
	},
}

func TestEventDecoder_Decode(t *testing.T) {
	d := NewEventDecoder()
	d.Register(TransferInfo)

	var transfer Transfer
	e := event.Entry {
		Name: []byte("Transfer"),
		Value: []byte(`{"from":[147,48,98,172,145,177,43,39,168,143,13,143,0,190,158,208,81,63,13,20],"to":[184,50,15,186,149,14,223,23,8,253,230,160,55,12,232,148,209,89,99,57],"tokenID":18}`),
	}
	if err := d.Decode(&transfer, &e); err != nil {
		t.Error(err)
	}
	if transfer.TokenID != 18 {
		t.Fail()
	}
}
