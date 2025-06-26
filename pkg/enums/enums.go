package enums

type TaskType string

const (
	// ------------k28常量-------------------------
	Move     TaskType = "Move"     // 移动
	XunGuai  TaskType = "XunGuai"  // 寻怪任务
	JieSuan  TaskType = "JieSuan"  // 匹配结算
	Boss     TaskType = "Boss"     // 刷boss任务
	BaoXiang TaskType = "BaoXiang" // 寻宝箱任务
	ZhangJie TaskType = "ZhangJie" // 章节任务
	JinRu    TaskType = "JinRu"    // 进入任务
	// ------------突破常量-------------------------

	BreakerPlayer TaskType = "breaker_player_" // 寻找对手
	BreakerAttack TaskType = "breaker_attack"  // 进攻
	BreakerReward TaskType = "breaker_reward"  // 奖励
	BreakerFail   TaskType = "breaker_fail"    // 突破失败
	BreakerWin    TaskType = "breaker_win"     // 突破胜利

)
