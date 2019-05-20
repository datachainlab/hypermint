package contract

import (
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

			b := makeEventData(cs.ev)
			ev, err := ParseEventData(b)
			assert.NoError(err)
			assert.EqualValues(cs.ev, ev)
		})
	}
}
