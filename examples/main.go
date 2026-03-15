package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/yourusername/gop"
)

func main() {
	// 示例 1: 基础初始化
	exampleBasic()

	// 示例 2: 窗口操作
	exampleWindow()

	// 示例 3: 鼠标操作
	exampleMouse()

	// 示例 4: 键盘操作
	exampleKeyboard()

	// 示例 5: 图色操作
	exampleImage()

	// 示例 6: OCR 识别
	exampleOCR()

	// 示例 7: 使用 Service 模式
	exampleService()
}

// 示例 1: 基础初始化与配置
func exampleBasic() {
	fmt.Println("========== 示例 1: 基础初始化 ==========")

	// 创建 OP 实例
	op, err := gop.NewOP("C:\\path\\to\\op_x64.dll")
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer op.Release()

	// 获取版本号
	version := op.Ver()
	fmt.Printf("OP 插件版本: %s\n", version)

	// 设置图片路径
	op.SetPath("C:\\images")

	// 获取当前路径
	path := op.GetPath()
	fmt.Printf("当前路径: %s\n", path)

	// 禁用错误弹窗
	op.SetShowErrorMsg(0)

	fmt.Println()
}

// 示例 2: 窗口操作
func exampleWindow() {
	fmt.Println("========== 示例 2: 窗口操作 ==========")

	op, err := gop.NewOP("op_x64.dll")
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer op.Release()

	// 查找窗口
	hwnd := op.FindWindow("Notepad", "")
	if hwnd == 0 {
		hwnd = op.FindWindow("", "无标题 - 记事本")
	}

	if hwnd != 0 {
		fmt.Printf("找到窗口，句柄: %d\n", hwnd)

		// 获取窗口信息
		title := op.GetWindowTitle(hwnd)
		className := op.GetWindowClass(hwnd)
		fmt.Printf("窗口标题: %s, 类名: %s\n", title, className)

		// 获取窗口位置和大小
		x1, y1, x2, y2 := op.GetWindowRect(hwnd)
		fmt.Printf("窗口位置: (%d,%d) - (%d,%d)\n", x1, y1, x2, y2)
		fmt.Printf("窗口大小: %dx%d\n", x2-x1, y2-y1)

		// 获取客户区大小
		width, height := op.GetClientSize(hwnd)
		fmt.Printf("客户区大小: %dx%d\n", width, height)
	} else {
		fmt.Println("未找到记事本窗口")
	}

	fmt.Println()
}

// 示例 3: 鼠标操作
func exampleMouse() {
	fmt.Println("========== 示例 3: 鼠标操作 ==========")

	op, err := gop.NewOP("op_x64.dll")
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer op.Release()

	// 移动鼠标到指定位置
	op.MoveTo(500, 300)
	op.Sleep(500)

	// 获取当前鼠标位置
	x, y := op.GetCursorPos()
	fmt.Printf("当前鼠标位置: (%d, %d)\n", x, y)

	// 左键单击
	op.LeftClick()
	op.Sleep(200)

	// 左键双击
	op.LeftDoubleClick()
	op.Sleep(200)

	// 右键单击
	op.RightClick()
	op.Sleep(200)

	// 滚轮操作
	op.WheelDown()
	op.Sleep(100)
	op.WheelUp()

	fmt.Println("鼠标操作完成")
	fmt.Println()
}

// 示例 4: 键盘操作
func exampleKeyboard() {
	fmt.Println("========== 示例 4: 键盘操作 ==========")

	op, err := gop.NewOP("op_x64.dll")
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer op.Release()

	// 使用虚拟键码按键
	op.KeyPress(0x41) // 'A' 键
	op.Sleep(100)

	// 使用字符串形式按键
	op.KeyPressChar("a")
	op.Sleep(100)
	op.KeyPressChar("enter")
	op.Sleep(100)
	op.KeyPressChar("space")
	op.Sleep(100)

	// 组合键操作: Ctrl+C 复制
	op.KeyDownChar("ctrl")
	op.KeyPressChar("c")
	op.KeyUpChar("ctrl")
	op.Sleep(100)

	// 组合键操作: Ctrl+V 粘贴
	op.KeyDownChar("ctrl")
	op.KeyPressChar("v")
	op.KeyUpChar("ctrl")

	fmt.Println("键盘操作完成")
	fmt.Println()
}

// 示例 5: 图色操作
func exampleImage() {
	fmt.Println("========== 示例 5: 图色操作 ==========")

	op, err := gop.NewOP("op_x64.dll")
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer op.Release()

	// 设置图片路径
	op.SetPath("C:\\images")

	// 截取屏幕
	op.Capture(0, 0, 1920, 1080, "screenshot.bmp")
	fmt.Println("截图已保存到 screenshot.bmp")

	// 获取指定坐标颜色
	color := op.GetColor(500, 300)
	fmt.Printf("坐标 (500,300) 的颜色: %s\n", color)

	// 查找颜色
	x, y, found := op.FindColor(0, 0, 1920, 1080, "FF0000", 0.9, 0)
	if found {
		fmt.Printf("找到红色，位置: (%d, %d)\n", x, y)
	} else {
		fmt.Println("未找到红色")
	}

	// 查找图片
	x, y, found = op.FindPic(0, 0, 1920, 1080, "button.bmp", "", 0.9, 0)
	if found {
		fmt.Printf("找到图片，位置: (%d, %d)\n", x, y)

		// 点击找到的图片
		op.MoveTo(x, y)
		op.LeftClick()
	} else {
		fmt.Println("未找到图片")
	}

	// 查找多个图片
	result := op.FindPicEx(0, 0, 1920, 1080, "icon.bmp", "", 0.9, 0)
	if result != "" {
		matches := strings.Split(result, "|")
		fmt.Printf("找到 %d 个匹配\n", len(matches))
	}

	fmt.Println()
}

// 示例 6: OCR 文字识别
func exampleOCR() {
	fmt.Println("========== 示例 6: OCR 文字识别 ==========")

	op, err := gop.NewOP("op_x64.dll")
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer op.Release()

	// 设置字库
	result := op.SetDict(0, "C:\\dict\\standard.txt")
	if result != 0 {
		fmt.Printf("设置字库失败，错误码: %d\n", op.GetLastError())
		return
	}

	// 选择使用哪个字库
	op.UseDict(0)

	// 识别屏幕区域的文字
	text := op.Ocr(100, 100, 500, 200, "FFFFFF-000000", 0.9)
	fmt.Printf("识别结果: %s\n", text)

	// 查找字符串
	idx, x, y := op.FindStr(0, 0, 1920, 1080, "确定|取消", "FFFFFF-000000", 0.9)
	if idx >= 0 {
		fmt.Printf("找到第 %d 个字符串，位置: (%d, %d)\n", idx, x, y)

		// 点击找到的字符串
		op.MoveTo(x, y)
		op.LeftClick()
	} else {
		fmt.Println("未找到字符串")
	}

	fmt.Println()
}

// 示例 7: 使用 Service 模式
func exampleService() {
	fmt.Println("========== 示例 7: 使用 Service 模式 ==========")

	// 创建 Service
	svc := gop.NewService("op_x64.dll")

	// 初始化
	if err := svc.Initialize(); err != nil {
		log.Printf("初始化失败: %v", err)
		return
	}
	defer svc.Close()

	// 获取版本
	version, err := svc.GetVersion()
	if err != nil {
		log.Printf("获取版本失败: %v", err)
		return
	}
	fmt.Printf("版本: %s\n", version)

	// 获取屏幕大小
	width, height, err := svc.GetScreenSize()
	if err != nil {
		log.Printf("获取屏幕大小失败: %v", err)
		return
	}
	fmt.Printf("屏幕大小: %dx%d\n", width, height)

	// 查找窗口
	hwnd, err := svc.FindWindow("", "记事本")
	if err != nil {
		log.Printf("查找窗口失败: %v", err)
		return
	}
	fmt.Printf("窗口句柄: %d\n", hwnd)

	// 获取窗口标题
	if hwnd != 0 {
		title, err := svc.GetWindowText(hwnd)
		if err != nil {
			log.Printf("获取窗口标题失败: %v", err)
			return
		}
		fmt.Printf("窗口标题: %s\n", title)
	}

	// 获取状态
	status, err := svc.GetStatus()
	if err != nil {
		log.Printf("获取状态失败: %v", err)
		return
	}
	fmt.Printf("状态: IsReady=%v, Version=%s\n", status.IsReady, status.Version)

	fmt.Println()
}
