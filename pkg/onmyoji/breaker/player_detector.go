package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"time"
)

// PlayerDetector
type PlayerDetector struct {
	ImgTemplates []onmyoji.ImgInfo // 模板图片信息
	window.Window
}

func newPlayerDetector(window window.Window, templates []onmyoji.ImgInfo) *PlayerDetector {
	logger.Info("创建玩家检测器，模板数量: %d", len(templates))
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
		logger.Info("开始检测玩家模板: %d ", string(img.Path))
		clicked, err := t.ClickAtTemplatePositionWithRandomOffset(img.Image, 0.8)
		if err != nil {
			logger.Error("点击玩家模板失败: %v", err)
			return err
		}

		if clicked {
			logger.Info("玩家检测成功，点击玩家模板: %d", img.Path)
			// 成功点击后，切换到攻击检测状态
			time.Sleep(500 * time.Millisecond) // 等待500毫秒以确保玩家界面加载完成
			logger.Info("玩家检测完成，切换到攻击检测状态")
			controller.Next(enums.BreakerAttack) // 切换到玩家检测完成状态
			break                                // 成功点击后跳出循环
		}

	}
	logger.Info("玩家检测失败，未找到匹配的玩家模板")

	return nil
}
