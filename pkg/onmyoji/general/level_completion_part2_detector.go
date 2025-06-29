package general

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

type LevelCompletionPart2Detector struct {
	ImgTemplate onmyoji.ImgInfo // 模板图片信息
	window.Window
	conf onmyoji.BreakerConfig
}

func newLevelCompletionPart2Detector(window window.Window, info onmyoji.ImgInfo, config onmyoji.Config) *LevelCompletionPart2Detector {
	return &LevelCompletionPart2Detector{
		ImgTemplate: info,
		Window:      window,
		conf:        config.Breaker,
	}
}
func (t *LevelCompletionPart2Detector) Name() enums.TaskType {
	return enums.LevelCompletionPart2
}
func (t *LevelCompletionPart2Detector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	logger.Info("开始执行关卡完成检测任务，使用模板: %s", t.ImgTemplate.Path)
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		logger.Info("攻击关卡通关检测成功")

	}
	controller.Next(enums.Challenge) // 切换到任务完成状态
	return nil
}
