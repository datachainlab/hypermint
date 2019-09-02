package contract

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeArgs(t *testing.T) {
	var cases = []struct {
		bs      []byte
		expects []string
	}{
		{
			[]byte{0, 0, 0, 0},
			[]string{},
		},
		{
			[]byte{0, 0, 0, 1, 0, 0, 0, 1, 49},
			[]string{"1"},
		},
		{
			[]byte{0, 0, 0, 2, 0, 0, 0, 1, 49, 0, 0, 0, 1, 49},
			[]string{"1", "1"},
		},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			args, err := DeserializeArgs(cs.bs)
			assert.NoError(t, err)
			assert.Equal(t, len(cs.expects), args.Len())
			var as = []string{}
			for i := 0; i < args.Len(); i++ {
				arg, ok := args.Get(i)
				assert.True(t, ok)
				as = append(as, string(arg))
			}
			assert.Equal(t, cs.expects, as)
		})
	}
}
