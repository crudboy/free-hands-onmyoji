package k28

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/types"
	"time"
)

type TuiJing struct {
	window.Window // 嵌入公共字段
}

func (t *TuiJing) Name() enums.TaskType {
	return enums.TuiJing
}

func (t *TuiJing) Execute(controller statemachine.TaskController) error {
	t.Count++ // 增加执行次数
	// 等待一段时间，让角色移动
	moveDelay := 800
	attribute, _ := controller.GetAttribute(types.MoveCount)
	// 等待移动完成
	time.Sleep(time.Duration(moveDelay) * time.Millisecond)
	_, _ = t.ClickFloor(165, 0)
	// 如果设置了最大移动次数并且已达到，则切换到下一个任务
	if attribute.(int) > 0 && t.Count >= attribute.(int) {
		logger.Info("已达到最大移动次数 %d，准备切换任务", attribute.(int))
		err := controller.SetAttribute(types.Move, true) // 设置移动完成标志
		if err != nil {
			logger.Error("设置移动完成标志失败: %v", err)
		} else {
			logger.Debug("成功设置移动完成标志 types.Move = true")
		}
		time.Sleep(300 * time.Millisecond) // 移动结束
		logger.Info("移动结束，切换到寻怪任务")
		controller.Next(enums.XunGuai) // 切换到下一个任务
	}

	return nil

}
