package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"os"
)

type FailDetector struct {
	ImgTemplate onmyoji.ImgInfo // 模板图片信息
	window.Window
}

func newBreakerFailDetector(window window.Window, info onmyoji.ImgInfo) *FailDetector {
	return &FailDetector{
		ImgTemplate: info,
		Window:      window,
	}
}
func (t *FailDetector) Name() enums.TaskType {
	return enums.BreakerFail
}
func (t *FailDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		window.AlertNotify("突破失败", "检测到突破失败状态，请检查游戏状态或重试。")
		os.Exit(1) // 退出程序或执行其他失败处理逻辑
	}
	logger.Info("未检测到突破失败，检测是否成功")
	controller.Next(enums.BreakerWin) // 切换到成功状态

	return nil
}
