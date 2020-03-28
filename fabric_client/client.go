package fabric_client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	feed "github.com/hello2mao/go-common/event"
)

// FabricClient defines the fabric sdk fabric-client
// NOTICE: FabricClient is just for only one org-channel-user, can not reuse.
type FabricClient struct {
	FabricClientConfig *FabricClientConfig
	Sdk                *fabsdk.FabricSDK // Fabric SDK 实例

	// all kinds of client
	ResMgmtClient *resmgmt.Client // 资源管理客户端
	ChannelClient *channel.Client // 通道客户端
	EventClient   *event.Client   // 事件客户端
	MspClient     *msp.Client     // MSP客户端
	LedgerClient  *ledger.Client

	// internal param
	channelProvider contextApi.ChannelProvider
	transactor      *transactor
	asyncExecOpts   []channel.RequestOption
	asyncQueryOpts  []channel.RequestOption
	execOpts        []channel.RequestOption
	queryOpts       []channel.RequestOption

	// async feed
	receiptFeed feed.Feed
	scope       feed.SubscriptionScope
}

// NewFabricClient returns the FabricClient instance
func NewFabricClient(fabricClientConfig *FabricClientConfig) (*FabricClient, error) {
	if fabricClientConfig == nil {
		return nil, fmt.Errorf("FabricClientConfig must not nil")
	}
	if fabricClientConfig.ConfigFile == "" {
		return nil, fmt.Errorf("FabricClientConfig ConfigFile must not empty")
	}
	// check config file
	if !FileOrDirExist(fabricClientConfig.ConfigFile) {
		return nil, fmt.Errorf("FabricClientConfig ConfigFile not exist")
	}

	fabricClient := &FabricClient{FabricClientConfig: fabricClientConfig}

	// Sdk
	sdk, err := fabsdk.New(config.FromFile(fabricClientConfig.ConfigFile))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create fabric SDK")
	}
	fabricClient.Sdk = sdk

	log.Infof("NewFabricClient with config: %+v", fabricClientConfig)
	return fabricClient, err
}

// Close close sdk
func (c *FabricClient) Close() {
	c.Sdk.Close()
}

func (c *FabricClient) createResMgmtClientIfNotExist() error {
	if c.ResMgmtClient != nil {
		return nil
	}
	if c.FabricClientConfig.OrgAdmin == "" || c.FabricClientConfig.OrgName != "" {
		return fmt.Errorf("OrgAdmin or OrgName must not empty")
	}
	resCliProvider := c.Sdk.Context(fabsdk.WithUser(c.FabricClientConfig.OrgAdmin), fabsdk.WithOrg(c.FabricClientConfig.OrgName))
	ResMgmtClient, err := resmgmt.New(resCliProvider)
	if err != nil {
		return errors.WithMessage(err, "failed to create fabric ResMgmtCli")
	}
	c.ResMgmtClient = ResMgmtClient
	return nil
}
