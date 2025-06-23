package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/utils"
	"image"
	"os"
	"path/filepath"
)

var templateMap = make(map[string]entity.ImgInfo)

func initImages() {
	logger.Info("加载困28任务模板图片")
	// 如果模板图片已经加载，则不需要重新加载
	if len(templateMap) > 0 {
		return
	}

	imgPath := "./k28/"
	files, err := os.ReadDir(imgPath)
	if err != nil {
		panic(fmt.Errorf("读取图片目录失败: %v", err))
	}
	for _, imgFile := range files {
		reader, err := os.Open(filepath.Join(imgPath, imgFile.Name()))
		if err != nil {
			continue
		}

		// 使用闭包确保每次迭代结束时关闭文件
		func(r *os.File) {
			defer r.Close()

			if r.Name() == filepath.Join(imgPath, ".DS_Store") {
				return
			}
			im, _, err := image.DecodeConfig(r)
			if err != nil {
				logger.Error("图片解码错误 %s: %v", imgFile.Name(), err)
				return // 跳过错误图片而不是终止程序
			}
			// 不要文件后缀
			fileName := imgFile.Name()[:len(imgFile.Name())-len(filepath.Ext(imgFile.Name()))]

			templateMap[fileName] = entity.ImgInfo{
				Path:    imgPath + imgFile.Name(),
				ImgMaxX: im.Width,
				ImgMaxY: im.Height,
				Image:   utils.ReadPic(imgPath + imgFile.Name()),
			}
		}(reader)
	}

	// 检查是否成功加载了所有必要的图片
	requiredImages := []string{
		string(enums.ZhangJie),
		string(enums.JinRu),
		string(enums.XunGuai),
		string(enums.JieSuan),
		string(enums.Boss),
		string(enums.BaoXiang),
	}

	for _, imgName := range requiredImages {
		if _, exists := templateMap[imgName]; !exists {
			panic(fmt.Errorf("模板图片 '%s' 未找到，请检查图片目录: %s", imgName, imgPath))
		}
	}
}
func Registration(machine *statemachine.StateMachine, info entity.WindowInfo, config onmyoji.Config) {
	initImages() // 初始化模板图片
	w := window.Window{
		WindowX: info.WindowX,
		WindowY: info.WindowY,
		WindowH: info.WindowH,
		WindowW: info.WindowW,
	}
	_registration(machine, newZhangJieTask(config, w, templateMap[string(enums.ZhangJie)]))
	_registration(machine, newJinRuTask(w, templateMap[string(enums.JinRu)]))
	_registration(machine, newXunGuaiTask(config, w, templateMap[string(enums.XunGuai)]))
	_registration(machine, newMoveTask(w))
	_registration(machine, newJieSuanTask(config, w, templateMap[string(enums.JieSuan)]))
	_registration(machine, newBossTask(w, templateMap[string(enums.Boss)]))
	_registration(machine, newBaoXiangTask(config, w, templateMap[string(enums.BaoXiang)]))
}

func _registration(machine *statemachine.StateMachine, task statemachine.NamedTask) {
	logger.Info("注册任务: %s", task.Name())
	machine.AddTask(task)
}
