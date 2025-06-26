package utils

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"gocv.io/x/gocv"
)

func TestFindImg(t *testing.T) {
	//ab x:449, y:242, num:0.99 //没打马赛克
	// bb x:449, y:242, num:0.992887 //打了马赛克

	img := ReadPic("/Users/crudboy/ab.jpg")
	template := ReadPic("/Users/crudboy/cd.jpg")
	x, y, num := FindTempPos(template, img)
	t.Logf("x:%d, y:%d, num:%f", x, y, num)

}
func TestFindImgWithMask(t *testing.T) {
	// 加载主图像和模板图像
	src := gocv.IMRead("/Users/crudboy/Debug.jpg", gocv.IMReadGrayScale)                     // 主图
	tmpl := gocv.IMRead("/Users/crudboy/Downloads/20250625224813.jpg", gocv.IMReadGrayScale) // 模板图

	if src.Empty() || tmpl.Empty() {
		fmt.Println("无法加载图像，请检查路径是否正确")
		return
	}

	// 创建结果矩阵 (rows = H - h + 1, cols = W - w + 1)
	result := gocv.NewMat()
	defer result.Close()

	// 执行模板匹配
	gocv.MatchTemplate(src, tmpl, &result, gocv.TmCcoeffNormed, gocv.NewMat())
	// 获取匹配结果的最大值、最小值以及对应的位置
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	fmt.Printf("最大匹配值: %.2f, 位置: %v\n", maxVal, maxLoc)

	// 设置匹配阈值（根据实际情况调整，0.8 表示较高相似度）
	threshold := float32(0.8)
	if maxVal < threshold {
		fmt.Println("未找到匹配区域或匹配度不足")
		return
	}

	// 如果匹配成功，绘制矩形标记模板区域
	tmplW := tmpl.Cols()
	tmplH := tmpl.Rows()
	matchRect := image.Rect(maxLoc.X, maxLoc.Y, maxLoc.X+tmplW, maxLoc.Y+tmplH)
	dst := gocv.IMRead("/Users/crudboy/Debug.jpg", gocv.IMReadColor) // 重新读取彩色图用于显示
	gocv.Rectangle(&dst, matchRect, color.RGBA{0, 255, 0, 255}, 2)

	// 保存结果图像
	gocv.IMWrite("output.jpg", dst)

	fmt.Printf("模板图像位于主图中的位置: (%d, %d)\n", maxLoc.X, maxLoc.Y)
	//ab x:449, y:242, num:0.99 //没打
}
