package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

type AttackDetector struct {
	ImgTemplate entity.ImgInfo // 模板图片信息
	window.Window
	conf onmyoji.BreakerConfig
}

func newAttackDetector(window window.Window, info entity.ImgInfo, config onmyoji.Config) *AttackDetector {
	return &AttackDetector{
		ImgTemplate: info,
		Window:      window,
		conf:        config.Breaker,
	}
}
func (t *AttackDetector) Name() enums.TaskType {
	return enums.BreakerAttack
}
func (t *AttackDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
	if err != nil {
		return err
	}

	if clicked {
		logger.Info("攻击检测成功，开始执行攻击操作")
		controller.Next(enums.BreakerWin) // 切换到任务完成状态

	}
	return nil
}
