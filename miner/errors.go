package miner

import "errors"

var (
	errDecodeFailed               = errors.New("fail to decode worker message")
	ErrInvalidSigner              = errors.New("message not signed by the sender")
	ErrUnauthorizedAddress        = errors.New("unauthorized address")
	ErrFailedDecodeOnlineProof    = errors.New("failed to decode online-proof message")
	ErrFailedDecodeOnlineQuestion = errors.New("failed to decode online-question message")
	ErrInvalidHeight              = errors.New("height inconsistency error")
	ErrInvalidProposer            = errors.New("err Not the proposer of this height")
	ErrInvalidProof               = errors.New("err Invalid proof")
	ErrInvalidValidator           = errors.New("err Not the validator of this height")
)
