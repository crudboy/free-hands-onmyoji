package onmyoji

import (
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

type TaskRegistrator interface {
	// Registration RegisterTasks 注册所有特定模式的任务到状态机
	Registration(*statemachine.StateMachine, window.Window, Config, map[string]entity.ImgInfo) error

	// LoadImageTemplates 加载图片模板
	LoadImageTemplates() (map[string]entity.ImgInfo, error)
}
type Registrator struct {
	machine    *statemachine.StateMachine // 状态机
	config     Config                     // 任务配置
	windowInfo window.Window
	imgMap     map[string]entity.ImgInfo // 模板图片信息

}

func NewRegistrator(machine *statemachine.StateMachine, windowInfo window.Window, config Config) *Registrator {
	return &Registrator{
		machine:    machine,
		config:     config,
		windowInfo: windowInfo,
	}
}
func (r *Registrator) Registration(registrator TaskRegistrator) error {
	if registrator == nil {
		return nil
	}
	if r.imgMap == nil {
		var err error
		r.imgMap, err = registrator.LoadImageTemplates()
		if err != nil {
			return err
		}
	}
	return registrator.Registration(r.machine, r.windowInfo, r.config, r.imgMap)
}
func Registration(machine *statemachine.StateMachine, task statemachine.NamedTask) {
	logger.Info("注册任务: %s", task.Name())
	machine.AddTask(task)
}
