package fabric_client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
)

func (c *FabricClient) createLedgerClientIfNotExist() error {
	if c.LedgerClient != nil {
		return nil
	}
	// createChannelClient because of channelProvider
	if err := c.createChannelClientIfNotExist(); err != nil {
		return err
	}
	// ledgerClient
	ledgerClient, err := ledger.New(c.channelProvider)
	if err != nil {
		return errors.WithMessage(err, "failed to create new ledger client")
	}
	c.LedgerClient = ledgerClient

	return nil
}

// GetBlockNumberByTxID get block number by txId
func (c *FabricClient) GetBlockNumberByTxID(txID string) (number uint64, err error) {
	block, err := c.LedgerClient.QueryBlockByTxID(fab.TransactionID(txID))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return block.Header.Number, nil
}
