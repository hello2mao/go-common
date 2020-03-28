package fabric_client

import (
	"testing"
)

const (
	testConfigFile   = "./test/config/config.yaml"
	testSDKPath      = "./test/sdk"
	testKeyStorePath = "/tmp/msp/keystore"
	testChannelID    = "scf-tongdao"
	testChaincodeID  = "test-cc"
	testUserName     = "Admin"
	testOrgName      = "bank"
)

func TestNewFabricClient(t *testing.T) {
	_, err := NewFabricClient(&FabricClientConfig{
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
}
