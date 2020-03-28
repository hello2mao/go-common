package fabric_client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type endorsementOpt struct {
	Nonce   []byte
	Creator []byte
}

//specificQueryHandler returns query handler with chain of ProposalProcessorHandler, EndorsementHandler, EndorsementValidationHandler and SignatureValidationHandler
func specificQueryHandler(endorsementOpt *endorsementOpt, next ...invoke.Handler) invoke.Handler {
	provider := func() []fab.TxnHeaderOpt {
		if endorsementOpt == nil {
			return nil
		}
		var TxnHeaderOpts []fab.TxnHeaderOpt
		TxnHeaderOpts = append(TxnHeaderOpts, fab.WithNonce(endorsementOpt.Nonce))
		TxnHeaderOpts = append(TxnHeaderOpts, fab.WithCreator(endorsementOpt.Creator))
		return TxnHeaderOpts
	}
	return invoke.NewProposalProcessorHandler(
		invoke.NewEndorsementHandlerWithOpts(
			invoke.NewEndorsementValidationHandler(
				invoke.NewSignatureValidationHandler(next...),
			),
			provider,
		),
	)
}

//specificExecuteHandler returns execute handler with chain of SelectAndEndorseHandler, EndorsementValidationHandler, SignatureValidationHandler and CommitHandler
func specificExecuteHandler(endorsementOpt *endorsementOpt, next ...invoke.Handler) invoke.Handler {
	provider := func() []fab.TxnHeaderOpt {
		if endorsementOpt == nil {
			return nil
		}
		var TxnHeaderOpts []fab.TxnHeaderOpt
		TxnHeaderOpts = append(TxnHeaderOpts, fab.WithNonce(endorsementOpt.Nonce))
		TxnHeaderOpts = append(TxnHeaderOpts, fab.WithCreator(endorsementOpt.Creator))
		return TxnHeaderOpts
	}
	endorsementHandler := invoke.NewEndorsementHandlerWithOpts(nil, provider)
	retHandler := invoke.NewSelectAndEndorseHandler(
		invoke.NewEndorsementValidationHandler(
			invoke.NewSignatureValidationHandler(invoke.NewCommitHandler(next...)),
		),
	)
	selectAndEndorseHandler := retHandler.(*invoke.SelectAndEndorseHandler)
	selectAndEndorseHandler.EndorsementHandler = endorsementHandler
	return retHandler
}
