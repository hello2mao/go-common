package fabric_client

import (
	"strconv"
	"testing"
	"time"
)

func TestRegisterAndEnrolledUser(t *testing.T) {
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
	newUserName := strconv.FormatInt(time.Now().Unix(), 10)
	err = fc.RegisterAndEnrollUser(newUserName, testSDKPath, testKeyStorePath)
	if err != nil {
		t.Errorf("RegisterAndEnrolledUser err: %v", err)
		return
	}
}
