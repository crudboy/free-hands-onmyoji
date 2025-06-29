package general

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/utils"
	"image"
	"os"
	"path/filepath"
)

type Registrator struct {
	Path string
}

func (r Registrator) Registration(machine *statemachine.StateMachine, w window.Window, config onmyoji.Config, imgMap map[string]onmyoji.ImgInfo) error {
	onmyoji.Registration(machine, newChallengeDetector(w, imgMap[string(enums.Challenge)], config))
	onmyoji.Registration(machine, newLevelCompletionDetector(w, imgMap[string(enums.LevelCompletion)], config))
	onmyoji.Registration(machine, newLevelCompletionPart2Detector(w, imgMap[string(enums.LevelCompletionPart2)], config))
	return nil
}
func (r Registrator) LoadImageTemplates() (map[string]onmyoji.ImgInfo, error) {

	logger.Info("加载困28任务模板图片")

	// 初始化模板图片map
	imgMap := make(map[string]onmyoji.ImgInfo)

	imgPath := r.Path
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
		func(file *os.File) {
			defer file.Close()

			if file.Name() == filepath.Join(imgPath, ".DS_Store") {
				return
			}
			im, _, err := image.DecodeConfig(file)
			if err != nil {
				logger.Error("图片解码错误 %s: %v", imgFile.Name(), err)
				return // 跳过错误图片而不是终止程序
			}
			// 不要文件后缀
			fileName := imgFile.Name()[:len(imgFile.Name())-len(filepath.Ext(imgFile.Name()))]

			imgMap[fileName] = onmyoji.ImgInfo{
				Path:    imgPath + imgFile.Name(),
				ImgMaxX: im.Width,
				ImgMaxY: im.Height,
				Image:   utils.ReadPic(imgPath + imgFile.Name()),
			}
		}(reader)
	}

	// 检查是否成功加载了所有必要的图片
	requiredImages := []string{
		string(enums.Challenge),
		string(enums.LevelCompletion),
		string(enums.LevelCompletionPart2),
	}

	for _, imgName := range requiredImages {
		if _, exists := imgMap[imgName]; !exists {
			panic(fmt.Errorf("模板图片 '%s' 未找到，请检查图片目录: %s", imgName, imgPath))
		}
	}
	return imgMap, nil
}
