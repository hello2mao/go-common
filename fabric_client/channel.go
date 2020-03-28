package fabric_client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

func (c *FabricClient) createChannelClientIfNotExist() error {
	if c.ChannelClient != nil {
		return nil
	}
	if c.FabricClientConfig.ChannelID == "" || c.FabricClientConfig.OrgName == "" || c.FabricClientConfig.UserName == "" {
		return fmt.Errorf("OrgAdmin or UserName or ChannelID must not empty")
	}
	clientContext := c.Sdk.ChannelContext(c.FabricClientConfig.ChannelID, fabsdk.WithUser(c.FabricClientConfig.UserName), fabsdk.WithOrg(c.FabricClientConfig.OrgName))
	channelClient, err := channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel fabric-client")
	}
	c.ChannelClient = channelClient

	channelContext, err := clientContext()
	if err != nil {
		return errors.WithMessage(err, "failed to create channel context")
	}
	c.transactor = newTransactor(channelContext, c.FabricClientConfig.EnableGM)
	c.asyncExecOpts = newAsyncExecOpts(channelContext)
	c.asyncQueryOpts = newAsyncQueryOpts(channelContext)
	c.execOpts = newExecOpts(channelContext)
	c.queryOpts = newQueryOpts(channelContext)

	return nil
}
