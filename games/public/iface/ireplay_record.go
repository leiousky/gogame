package iface

import "games/public/define"

/// <summary>
/// 游戏记录接口
/// </summary>
type IReplayRecord interface {
	SaveReplay(record define.ReplayRecord)
}
