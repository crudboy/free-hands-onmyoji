package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

// PlayerDetector 玩家检测器
type PlayerDetector struct {
	ImgTemplates []entity.ImgInfo // 模板图片信息
	window.Window
}

func newPlayerDetector(window window.Window, templates []entity.ImgInfo) *PlayerDetector {
	return &PlayerDetector{
		ImgTemplates: templates,
		Window:       window,
	}
}
func (t *PlayerDetector) Name() enums.TaskType {
	return enums.BreakerPlayer
}
func (t *PlayerDetector) Execute(controller statemachine.TaskController) error {
	for _, img := range t.ImgTemplates {
		clicked, err := t.ClickAtTemplatePositionWithRandomOffset(img.Image, 0.8)
		if err != nil {
			return err
		}

		if clicked {
			controller.Next(enums.BreakerAttack) // 切换到玩家检测完成状态
		} else {
			logger.Info("玩家检测未匹配到，继续尝试下一个模板")
		}
	}
	return nil
}
