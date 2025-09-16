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
	success := make(chan struct{})

	for i := 0; i < len(t.ImgTemplate); i++ {
		go func(idx int) {
			logger.Info("检测奖励模板: %d", t.ImgTemplate[idx].Path)
			clicked, err := t.ClickAtTemplatePositionWithOffset(t.ImgTemplate[idx].Image, 0.8, 0, 100)
			if err == nil && clicked {
				select {
				case success <- struct{}{}:
				default:
				}
			}
		}(i)
	}

	select {
	case <-success:
		logger.Info("奖励点击成功")
		time.Sleep(2000 * time.Millisecond)
		controller.Next(enums.BreakerPlayer)
		return nil
	case <-time.After(2000 * time.Millisecond):
		logger.Info("奖励检测失败，未找到匹配的奖励模板开始检测玩家")
		controller.Next(enums.BreakerPlayer)
		return nil
	}
}
