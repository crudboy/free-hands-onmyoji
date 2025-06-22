package entity

import "image"

type ImgInfo struct {
	Path    string
	ImgMaxX int
	ImgMaxY int
	Image   image.Image
}
type WindowInfo struct {
	WindowX int // 程序窗口在屏幕上的起始位置（偏移量）
	WindowY int // 程序窗口在屏幕上的起始位置（偏移量）
	WindowH int // 截图区域的高度
	WindowW int // 截图区域的宽度

}
