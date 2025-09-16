package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
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
func (t *BreakerWinDetector) Name() enums.TaskType {
	return enums.BreakerWin
}
func (t *BreakerWinDetector) Execute(controller statemachine.TaskController) error {
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplates.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		logger.Info("检测到突破胜利，检测是否有奖励")
		time.Sleep(1500 * time.Millisecond)  // 等待奖励出现
		controller.Next(enums.BreakerReward) //  检测一次奖励状态
		return nil
	}
	logger.Info("未检测到突破胜利， 尝试检测是否失败")
	controller.Next(enums.BreakerFail) // 切换到失败状态
	return nil
}
