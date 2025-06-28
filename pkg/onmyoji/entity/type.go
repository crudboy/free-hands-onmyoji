package entity

import "image"

type ImgInfo struct {
	Path    string
	ImgMaxX int
	ImgMaxY int
	Image   image.Image
}
