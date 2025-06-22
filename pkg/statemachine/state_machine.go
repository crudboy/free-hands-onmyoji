package statemachine

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/types"
)

// TaskController 接口暴露任务执行时需要的控制方法和任务间信息传递方法
type TaskController interface {
	Next(key enums.TaskType) error                              // 切换到指定名称的任务
	NextIndex()                                                 // 按索引切换到下一个任务
	GetAttribute(name types.TaskState) (interface{}, error)     // 获取任务属性，用于任务间信息传递
	SetAttribute(name types.TaskState, value interface{}) error // 设置任务属性，用于任务间信息传递
	ClearAttributes()                                           // 清除所有任务属性
}

// NamedTask 接口扩展了基本的Task接口，添加了名称功能
type NamedTask interface {
	Name() enums.TaskType
	Execute(controller TaskController) error
}

type Task interface {
	Execute(controller TaskController) error
}

type StateMachine struct {
	tasks      []NamedTask
	current    int
	attributes map[types.TaskState]interface{} // 用于存储任务间共享的信息
}

func NewStateMachine() *StateMachine {

	return &StateMachine{
		tasks:      []NamedTask{},
		current:    0,
		attributes: make(map[types.TaskState]interface{}), // 初始化属性映射
	}
}

// Next 根据任务名称切换到指定任务
func (sm *StateMachine) Next(key enums.TaskType) error {
	for i, task := range sm.tasks {
		if task.Name() == key {
			sm.current = i
			logger.Debug("切换到任务: %s (索引: %d)", key, i)
			return nil
		}
	}
	panic(fmt.Sprintf("未找到名为 '%s' 的任务", key))
}

// NextIndex 按索引切换到下一个任务
func (sm *StateMachine) NextIndex() {
	if len(sm.tasks) == 0 {
		return
	}
	sm.current = (sm.current + 1) % len(sm.tasks)
	logger.Debug("切换到任务: %s (索引: %d)", sm.tasks[sm.current].Name(), sm.current)
}

func (sm *StateMachine) Run() error {
	if len(sm.tasks) == 0 {
		return nil
	}
	logger.Debug("当前任务: %s (索引: %d)", sm.tasks[sm.current].Name(), sm.current)
	return sm.tasks[sm.current].Execute(sm)
}

func (sm *StateMachine) AddTask(task NamedTask) {
	sm.tasks = append(sm.tasks, task)
}

// GetCurrentTask 获取当前任务
func (sm *StateMachine) GetCurrentTask() NamedTask {
	if len(sm.tasks) == 0 {
		return nil
	}
	return sm.tasks[sm.current]
}
func (sm *StateMachine) ClearAttributes() {
	if sm.attributes != nil {
		sm.attributes = make(map[types.TaskState]interface{})
	}

}

// GetTaskByName 根据名称获取任务
func (sm *StateMachine) GetTaskByName(name enums.TaskType) (NamedTask, error) {
	for _, task := range sm.tasks {
		if task.Name() == name {
			return task, nil
		}
	}
	return nil, fmt.Errorf("未找到名为 '%s' 的任务", name)
}

// GetAttribute 获取属性值，用于任务间信息传递
func (sm *StateMachine) GetAttribute(name types.TaskState) (interface{}, error) {
	if sm.attributes == nil {
		return nil, fmt.Errorf("属性 '%v' 不存在", name)
	}

	value, exists := sm.attributes[name]
	if !exists {
		return nil, fmt.Errorf("属性 '%v' 不存在", name)
	}

	return value, nil
}

// SetAttribute 设置属性值，用于任务间信息传递
func (sm *StateMachine) SetAttribute(name types.TaskState, value interface{}) error {
	if sm.attributes == nil {
		sm.attributes = make(map[types.TaskState]interface{})
	}

	sm.attributes[name] = value
	return nil
}
