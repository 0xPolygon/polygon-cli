package claim

import (
	"errors"
)

var (
	ErrNotReadyForClaim      = errors.New("the claim transaction is not yet ready to be claimed, try again in a few blocks")
	ErrDepositAlreadyClaimed = errors.New("the claim transaction has already been claimed")
)
