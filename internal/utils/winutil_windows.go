//go:build windows
// +build windows

package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32dll               = syscall.NewLazyDLL("user32.dll")
	procEnumWindows         = user32dll.NewProc("EnumWindows")
	procGetWindowTextW      = user32dll.NewProc("GetWindowTextW")
	procGetWindowRect       = user32dll.NewProc("GetWindowRect")
	procSetProcessDPIAwareW = user32dll.NewProc("SetProcessDPIAware")
	procShowWindow          = user32dll.NewProc("ShowWindow")
	procSetForegroundWindow = user32dll.NewProc("SetForegroundWindow")
)

const (
	swRestore = 9
)

// EnableProcessDPIAware 关闭DPI缩放虚拟化，使用物理像素坐标
func EnableProcessDPIAware() error {
	r1, _, err := procSetProcessDPIAwareW.Call()
	if r1 == 0 { // 调用失败
		if err != syscall.Errno(0) {
			return fmt.Errorf("SetProcessDPIAware 调用失败: %v", err)
		}
		return fmt.Errorf("SetProcessDPIAware 调用失败")
	}
	return nil
}

// FindHwndByTitle 根据窗口标题获取窗口句柄（精确匹配）
func FindHwndByTitle(title string) (syscall.Handle, error) {
	var (
		found     bool
		hwndFound syscall.Handle
	)

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		buf := make([]uint16, 512)
		procGetWindowTextW.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)),
		)
		current := syscall.UTF16ToString(buf)
		if current == title {
			hwndFound = hwnd
			found = true
			return 0 // 停止枚举
		}
		return 1 // 继续枚举
	})

	r1, _, err := procEnumWindows.Call(cb, 0)
	if r1 == 0 && !found {
		// 如果不是被我们主动中断（found），视为失败
		if err != syscall.Errno(0) {
			return 0, fmt.Errorf("EnumWindows 调用失败: %v", err)
		}
		return 0, fmt.Errorf("未找到标题为 '%s' 的窗口", title)
	}

	if !found {
		return 0, fmt.Errorf("未找到标题为 '%s' 的窗口", title)
	}
	return hwndFound, nil
}

type winRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// GetWindowBounds 根据句柄获取窗口的屏幕位置与大小（x,y,w,h）
func GetWindowBounds(hwnd syscall.Handle) (int, int, int, int, error) {
	var r winRect
	ret, _, err := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	if ret == 0 {
		if err != syscall.Errno(0) {
			return 0, 0, 0, 0, fmt.Errorf("GetWindowRect 失败: %v", err)
		}
		return 0, 0, 0, 0, fmt.Errorf("GetWindowRect 失败")
	}

	x := int(r.Left)
	y := int(r.Top)
	w := int(r.Right - r.Left)
	h := int(r.Bottom - r.Top)
	return x, y, w, h, nil
}

// ActivateWindow 还原并激活指定窗口到前台
func ActivateWindow(hwnd syscall.Handle) error {
	if hwnd == 0 {
		return fmt.Errorf("无效的窗口句柄")
	}
	// 先尝试还原窗口（处理最小化等情况）
	procShowWindow.Call(uintptr(hwnd), uintptr(uint32(swRestore)))

	// 设置到前台
	r1, _, err := procSetForegroundWindow.Call(uintptr(hwnd))
	if r1 == 0 {
		if err != syscall.Errno(0) {
			return fmt.Errorf("SetForegroundWindow 失败: %v", err)
		}
		return fmt.Errorf("SetForegroundWindow 失败")
	}
	return nil
}
