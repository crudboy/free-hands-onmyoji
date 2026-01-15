package window

// GetWindowPosition Windows 实现：获取特定应用窗口的位置和大小信息
func (w *MacPlatform) GetWindowPosition() (Window, error) {
	return Window{}, nil
}

// GetWindowPositionOnSecondDisplay Windows 实现：获取第二显示器上的窗口位置
func (w *MacPlatform) GetWindowPositionOnSecondDisplay(displayId int) (Window, error) {
	return Window{}, nil
}

// SetScaleWindow Windows 实现：设置窗口大小
func (w *MacPlatform) SetScaleWindow(width, height int) {}

// CloseApplication Windows 实现：关闭应用
func (w *MacPlatform) CloseApplication() error {
	return nil
}
