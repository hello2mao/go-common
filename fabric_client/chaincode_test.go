package fabric_client

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func TestExecute(t *testing.T) {
	fc, err := NewFabricClient(&FabricClientConfig{
		ConfigFile:  testConfigFile,
		ChannelID:   testChannelID,
		ChaincodeID: testChaincodeID,
		UserName:    testUserName,
		OrgName:     testOrgName,
		OrgAdmin:    "Admin",
		EnableGM:    true,
	})
	if err != nil {
		t.Errorf("NewFabricClient err: %v", err)
		return
	}
	name, _ := GetRandomBytes(10)
	jsonContent, err := json.Marshal(
		struct {
			Name string `json:"name"`
		}{
			Name: string(name),
		})
	if err != nil {
		t.Errorf("json.Marshal err: %v\n", err)
		return
	}
	response, err := fc.Execute(channel.Request{
		ChaincodeID: testChaincodeID,
		Fcn:         "setKey",
		Args:        [][]byte{[]byte("testKey"), jsonContent},
	})
	if err != nil {
		t.Errorf("failed to Execute: %v\n", err)
		return
	}
	if response.ChaincodeStatus != 200 {
		t.Errorf("ChaincodeStatus != 200")
		return
	}
}

func TestComputeTxId(t *testing.T) {
	fc, err := NewFabricClient(&FabricClientConfig{
		ConfigFile:  testConfigFile,
		ChannelID:   testChannelID,
		ChaincodeID: testChaincodeID,
		UserName:    testUserName,
		OrgName:     testOrgName,
		OrgAdmin:    "Admin",
		EnableGM:    true,
	})
	if err != nil {
		t.Errorf("NewFabricClient err: %v", err)
		return
	}
	txId, nonce, err := fc.ComputeTxId()
	if err != nil {
		t.Errorf("ComputeTxId err: %v", err)
		return
	}
	t.Logf("txId: %s\n", txId)
	t.Logf("nonce: %s\n", hex.EncodeToString(nonce))

}

func TestExecuteWithNonce(t *testing.T) {
	fc, err := NewFabricClient(&FabricClientConfig{
		ConfigFile:  testConfigFile,
		ChannelID:   testChannelID,
		ChaincodeID: testChaincodeID,
		UserName:    testUserName,
		OrgName:     testOrgName,
		OrgAdmin:    "Admin",
		EnableGM:    true,
	})
	if err != nil {
		t.Errorf("NewFabricClient err: %v", err)
		return
	}
	name, _ := GetRandomBytes(10)
	jsonContent, err := json.Marshal(
		struct {
			Name string `json:"name"`
		}{
			Name: string(name),
		})
	if err != nil {
		t.Errorf("json.Marshal err: %v\n", err)
		return
	}
	txId, nonce, err := fc.ComputeTxId()
	if err != nil {
		t.Errorf("ComputeTxId err: %v", err)
		return
	}
	t.Logf("ComputeTxId: %s\n", txId)
	response, err := fc.ExecuteWithNonce(channel.Request{
		ChaincodeID: testChaincodeID,
		Fcn:         "setKey",
		Args:        [][]byte{[]byte("testKey"), jsonContent},
	}, nonce)
	if err != nil {
		t.Errorf("failed to Execute: %v\n", err)
		return
	}
	if response.ChaincodeStatus != 200 {
		t.Errorf("ChaincodeStatus != 200")
		return
	}
	t.Logf("ExecuteWithNonce TxId: %s\n", string(response.TransactionID))

	if txId != string(response.TransactionID) {
		t.Errorf("TransactionID not equal")
		return
	}
}

func TestQuery(t *testing.T) {
	fc, err := NewFabricClient(&FabricClientConfig{
		ConfigFile:  testConfigFile,
		ChannelID:   testChannelID,
		ChaincodeID: testChaincodeID,
		UserName:    testUserName,
		OrgName:     testOrgName,
		OrgAdmin:    "Admin",
		EnableGM:    true,
	})
	if err != nil {
		t.Errorf("NewFabricClient err: %v", err)
		return
	}
	response, err := fc.Query(channel.Request{
		ChaincodeID: testChaincodeID,
		Fcn:         "getKey",
		Args:        [][]byte{[]byte("testKey")},
	})
	if err != nil {
		t.Errorf("failed to Execute: %v\n", err)
		return
	}
	t.Logf("get value: %s\n", string(response.Payload))
}
