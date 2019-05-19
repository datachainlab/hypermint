package contract

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bluele/hypermint/pkg/logger"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestKeccak256(t *testing.T) {
	var cases = []struct {
		msg    []byte
		expect string
	}{
		{[]byte(""), "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},
		{[]byte("1"), "c89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6"},
		{[]byte("2"), "ad7c5bef027816a800da1736444fb58a807ef4c9603b7848673f7e3a68eb14a5"},
		{[]byte("3"), "2a80e1ef1d7842f27f2e6be0972bb708b9a135c38860dbe73c27c3486c34f4de"},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)

			e, err := hex.DecodeString(cs.expect)
			assert.NoError(err)
			b, err := util.Keccak256(cs.msg)
			assert.NoError(err)
			assert.Equal(e, b)

			buf := bytes.NewBuffer(nil)
			buf.Write(cs.msg)
			buf.Write(make([]byte, 32))
			mem := buf.Bytes()

			ps := newMockProcess()
			r := NewReader(mem, 0, int64(len(cs.msg)))
			w := NewWriter(mem, int64(len(cs.msg)), 32)
			assert.Equal(32, Keccak256(ps, r, w))
			assert.Equal(e, w.(Reader).Read())
		})
	}
}

type mockProcess struct {
	Process
}

func (p *mockProcess) Logger() logger.Logger {
	if testing.Verbose() {
		return logger.GetDefaultLogger("*:debug")
	}
	return logger.GetDefaultLogger("*:none")
}

func newMockProcess() Process {
	return &mockProcess{}
}
