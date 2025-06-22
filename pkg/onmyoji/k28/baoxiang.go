package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"time"
)

type BaoXiang struct {
	TemplateImg   entity.ImgInfo // 模板图片信息
	window.Window                // 嵌入公共字段
}

func (t *BaoXiang) Name() enums.TaskType {
	return enums.BaoXiang
}

func (t *BaoXiang) Execute(controller statemachine.TaskController) error {
	logger.Info("第 %d 次执行", t.Count)
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.TemplateImg.Image, 0.8)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}

	if clicked {
		logger.Info("宝箱点击成功")
		time.Sleep(200 * time.Millisecond) // 等待点击操作完成
		// 点击地板
		t.ClickFloor(0, 50)
		time.Sleep(200 * time.Millisecond) // 等待地板点击完成
	} else {
		//宝箱寻找完成
		time.Sleep(500 * time.Millisecond) // 等待切换完成
		controller.Next(enums.JinRu)       // 切换到进入任务
		logger.Info("宝箱没有匹配到，切换到进入任务")
	}
	return nil

}
