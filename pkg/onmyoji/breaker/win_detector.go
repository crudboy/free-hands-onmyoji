package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"time"
)

// BreakerWinDetector 突破胜利检测器
type BreakerWinDetector struct {
	ImgTemplates entity.ImgInfo // 模板图片信息
	window.Window
}

func newBreakerWinDetector(window window.Window, template entity.ImgInfo) *BreakerWinDetector {
	return &BreakerWinDetector{
		ImgTemplates: template,
		Window:       window,
	}
}
func (t *BreakerWinDetector) Name() enums.TaskType {
	return enums.BreakerWin
}
func (t *BreakerWinDetector) Execute(controller statemachine.TaskController) error {
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplates.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		logger.Info("检测到突破胜利，开始奖励检测")
		time.Sleep(500 * time.Millisecond)   // 等待500毫秒以确保奖励界面加载完成
		controller.Next(enums.BreakerReward) //  检测一次奖励状态
	}
	return nil
}
