package validator

import (
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto"
)

const (
	STATUS_NONE = iota
	STATUS_CANDIDATE
	STATUS_ACTIVE
	STATUS_REMOVED
)

type Validator struct {
	Address     common.Address
	PubKey      crypto.PubKey
	VotingPower uint64
	Status      uint8
	Commitment  []byte
}

func MakeValidatorFromTx(tx *transaction.ValidatorAddTx) *Validator {
	return &Validator{
		Address:     tx.From,
		PubKey:      tx.PubKey,
		Commitment:  tx.Commitment,
		VotingPower: tx.Amount,
		Status:      STATUS_CANDIDATE,
	}
}
