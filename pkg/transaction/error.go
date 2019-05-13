package transaction

import (
	"github.com/bluele/hypermint/pkg/abci/types"
)

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace types.CodespaceType = "2"

	CodeInvalidTx       types.CodeType = 101
	CodeInvalidTransfer types.CodeType = 102
	CodeFailTransfer    types.CodeType = 103
	CodeInvalidDeploy   types.CodeType = 104
	CodeInvalidCall     types.CodeType = 105
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code types.CodeType) string {
	switch code {
	case CodeInvalidTransfer:
		return "invalid transfer"
	default:
		return types.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidTx(codespace types.CodespaceType, msg string) types.Error {
	return newError(codespace, CodeInvalidTx, msg)
}

func ErrInvalidTransfer(codespace types.CodespaceType, msg string) types.Error {
	return newError(codespace, CodeInvalidTransfer, msg)
}

func ErrFailTransfer(codespace types.CodespaceType, msg string) types.Error {
	return newError(codespace, CodeFailTransfer, msg)
}

func ErrInvalidDeploy(codespace types.CodespaceType, msg string) types.Error {
	return newError(codespace, CodeInvalidDeploy, msg)
}

func ErrInvalidCall(codespace types.CodespaceType, msg string) types.Error {
	return newError(codespace, CodeInvalidCall, msg)
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
