package enums

type TaskType string

const (
	TuiJing  TaskType = "Move"     // 移动
	XunGuai  TaskType = "XunGuai"  // 寻怪任务
	JieSuan  TaskType = "JieSuan"  // 匹配结算
	Boss     TaskType = "Boss"     // 刷boss任务
	BaoXiang TaskType = "BaoXiang" // 寻宝箱任务
	ZhangJie TaskType = "ZhangJie" // 章节任务
	JinRu    TaskType = "JinRu"    // 进入任务

)
