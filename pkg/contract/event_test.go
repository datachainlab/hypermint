package contract

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"
)

func TestEventData(t *testing.T) {
	var cases = []struct {
		ev       *Event
		hasError bool
	}{
		{&Event{common.RandBytes(32), common.RandBytes(1024)}, false},
		{&Event{common.RandBytes(33), common.RandBytes(1024)}, true},
		{&Event{common.RandBytes(32), common.RandBytes(1025)}, true},
		{&Event{nil, common.RandBytes(1024)}, true},
		{&Event{common.RandBytes(32), nil}, true},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)

			err := validateEvent(cs.ev)
			if cs.hasError {
				assert.Error(err)
				return
			}
			assert.NoError(err)

			ev, err := ParseEventData(cs.ev.Bytes())
			assert.NoError(err)
			assert.EqualValues(cs.ev, ev)
		})
	}
}

func TestMakeEvent(t *testing.T) {
	assert := assert.New(t)
	name, value := "name", "value"
	ev1, err := MakeEvent(name, value)
	assert.NoError(err)

	hv := hex.EncodeToString([]byte(value))
	ev2, err := MakeEvent(name, "0x"+hv)
	assert.NoError(err)
	assert.Equal(ev1, ev2)
}
