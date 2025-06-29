package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"time"
)

type RewardDetector struct {
	ImgTemplate []onmyoji.ImgInfo // 模板图片信息
	window.Window
}

func newRewardDetector(window window.Window, info []onmyoji.ImgInfo) *RewardDetector {
	logger.Info("创建奖励检测器，模板数量: %d", len(info))
	return &RewardDetector{
		ImgTemplate: info,
		Window:      window,
	}
}
func (t *RewardDetector) Name() enums.TaskType {
	return enums.BreakerReward
}
func (t *RewardDetector) Execute(controller statemachine.TaskController) error {
	for i := 0; i < len(t.ImgTemplate); i++ {
		logger.Info("检测奖励模板: %d", t.ImgTemplate[i].Path)
		clicked, err := t.ClickAtTemplatePositionWithOffset(t.ImgTemplate[i].Image, 0.8, 0, 100)
		time.Sleep(500 * time.Millisecond) // 等待奖励页面消失
		if err != nil {
			logger.Error("点击奖励模板失败: %v", err)
			return err
		}
		if clicked {
			logger.Info("奖励点击成功")
			// 成功点击后，切换到玩家检测状态
			return nil // 成功点击后返回
		}
	}
	logger.Info("奖励检测失败，未找到匹配的奖励模板开始检测玩家")
	controller.Next(enums.BreakerPlayer) // 切换到玩家检测状态
	return nil
}
