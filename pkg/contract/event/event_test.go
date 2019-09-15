package event

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"
)

func TestEventData(t *testing.T) {
	var cases = []struct {
		e        *Entry
		hasError bool
	}{
		{&Entry{common.RandBytes(32), common.RandBytes(1024)}, false},
		{&Entry{common.RandBytes(33), common.RandBytes(1024)}, true},
		{&Entry{common.RandBytes(32), common.RandBytes(1025)}, true},
		{&Entry{nil, common.RandBytes(1024)}, true},
		{&Entry{common.RandBytes(32), nil}, true},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)

			err := cs.e.Validate()
			if cs.hasError {
				assert.Error(err)
				return
			}
			assert.NoError(err)

			e, err := ParseEntry(cs.e.Bytes())
			assert.NoError(err)
			assert.EqualValues(cs.e, e)
		})
	}
}

func TestMakeEntry(t *testing.T) {
	assert := assert.New(t)
	name, value := "name", "value"
	e1, err := MakeEntry(name, value)
	assert.NoError(err)

	hv := hex.EncodeToString([]byte(value))
	e2, err := MakeEntry(name, "0x"+hv)
	assert.NoError(err)
	assert.Equal(e1, e2)
}
