package fabric_client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

func (c *FabricClient) createEventClientIfNotExist() error {
	if c.EventClient != nil {
		return nil
	}
	if c.FabricClientConfig.ChannelID == "" || c.FabricClientConfig.OrgName == "" || c.FabricClientConfig.UserName == "" {
		return fmt.Errorf("OrgAdmin or UserName or ChannelID must not empty")
	}
	clientContext := c.Sdk.ChannelContext(c.FabricClientConfig.ChannelID, fabsdk.WithUser(c.FabricClientConfig.UserName), fabsdk.WithOrg(c.FabricClientConfig.OrgName))
	eventClient, err := event.New(clientContext, event.WithBlockEvents())
	if err != nil {
		return errors.WithMessage(err, "failed to create event fabric-client")
	}
	c.EventClient = eventClient
	return nil
}
