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
	DuiSou      TaskType = "DuiSou"      // 寻找对手
	JiangLi     TaskType = "JiangLi"     // 领取奖励
	ToPoJieSuan TaskType = "ToPoJieSuan" // 突破结算
	ToPoJieRun  TaskType = "ToPoJieRun"  // 突破进入
	ToPoShiBai  TaskType = "ToPoShiBai"  // 突破失败
	ToPoShengLi TaskType = "ToPoShengLi" // 突破胜利

)
