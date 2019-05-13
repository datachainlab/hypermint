package errors

import "github.com/bluele/hypermint/pkg/abci/types"

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace types.CodespaceType = "cool"

	CodeInvalidSignature types.CodeType = 101
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code types.CodeType) string {
	switch code {
	case CodeInvalidSignature:
		return "invalid signature"
	default:
		return types.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidSignature(codespace types.CodespaceType, msg string) types.Error {
	return newError(codespace, CodeInvalidSignature, msg)
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code types.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace types.CodespaceType, code types.CodeType, msg string) types.Error {
	msg = msgOrDefaultMsg(msg, code)
	return types.NewError(codespace, code, msg)
}
