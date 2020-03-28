package fabric_client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/filter"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
)

// Query invoke query on chaincode
func (c *FabricClient) Query(request channel.Request, targets ...string) (channel.Response, error) {
	if err := c.createChannelClientIfNotExist(); err != nil {
		return channel.Response{}, err
	}
	if request.ChaincodeID != c.FabricClientConfig.ChaincodeID {
		return channel.Response{}, fmt.Errorf("ChaincodeID mismatch")
	}
	opts := c.queryOpts
	if targets != nil {
		opts = append(opts, channel.WithTargetEndpoints(targets...))
	}
	response, err := c.ChannelClient.Query(request, opts...)
	if err != nil {
		return channel.Response{}, errors.WithMessage(err, "failed to query chaincode")
	}
	return response, nil
}

// Execute invoke exec on chaincode
func (c *FabricClient) Execute(request channel.Request, targets ...string) (channel.Response, error) {
	if err := c.createChannelClientIfNotExist(); err != nil {
		return channel.Response{}, err
	}
	if request.ChaincodeID != c.FabricClientConfig.ChaincodeID {
		return channel.Response{}, fmt.Errorf("ChaincodeID mismatch")
	}
	opts := c.execOpts
	if targets != nil {
		opts = append(opts, channel.WithTargetEndpoints(targets...))
	}
	response, err := c.ChannelClient.Execute(request, opts...)
	if err != nil {
		return channel.Response{}, errors.WithMessage(err, "failed to execute chaincode")
	}
	return response, nil
}

// ExecuteWithNonce invoke exec on chaincode with target nonce
func (c *FabricClient) ExecuteWithNonce(request channel.Request, nonce []byte, targets ...string) (channel.Response, error) {
	if err := c.createChannelClientIfNotExist(); err != nil {
		return channel.Response{}, err
	}
	if request.ChaincodeID != c.FabricClientConfig.ChaincodeID {
		return channel.Response{}, fmt.Errorf("ChaincodeID mismatch")
	}
	invoker := specificExecuteHandler(&endorsementOpt{
		Nonce: nonce,
	})
	opts := c.asyncExecOpts
	if targets != nil {
		opts = append(opts, channel.WithTargetEndpoints(targets...))
	}
	response, err := c.ChannelClient.InvokeHandler(invoker, request, opts...)
	if err != nil {
		return channel.Response{}, errors.WithMessage(err, "failed to execute chaincode")
	}
	return response, nil
}

// ComputeTxId compute tx id
func (c *FabricClient) ComputeTxId() (txId string, nonce []byte, err error) {
	if err := c.createChannelClientIfNotExist(); err != nil {
		return "", nil, err
	}
	return c.transactor.createTxid()
}

func newAsyncExecOpts(channelContext context.Channel) []channel.RequestOption {
	return []channel.RequestOption{
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetFilter(filter.NewEndpointFilter(channelContext, filter.EndorsingPeer)),
		channel.WithTimeout(fab.Execute, channelContext.EndpointConfig().Timeout(fab.Execute)),
	}
}

func newAsyncQueryOpts(channelContext context.Channel) []channel.RequestOption {
	return []channel.RequestOption{
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetFilter(filter.NewEndpointFilter(channelContext, filter.ChaincodeQuery)),
		channel.WithTimeout(fab.Query, channelContext.EndpointConfig().Timeout(fab.Query)),
	}
}

func newExecOpts(channelContext context.Channel) []channel.RequestOption {
	return []channel.RequestOption{
		channel.WithRetry(retry.DefaultChannelOpts),
	}
}

func newQueryOpts(channelContext context.Channel) []channel.RequestOption {
	return []channel.RequestOption{
		channel.WithRetry(retry.DefaultChannelOpts),
	}
}
