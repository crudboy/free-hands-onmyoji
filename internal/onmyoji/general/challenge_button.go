package general

import (
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
)

type ChallengeDetector struct {
	ImgTemplate onmyoji.ImgInfo // 模板图片信息
	window.Window
	conf onmyoji.BreakerConfig
}

func newChallengeDetector(window window.Window, info onmyoji.ImgInfo, config onmyoji.Config) *ChallengeDetector {
	return &ChallengeDetector{
		ImgTemplate: info,
		Window:      window,
		conf:        config.Breaker,
	}
}
func (t *ChallengeDetector) Name() tasks.TaskType {
	return tasks.Challenge
}
func (t *ChallengeDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	logger.Info("开始执行挑战按钮检测任务，使用模板: %s", t.ImgTemplate.Path)
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		logger.Info("挑战按钮点击成功")
	}
	controller.Next(tasks.LevelCompletion) // 切换到任务完成状态

	return nil
}
