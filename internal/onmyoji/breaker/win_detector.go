package breaker

import (
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
	"time"
)

// BreakerWinDetector 突破胜利检测器
type BreakerWinDetector struct {
	ImgTemplates onmyoji.ImgInfo // 模板图片信息
	window.Window
}

func newBreakerWinDetector(window window.Window, template onmyoji.ImgInfo) *BreakerWinDetector {
	return &BreakerWinDetector{
		ImgTemplates: template,
		Window:       window,
	}
}
func (t *BreakerWinDetector) Name() tasks.TaskType {
	return tasks.BreakerWin
}
func (t *BreakerWinDetector) Execute(controller statemachine.TaskController) error {
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplates.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		logger.Info("检测到突破胜利，检测是否有奖励")
		time.Sleep(1500 * time.Millisecond)  // 等待奖励出现
		controller.Next(tasks.BreakerReward) //  检测一次奖励状态
		return nil
	}
	logger.Info("未检测到突破胜利， 尝试检测是否失败")
	controller.Next(tasks.BreakerFail) // 切换到失败状态
	return nil
}
