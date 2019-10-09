package bind

import (
	"bytes"
	"encoding/gob"
)

type Arg interface {
	Bytes() []byte
}

func Args(args ...Arg) [][]byte {
	var bss [][]byte
	for _, arg := range args {
		bss = append(bss, arg.Bytes())
	}
	return bss
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

