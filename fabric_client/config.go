package fabric_client

// FabricClientConfig defines the fabric-client config
type FabricClientConfig struct {
	ConfigFile string // must: sdk的配置文件路径

	ChannelID   string // opt: 通道id
	ChaincodeID string // opt: 链码名称
	UserName    string // opt: 组织的普通用户
	OrgName     string // opt: 组织的名称
	OrgAdmin    string // opt: 组织的管理员用户

	EnableGM bool
}
