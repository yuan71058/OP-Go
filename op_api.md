# OP 插件 Go 语言封装 API 文档

## 概述

OP (Operator & Open) 是一款专为 Windows 设计的开源自动化插件，提供屏幕读取、输入模拟、图像处理、OCR 等功能。

项目地址：https://github.com/WallBreaker2/op

---

## 目录

- [初始化与基础](#初始化与基础)
- [后台绑定](#后台绑定)
- [窗口操作](#窗口操作)
- [鼠标操作](#鼠标操作)
- [键盘操作](#键盘操作)
- [图色操作](#图色操作)
- [OCR 文字识别](#ocr-文字识别)
- [系统命令](#系统命令)

---

## 初始化与基础

### NewOP

创建一个新的 OP 插件实例

```go
func NewOP(dllPath string) (*OP, error)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| dllPath | string | OP 插件 DLL 文件的完整路径 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| op | *OP | OP 插件实例 |
| err | error | 错误信息 |

**使用方式：**

1. **免注册方式**：需要 tools_64.dll 或 tools.dll 文件
2. **注册方式**：需要先使用 `regsvr32 op_x64.dll` 注册

**示例：**

```go
op, err := op.NewOP("C:\\path\\to\\op_x64.dll")
if err != nil {
    log.Fatal(err)
}
defer op.Release()
```

---

### Release

释放 OP 插件实例

```go
func (o *OP) Release()
```

**说明：**

- 必须在不再使用 OP 实例时调用，以释放 COM 资源
- 建议在创建后立即使用 `defer op.Release()`

---

### Ver

获取 OP 插件的版本号

```go
func (o *OP) Ver() string
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| version | string | 插件版本号 |

---

### SetPath

设置全局路径

```go
func (o *OP) SetPath(path string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | 全局路径 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功，非 0 表示失败 |

---

### GetPath

获取全局路径

```go
func (o *OP) GetPath() string
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| path | string | 当前全局路径 |

---

### GetBasePath

获取插件所在目录

```go
func (o *OP) GetBasePath() string
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| path | string | 插件所在目录路径 |

---

### GetID

获取当前对象的 ID 值

```go
func (o *OP) GetID() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| id | int | 对象 ID |

---

### GetLastError

获取最后的错误

```go
func (o *OP) GetLastError() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| errorCode | int | 错误代码 |

---

### SetShowErrorMsg

设置是否弹出错误信息

```go
func (o *OP) SetShowErrorMsg(show int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| show | int | 0 表示不弹出，1 表示弹出 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### Sleep

休眠指定时间

```go
func (o *OP) Sleep(ms int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| ms | int | 休眠时间，单位毫秒 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### EnablePicCache

设置是否开启或关闭插件内部的图片缓存机制

```go
func (o *OP) EnablePicCache(enable int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int | 0 表示关闭，1 表示开启 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

## 后台绑定

### BindWindow

绑定指定的窗口

```go
func (o *OP) BindWindow(hwnd int, display, mouse, keypad string, mode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |
| display | string | 显示模式："normal", "gdi", "gdi2", "dx", "dx2" 等 |
| mouse | string | 鼠标模式："normal", "windows", "dx" 等 |
| keypad | string | 键盘模式："normal", "windows", "dx" 等 |
| mode | int | 绑定模式：0 表示正常模式，其他值参考文档 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功，非 0 表示失败 |

**显示模式说明：**

| 模式 | 说明 |
|------|------|
| normal | 正常模式，使用 Windows API 截图 |
| gdi | GDI 模式 |
| gdi2 | GDI 模式2 |
| dx | DirectX 模式 |
| dx2 | DirectX 模式2 |

**示例：**

```go
// 绑定窗口，使用 DX 模式
hwnd := op.FindWindow("", "记事本")
result := op.BindWindow(hwnd, "dx", "dx", "dx", 0)
if result == 0 {
    fmt.Println("绑定成功")
}
```

---

### UnBindWindow

解除绑定窗口

```go
func (o *OP) UnBindWindow() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### GetBindWindow

获取当前对象已经绑定的窗口句柄

```go
func (o *OP) GetBindWindow() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 绑定的窗口句柄，0 表示未绑定 |

---

### IsBind

判断当前对象是否已绑定窗口

```go
func (o *OP) IsBind() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示未绑定，1 表示已绑定 |

---

## 窗口操作

### EnumWindow

根据指定条件枚举系统中符合条件的窗口

```go
func (o *OP) EnumWindow(parent int, title, className string, filter int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| parent | int | 父窗口句柄，0 表示桌面窗口 |
| title | string | 窗口标题，空字符串表示匹配所有 |
| className | string | 窗口类名，空字符串表示匹配所有 |
| filter | int | 过滤条件：1 表示匹配标题，2 表示匹配类名，3 表示同时匹配 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| handles | string | 窗口句柄列表，格式为 "hwnd1,hwnd2,hwnd3" |

---

### FindWindow

查找符合类名或者标题名的顶层可见窗口

```go
func (o *OP) FindWindow(className, title string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| className | string | 窗口类名，空字符串表示匹配所有 |
| title | string | 窗口标题，空字符串表示匹配所有 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄，0 表示未找到 |

**示例：**

```go
// 查找记事本窗口
hwnd := op.FindWindow("Notepad", "")
if hwnd != 0 {
    fmt.Printf("找到记事本窗口，句柄: %d\n", hwnd)
}

// 查找标题包含 "Excel" 的窗口
hwnd = op.FindWindow("", "Excel")
```

---

### FindWindowEx

查找符合类名或者标题名的窗口（支持子窗口）

```go
func (o *OP) FindWindowEx(parent int, className, title string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| parent | int | 父窗口句柄，0 表示桌面窗口 |
| className | string | 窗口类名 |
| title | string | 窗口标题 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄，0 表示未找到 |

---

### GetWindowTitle

获取窗口的标题

```go
func (o *OP) GetWindowTitle(hwnd int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| title | string | 窗口标题 |

---

### GetWindowClass

获取窗口的类名

```go
func (o *OP) GetWindowClass(hwnd int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| className | string | 窗口类名 |

---

### GetWindowRect

获取窗口在屏幕上的位置

```go
func (o *OP) GetWindowRect(hwnd int) (x1, y1, x2, y2 int)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| x1 | int | 窗口左上角 X 坐标 |
| y1 | int | 窗口左上角 Y 坐标 |
| x2 | int | 窗口右下角 X 坐标 |
| y2 | int | 窗口右下角 Y 坐标 |

**示例：**

```go
hwnd := op.FindWindow("Notepad", "")
x1, y1, x2, y2 := op.GetWindowRect(hwnd)
fmt.Printf("窗口位置: (%d,%d) - (%d,%d)\n", x1, y1, x2, y2)
fmt.Printf("窗口大小: %dx%d\n", x2-x1, y2-y1)
```

---

### GetClientSize

获取窗口客户区的宽度和高度

```go
func (o *OP) GetClientSize(hwnd int) (width, height int)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| width | int | 客户区宽度 |
| height | int | 客户区高度 |

---

### GetWindowState

获取指定窗口的一些属性

```go
func (o *OP) GetWindowState(hwnd, flag int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |
| flag | int | 要获取的状态类型 |

**flag 取值：**

| 值 | 说明 |
|----|------|
| 0 | 判断窗口是否存在 |
| 1 | 判断窗口是否可见 |
| 2 | 判断窗口是否可用 |
| 3 | 判断窗口是否最大化 |
| 4 | 判断窗口是否最小化 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示否，1 表示是 |

---

### SetWindowState

设置窗口的状态

```go
func (o *OP) SetWindowState(hwnd, flag int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |
| flag | int | 状态标志 |

**flag 取值：**

| 值 | 说明 |
|----|------|
| 0 | 关闭窗口 |
| 1 | 激活窗口 |
| 2 | 最小化窗口 |
| 3 | 最大化窗口 |
| 4 | 恢复窗口 |
| 5 | 隐藏窗口 |
| 6 | 显示窗口 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### MoveWindow

移动指定窗口到指定位置

```go
func (o *OP) MoveWindow(hwnd, x, y int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |
| x | int | 新的 X 坐标 |
| y | int | 新的 Y 坐标 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### SendString

向指定窗口发送文本数据

```go
func (o *OP) SendString(hwnd int, str string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |
| str | string | 要发送的字符串 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### SendStringIme

向指定窗口发送文本数据（输入法方式）

```go
func (o *OP) SendStringIme(hwnd int, str string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄 |
| str | string | 要发送的字符串 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

**说明：**

- 使用输入法方式发送字符串，可以发送中文等特殊字符

---

## 鼠标操作

### MoveTo

把鼠标移动到目的点

```go
func (o *OP) MoveTo(x, y int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int | 目标 X 坐标 |
| y | int | 目标 Y 坐标 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### LeftClick

按下鼠标左键

```go
func (o *OP) LeftClick() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### LeftDoubleClick

双击鼠标左键

```go
func (o *OP) LeftDoubleClick() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### LeftDown

按住鼠标左键

```go
func (o *OP) LeftDown() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### LeftUp

弹起鼠标左键

```go
func (o *OP) LeftUp() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### RightClick

按下鼠标右键

```go
func (o *OP) RightClick() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### RightDown

按住鼠标右键

```go
func (o *OP) RightDown() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### RightUp

弹起鼠标右键

```go
func (o *OP) RightUp() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### WheelDown

滚轮向下滚

```go
func (o *OP) WheelDown() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### WheelUp

滚轮向上滚

```go
func (o *OP) WheelUp() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### GetCursorPos

获取鼠标位置

```go
func (o *OP) GetCursorPos() (x, y int)
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| x | int | 鼠标 X 坐标 |
| y | int | 鼠标 Y 坐标 |

**示例：**

```go
x, y := op.GetCursorPos()
fmt.Printf("鼠标位置: (%d, %d)\n", x, y)
```

---

## 键盘操作

### KeyPress

按下并弹起指定的虚拟键码

```go
func (o *OP) KeyPress(vkCode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| vkCode | int | 虚拟键码，如 0x41 表示 'A' |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

**常用虚拟键码：**

| 键码 | 值 | 说明 |
|------|----|------|
| VK_ESCAPE | 0x1B | Esc 键 |
| VK_SPACE | 0x20 | 空格键 |
| VK_RETURN | 0x0D | 回车键 |
| VK_TAB | 0x09 | Tab 键 |
| VK_BACK | 0x08 | Backspace 键 |
| VK_DELETE | 0x2E | Delete 键 |
| VK_LEFT | 0x25 | 左箭头 |
| VK_UP | 0x26 | 上箭头 |
| VK_RIGHT | 0x27 | 右箭头 |
| VK_DOWN | 0x28 | 下箭头 |
| VK_F1 - VK_F12 | 0x70-0x7B | F1-F12 |
| VK_LSHIFT | 0xA0 | 左 Shift |
| VK_RSHIFT | 0xA1 | 右 Shift |
| VK_LCONTROL | 0xA2 | 左 Ctrl |
| VK_RCONTROL | 0xA3 | 右 Ctrl |
| VK_LMENU | 0xA4 | 左 Alt |
| VK_RMENU | 0xA5 | 右 Alt |

---

### KeyPressChar

按下并弹起指定的虚拟键码（字符串形式）

```go
func (o *OP) KeyPressChar(keyStr string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| keyStr | string | 键名字符串，如 "a", "enter", "ctrl", "alt" 等 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### KeyDown

按住指定的虚拟键码

```go
func (o *OP) KeyDown(vkCode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| vkCode | int | 虚拟键码 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### KeyDownChar

按住指定的虚拟键码（字符串形式）

```go
func (o *OP) KeyDownChar(keyStr string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| keyStr | string | 键名字符串 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### KeyUp

弹起虚拟键

```go
func (o *OP) KeyUp(vkCode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| vkCode | int | 虚拟键码 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### KeyUpChar

弹起虚拟键（字符串形式）

```go
func (o *OP) KeyUpChar(keyStr string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| keyStr | string | 键名字符串 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

## 图色操作

### Capture

抓取指定区域的图像，保存为文件

```go
func (o *OP) Capture(x1, y1, x2, y2 int, file string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int | 区域左上角 X 坐标 |
| y1 | int | 区域左上角 Y 坐标 |
| x2 | int | 区域右下角 X 坐标 |
| y2 | int | 区域右下角 Y 坐标 |
| file | string | 保存的文件路径 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

**示例：**

```go
// 截图整个屏幕
op.Capture(0, 0, 1920, 1080, "screenshot.bmp")
```

---

### FindPic

查找图片

```go
func (o *OP) FindPic(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) (x, y int, found bool)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int | 查找区域左上角 X 坐标 |
| y1 | int | 查找区域左上角 Y 坐标 |
| x2 | int | 查找区域右下角 X 坐标 |
| y2 | int | 查找区域右下角 Y 坐标 |
| picName | string | 图片文件名，多个图片用 \| 分隔 |
| deltaColor | string | 偏色，格式为 "RRGGBB"，空字符串表示无偏色 |
| sim | float64 | 相似度，取值范围 0.0-1.0 |
| dir | int | 查找方向：0 表示从左到右从上到下 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| x | int | 找到的图片左上角 X 坐标 |
| y | int | 找到的图片左上角 Y 坐标 |
| found | bool | 是否找到 |

**示例：**

```go
// 查找单个图片
x, y, found := op.FindPic(0, 0, 1920, 1080, "button.bmp", "", 0.9, 0)
if found {
    fmt.Printf("找到图片，位置: (%d, %d)\n", x, y)
}

// 查找多个图片
x, y, found = op.FindPic(0, 0, 1920, 1080, "btn1.bmp|btn2.bmp|btn3.bmp", "", 0.9, 0)
if found {
    fmt.Printf("找到第 %d 个图片，位置: (%d, %d)\n", x, y)
}
```

---

### FindPicEx

多图查找，返回所有找到的图片信息

```go
func (o *OP) FindPicEx(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) string
```

**参数说明：**

与 FindPic 相同

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | string | 找到的图片信息，格式为 "idx,x,y\|idx,x,y\|..." |

---

### FindColor

查找指定区域内的颜色

```go
func (o *OP) FindColor(x1, y1, x2, y2 int, color string, sim float64, dir int) (x, y int, found bool)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| color | string | 要查找的颜色，格式为 "RRGGBB"，如 "FF0000" 表示红色 |
| sim | float64 | 相似度，取值范围 0.0-1.0，1.0 表示完全匹配 |
| dir | int | 查找方向，0 表示从左到右从上到下 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| x | int | 找到的颜色 X 坐标 |
| y | int | 找到的颜色 Y 坐标 |
| found | bool | 是否找到 |

**示例：**

```go
// 查找红色
x, y, found := op.FindColor(0, 0, 1920, 1080, "FF0000", 1.0, 0)
if found {
    fmt.Printf("找到红色，位置: (%d, %d)\n", x, y)
}
```

---

### FindColorEx

多点找色

```go
func (o *OP) FindColorEx(x1, y1, x2, y2 int, color string, sim float64, dir int) (x, y int, found bool)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| color | string | 颜色描述，格式为 "颜色\|偏移X\|偏移Y,颜色\|偏移X\|偏移Y,..." |
| sim | float64 | 相似度，取值范围 0.0-1.0 |
| dir | int | 查找方向 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| x | int | 找到的第一个点的 X 坐标 |
| y | int | 找到的第一个点的 Y 坐标 |
| found | bool | 是否找到 |

**示例：**

```go
// 多点找色，查找一个由多个颜色点组成的图案
// 格式: "主颜色|相对X|相对Y,次颜色|相对X|相对Y,..."
colorStr := "FF0000|0|0,00FF00|10|0,0000FF|20|0"
x, y, found := op.FindColorEx(0, 0, 1920, 1080, colorStr, 0.9, 0)
```

---

### FindColorEx2

多点找色，返回所有找到的点

```go
func (o *OP) FindColorEx2(x1, y1, x2, y2 int, color string, sim float64, dir int) string
```

**参数说明：**

与 FindColorEx 相同

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | string | 所有找到的坐标字符串，格式为 "x,y\|x,y\|..." |

---

### GetColor

获取指定坐标的颜色

```go
func (o *OP) GetColor(x, y int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int | X 坐标 |
| y | int | Y 坐标 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| color | string | 颜色值，格式为 "RRGGBB" |

**示例：**

```go
color := op.GetColor(100, 100)
fmt.Printf("坐标 (100,100) 的颜色: %s\n", color)
```

---

### CmpColor

比较指定坐标点的颜色

```go
func (o *OP) CmpColor(x, y int, color string, sim float64) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int | X 坐标 |
| y | int | Y 坐标 |
| color | string | 要比较的颜色，格式为 "RRGGBB" |
| sim | float64 | 相似度 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示匹配，1 表示不匹配 |

---

## OCR 文字识别

### SetDict

设置字库文件

```go
func (o *OP) SetDict(index int, file string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int | 字库索引，范围 0-9 |
| file | string | 字库文件路径 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### UseDict

选择使用哪个字库文件进行识别

```go
func (o *OP) UseDict(index int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int | 字库索引，范围 0-9 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

### Ocr

识别屏幕范围内的字符串

```go
func (o *OP) Ocr(x1, y1, x2, y2 int, colorFormat string, sim float64) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 识别区域 |
| colorFormat | string | 颜色格式，格式为 "RRGGBB-RRGGBB" |
| sim | float64 | 相似度 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| text | string | 识别出的文字 |

**示例：**

```go
// 设置字库
op.SetDict(0, "dict.txt")
op.UseDict(0)

// 识别文字
text := op.Ocr(100, 100, 500, 200, "FFFFFF-000000", 0.9)
fmt.Printf("识别结果: %s\n", text)
```

---

### FindStr

在屏幕范围内查找字符串

```go
func (o *OP) FindStr(x1, y1, x2, y2 int, str, colorFormat string, sim float64) (ret, x, y int)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| str | string | 要查找的字符串，支持多字符串查找，用 \| 分隔 |
| colorFormat | string | 颜色格式，格式为 "RRGGBB-RRGGBB" |
| sim | float64 | 相似度，取值范围 0.0-1.0 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| ret | int | 找到的字符串索引（从 0 开始），未找到返回 -1 |
| x | int | 找到的字符串左上角 X 坐标 |
| y | int | 找到的字符串左上角 Y 坐标 |

**示例：**

```go
// 查找单个字符串
idx, x, y := op.FindStr(0, 0, 1920, 1080, "确定", "FFFFFF-000000", 0.9)
if idx >= 0 {
    fmt.Printf("找到文字，位置: (%d, %d)\n", x, y)
}

// 查找多个字符串
idx, x, y = op.FindStr(0, 0, 1920, 1080, "确定|取消|关闭", "FFFFFF-000000", 0.9)
if idx >= 0 {
    fmt.Printf("找到第 %d 个文字，位置: (%d, %d)\n", idx, x, y)
}
```

---

### FindStrEx

在屏幕范围内查找字符串，返回所有找到的

```go
func (o *OP) FindStrEx(x1, y1, x2, y2 int, str, colorFormat string, sim float64) string
```

**参数说明：**

与 FindStr 相同

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | string | 所有找到的字符串信息，格式为 "idx,x,y,str\|idx,x,y,str\|..." |

---

## 系统命令

### GetScreenWidth

获取屏幕宽度

```go
func (o *OP) GetScreenWidth() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| width | int | 屏幕宽度 |

---

### GetScreenHeight

获取屏幕高度

```go
func (o *OP) GetScreenHeight() int
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| height | int | 屏幕高度 |

---

### GetClipboard

获取剪贴板内容

```go
func (o *OP) GetClipboard() string
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| text | string | 剪贴板内容 |

---

### SetClipboard

设置剪贴板内容

```go
func (o *OP) SetClipboard(str string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | 要设置的剪贴板内容 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| result | int | 0 表示成功 |

---

## 常量定义

### 虚拟键码

```go
const (
    VK_LBUTTON    = 0x01 // 鼠标左键
    VK_RBUTTON    = 0x02 // 鼠标右键
    VK_CANCEL     = 0x03 // Cancel
    VK_MBUTTON    = 0x04 // 鼠标中键
    VK_XBUTTON1   = 0x05 // 鼠标侧键1
    VK_XBUTTON2   = 0x06 // 鼠标侧键2
    VK_BACK       = 0x08 // Backspace
    VK_TAB        = 0x09 // Tab
    VK_CLEAR      = 0x0C // Clear
    VK_RETURN     = 0x0D // Enter
    VK_SHIFT      = 0x10 // Shift
    VK_CONTROL    = 0x11 // Ctrl
    VK_MENU       = 0x12 // Alt
    VK_PAUSE      = 0x13 // Pause
    VK_CAPITAL    = 0x14 // Caps Lock
    VK_ESCAPE     = 0x1B // Esc
    VK_SPACE      = 0x20 // 空格
    VK_PRIOR      = 0x21 // Page Up
    VK_NEXT       = 0x22 // Page Down
    VK_END        = 0x23 // End
    VK_HOME       = 0x24 // Home
    VK_LEFT       = 0x25 // 左箭头
    VK_UP         = 0x26 // 上箭头
    VK_RIGHT      = 0x27 // 右箭头
    VK_DOWN       = 0x28 // 下箭头
    VK_SELECT     = 0x29 // Select
    VK_PRINT      = 0x2A // Print
    VK_EXECUTE    = 0x2B // Execute
    VK_SNAPSHOT   = 0x2C // Print Screen
    VK_INSERT     = 0x2D // Insert
    VK_DELETE     = 0x2E // Delete
    VK_HELP       = 0x2F // Help
    VK_0          = 0x30 // 0
    VK_1          = 0x31 // 1
    VK_2          = 0x32 // 2
    VK_3          = 0x33 // 3
    VK_4          = 0x34 // 4
    VK_5          = 0x35 // 5
    VK_6          = 0x36 // 6
    VK_7          = 0x37 // 7
    VK_8          = 0x38 // 8
    VK_9          = 0x39 // 9
    VK_A          = 0x41 // A
    VK_B          = 0x42 // B
    VK_C          = 0x43 // C
    VK_D          = 0x44 // D
    VK_E          = 0x45 // E
    VK_F          = 0x46 // F
    VK_G          = 0x47 // G
    VK_H          = 0x48 // H
    VK_I          = 0x49 // I
    VK_J          = 0x4A // J
    VK_K          = 0x4B // K
    VK_L          = 0x4C // L
    VK_M          = 0x4D // M
    VK_N          = 0x4E // N
    VK_O          = 0x4F // O
    VK_P          = 0x50 // P
    VK_Q          = 0x51 // Q
    VK_R          = 0x52 // R
    VK_S          = 0x53 // S
    VK_T          = 0x54 // T
    VK_U          = 0x55 // U
    VK_V          = 0x56 // V
    VK_W          = 0x57 // W
    VK_X          = 0x58 // X
    VK_Y          = 0x59 // Y
    VK_Z          = 0x5A // Z
    VK_LWIN       = 0x5B // 左Win
    VK_RWIN       = 0x5C // 右Win
    VK_APPS       = 0x5D // Applications
    VK_SLEEP      = 0x5F // Sleep
    VK_NUMPAD0    = 0x60 // 小键盘0
    VK_NUMPAD1    = 0x61 // 小键盘1
    VK_NUMPAD2    = 0x62 // 小键盘2
    VK_NUMPAD3    = 0x63 // 小键盘3
    VK_NUMPAD4    = 0x64 // 小键盘4
    VK_NUMPAD5    = 0x65 // 小键盘5
    VK_NUMPAD6    = 0x66 // 小键盘6
    VK_NUMPAD7    = 0x67 // 小键盘7
    VK_NUMPAD8    = 0x68 // 小键盘8
    VK_NUMPAD9    = 0x69 // 小键盘9
    VK_MULTIPLY   = 0x6A // 小键盘*
    VK_ADD        = 0x6B // 小键盘+
    VK_SEPARATOR  = 0x6C // Separator
    VK_SUBTRACT   = 0x6D // 小键盘-
    VK_DECIMAL    = 0x6E // 小键盘.
    VK_DIVIDE     = 0x6F // 小键盘/
    VK_F1         = 0x70 // F1
    VK_F2         = 0x71 // F2
    VK_F3         = 0x72 // F3
    VK_F4         = 0x73 // F4
    VK_F5         = 0x74 // F5
    VK_F6         = 0x75 // F6
    VK_F7         = 0x76 // F7
    VK_F8         = 0x77 // F8
    VK_F9         = 0x78 // F9
    VK_F10        = 0x79 // F10
    VK_F11        = 0x7A // F11
    VK_F12        = 0x7B // F12
    VK_NUMLOCK    = 0x90 // Num Lock
    VK_SCROLL     = 0x91 // Scroll Lock
    VK_LSHIFT     = 0xA0 // 左Shift
    VK_RSHIFT     = 0xA1 // 右Shift
    VK_LCONTROL   = 0xA2 // 左Ctrl
    VK_RCONTROL   = 0xA3 // 右Ctrl
    VK_LMENU      = 0xA4 // 左Alt
    VK_RMENU      = 0xA5 // 右Alt
)
```

---

## 错误处理

OP 插件的函数通常返回整数结果：

- **0**：表示成功
- **非0**：表示失败，具体错误码可以通过 `GetLastError()` 获取

**建议的错误处理方式：**

```go
// 基础函数检查
result := op.SetPath("C:\\images")
if result != 0 {
    errCode := op.GetLastError()
    fmt.Printf("设置路径失败，错误码: %d\n", errCode)
}

// 查找函数检查
x, y, found := op.FindPic(0, 0, 1920, 1080, "btn.bmp", "", 0.9, 0)
if !found {
    fmt.Println("未找到图片")
} else {
    fmt.Printf("找到图片，位置: (%d, %d)\n", x, y)
}
```

---

## 补充函数（新增46个）

以下函数为根据官方文档补充的函数：

### Base 基础函数补充

#### InjectDll

将指定的DLL注入到指定的进程中

```go
func (o *OP) InjectDll(processName, dllName string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| processName | string | 指定要注入DLL的进程名称 |
| dllName | string | 注入的DLL名称 |

**返回值：** 0表示失败，1表示成功

---

#### CapturePre

取上次操作的图色区域，保存为文件(24位位图)

```go
func (o *OP) CapturePre(file string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | 设置保存文件名 |

**返回值：** 0表示失败，1表示成功

---

#### SetScreenDataMode

设置屏幕数据模式

```go
func (o *OP) SetScreenDataMode(mode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| mode | int | 0表示从上到下(默认)，1表示从下到上 |

**返回值：** 0表示失败，1表示成功

---

### Window 窗口操作补充

#### EnumWindowByProcess

根据指定进程以及其它条件，枚举系统中符合条件的窗口

```go
func (o *OP) EnumWindowByProcess(processName, title, className string, filter int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| processName | string | 进程名称 |
| title | string | 窗口的标题 |
| className | string | 窗口的类名 |
| filter | int | 过滤条件，1匹配标题，2匹配类名，4只匹配第一层子窗口，8匹配顶级窗口，16匹配可见窗口，32按打开顺序排列 |

**返回值：** 返回所有匹配到的窗口句柄字符串，格式为"hwnd1,hwnd2,..."

---

#### EnumProcess

根据指定进程名，枚举系统中符合条件的进程PID

```go
func (o *OP) EnumProcess(name string) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| name | string | 进程名称 |

**返回值：** 返回所有匹配的进程PID，格式为"pid1,pid2,..."

---

#### ClientToScreen

把窗口坐标转换为屏幕坐标

```go
func (o *OP) ClientToScreen(hwnd int, x, y *int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| x | *int | 窗口X坐标（变参，会被修改为屏幕坐标） |
| y | *int | 窗口Y坐标（变参，会被修改为屏幕坐标） |

**返回值：** 0表示失败，1表示成功

---

#### FindWindowByProcess

根据指定的进程名字，来查找可见窗口

```go
func (o *OP) FindWindowByProcess(processName, className, title string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| processName | string | 进程名，如"notepad.exe"，精确匹配但不区分大小写 |
| className | string | 窗口类名，模糊匹配 |
| title | string | 窗口标题，模糊匹配 |

**返回值：** 窗口句柄，没找到返回0

---

#### FindWindowByProcessId

根据指定的进程Id，来查找可见窗口

```go
func (o *OP) FindWindowByProcessId(processId int, className, title string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| processId | int | 进程id |
| className | string | 窗口类名，模糊匹配 |
| title | string | 窗口标题，模糊匹配 |

**返回值：** 窗口句柄，没找到返回0

---

#### GetClientRect

获取窗口客户区域在屏幕上的位置

```go
func (o *OP) GetClientRect(hwnd int) (x1, y1, x2, y2 int)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |

**返回值：** x1, y1, x2, y2 客户区坐标

---

#### GetForegroundFocus

获取顶层活动窗口中具有输入焦点的窗口句柄

```go
func (o *OP) GetForegroundFocus() int
```

**返回值：** 窗口句柄

---

#### GetForegroundWindow

获取顶层活动窗口

```go
func (o *OP) GetForegroundWindow() int
```

**返回值：** 窗口句柄

---

#### GetMousePointWindow

获取鼠标指向的可见窗口句柄

```go
func (o *OP) GetMousePointWindow() int
```

**返回值：** 窗口句柄

---

#### GetPointWindow

获取给定坐标的可见窗口句柄

```go
func (o *OP) GetPointWindow(x, y int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int | 屏幕X坐标 |
| y | int | 屏幕Y坐标 |

**返回值：** 窗口句柄

---

#### GetProcessInfo

根据指定的pid获取进程详细信息

```go
func (o *OP) GetProcessInfo(pid int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| pid | int | 进程pid |

**返回值：** 返回格式"进程名|进程路径|cpu|内存"

---

#### GetSpecialWindow

获取特殊窗口

```go
func (o *OP) GetSpecialWindow(flag int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| flag | int | 0获取桌面窗口，1获取任务栏窗口 |

**返回值：** 窗口句柄

---

#### GetWindow

获取给定窗口相关的窗口句柄

```go
func (o *OP) GetWindow(hwnd, flag int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| flag | int | 0获取父窗口，1获取第一个子窗口，2获取First窗口，3获取Last窗口，4获取下一个窗口，5获取上一个窗口，6获取拥有者窗口，7获取顶层窗口 |

**返回值：** 窗口句柄

---

#### GetWindowProcessId

获取指定窗口所在的进程ID

```go
func (o *OP) GetWindowProcessId(hwnd int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |

**返回值：** 进程ID

---

#### GetWindowProcessPath

获取指定窗口所在的进程的exe文件全路径

```go
func (o *OP) GetWindowProcessPath(hwnd int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |

**返回值：** 进程所在的全路径

---

#### ScreenToClient

把屏幕坐标转换为窗口坐标

```go
func (o *OP) ScreenToClient(hwnd int, x, y *int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| x | *int | 屏幕X坐标（变参，会被修改为窗口坐标） |
| y | *int | 屏幕Y坐标（变参，会被修改为窗口坐标） |

**返回值：** 0表示失败，1表示成功

---

#### SetClientSize

设置窗口客户区域的宽度和高度

```go
func (o *OP) SetClientSize(hwnd, width, height int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| width | int | 宽度 |
| height | int | 高度 |

**返回值：** 0表示失败，1表示成功

---

#### SetWindowSize

设置窗口的大小

```go
func (o *OP) SetWindowSize(hwnd, width, height int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| width | int | 宽度 |
| height | int | 高度 |

**返回值：** 0表示失败，1表示成功

---

#### SetWindowText

设置窗口的标题

```go
func (o *OP) SetWindowText(hwnd int, title string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| title | string | 标题 |

**返回值：** 0表示失败，1表示成功

---

#### SetWindowTransparent

设置窗口的透明度

```go
func (o *OP) SetWindowTransparent(hwnd, trans int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |
| trans | int | 透明度取值(0-255)，越小透明度越大，0为完全透明，255为完全不透明 |

**返回值：** 0表示失败，1表示成功

---

#### SendPaste

向指定窗口发送粘贴命令

```go
func (o *OP) SendPaste(hwnd int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 指定的窗口句柄 |

**返回值：** 0表示失败，1表示成功

---

### Mouse 鼠标操作补充

#### MoveR

鼠标相对于上次的位置移动rx,ry

```go
func (o *OP) MoveR(rx, ry int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| rx | int | 相对于上次的X偏移 |
| ry | int | 相对于上次的Y偏移 |

**返回值：** 0表示失败，1表示成功

---

#### MoveToEx

把鼠标移动到目的范围内的任意一点

```go
func (o *OP) MoveToEx(x, y, w, h int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int | X坐标 |
| y | int | Y坐标 |
| w | int | 宽度(从x计算起) |
| h | int | 高度(从y计算起) |

**返回值：** 返回要移动到的目标点，格式为"x,y"

---

#### MiddleClick

按下鼠标中键

```go
func (o *OP) MiddleClick() int
```

**返回值：** 0表示失败，1表示成功

---

#### MiddleDown

按住鼠标中键

```go
func (o *OP) MiddleDown() int
```

**返回值：** 0表示失败，1表示成功

---

#### MiddleUp

弹起鼠标中键

```go
func (o *OP) MiddleUp() int
```

**返回值：** 0表示失败，1表示成功

---

#### SetMouseDelay

设置鼠标单击或双击时，鼠标按下和弹起之间的时间间隔

```go
func (o *OP) SetMouseDelay(mouseType string, delay int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| mouseType | string | 鼠标类型，取值: "normal" \| "windows" \| "dx" |
| delay | int | 指定鼠标按下和弹起之间的时间间隔，单位毫秒 |

**默认值：**
- "normal": 30ms
- "windows": 10ms
- "dx": 40ms

**返回值：** 0表示失败，1表示成功

---

### Keypad 键盘操作补充

#### GetKeyState

获取指定的按键状态（前台信息，不是后台）

```go
func (o *OP) GetKeyState(vkCode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| vkCode | int | 虚拟按键码 |

**返回值：** 0表示失败，1表示成功

---

#### WaitKey

等待指定的按键按下（前台，不是后台）

```go
func (o *OP) WaitKey(vkCode int, timeOut string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| vkCode | int | 虚拟按键码，当此值为0时表示等待任意按键。鼠标左键是1，鼠标右键是2，鼠标中键是4 |
| timeOut | string | 等待多久，单位毫秒。如果是0，表示一直等待。注意：官方文档中此参数为string类型 |

**返回值：** 0表示失败，1表示成功

---

#### SetKeypadDelay

设置按键时，键盘按下和弹起之间的时间间隔

```go
func (o *OP) SetKeypadDelay(keypadType string, delay int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| keypadType | string | 键盘类型，取值: "normal" \| "normal2" \| "windows" \| "dx" |
| delay | int | 指定键盘按下和弹起之间的时间间隔，单位毫秒 |

**默认值：**
- "normal": 30ms
- "normal2": 30ms
- "windows": 10ms
- "dx": 50ms

**返回值：** 0表示失败，1表示成功

---

### OCR 文字识别补充

#### SetMemDict

设置内存字库文件

```go
func (o *OP) SetMemDict(index, data, size int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int | 字库的序号，范围0-9 |
| data | int | 字库内容数据的内存地址（整数指针） |
| size | int | 字库大小 |

**返回值：** 0表示失败，1表示成功

---

#### OcrEx

识别屏幕范围内的字符串，返回识别到的字符串以及每个字符的坐标

```go
func (o *OP) OcrEx(x1, y1, x2, y2 int, colorFormat string, sim float64) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 识别区域 |
| colorFormat | string | 颜色格式串，比如"FFFFFF-000000\|CCCCCC-000000"每种颜色用"\|"分割 |
| sim | float64 | 相似度，取值范围0.1-1.0 |

**返回值：** 返回识别到的字符串以及坐标，格式为"char$x$y\|char$x$y\|..."

---

#### OcrAuto

识别屏幕范围内的字符串，自动二值化，无需指定颜色

```go
func (o *OP) OcrAuto(x1, y1, x2, y2 int, sim float64) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 识别区域 |
| sim | float64 | 相似度，取值范围0.1-1.0 |

**返回值：** 返回识别到的字符串

**说明：** 适用于字体颜色和背景相差较大的场合

---

#### OcrFromFile

从文件中识别图片

```go
func (o *OP) OcrFromFile(fileName, colorFormat string, sim float64) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| fileName | string | 文件名 |
| colorFormat | string | 颜色格式串 |
| sim | float64 | 相似度，取值范围0.1-1.0 |

**返回值：** 返回识别到的字符串

---

#### OcrAutoFromFile

从文件中识别图片，自动二值化，无需指定颜色

```go
func (o *OP) OcrAutoFromFile(fileName string, sim float64) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| fileName | string | 文件名 |
| sim | float64 | 相似度，取值范围0.1-1.0 |

**返回值：** 返回识别到的字符串

---

#### FindLine

在指定的屏幕坐标范围内，查找指定颜色的直线

```go
func (o *OP) FindLine(x1, y1, x2, y2 int, colorFormat string, sim float64) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| colorFormat | string | 颜色格式串 |
| sim | float64 | 相似度，取值范围0.1-1.0 |

**返回值：** 返回识别到的结果

---

### System 系统命令补充

#### RunApp

运行可执行文件，可指定模式

```go
func (o *OP) RunApp(appPath string, mode int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| appPath | string | 指定的可执行程序全路径 |
| mode | int | 取值0表示普通模式，1表示加强模式 |

**返回值：** 0表示失败，1表示成功

---

#### WinExec

运行可执行文件，可指定显示模式

```go
func (o *OP) WinExec(cmdLine string, cmdShow int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| cmdLine | string | 指定的可执行程序全路径 |
| cmdShow | int | 取值0表示隐藏，1表示用最近的大小和位置显示并激活 |

**返回值：** 0表示失败，1表示成功

---

#### GetCmdStr

运行命令行并返回结果

```go
func (o *OP) GetCmdStr(cmdLine string, milliseconds int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| cmdLine | string | 指定的命令行 |
| milliseconds | int | 等待的时间（毫秒） |

**返回值：** cmd输出的字符

---

#### Delay

实现一个指定毫秒数的延迟，同时确保在此期间不会阻塞用户界面（UI）操作

```go
func (o *OP) Delay(ms int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| ms | int | 指定延迟的时间，单位为毫秒 |

**返回值：** 0表示失败，1表示成功

---

#### Delays

实现一个指定毫秒数的延迟，同时确保在此期间不会阻塞用户界面（UI）操作，随机选择延迟时间

```go
func (o *OP) Delays(msMin, msMax int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| msMin | int | 指定延迟时间的最小值，单位为毫秒 |
| msMax | int | 指定延迟时间的最大值，单位为毫秒 |

**返回值：** 0表示失败，1表示成功

**说明：** 函数将随机选择一个介于msMin和msMax之间的延迟时间

---

## 参数修正说明

以下函数的参数已根据官方文档进行修正：

### WaitKey

**修正前：** `func (o *OP) WaitKey(vkCode, timeOut int) int`

**修正后：** `func (o *OP) WaitKey(vkCode int, timeOut string) int`

**说明：** `timeOut` 参数类型从 `int` 改为 `string`，与官方文档一致。

---

### WriteData

**修正前：** `func (o *OP) WriteData(hwnd int, addr int64, data string) int`

**修正后：** `func (o *OP) WriteData(hwnd int, address, data string, size int) int`

**说明：**
- `address` 参数类型从 `int64` 改为 `string`
- 新增 `size` 参数

---

### ReadData

**修正前：** `func (o *OP) ReadData(hwnd int, addr int64, len int) string`

**修正后：** `func (o *OP) ReadData(hwnd int, address string, size int) string`

**说明：**
- `address` 参数类型从 `int64` 改为 `string`
- `len` 参数名改为 `size`

---

### AStarFindPath

**修正前：** `func (o *OP) AStarFindPath(mapData string, startX, startY, endX, endY int) string`

**修正后：** `func (o *OP) AStarFindPath(mapWidth, mapHeight int, disablePoints string, beginX, beginY, endX, endY int) string`

**说明：** 参数完全重新设计，与官方文档一致：
- `mapWidth`, `mapHeight`: 地图宽度和高度
- `disablePoints`: 不可通行的坐标，以"|"分割
- `beginX`, `beginY`: 源坐标
- `endX`, `endY`: 目的坐标

---

### FindNearestPos

**修正前：** `func (o *OP) FindNearestPos(mapData string, curX, curY, targetX, targetY int) string`

**修正后：** `func (o *OP) FindNearestPos(allPos string, posType, x, y int) string`

**说明：** 参数完全重新设计，与官方文档一致：
- `allPos`: 位置集合，以"|"分割的坐标列表
- `posType`: 类型
- `x`, `y`: 参考坐标

---

### FindColorBlock

**修正前：** `func (o *OP) FindColorBlock(x1, y1, x2, y2 int, color string, sim float64, w, h, dir int) (x, y int, found bool)`

**修正后：** `func (o *OP) FindColorBlock(x1, y1, x2, y2 int, color string, sim float64, count, width, height int) (x, y int, found bool)`

**说明：**
- 参数从 `w, h, dir` 改为 `count, width, height`
- 新增 `count` 参数（在宽度为width,高度为height的颜色块中，符合color颜色的最小数量）
- 移除 `dir` 参数
- 返回值改为通过变参指针获取

---

### FindColorBlockEx

**修正前：** `func (o *OP) FindColorBlockEx(x1, y1, x2, y2 int, color string, sim float64, w, h, dir int) string`

**修正后：** `func (o *OP) FindColorBlockEx(x1, y1, x2, y2 int, color string, sim float64, count, width, height int) string`

**说明：**
- 参数从 `w, h, dir` 改为 `count, width, height`
- 新增 `count` 参数
- 移除 `dir` 参数

---

### SetDisplayInput

**修正前：** `func (o *OP) SetDisplayInput(mode int) int`

**修正后：** `func (o *OP) SetDisplayInput(mode string) int`

**说明：** `mode` 参数类型从 `int` 改为 `string`，与官方文档一致。

---

### GetScreenData

**修正前：** `func (o *OP) GetScreenData(x1, y1, x2, y2 int) string`

**修正后：** `func (o *OP) GetScreenData(x1, y1, x2, y2 int) int`

**说明：** 返回值从 `string` 改为 `int`（数据指针）。

---

### GetScreenDataBmp

**修正前：** `func (o *OP) GetScreenDataBmp(x1, y1, x2, y2 int) string`

**修正后：** `func (o *OP) GetScreenDataBmp(x1, y1, x2, y2 int) (data, size int, ret int)`

**说明：**
- 返回值从 `string` 改为多返回值 `(data, size, ret)`
- `data`: 图片数据指针
- `size`: 图片数据长度
- `ret`: 0表示失败，1表示成功

---

### SetMemDict

**修正前：** `func (o *OP) SetMemDict(index int, data []byte, size int) int`

**修正后：** `func (o *OP) SetMemDict(index, data, size int) int`

**说明：** `data` 参数类型从 `[]byte` 改为 `int`（内存地址指针）。

---

## 新增补充函数（第二次补充）

以下函数为根据完整函数列表再次补充的函数：

### ImageProc 图片处理补充

#### FindMultiColor

根据指定的多点颜色查找颜色坐标

```go
func (o *OP) FindMultiColor(x1, y1, x2, y2 int, firstColor, offsetColor string, sim float64, dir int) (x, y int, found bool)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| firstColor | string | 第一个点的颜色，格式为"RRGGBB" |
| offsetColor | string | 其他点的颜色偏移，格式为"x1\|y1\|RRGGBB,x2\|y2\|RRGGBB,..." |
| sim | float64 | 相似度，取值范围 0.0-1.0 |
| dir | int | 查找方向，0 表示从左到右从上到下 |

**返回值：** x, y 找到的颜色坐标，found 是否找到

---

#### FindMultiColorEx

根据指定的多点颜色查找所有颜色坐标

```go
func (o *OP) FindMultiColorEx(x1, y1, x2, y2 int, firstColor, offsetColor string, sim float64, dir int) string
```

**参数说明：** 同 FindMultiColor

**返回值：** 所有找到的坐标字符串，格式为"x,y\|x,y\|..."

---

#### FindPicExS

查找多个图片，同 FindPicEx

```go
func (o *OP) FindPicExS(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) string
```

**参数说明：** 同 FindPic

**返回值：** 所有找到的图片信息，格式为"idx,x,y\|idx,x,y\|..."

---

#### FindColorBlock

查找指定区域内的颜色块

```go
func (o *OP) FindColorBlock(x1, y1, x2, y2 int, color string, sim float64, count, width, height int) (x, y int, found bool)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| color | string | 颜色格式串，比如"FFFFFF-000000\|CCCCCC-000000"每种颜色用"\|"分割 |
| sim | float64 | 相似度，取值范围 0.1-1.0 |
| count | int | 在宽度为width,高度为height的颜色块中，符合color颜色的最小数量 |
| width | int | 颜色块的宽度 |
| height | int | 颜色块的高度 |

**返回值：** x, y 找到的颜色块坐标，found 是否找到

---

#### FindColorBlockEx

查找指定区域内的所有颜色块

```go
func (o *OP) FindColorBlockEx(x1, y1, x2, y2 int, color string, sim float64, count, width, height int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 查找区域 |
| color | string | 颜色格式串，比如"FFFFFF-000000\|CCCCCC-000000"每种颜色用"\|"分割 |
| sim | float64 | 相似度，取值范围 0.1-1.0 |
| count | int | 在宽度为width,高度为height的颜色块中，符合color颜色的最小数量 |
| width | int | 颜色块的宽度 |
| height | int | 颜色块的高度 |

**返回值：** 所有找到的坐标字符串，格式为"x,y\|x,y\|..."

---

#### SetDisplayInput

设置图像输入方式

```go
func (o *OP) SetDisplayInput(mode string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| mode | string | 图色输入模式，取值:<br>- "screen": 默认的模式，表示使用显示器或者后台窗口<br>- "pic:文件名": 指定输入模式为指定的图片<br>- "mem:addr": 指定输入模式为内存中的图片 |

**返回值：** 0表示失败，1表示成功

---

#### LoadPic

预加载图片

```go
func (o *OP) LoadPic(picName string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| picName | string | 图片文件名 |

**返回值：** 0表示失败，1表示成功

---

#### FreePic

释放图片

```go
func (o *OP) FreePic(picName string) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| picName | string | 图片文件名，空字符串表示释放所有图片 |

**返回值：** 0表示失败，1表示成功

---

#### GetScreenData

获取指定区域的图像数据

```go
func (o *OP) GetScreenData(x1, y1, x2, y2 int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 区域坐标 |

**返回值：** 返回的是指定区域的二进制图片颜色数据指针，每个颜色是4个字节,表示方式为(00RRGGBB)

---

#### GetScreenDataBmp

获取指定区域的图像（BMP格式）

```go
func (o *OP) GetScreenDataBmp(x1, y1, x2, y2 int) (data, size int, ret int)
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1, y1, x2, y2 | int | 区域坐标 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| data | int | 返回图片的数据指针 |
| size | int | 返回图片的数据长度 |
| ret | int | 0表示失败，1表示成功 |

---

#### MatchPicName

使用通配符并获取文件集合

```go
func (o *OP) MatchPicName(picName string) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| picName | string | 图片文件名，支持通配符如 "*.bmp" |

**返回值：** 匹配的文件名列表，格式为 "file1\|file2\|file3"

---

#### LoadMemPic

从内存中加载图片，并将加载结果返回

```go
func (o *OP) LoadMemPic(fileName string, data, size int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| fileName | string | 图片的文件名 |
| data | int | 图像数据的内存地址（整数指针） |
| size | int | 图像数据的大小 |

**返回值：** 0表示失败，1表示成功

---

### Memory 内存操作

#### WriteData

向某进程写入数据

```go
func (o *OP) WriteData(hwnd int, address, data string, size int) int
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄，用于指定要在哪个窗口内写入数据 |
| address | string | 写入数据的地址（字符串类型） |
| data | string | 写入的数据 |
| size | int | 写入的数据的大小 |

**返回值：** 0表示失败，1表示成功

---

#### ReadData

读取数据

```go
func (o *OP) ReadData(hwnd int, address string, size int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int | 窗口句柄，用于指定要从哪个窗口内读取数据 |
| address | string | 表示要读取数据的地址（字符串类型） |
| size | int | 要读取的数据的大小 |

**返回值：** 读取到的数值

---

### Algorithm 算法

#### AStarFindPath

A星算法

```go
func (o *OP) AStarFindPath(mapWidth, mapHeight int, disablePoints string, beginX, beginY, endX, endY int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| mapWidth | int | 地图宽度 |
| mapHeight | int | 地图高度 |
| disablePoints | string | 不可通行的坐标，以"\|"分割，例如:"10,15\|20,30" |
| beginX, beginY | int | 源坐标 |
| endX, endY | int | 目的坐标 |

**返回值：** 找到的路径结果，格式为"x1,y1\|x2,y2\|..."

---

#### FindNearestPos

查找最近的位置

```go
func (o *OP) FindNearestPos(allPos string, posType, x, y int) string
```

**参数说明：**

| 参数名 | 类型 | 说明 |
|--------|------|------|
| allPos | string | 位置集合，以"\|"分割的坐标列表，例如:"10,15\|20,30\|50,60" |
| posType | int | 类型 |
| x, y | int | 参考坐标 |

**返回值：** 最接近指定坐标 (x, y) 的位置，格式为"x,y"

---

## 注意事项

1. **资源释放**：使用完 OP 实例后必须调用 `Release()` 释放资源
2. **线程安全**：OP 插件不是线程安全的，不要在多个 goroutine 中同时使用同一个实例
3. **管理员权限**：某些操作需要管理员权限才能执行
4. **路径设置**：使用图片查找功能前，确保设置了正确的图片路径
5. **字库设置**：使用 OCR 功能前，需要先设置并选择字库
