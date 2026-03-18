================================================================================
                         OP-Go - OP 插件 Go 语言封装库
================================================================================

简介
----
OP-Go 是 OP (Operator & Open) 插件的 Go 语言封装库，专为 Windows 平台设计，
提供屏幕读取、输入模拟、图像处理、OCR 文字识别等自动化功能。

项目地址: https://github.com/yuan71058/OP-Go
OP 官方: https://github.com/WallBreaker2/op


功能特性
--------
[√] 窗口操作 - 查找窗口、获取窗口信息、移动窗口、设置窗口状态
[√] 鼠标操作 - 移动、点击、滚轮、拖拽等
[√] 键盘操作 - 按键、组合键、字符串输入
[√] 图色操作 - 截图、找图、找色、取色、比色
[√] OCR 识别 - 文字识别、字库支持
[√] 后台绑定 - 支持后台窗口操作（DX 模式）
[√] 内存操作 - 进程内存读写
[√] 系统命令 - 剪贴板操作、运行程序等


安装要求
--------
1. Windows 操作系统（仅支持 Windows）
2. Go 1.21 或更高版本
3. OP 插件 DLL 文件:
   - op_x64.dll 或 op_x86.dll（根据系统架构选择）
   - tools_64.dll 或 tools.dll（免注册方式需要）


安装方法
--------
    go get github.com/yuan71058/OP-Go


快速开始
--------

### 基础使用示例

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

        // 查找窗口
        hwnd := opInst.FindWindow("", "记事本")
        if hwnd != 0 {
            fmt.Printf("找到窗口，句柄: %d\n", hwnd)
        }
    }


### Service 模式（推荐）

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
    }


目录结构
--------
    E:\SRC\gop\
    ├── examples/           # 示例代码
    │   └── main.go        # 使用示例
    ├── op/                # OP 插件文件
    │   ├── op-0.4.5.zip              # OP 插件主文件
    │   ├── OPTestTool-0.4.5.zip      # OP 测试工具
    │   └── op-wiki-离线文档0.45.exe  # OP 离线文档
    ├── op.go              # OP 核心封装代码
    ├── service.go         # Service 模式封装
    ├── op_api.md          # 详细 API 文档
    ├── README.md          # Markdown 格式说明
    └── README.txt         # 本文件（文本格式说明）


主要 API 概览
-------------

### 基础函数
    NewOP(dllPath string) (*OP, error)     - 创建 OP 实例
    (*OP) Release()                         - 释放 OP 实例
    (*OP) Ver() string                      - 获取版本号
    (*OP) SetPath(path string) int          - 设置全局路径
    (*OP) GetPath() string                  - 获取全局路径

### 窗口操作
    (*OP) FindWindow(className, title string) int                    - 查找窗口
    (*OP) GetWindowTitle(hwnd int) string                            - 获取窗口标题
    (*OP) GetWindowRect(hwnd int) (x1, y1, x2, y2 int)               - 获取窗口位置
    (*OP) GetClientSize(hwnd int) (width, height int)                - 获取客户区大小
    (*OP) SetWindowState(hwnd, flag int) int                         - 设置窗口状态

### 鼠标操作
    (*OP) MoveTo(x, y int) int              - 移动鼠标
    (*OP) LeftClick() int                   - 左键单击
    (*OP) LeftDoubleClick() int             - 左键双击
    (*OP) RightClick() int                  - 右键单击
    (*OP) GetCursorPos() (x, y int)         - 获取鼠标位置

### 键盘操作
    (*OP) KeyPress(vkCode int) int                  - 按键（虚拟键码）
    (*OP) KeyPressChar(keyStr string) int           - 按键（字符）
    (*OP) KeyDownChar(keyStr string) int            - 按住键
    (*OP) KeyUpChar(keyStr string) int              - 弹起键

### 图色操作
    (*OP) Capture(x1, y1, x2, y2 int, file string) int              - 截图
    (*OP) GetColor(x, y int) string                                 - 取色
    (*OP) FindColor(x1, y1, x2, y2 int, color string, sim float64, dir int) (x, y int, found bool)  - 找色
    (*OP) FindPic(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) (x, y int, found bool)  - 找图

### OCR 文字识别
    (*OP) SetDict(index int, file string) int                       - 设置字库
    (*OP) UseDict(index int) int                                    - 选择字库
    (*OP) Ocr(x1, y1, x2, y2 int, colorFormat string, sim float64) string  - 识别文字
    (*OP) FindStr(x1, y1, x2, y2 int, str, colorFormat string, sim float64) (ret, x, y int)  - 查找字符串


虚拟键码常量
------------
    VK_LBUTTON  = 0x01     // 鼠标左键
    VK_RBUTTON  = 0x02     // 鼠标右键
    VK_MBUTTON  = 0x04     // 鼠标中键
    VK_BACK     = 0x08     // Backspace
    VK_TAB      = 0x09     // Tab
    VK_RETURN   = 0x0D     // Enter
    VK_SHIFT    = 0x10     // Shift
    VK_CONTROL  = 0x11     // Ctrl
    VK_MENU     = 0x12     // Alt
    VK_ESCAPE   = 0x1B     // Esc
    VK_SPACE    = 0x20     // 空格
    VK_LEFT     = 0x25     // 左箭头
    VK_UP       = 0x26     // 上箭头
    VK_RIGHT    = 0x27     // 右箭头
    VK_DOWN     = 0x28     // 下箭头


注意事项
--------
1. 仅支持 Windows 平台
2. 某些功能可能需要管理员权限
3. 确保 DLL 文件路径正确
4. 使用完毕后务必调用 Release() 释放资源
5. 免注册方式需要将 tools_64.dll 或 tools.dll 放在 DLL 同目录


相关链接
--------
- 项目主页: https://github.com/yuan71058/OP-Go
- OP 官方 GitHub: https://github.com/WallBreaker2/op
- OP 官方文档: https://github.com/WallBreaker2/op/wiki
- Go-OLE 库: https://github.com/go-ole/go-ole


许可证
------
MIT License


作者
----
GitHub: https://github.com/yuan71058

================================================================================
                         感谢使用 OP-Go！
================================================================================
