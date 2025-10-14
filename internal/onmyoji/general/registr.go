package general

import (
	"fmt"
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
	"free-hands-onmyoji/internal/utils"
	"image"
	"os"
	"path/filepath"
)

type Registrator struct {
	Path string
}

func (r Registrator) Registration(machine *statemachine.StateMachine, w window.Window, config onmyoji.Config, imgMap map[string]onmyoji.ImgInfo) error {
	onmyoji.Registration(machine, newChallengeDetector(w, imgMap[string(tasks.Challenge)], config))
	onmyoji.Registration(machine, newLevelCompletionDetector(w, imgMap[string(tasks.LevelCompletion)], config))
	onmyoji.Registration(machine, newLevelCompletionPart2Detector(w, imgMap[string(tasks.LevelCompletionPart2)], config))
	return nil
}
func (r Registrator) LoadImageTemplates() (map[string]onmyoji.ImgInfo, error) {

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
		string(tasks.Challenge),
		string(tasks.LevelCompletion),
		string(tasks.LevelCompletionPart2),
	}

	for _, imgName := range requiredImages {
		if _, exists := imgMap[imgName]; !exists {
			panic(fmt.Errorf("模板图片 '%s' 未找到，请检查图片目录: %s", imgName, imgPath))
		}
	}
	return imgMap, nil
}
