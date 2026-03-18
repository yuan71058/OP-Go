package main

import (
	"fmt"
	"log"
	"strings"

	op "github.com/yuan71058/OP-Go"
)

const dllPath = "E:\\SRC\\gop\\examples\\op_x86.dll"

func main() {
	fmt.Println("========== OP-Go 32位示例 ==========")
	fmt.Println()

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
	opInst, err := op.NewOP(dllPath)
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer opInst.Release()

	// 获取版本号
	version := opInst.Ver()
	fmt.Printf("OP 插件版本: %s\n", version)

	// 设置图片路径
	opInst.SetPath("C:\\images")

	// 获取当前路径
	path := opInst.GetPath()
	fmt.Printf("当前路径: %s\n", path)

	// 禁用错误弹窗
	opInst.SetShowErrorMsg(0)

	fmt.Println()
}

// 示例 2: 窗口操作
func exampleWindow() {
	fmt.Println("========== 示例 2: 窗口操作 ==========")

	opInst, err := op.NewOP(dllPath)
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer opInst.Release()

	// 查找窗口
	hwnd := opInst.FindWindow("Notepad", "")
	if hwnd == 0 {
		hwnd = opInst.FindWindow("", "无标题 - 记事本")
	}

	if hwnd != 0 {
		fmt.Printf("找到窗口，句柄: %d\n", hwnd)

		// 获取窗口信息
		title := opInst.GetWindowTitle(hwnd)
		className := opInst.GetWindowClass(hwnd)
		fmt.Printf("窗口标题: %s, 类名: %s\n", title, className)

		// 获取窗口位置和大小
		x1, y1, x2, y2 := opInst.GetWindowRect(hwnd)
		fmt.Printf("窗口位置: (%d,%d) - (%d,%d)\n", x1, y1, x2, y2)
		fmt.Printf("窗口大小: %dx%d\n", x2-x1, y2-y1)

		// 获取客户区大小
		width, height := opInst.GetClientSize(hwnd)
		fmt.Printf("客户区大小: %dx%d\n", width, height)
	} else {
		fmt.Println("未找到记事本窗口")
	}

	fmt.Println()
}

// 示例 3: 鼠标操作
func exampleMouse() {
	fmt.Println("========== 示例 3: 鼠标操作 ==========")

	opInst, err := op.NewOP(dllPath)
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer opInst.Release()

	// 移动鼠标到指定位置
	opInst.MoveTo(500, 300)
	opInst.Sleep(500)

	// 获取当前鼠标位置
	x, y := opInst.GetCursorPos()
	fmt.Printf("当前鼠标位置: (%d, %d)\n", x, y)

	// 左键单击
	opInst.LeftClick()
	opInst.Sleep(200)

	// 左键双击
	opInst.LeftDoubleClick()
	opInst.Sleep(200)

	// 右键单击
	opInst.RightClick()
	opInst.Sleep(200)

	// 滚轮操作
	opInst.WheelDown()
	opInst.Sleep(100)
	opInst.WheelUp()

	fmt.Println("鼠标操作完成")
	fmt.Println()
}

// 示例 4: 键盘操作
func exampleKeyboard() {
	fmt.Println("========== 示例 4: 键盘操作 ==========")

	opInst, err := op.NewOP(dllPath)
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer opInst.Release()

	// 使用虚拟键码按键
	opInst.KeyPress(0x41) // 'A' 键
	opInst.Sleep(100)

	// 使用字符串形式按键
	opInst.KeyPressChar("a")
	opInst.Sleep(100)
	opInst.KeyPressChar("enter")
	opInst.Sleep(100)
	opInst.KeyPressChar("space")
	opInst.Sleep(100)

	// 组合键操作: Ctrl+C 复制
	opInst.KeyDownChar("ctrl")
	opInst.KeyPressChar("c")
	opInst.KeyUpChar("ctrl")
	opInst.Sleep(100)

	// 组合键操作: Ctrl+V 粘贴
	opInst.KeyDownChar("ctrl")
	opInst.KeyPressChar("v")
	opInst.KeyUpChar("ctrl")

	fmt.Println("键盘操作完成")
	fmt.Println()
}

// 示例 5: 图色操作
func exampleImage() {
	fmt.Println("========== 示例 5: 图色操作 ==========")

	opInst, err := op.NewOP(dllPath)
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer opInst.Release()

	// 设置图片路径
	opInst.SetPath("C:\\images")

	// 截取屏幕
	opInst.Capture(0, 0, 1920, 1080, "screenshot.bmp")
	fmt.Println("截图已保存到 screenshot.bmp")

	// 获取指定坐标颜色
	color := opInst.GetColor(500, 300)
	fmt.Printf("坐标 (500,300) 的颜色: %s\n", color)

	// 查找颜色
	x, y, found := opInst.FindColor(0, 0, 1920, 1080, "FF0000", 0.9, 0)
	if found {
		fmt.Printf("找到红色，位置: (%d, %d)\n", x, y)
	} else {
		fmt.Println("未找到红色")
	}

	// 查找图片
	x, y, found = opInst.FindPic(0, 0, 1920, 1080, "button.bmp", "", 0.9, 0)
	if found {
		fmt.Printf("找到图片，位置: (%d, %d)\n", x, y)

		// 点击找到的图片
		opInst.MoveTo(x, y)
		opInst.LeftClick()
	} else {
		fmt.Println("未找到图片")
	}

	// 查找多个图片
	result := opInst.FindPicEx(0, 0, 1920, 1080, "icon.bmp", "", 0.9, 0)
	if result != "" {
		matches := strings.Split(result, "|")
		fmt.Printf("找到 %d 个匹配\n", len(matches))
	}

	fmt.Println()
}

// 示例 6: OCR 文字识别
func exampleOCR() {
	fmt.Println("========== 示例 6: OCR 文字识别 ==========")

	opInst, err := op.NewOP(dllPath)
	if err != nil {
		log.Printf("创建 OP 实例失败: %v", err)
		return
	}
	defer opInst.Release()

	// 设置字库
	result := opInst.SetDict(0, "C:\\dict\\standard.txt")
	if result != 0 {
		fmt.Printf("设置字库失败，错误码: %d\n", opInst.GetLastError())
		return
	}

	// 选择使用哪个字库
	opInst.UseDict(0)

	// 识别屏幕区域的文字
	text := opInst.Ocr(100, 100, 500, 200, "FFFFFF-000000", 0.9)
	fmt.Printf("识别结果: %s\n", text)

	// 查找字符串
	idx, x, y := opInst.FindStr(0, 0, 1920, 1080, "确定|取消", "FFFFFF-000000", 0.9)
	if idx >= 0 {
		fmt.Printf("找到第 %d 个字符串，位置: (%d, %d)\n", idx, x, y)

		// 点击找到的字符串
		opInst.MoveTo(x, y)
		opInst.LeftClick()
	} else {
		fmt.Println("未找到字符串")
	}

	fmt.Println()
}

// 示例 7: 使用 Service 模式
func exampleService() {
	fmt.Println("========== 示例 7: 使用 Service 模式 ==========")

	// 创建 Service
	svc := op.NewService(dllPath)

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
