package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	op "github.com/yuan71058/OP-Go"
)

const dllPath = `E:\SRC\gop\examples\op_x86.dll`

func main() {
	fmt.Println("========== 多线程记事本操作示例 ==========")

	// 创建 OP 主对象
	mainOP, err := op.NewOP(dllPath)
	if err != nil {
		log.Fatalf("创建 OP 主对象失败: %v", err)
	}
	defer mainOP.Release()

	// 禁用错误弹窗
	mainOP.SetShowErrorMsg(0)
	fmt.Printf("OP 插件版本: %s\n\n", mainOP.Ver())

	// 启动 3 个记事本并获取窗口句柄
	const windowCount = 3
	hwnds := make([]int, windowCount)
	editHwnds := make([]int, windowCount)

	fmt.Println("启动记事本...")
	for i := 0; i < windowCount; i++ {
		mainOP.WinExec("notepad", 1)
		mainOP.Sleep(500)
	}

	mainOP.Sleep(1000) // 等待所有记事本启动完成

	// 查找所有记事本窗口
	for i := 0; i < windowCount; i++ {
		hwnds[i] = mainOP.FindWindow("Notepad", "")
		if hwnds[i] == 0 {
			log.Fatalf("未找到第 %d 个记事本窗口", i+1)
		}
		// 查找编辑框控件 (Edit 类)
		editHwnds[i] = mainOP.FindWindowEx(hwnds[i], "Edit", "")
		if editHwnds[i] == 0 {
			log.Fatalf("未找到第 %d 个记事本的编辑框", i+1)
		}
		fmt.Printf("记事本 %d: 窗口句柄=%d, 编辑框句柄=%d\n", i+1, hwnds[i], editHwnds[i])
	}
	fmt.Println()

	// 创建子对象并绑定窗口
	subOPs := make([]*op.OP, windowCount)
	for i := 0; i < windowCount; i++ {
		subOPs[i], err = op.NewOP(dllPath)
		if err != nil {
			log.Fatalf("创建子对象 %d 失败: %v", i+1, err)
		}
		// 绑定编辑框窗口: gdi模式, 鼠标键盘Windows模式
		ret := subOPs[i].BindWindow(editHwnds[i], "gdi", "windows", "windows", 0)
		if ret == 0 {
			log.Fatalf("绑定窗口 %d 失败", i+1)
		}
		fmt.Printf("子对象 %d 绑定成功\n", i+1)
	}
	fmt.Println()

	// 多线程输入文字
	var wg sync.WaitGroup
	inputChars := []string{"1", "2", "3"}

	fmt.Println("开始多线程输入...")
	startTime := time.Now()

	for i := 0; i < windowCount; i++ {
		wg.Add(1)
		go func(index int, char string, subOP *op.OP) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				subOP.SendString(editHwnds[index], char)
				subOP.Sleep(100)
			}
			fmt.Printf("窗口 %d 输入完成 (输入了200个'%s')\n", index+1, char)
		}(i, inputChars[i], subOPs[i])
	}

	wg.Wait()
	elapsed := time.Since(startTime)
	fmt.Printf("\n所有窗口输入完成，耗时: %v\n\n", elapsed)

	// 延时5秒
	fmt.Println("延时5秒...")
	mainOP.Sleep(5000)

	// 解绑窗口并关闭
	fmt.Println("\n解绑窗口并关闭记事本...")
	for i := 0; i < windowCount; i++ {
		subOPs[i].UnBindWindow()
		// 关闭窗口
		mainOP.SendMessage(hwnds[i], 0x0010, 0, 0) // WM_CLOSE = 0x0010
		fmt.Printf("窗口 %d 已解绑并关闭\n", i+1)
	}

	// 释放子对象
	for i := 0; i < windowCount; i++ {
		subOPs[i].Release()
	}

	fmt.Println("\n资源释放完成，示例结束")
}
