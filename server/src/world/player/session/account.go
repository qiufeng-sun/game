package session

// 帐号信息
type AccInfo struct {
	AccId          int // 帐号id
	CreateServerId int // 创建角色的服务器id
	Platform       int // 平台类型
	GmLv           int // gm等级
}

// 创建新的帐号信息
func NewAccInfo(accId, createServerId, platform, gmLv int) *AccInfo {
	return &AccInfo{accId, createServerId, platform, gmLv}
}
