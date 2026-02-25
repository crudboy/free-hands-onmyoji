package k28

import (
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
)

// CompletionDetector 结算按钮检测任务
type CompletionDetector struct {
	ImgTemplate   onmyoji.ImgInfo
	window.Window // 嵌入公共字段
}

func newCompletionDetectorTask(window window.Window, info onmyoji.ImgInfo) *CompletionDetector {
	return &CompletionDetector{
		ImgTemplate: info,
		Window:      window,
	}
}
func (t *CompletionDetector) Name() tasks.TaskType {
	return tasks.Win
}
func (t *CompletionDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并点击
	clicked, err := t.ClickAtTemplatePosition(t.ImgTemplate.Image, 0.8)

	if err != nil {
		return err
	}

	if clicked {
		logger.Info("检测胜利界面点击成功，切换到结算任务")
		controller.Next(tasks.JieSuan)
		return nil
	}
	logger.Info("未找到胜利界面或相似度不足，跳过点击操作")
	return nil

}
