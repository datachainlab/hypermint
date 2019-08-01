package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	type V = Version
	var cases = []struct {
		version Version
	}{
		{V{0, 0}},
		{V{1, 1}},
		{V{2, 2}},
		{V{1, 2}},
		{V{2, 1}},
	}
	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)

			b := cs.version.Bytes()
			v, err := MakeVersion(b)
			assert.NoError(err)
			assert.Equal(cs.version, v)
		})
	}
}

func TestValueObject(t *testing.T) {
	type V = Version
	var cases = []struct {
		value   []byte
		version Version
	}{
		{[]byte{}, V{0, 0}},
		{[]byte{}, V{1, 1}},
		{[]byte("test"), V{2, 2}},
		{[]byte("10000"), V{2, 2}},
		{[]byte("9990"), V{2, 2}},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)

			vo := &ValueObject{Value: cs.value, Version: cs.version}
			b := vo.Marshal()
			nv := new(ValueObject)
			assert.NoError(nv.Unmarshal(b))
			assert.Equal(vo, nv)
		})
	}
}

func TestRWSet(t *testing.T) {
	type V = Version
	type ri struct {
		key     string
		version Version
	}
	type wi struct {
		key   string
		value string
	}
	type iri ri // ignore item
	type iwi wi // ignore item

	var cases = []struct {
		ops []interface{}
	}{
		{[]interface{}{ri{"z", V{1, 1}}, ri{"c", V{1, 1}}, ri{"b", V{1, 1}}, ri{"a", V{1, 1}}}},
		{[]interface{}{ri{"z", V{1, 1}}, ri{"a", V{1, 2}}, ri{"b", V{1, 3}}}},
		{[]interface{}{ri{"z", V{1, 1}}, wi{"a", "A"}, wi{"b", "B"}}},
		{[]interface{}{wi{"a", "A"}, iwi{"a", "A"}, ri{"a", V{1, 1}}, wi{"b", "B"}, iri{"a", V{1, 1}}}},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)
			rwm := NewRWSetMap()
			var rs []ri
			var ws []wi
			for _, op := range cs.ops {
				switch op := op.(type) {
				case ri:
					rwm.AddRead([]byte(op.key), op.version)
					rs = append(rs, op)
				case wi:
					rwm.AddWrite([]byte(op.key), []byte(op.value))
					ws = append(ws, op)
				case iri:
					rwm.AddRead([]byte(op.key), op.version)
				case iwi:
					rwm.AddWrite([]byte(op.key), []byte(op.value))
				default:
					t.Fatalf("unknown type %T", op)
				}
			}

			rws := rwm.ToSet()
			assert.Equal(len(rs), len(rws.ReadSet))
			for i, r := range rs {
				actual := rws.ReadSet[i]
				assert.Equal([]byte(r.key), actual.Key)
				assert.Equal(r.version, actual.Version)
			}
			assert.Equal(len(ws), len(rws.WriteSet))
			for i, w := range ws {
				actual := rws.WriteSet[i]
				assert.Equal([]byte(w.key), actual.Key)
				assert.Equal([]byte(w.value), actual.Value)
			}
		})
	}
}
