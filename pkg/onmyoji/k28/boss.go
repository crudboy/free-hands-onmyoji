package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/types"
)

type Boss struct {
	TemplateImg   entity.ImgInfo // 模板图片信息
	window.Window                // 嵌入公共字段
}

func (t *Boss) Name() enums.TaskType {
	return enums.Boss
}
func (t *Boss) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.TemplateImg.Image, 0.8)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}

	if clicked {
		controller.SetAttribute(types.Boss, true) // 设置点击完成标志
		controller.Next(enums.JieSuan)            // 切换到刷怪任务
	} else {
		logger.Warn("模板匹配相似度过低，跳过点击操作")
		controller.Next(enums.XunGuai) // 切换到寻怪任务
		logger.Info("Boss没有匹配到切换到寻怪任务")
	}
	return nil

}
