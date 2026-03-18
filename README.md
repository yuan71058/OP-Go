# 🎮 OP-Go - OP 插件 Go 语言封装库

<p align="center">
  <img src="https://img.shields.io/badge/Go-%3E%3D1.21-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Windows-0078D6?style=flat-square&logo=windows&logoColor=white" alt="Windows">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/Platform-Windows%20Only-blueviolet?style=flat-square" alt="Platform">
</p>

<p align="center">
  <b>🚀 专为 Windows 平台设计的自动化操作库</b>
</p>

<p align="center">
  <a href="#-功能特性">功能特性</a> •
  <a href="#-快速开始">快速开始</a> •
  <a href="#-安装">安装</a> •
  <a href="#-api-文档">API 文档</a> •
  <a href="#-使用示例">示例</a>
</p>

---

## ✨ 功能特性

| 功能模块 | 说明 | 状态 |
|---------|------|------|
| 🪟 **窗口操作** | 查找窗口、获取信息、移动、设置状态 | ✅ |
| 🖱️ **鼠标操作** | 移动、点击、滚轮、拖拽 | ✅ |
| ⌨️ **键盘操作** | 按键、组合键、字符串输入 | ✅ |
| 🎨 **图色操作** | 截图、找图、找色、取色、比色 | ✅ |
| 🔤 **OCR 识别** | 文字识别、字库支持 | ✅ |
| 🎯 **后台绑定** | 支持后台窗口操作（DX 模式） | ✅ |
| 💾 **内存操作** | 进程内存读写 | ✅ |
| ⚙️ **系统命令** | 剪贴板操作、运行程序等 | ✅ |

---

## 📦 安装

```bash
go get github.com/yuan71058/OP-Go
```

### 前置要求

1. **Windows 操作系统**（仅支持 Windows）
2. **OP 插件 DLL 文件**：
   - `op_x64.dll` 或 `op_x86.dll`（根据系统架构选择）
   - `tools_64.dll` 或 `tools.dll`（免注册方式需要）

> 📥 从 [OP 官方 GitHub](https://github.com/WallBreaker2/op) 下载插件文件

---

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "fmt"
    "log"
    
    op "github.com/yuan71058/OP-Go"
)

func main() {
    // 创建 OP 实例
    opInst, err := op.NewOP("C:\\path\\to\\op_x64.dll")
    if err != nil {
        log.Fatal(err)
    }
    defer opInst.Release()

    // 获取版本号
    version := opInst.Ver()
    fmt.Printf("OP 版本: %s\n", version)

    // 设置图片路径
    opInst.SetPath("C:\\images")

    // 查找窗口
    hwnd := opInst.FindWindow("", "记事本")
    if hwnd != 0 {
        fmt.Printf("找到窗口，句柄: %d\n", hwnd)
    }
}
```

### Service 模式（推荐）

```go
package main

import (
    "log"
    
    op "github.com/yuan71058/OP-Go"
)

func main() {
    // 创建 Service
    svc := op.NewService("C:\\path\\to\\op_x64.dll")

    // 初始化
    if err := svc.Initialize(); err != nil {
        log.Fatal(err)
    }
    defer svc.Close()

    // 获取版本
    version, err := svc.GetVersion()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("版本: %s", version)

    // 查找窗口
    hwnd, err := svc.FindWindow("", "记事本")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("窗口句柄: %d", hwnd)
}
```

---

## 📚 API 文档

### 🔧 基础函数

| 函数 | 说明 |
|------|------|
| `NewOP(dllPath string) (*OP, error)` | 创建 OP 实例 |
| `(*OP) Release()` | 释放 OP 实例 |
| `(*OP) Ver() string` | 获取版本号 |
| `(*OP) SetPath(path string) int` | 设置全局路径 |
| `(*OP) GetPath() string` | 获取全局路径 |
| `(*OP) GetLastError() int` | 获取最后错误码 |
| `(*OP) SetShowErrorMsg(show int) int` | 设置是否显示错误弹窗 |
| `(*OP) Sleep(ms int) int` | 休眠指定毫秒 |

### 🪟 窗口操作

| 函数 | 说明 |
|------|------|
| `(*OP) FindWindow(className, title string) int` | 查找窗口 |
| `(*OP) FindWindowEx(parent int, className, title string) int` | 查找子窗口 |
| `(*OP) GetWindowTitle(hwnd int) string` | 获取窗口标题 |
| `(*OP) GetWindowClass(hwnd int) string` | 获取窗口类名 |
| `(*OP) GetWindowRect(hwnd int) (x1, y1, x2, y2 int)` | 获取窗口位置 |
| `(*OP) GetClientSize(hwnd int) (width, height int)` | 获取客户区大小 |
| `(*OP) SetWindowState(hwnd, flag int) int` | 设置窗口状态 |
| `(*OP) MoveWindow(hwnd, x, y int) int` | 移动窗口 |
| `(*OP) EnumWindow(parent int, title, className string, filter int) string` | 枚举窗口 |

### 🎯 后台绑定

| 函数 | 说明 |
|------|------|
| `(*OP) BindWindow(hwnd int, display, mouse, keypad string, mode int) int` | 绑定窗口 |
| `(*OP) UnBindWindow() int` | 解绑窗口 |
| `(*OP) IsBind() int` | 判断是否已绑定 |

**绑定模式说明**：
- `display`: `normal`, `gdi`, `gdi2`, `dx`, `dx2`
- `mouse`: `normal`, `windows`, `dx`
- `keypad`: `normal`, `windows`, `dx`

### 🖱️ 鼠标操作

| 函数 | 说明 |
|------|------|
| `(*OP) MoveTo(x, y int) int` | 移动鼠标到指定位置 |
| `(*OP) LeftClick() int` | 左键单击 |
| `(*OP) LeftDoubleClick() int` | 左键双击 |
| `(*OP) LeftDown() / LeftUp() int` | 左键按下/弹起 |
| `(*OP) RightClick() int` | 右键单击 |
| `(*OP) RightDown() / RightUp() int` | 右键按下/弹起 |
| `(*OP) MiddleClick() int` | 中键单击 |
| `(*OP) WheelUp() / WheelDown() int` | 滚轮向上/向下 |
| `(*OP) GetCursorPos() (x, y int)` | 获取鼠标位置 |

### ⌨️ 键盘操作

| 函数 | 说明 |
|------|------|
| `(*OP) KeyPress(vkCode int) int` | 按下并弹起虚拟键码 |
| `(*OP) KeyPressChar(keyStr string) int` | 按下并弹起字符键 |
| `(*OP) KeyDown(vkCode int) / KeyUp(vkCode int) int` | 按住/弹起虚拟键码 |
| `(*OP) KeyDownChar(keyStr string) / KeyUpChar(keyStr string) int` | 按住/弹起字符键 |
| `(*OP) GetKeyState(vkCode int) int` | 获取按键状态 |

### 🎨 图色操作

| 函数 | 说明 |
|------|------|
| `(*OP) Capture(x1, y1, x2, y2 int, file string) int` | 截图保存为文件 |
| `(*OP) GetColor(x, y int) string` | 获取指定坐标颜色 |
| `(*OP) CmpColor(x, y int, color string, sim float64) int` | 比较颜色 |
| `(*OP) FindColor(x1, y1, x2, y2 int, color string, sim float64, dir int) (x, y int, found bool)` | 查找颜色 |
| `(*OP) FindPic(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) (x, y int, found bool)` | 查找图片 |
| `(*OP) FindPicEx(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) string` | 查找所有图片 |

### 🔤 OCR 文字识别

| 函数 | 说明 |
|------|------|
| `(*OP) SetDict(index int, file string) int` | 设置字库文件 |
| `(*OP) UseDict(index int) int` | 选择字库 |
| `(*OP) Ocr(x1, y1, x2, y2 int, colorFormat string, sim float64) string` | 识别文字 |
| `(*OP) OcrAuto(x1, y1, x2, y2 int, sim float64) string` | 自动识别文字 |
| `(*OP) FindStr(x1, y1, x2, y2 int, str, colorFormat string, sim float64) (ret, x, y int)` | 查找字符串 |

---

## 💡 使用示例

### 示例 1：窗口自动化

```go
package main

import (
    "fmt"
    "log"
    
    op "github.com/yuan71058/OP-Go"
)

func main() {
    opInst, err := op.NewOP("op_x64.dll")
    if err != nil {
        log.Fatal(err)
    }
    defer opInst.Release()

    // 查找记事本窗口
    hwnd := opInst.FindWindow("Notepad", "")
    if hwnd == 0 {
        log.Fatal("未找到记事本窗口")
    }

    // 激活窗口
    opInst.SetWindowState(hwnd, 1)
    opInst.Sleep(500)

    // 输入文本
    opInst.SendString(hwnd, "Hello, World!")
}
```

### 示例 2：找图点击

```go
package main

import (
    "fmt"
    "log"
    
    op "github.com/yuan71058/OP-Go"
)

func main() {
    opInst, err := op.NewOP("op_x64.dll")
    if err != nil {
        log.Fatal(err)
    }
    defer opInst.Release()

    // 设置图片路径
    opInst.SetPath("C:\\images")

    // 查找图片
    x, y, found := opInst.FindPic(0, 0, 1920, 1080, "button.bmp", "", 0.9, 0)
    if found {
        fmt.Printf("找到图片，位置: (%d, %d)\n", x, y)
        opInst.MoveTo(x, y)
        opInst.LeftClick()
    } else {
        fmt.Println("未找到图片")
    }
}
```

### 示例 3：OCR 识别

```go
package main

import (
    "fmt"
    "log"
    
    op "github.com/yuan71058/OP-Go"
)

func main() {
    opInst, err := op.NewOP("op_x64.dll")
    if err != nil {
        log.Fatal(err)
    }
    defer opInst.Release()

    // 设置字库
    opInst.SetDict(0, "C:\\dict\\standard.txt")
    opInst.UseDict(0)

    // 识别屏幕区域的文字
    text := opInst.Ocr(100, 100, 500, 200, "FFFFFF-000000", 0.9)
    fmt.Printf("识别结果: %s\n", text)
}
```

### 示例 4：后台绑定模式

```go
package main

import (
    "fmt"
    "log"
    
    op "github.com/yuan71058/OP-Go"
)

func main() {
    opInst, err := op.NewOP("op_x64.dll")
    if err != nil {
        log.Fatal(err)
    }
    defer opInst.Release()

    // 查找游戏窗口
    hwnd := opInst.FindWindow("", "游戏窗口")
    if hwnd == 0 {
        log.Fatal("未找到窗口")
    }

    // 绑定窗口（后台模式）
    ret := opInst.BindWindow(hwnd, "dx", "dx", "dx", 0)
    if ret == 0 {
        log.Fatal("绑定失败")
    }
    defer opInst.UnBindWindow()

    // 后台截图
    opInst.Capture(0, 0, 800, 600, "screenshot.bmp")

    // 后台点击
    opInst.MoveTo(400, 300)
    opInst.LeftClick()
}
```

### 示例 5：多线程操作（推荐）

```go
package main

import (
    "fmt"
    "log"
    "os/exec"
    "strconv"
    "strings"
    "sync"
    "time"

    op "github.com/yuan71058/OP-Go"
)

func main() {
    // 创建 OP 主对象
    mainOP, err := op.NewOP("op_x86.dll")
    if err != nil {
        log.Fatal(err)
    }
    defer mainOP.Release()

    mainOP.SetShowErrorMsg(0)

    // 启动 3 个记事本
    const windowCount = 3
    for i := 0; i < windowCount; i++ {
        mainOP.WinExec("notepad", 1)
        mainOP.Sleep(300)
    }
    mainOP.Sleep(1500)

    // 枚举所有记事本窗口
    hwndStr := mainOP.EnumWindow(0, "", "Notepad", 1+2+4+8+16)
    hwndList := parseIntList(hwndStr)
    hwnds := hwndList[:windowCount]

    // 查找编辑框句柄并获取进程ID
    editHwnds := make([]int, windowCount)
    pids := make([]int, windowCount)
    for i := 0; i < windowCount; i++ {
        editHwnds[i] = mainOP.FindWindowEx(hwnds[i], "Edit", "")
        pids[i] = mainOP.GetWindowProcessId(hwnds[i])
    }

    // 创建子对象并绑定窗口
    subOPs := make([]*op.OP, windowCount)
    for i := 0; i < windowCount; i++ {
        subOPs[i], _ = op.NewOP("op_x86.dll")
        subOPs[i].BindWindow(editHwnds[i], "gdi", "windows", "windows", 0)
    }

    // 智能排列窗口
    arrangeWindows(mainOP, hwnds)

    // 多线程输入文字
    var wg sync.WaitGroup
    inputChars := []string{"1", "2", "3"}

    for i := 0; i < windowCount; i++ {
        wg.Add(1)
        go func(index int, char string, subOP *op.OP) {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                subOP.SendString(editHwnds[index], char)
                subOP.Sleep(100)
            }
        }(i, inputChars[i], subOPs[i])
    }
    wg.Wait()

    // 解绑窗口
    for i := 0; i < windowCount; i++ {
        subOPs[i].UnBindWindow()
    }

    // 结束进程
    for i := 0; i < windowCount; i++ {
        cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pids[i]))
        cmd.Run()
        subOPs[i].Release()
    }
}
```

---

## ⌨️ 虚拟键码常量

```go
const (
    VK_LBUTTON  = 0x01 // 鼠标左键
    VK_RBUTTON  = 0x02 // 鼠标右键
    VK_MBUTTON  = 0x04 // 鼠标中键
    VK_BACK     = 0x08 // Backspace
    VK_TAB      = 0x09 // Tab
    VK_RETURN   = 0x0D // Enter
    VK_SHIFT    = 0x10 // Shift
    VK_CONTROL  = 0x11 // Ctrl
    VK_MENU     = 0x12 // Alt
    VK_ESCAPE   = 0x1B // Esc
    VK_SPACE    = 0x20 // 空格
    VK_LEFT     = 0x25 // 左箭头
    VK_UP       = 0x26 // 上箭头
    VK_RIGHT    = 0x27 // 右箭头
    VK_DOWN     = 0x28 // 下箭头
    VK_0        = 0x30 // 0
    // ... A-Z, F1-F12 等
)
```

---

## ⚠️ 注意事项

1. **仅支持 Windows 平台**：OP 插件是 Windows 专用插件
2. **管理员权限**：某些功能可能需要管理员权限
3. **DLL 文件路径**：确保 DLL 文件路径正确
4. **资源释放**：使用完毕后务必调用 `Release()` 释放资源
5. **免注册方式**：将 `tools_64.dll` 或 `tools.dll` 放在 DLL 同目录或系统目录

---

## 🔗 相关链接

- 📘 [OP 官方 GitHub](https://github.com/WallBreaker2/op)
- 📖 [OP 官方文档](https://github.com/WallBreaker2/op/wiki)
- 🛠️ [Go-OLE 库](https://github.com/go-ole/go-ole)

---

## 📄 许可证

MIT License © 2024

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

<p align="center">
  <sub>Built with ❤️ by <a href="https://github.com/yuan71058">yuan71058</a></sub>
</p>
