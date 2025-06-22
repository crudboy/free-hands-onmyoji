package types

type Task interface {
	Execute(sm interface{}) error // 使用 interface{} 避免循环依赖
}
type TaskState int

const (
	Move      TaskState = iota // 是否移动过
	MoveCount                  // 移动次数
	Boss                       // 是否刷过boss
	Finish                     // 是否完成任务
)
