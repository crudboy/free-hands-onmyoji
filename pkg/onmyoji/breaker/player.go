package breaker

import (
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

// PlayerDetector 玩家检测器
type PlayerDetector struct {
	ImgTemplates []entity.ImgInfo // 模板图片信息
	window.Window
}

func newPlayerDetector(window window.Window, templateImg []entity.ImgInfo) *PlayerDetector {
	return &PlayerDetector{
		ImgTemplates: templateImg,
		Window:       window,
	}
}
func (t *PlayerDetector) Name() string {
	return "PlayerDetector"
}
func (t *PlayerDetector) Execute(controller statemachine.TaskController) error {
	for _, img := range t.ImgTemplates {
		clicked, err := t.ClickAtTemplatePositionWithRandomOffset(img.Image, 0.8)
		if err != nil {
			return err
		}

		if clicked {
			controller.Next("PlayerDetected") // 切换到玩家检测完成状态
		} else {
			controller.Next("PlayerNotDetected") // 切换到玩家未检测状态
		}
	}
	return nil
}
