// Package op 提供了 OP 插件的 Go 语言封装
// OP (Operator & Open) 是一款专为 Windows 设计的开源自动化插件
// 提供屏幕读取、输入模拟、图像处理、OCR 等功能
// 项目地址：https://github.com/WallBreaker2/op
package op

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// myDISPPARAMS 自定义的 DISPPARAMS 结构体，字段可导出
type myDISPPARAMS struct {
	Rgvarg            uintptr
	RgdispidNamedArgs uintptr
	CArgs             uint32
	CNamedArgs        uint32
}

// opLogger OP 插件专用日志
var opLogger *log.Logger
var opLogFile *os.File

// initOPLog 初始化 OP 日志
func initOPLog() {
	if opLogFile != nil {
		return
	}

	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		return
	}
	execDir := filepath.Dir(execPath)
	logDir := filepath.Join(execDir, "config", "log")

	// 创建日志目录
	os.MkdirAll(logDir, 0755)

	// 创建日志文件
	logPath := filepath.Join(logDir, "op_debug.log")
	opLogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	opLogger = log.New(opLogFile, "", log.LstdFlags|log.Lmicroseconds)
}

// logf 记录 OP 日志
func logf(format string, v ...interface{}) {
	initOPLog()
	if opLogger != nil {
		opLogger.Printf(format, v...)
		opLogFile.Sync()
	}
	// 同时输出到控制台
	fmt.Printf("[OP] "+format+"\n", v...)
}

// 全局变量，用于免注册加载
var (
	toolsDLL *syscall.LazyDLL
	setupW   *syscall.LazyProc
)

// OP 表示一个 OP 插件实例
// 通过 COM 接口调用 OP 插件的功能
type OP struct {
	object *ole.IDispatch // COM 对象的 IDispatch 接口
}

// init 初始化 tools DLL
func init() {
	// 尝试加载 tools_64.dll 或 tools.dll（免注册方式）
	toolsDLL = syscall.NewLazyDLL("tools_64.dll")
	setupW = toolsDLL.NewProc("setupW")

	// 如果 tools_64.dll 不存在，尝试 tools.dll
	if setupW.Find() != nil {
		toolsDLL = syscall.NewLazyDLL("tools.dll")
		setupW = toolsDLL.NewProc("setupW")
	}
}

// SetupW 使用 Unicode 编码注册 OP 插件（免注册方式）
func SetupW(path string) bool {
	p0, _ := syscall.UTF16PtrFromString(path)
	ret, _, _ := setupW.Call(uintptr(unsafe.Pointer(p0)))
	return ret != 0
}

// NewOP 创建一个新的 OP 插件实例
// 参数 dllPath: OP 插件 DLL 文件的完整路径
// 返回值:
//   - *OP: OP 插件实例
//   - error: 错误信息
//
// 支持两种使用方式：
//  1. 免注册方式：需要 tools_64.dll 或 tools.dll 文件
//  2. 注册方式：需要先使用 regsvr32 op_x64.dll 注册
func NewOP(dllPath string) (*OP, error) {
	// 获取 DLL 所在目录
	dllDir := filepath.Dir(dllPath)

	// 尝试加载的 tools DLL 列表（按优先级）
	// tools_64.dll: 64位专用
	// tools_86.dll: 32位专用
	// tools.dll: 默认（通常为64位）
	toolsNames := []string{"tools_64.dll", "tools_86.dll", "tools.dll"}

	for _, toolsName := range toolsNames {
		toolsPath := filepath.Join(dllDir, toolsName)
		if _, err := os.Stat(toolsPath); err == nil {
			dll, err := syscall.LoadDLL(toolsPath)
			if err == nil {
				// 先尝试 setupA (ANSI 版本)
				proc, err := dll.FindProc("setupA")
				if err == nil {
					// 将路径转换为字节数组（ANSI）
					pathBytes := []byte(dllPath + "\x00")
					ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&pathBytes[0])))
					if ret != 0 {
						return newOPWithCOM()
					}
				}
				// 再尝试 setupW (Unicode 版本)
				proc, err = dll.FindProc("setupW")
				if err == nil {
					p0, _ := syscall.UTF16PtrFromString(dllPath)
					ret, _, _ := proc.Call(uintptr(unsafe.Pointer(p0)))
					if ret != 0 {
						return newOPWithCOM()
					}
				}
				dll.Release()
			}
		}
	}

	// 尝试全局 tools DLL（如果已加载）
	if setupW.Find() == nil {
		ok := SetupW(dllPath)
		if ok {
			return newOPWithCOM()
		}
	}

	// 免注册失败，尝试注册方式
	return newOPWithCOM()
}

// newOPWithCOM 使用 COM 方式创建 OP 插件实例
func newOPWithCOM() (*OP, error) {
	// 初始化 COM 库
	ole.CoInitialize(0)

	// 创建 OP 插件的 COM 对象
	unknown, err := oleutil.CreateObject("op.opsoft")
	if err != nil {
		ole.CoUninitialize()
		return nil, fmt.Errorf("创建 OP 对象失败：%v", err)
	}

	// 获取 IDispatch 接口
	dispatch, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		ole.CoUninitialize()
		return nil, fmt.Errorf("获取 IDispatch 失败：%v", err)
	}

	return &OP{object: dispatch}, nil
}

// Release 释放 OP 插件实例
// 必须在不再使用 OP 实例时调用，以释放 COM 资源
func (o *OP) Release() {
	if o.object != nil {
		o.object.Release()
		o.object = nil
	}
	ole.CoUninitialize()
}

// ==================== Base 基础函数 ====================

// Ver 获取 OP 插件的版本号
func (o *OP) Ver() string {
	ret, _ := oleutil.CallMethod(o.object, "Ver")
	return ret.ToString()
}

// SetPath 设置全局路径
func (o *OP) SetPath(path string) int {
	ret, _ := oleutil.CallMethod(o.object, "SetPath", path)
	return int(ret.Val)
}

// GetPath 获取全局路径
func (o *OP) GetPath() string {
	ret, _ := oleutil.CallMethod(o.object, "GetPath")
	return ret.ToString()
}

// GetBasePath 获取插件所在目录
func (o *OP) GetBasePath() string {
	ret, _ := oleutil.CallMethod(o.object, "GetBasePath")
	return ret.ToString()
}

// GetID 获取当前对象的 ID 值
func (o *OP) GetID() int {
	ret, _ := oleutil.CallMethod(o.object, "GetID")
	return int(ret.Val)
}

// GetLastError 获取最后的错误
func (o *OP) GetLastError() int {
	ret, _ := oleutil.CallMethod(o.object, "GetLastError")
	return int(ret.Val)
}

// SetShowErrorMsg 设置是否弹出错误信息
func (o *OP) SetShowErrorMsg(show int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetShowErrorMsg", show)
	return int(ret.Val)
}

// Sleep 休眠指定时间
func (o *OP) Sleep(ms int) int {
	ret, _ := oleutil.CallMethod(o.object, "Sleep", ms)
	return int(ret.Val)
}

// EnablePicCache 设置是否开启或关闭插件内部的图片缓存机制
func (o *OP) EnablePicCache(enable int) int {
	ret, _ := oleutil.CallMethod(o.object, "EnablePicCache", enable)
	return int(ret.Val)
}

// InjectDll 将指定的DLL注入到指定的进程中
// 参数:
//   - processName: 指定要注入DLL的进程名称
//   - dllName: 注入的DLL名称
//
// 返回值: 0表示失败，1表示成功
func (o *OP) InjectDll(processName, dllName string) int {
	ret, _ := oleutil.CallMethod(o.object, "InjectDll", processName, dllName)
	return int(ret.Val)
}

// CapturePre 取上次操作的图色区域，保存为文件(24位位图)
// 参数:
//   - file: 设置保存文件名，保存路径是SetPath设置的目录，也可以指定全路径
//
// 返回值: 0表示失败，1表示成功
func (o *OP) CapturePre(file string) int {
	ret, _ := oleutil.CallMethod(o.object, "CapturePre", file)
	return int(ret.Val)
}

// SetScreenDataMode 设置屏幕数据模式
// 参数:
//   - mode: 0表示从上到下(默认)，1表示从下到上
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetScreenDataMode(mode int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetScreenDataMode", mode)
	return int(ret.Val)
}

// ==================== Background 后台绑定 ====================

// BindWindow 绑定指定的窗口
func (o *OP) BindWindow(hwnd int, display, mouse, keypad string, mode int) int {
	ret, _ := oleutil.CallMethod(o.object, "BindWindow", hwnd, display, mouse, keypad, mode)
	return int(ret.Val)
}

// UnBindWindow 解除绑定窗口
func (o *OP) UnBindWindow() int {
	ret, _ := oleutil.CallMethod(o.object, "UnBindWindow")
	return int(ret.Val)
}

// GetBindWindow 获取当前对象已经绑定的窗口句柄
func (o *OP) GetBindWindow() int {
	ret, _ := oleutil.CallMethod(o.object, "GetBindWindow")
	return int(ret.Val)
}

// IsBind 判断当前对象是否已绑定窗口
func (o *OP) IsBind() int {
	ret, _ := oleutil.CallMethod(o.object, "IsBind")
	return int(ret.Val)
}

// ==================== Window 窗口操作 ====================

// EnumWindow 根据指定条件枚举系统中符合条件的窗口
func (o *OP) EnumWindow(parent int, title, className string, filter int) string {
	ret, _ := oleutil.CallMethod(o.object, "EnumWindow", parent, title, className, filter)
	return ret.ToString()
}

// FindWindow 查找符合类名或者标题名的顶层可见窗口
func (o *OP) FindWindow(className, title string) int {
	ret, _ := oleutil.CallMethod(o.object, "FindWindow", className, title)
	return int(ret.Val)
}

// FindWindowEx 查找符合类名或者标题名的窗口
func (o *OP) FindWindowEx(parent int, className, title string) int {
	ret, _ := oleutil.CallMethod(o.object, "FindWindowEx", parent, className, title)
	return int(ret.Val)
}

// GetWindowTitle 获取窗口的标题
func (o *OP) GetWindowTitle(hwnd int) string {
	ret, _ := oleutil.CallMethod(o.object, "GetWindowTitle", hwnd)
	return ret.ToString()
}

// GetWindowClass 获取窗口的类名
func (o *OP) GetWindowClass(hwnd int) string {
	ret, _ := oleutil.CallMethod(o.object, "GetWindowClass", hwnd)
	return ret.ToString()
}

// GetWindowRect 获取窗口在屏幕上的位置
// 参考大漠写法：使用 oleutil.CallMethod，调用后从 VARIANT.Val 读取值
func (o *OP) GetWindowRect(hwnd int) (x1, y1, x2, y2 int) {
	logf("[GetWindowRect] 开始调用, hwnd=%d", hwnd)

	// 创建输出参数 VARIANT
	x1Var := ole.NewVariant(ole.VT_I4, 0)
	y1Var := ole.NewVariant(ole.VT_I4, 0)
	x2Var := ole.NewVariant(ole.VT_I4, 0)
	y2Var := ole.NewVariant(ole.VT_I4, 0)

	logf("[GetWindowRect] 调用 COM 方法...")
	ret, err := oleutil.CallMethod(o.object, "GetWindowRect", hwnd, &x1Var, &y1Var, &x2Var, &y2Var)
	if err != nil {
		logf("[GetWindowRect] 调用失败: %v", err)
		return 0, 0, 0, 0
	}

	// 从 VARIANT.Val 读取值（参考大漠写法）
	x1Val := int32(x1Var.Val)
	y1Val := int32(y1Var.Val)
	x2Val := int32(x2Var.Val)
	y2Val := int32(y2Var.Val)

	// 清理 VARIANT
	x1Var.Clear()
	y1Var.Clear()
	x2Var.Clear()
	y2Var.Clear()

	logf("[GetWindowRect] 调用成功: ret=%d, x1=%d, y1=%d, x2=%d, y2=%d", int(ret.Val), x1Val, y1Val, x2Val, y2Val)
	logf("[GetWindowRect] 调用成功: hwnd=%d, x1=%d, y1=%d, x2=%d, y2=%d, width=%d, height=%d",
		hwnd, x1Val, y1Val, x2Val, y2Val, x2Val-x1Val, y2Val-y1Val)

	return int(x1Val), int(y1Val), int(x2Val), int(y2Val)
}

// GetClientSize 获取窗口客户区的宽度和高度
// 参考大漠写法：使用 oleutil.CallMethod，调用后从 VARIANT.Val 读取值
func (o *OP) GetClientSize(hwnd int) (width, height int) {
	logf("[GetClientSize] 开始调用, hwnd=%d", hwnd)

	// 创建输出参数 VARIANT
	widthVar := ole.NewVariant(ole.VT_I4, 0)
	heightVar := ole.NewVariant(ole.VT_I4, 0)

	logf("[GetClientSize] 调用 COM 方法...")
	ret, err := oleutil.CallMethod(o.object, "GetClientSize", hwnd, &widthVar, &heightVar)
	if err != nil {
		logf("[GetClientSize] 调用失败: %v", err)
		return 0, 0
	}

	// 从 VARIANT.Val 读取值（参考大漠写法）
	widthVal := int32(widthVar.Val)
	heightVal := int32(heightVar.Val)

	// 清理 VARIANT
	widthVar.Clear()
	heightVar.Clear()

	logf("[GetClientSize] 调用成功: ret=%d, width=%d, height=%d", int(ret.Val), widthVal, heightVal)
	logf("[GetClientSize] 调用成功: hwnd=%d, width=%d, height=%d", hwnd, widthVal, heightVal)

	return int(widthVal), int(heightVal)
}

// GetWindowState 获取指定窗口的一些属性
func (o *OP) GetWindowState(hwnd, flag int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetWindowState", hwnd, flag)
	return int(ret.Val)
}

// SetWindowState 设置窗口的状态
func (o *OP) SetWindowState(hwnd, flag int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetWindowState", hwnd, flag)
	return int(ret.Val)
}

// MoveWindow 移动指定窗口到指定位置
func (o *OP) MoveWindow(hwnd, x, y int) int {
	ret, _ := oleutil.CallMethod(o.object, "MoveWindow", hwnd, x, y)
	return int(ret.Val)
}

// SendString 向指定窗口发送文本数据
func (o *OP) SendString(hwnd int, str string) int {
	ret, _ := oleutil.CallMethod(o.object, "SendString", hwnd, str)
	return int(ret.Val)
}

// SendStringIme 向指定窗口发送文本数据（输入法方式）
func (o *OP) SendStringIme(hwnd int, str string) int {
	ret, _ := oleutil.CallMethod(o.object, "SendStringIme", hwnd, str)
	return int(ret.Val)
}

// EnumWindowByProcess 根据指定进程以及其它条件，枚举系统中符合条件的窗口
// 参数:
//   - processName: 进程名称
//   - title: 窗口的标题
//   - className: 窗口的类名
//   - filter: 过滤条件，1匹配标题，2匹配类名，4只匹配第一层子窗口，8匹配顶级窗口，16匹配可见窗口，32按打开顺序排列
//
// 返回值: 返回所有匹配到的窗口句柄字符串，格式为"hwnd1,hwnd2,..."
func (o *OP) EnumWindowByProcess(processName, title, className string, filter int) string {
	ret, _ := oleutil.CallMethod(o.object, "EnumWindowByProcess", processName, title, className, filter)
	return ret.ToString()
}

// EnumProcess 根据指定进程名，枚举系统中符合条件的进程PID
// 参数:
//   - name: 进程名称
//
// 返回值: 返回所有匹配的进程PID，格式为"pid1,pid2,..."
func (o *OP) EnumProcess(name string) string {
	ret, _ := oleutil.CallMethod(o.object, "EnumProcess", name)
	return ret.ToString()
}

// ClientToScreen 把窗口坐标转换为屏幕坐标
// 参数:
//   - hwnd: 指定的窗口句柄
//   - x: 窗口X坐标（变参，会被修改为屏幕坐标）
//   - y: 窗口Y坐标（变参，会被修改为屏幕坐标）
//
// 返回值: 0表示失败，1表示成功
func (o *OP) ClientToScreen(hwnd int, x, y *int) int {
	xVar := ole.NewVariant(ole.VT_I4, int64(*x))
	yVar := ole.NewVariant(ole.VT_I4, int64(*y))
	ret, _ := oleutil.CallMethod(o.object, "ClientToScreen", hwnd, &xVar, &yVar)
	*x = int(int32(xVar.Val))
	*y = int(int32(yVar.Val))
	xVar.Clear()
	yVar.Clear()
	return int(ret.Val)
}

// FindWindowByProcess 根据指定的进程名字，来查找可见窗口
// 参数:
//   - processName: 进程名，如"notepad.exe"，精确匹配但不区分大小写
//   - className: 窗口类名，模糊匹配
//   - title: 窗口标题，模糊匹配
//
// 返回值: 窗口句柄，没找到返回0
func (o *OP) FindWindowByProcess(processName, className, title string) int {
	ret, _ := oleutil.CallMethod(o.object, "FindWindowByProcess", processName, className, title)
	return int(ret.Val)
}

// FindWindowByProcessId 根据指定的进程Id，来查找可见窗口
// 参数:
//   - processId: 进程id
//   - className: 窗口类名，模糊匹配
//   - title: 窗口标题，模糊匹配
//
// 返回值: 窗口句柄，没找到返回0
func (o *OP) FindWindowByProcessId(processId int, className, title string) int {
	ret, _ := oleutil.CallMethod(o.object, "FindWindowByProcessId", processId, className, title)
	return int(ret.Val)
}

// GetClientRect 获取窗口客户区域在屏幕上的位置
// 参数:
//   - hwnd: 指定的窗口句柄
//
// 返回值: x1, y1, x2, y2 客户区坐标
func (o *OP) GetClientRect(hwnd int) (x1, y1, x2, y2 int) {
	x1Var := ole.NewVariant(ole.VT_I4, 0)
	y1Var := ole.NewVariant(ole.VT_I4, 0)
	x2Var := ole.NewVariant(ole.VT_I4, 0)
	y2Var := ole.NewVariant(ole.VT_I4, 0)
	ret, _ := oleutil.CallMethod(o.object, "GetClientRect", hwnd, &x1Var, &y1Var, &x2Var, &y2Var)
	x1 = int(int32(x1Var.Val))
	y1 = int(int32(y1Var.Val))
	x2 = int(int32(x2Var.Val))
	y2 = int(int32(y2Var.Val))
	x1Var.Clear()
	y1Var.Clear()
	x2Var.Clear()
	y2Var.Clear()
	_ = ret
	return
}

// GetForegroundFocus 获取顶层活动窗口中具有输入焦点的窗口句柄
// 返回值: 窗口句柄
func (o *OP) GetForegroundFocus() int {
	ret, _ := oleutil.CallMethod(o.object, "GetForegroundFocus")
	return int(ret.Val)
}

// GetForegroundWindow 获取顶层活动窗口
// 返回值: 窗口句柄
func (o *OP) GetForegroundWindow() int {
	ret, _ := oleutil.CallMethod(o.object, "GetForegroundWindow")
	return int(ret.Val)
}

// GetMousePointWindow 获取鼠标指向的可见窗口句柄
// 返回值: 窗口句柄
func (o *OP) GetMousePointWindow() int {
	ret, _ := oleutil.CallMethod(o.object, "GetMousePointWindow")
	return int(ret.Val)
}

// GetPointWindow 获取给定坐标的可见窗口句柄
// 参数:
//   - x: 屏幕X坐标
//   - y: 屏幕Y坐标
//
// 返回值: 窗口句柄
func (o *OP) GetPointWindow(x, y int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetPointWindow", x, y)
	return int(ret.Val)
}

// GetProcessInfo 根据指定的pid获取进程详细信息
// 参数:
//   - pid: 进程pid
//
// 返回值: 返回格式"进程名|进程路径|cpu|内存"
func (o *OP) GetProcessInfo(pid int) string {
	ret, _ := oleutil.CallMethod(o.object, "GetProcessInfo", pid)
	return ret.ToString()
}

// GetSpecialWindow 获取特殊窗口
// 参数:
//   - flag: 0获取桌面窗口，1获取任务栏窗口
//
// 返回值: 窗口句柄
func (o *OP) GetSpecialWindow(flag int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetSpecialWindow", flag)
	return int(ret.Val)
}

// GetWindow 获取给定窗口相关的窗口句柄
// 参数:
//   - hwnd: 指定的窗口句柄
//   - flag: 0获取父窗口，1获取第一个子窗口，2获取First窗口，3获取Last窗口，4获取下一个窗口，5获取上一个窗口，6获取拥有者窗口，7获取顶层窗口
//
// 返回值: 窗口句柄
func (o *OP) GetWindow(hwnd, flag int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetWindow", hwnd, flag)
	return int(ret.Val)
}

// GetWindowProcessId 获取指定窗口所在的进程ID
// 参数:
//   - hwnd: 指定的窗口句柄
//
// 返回值: 进程ID
func (o *OP) GetWindowProcessId(hwnd int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetWindowProcessId", hwnd)
	return int(ret.Val)
}

// GetWindowProcessPath 获取指定窗口所在的进程的exe文件全路径
// 参数:
//   - hwnd: 指定的窗口句柄
//
// 返回值: 进程所在的全路径
func (o *OP) GetWindowProcessPath(hwnd int) string {
	ret, _ := oleutil.CallMethod(o.object, "GetWindowProcessPath", hwnd)
	return ret.ToString()
}

// ScreenToClient 把屏幕坐标转换为窗口坐标
// 参数:
//   - hwnd: 指定的窗口句柄
//   - x: 屏幕X坐标（变参，会被修改为窗口坐标）
//   - y: 屏幕Y坐标（变参，会被修改为窗口坐标）
//
// 返回值: 0表示失败，1表示成功
func (o *OP) ScreenToClient(hwnd int, x, y *int) int {
	xVar := ole.NewVariant(ole.VT_I4, int64(*x))
	yVar := ole.NewVariant(ole.VT_I4, int64(*y))
	ret, _ := oleutil.CallMethod(o.object, "ScreenToClient", hwnd, &xVar, &yVar)
	*x = int(int32(xVar.Val))
	*y = int(int32(yVar.Val))
	xVar.Clear()
	yVar.Clear()
	return int(ret.Val)
}

// SetClientSize 设置窗口客户区域的宽度和高度
// 参数:
//   - hwnd: 指定的窗口句柄
//   - width: 宽度
//   - height: 高度
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetClientSize(hwnd, width, height int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetClientSize", hwnd, width, height)
	return int(ret.Val)
}

// SetWindowSize 设置窗口的大小
// 参数:
//   - hwnd: 指定的窗口句柄
//   - width: 宽度
//   - height: 高度
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetWindowSize(hwnd, width, height int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetWindowSize", hwnd, width, height)
	return int(ret.Val)
}

// SetWindowText 设置窗口的标题
// 参数:
//   - hwnd: 指定的窗口句柄
//   - title: 标题
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetWindowText(hwnd int, title string) int {
	ret, _ := oleutil.CallMethod(o.object, "SetWindowText", hwnd, title)
	return int(ret.Val)
}

// SetWindowTransparent 设置窗口的透明度
// 参数:
//   - hwnd: 指定的窗口句柄
//   - trans: 透明度取值(0-255)，越小透明度越大，0为完全透明，255为完全不透明
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetWindowTransparent(hwnd, trans int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetWindowTransparent", hwnd, trans)
	return int(ret.Val)
}

// SendPaste 向指定窗口发送粘贴命令
// 参数:
//   - hwnd: 指定的窗口句柄
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SendPaste(hwnd int) int {
	ret, _ := oleutil.CallMethod(o.object, "SendPaste", hwnd)
	return int(ret.Val)
}

// ==================== Mouse 鼠标操作 ====================

// MoveTo 把鼠标移动到目的点
func (o *OP) MoveTo(x, y int) int {
	ret, _ := oleutil.CallMethod(o.object, "MoveTo", x, y)
	return int(ret.Val)
}

// LeftClick 按下鼠标左键
func (o *OP) LeftClick() int {
	ret, _ := oleutil.CallMethod(o.object, "LeftClick")
	return int(ret.Val)
}

// LeftDoubleClick 双击鼠标左键
func (o *OP) LeftDoubleClick() int {
	ret, _ := oleutil.CallMethod(o.object, "LeftDoubleClick")
	return int(ret.Val)
}

// LeftDown 按住鼠标左键
func (o *OP) LeftDown() int {
	ret, _ := oleutil.CallMethod(o.object, "LeftDown")
	return int(ret.Val)
}

// LeftUp 弹起鼠标左键
func (o *OP) LeftUp() int {
	ret, _ := oleutil.CallMethod(o.object, "LeftUp")
	return int(ret.Val)
}

// RightClick 按下鼠标右键
func (o *OP) RightClick() int {
	ret, _ := oleutil.CallMethod(o.object, "RightClick")
	return int(ret.Val)
}

// RightDown 按住鼠标右键
func (o *OP) RightDown() int {
	ret, _ := oleutil.CallMethod(o.object, "RightDown")
	return int(ret.Val)
}

// RightUp 弹起鼠标右键
func (o *OP) RightUp() int {
	ret, _ := oleutil.CallMethod(o.object, "RightUp")
	return int(ret.Val)
}

// WheelDown 滚轮向下滚
func (o *OP) WheelDown() int {
	ret, _ := oleutil.CallMethod(o.object, "WheelDown")
	return int(ret.Val)
}

// WheelUp 滚轮向上滚
func (o *OP) WheelUp() int {
	ret, _ := oleutil.CallMethod(o.object, "WheelUp")
	return int(ret.Val)
}

// MoveR 鼠标相对于上次的位置移动rx,ry
// 参数:
//   - rx: 相对于上次的X偏移
//   - ry: 相对于上次的Y偏移
//
// 返回值: 0表示失败，1表示成功
func (o *OP) MoveR(rx, ry int) int {
	ret, _ := oleutil.CallMethod(o.object, "MoveR", rx, ry)
	return int(ret.Val)
}

// MoveToEx 把鼠标移动到目的范围内的任意一点
// 参数:
//   - x: X坐标
//   - y: Y坐标
//   - w: 宽度(从x计算起)
//   - h: 高度(从y计算起)
//
// 返回值: 返回要移动到的目标点，格式为"x,y"
func (o *OP) MoveToEx(x, y, w, h int) string {
	ret, _ := oleutil.CallMethod(o.object, "MoveToEx", x, y, w, h)
	return ret.ToString()
}

// MiddleClick 按下鼠标中键
// 返回值: 0表示失败，1表示成功
func (o *OP) MiddleClick() int {
	ret, _ := oleutil.CallMethod(o.object, "MiddleClick")
	return int(ret.Val)
}

// MiddleDown 按住鼠标中键
// 返回值: 0表示失败，1表示成功
func (o *OP) MiddleDown() int {
	ret, _ := oleutil.CallMethod(o.object, "MiddleDown")
	return int(ret.Val)
}

// MiddleUp 弹起鼠标中键
// 返回值: 0表示失败，1表示成功
func (o *OP) MiddleUp() int {
	ret, _ := oleutil.CallMethod(o.object, "MiddleUp")
	return int(ret.Val)
}

// SetMouseDelay 设置鼠标单击或双击时，鼠标按下和弹起之间的时间间隔
// 参数:
//   - mouseType: 鼠标类型，取值: "normal" | "windows" | "dx"
//   - delay: 指定鼠标按下和弹起之间的时间间隔，单位毫秒
//     当取值为"normal"，默认为30ms
//     当取值为"windows"，默认为10ms
//     当取值为"dx"，默认为40ms
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetMouseDelay(mouseType string, delay int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetMouseDelay", mouseType, delay)
	return int(ret.Val)
}

// GetCursorPos 获取鼠标位置
// 参考大漠写法：使用 oleutil.CallMethod，调用后从 VARIANT.Val 读取值
func (o *OP) GetCursorPos() (x, y int) {
	// 创建输出参数 VARIANT
	xVar := ole.NewVariant(ole.VT_I4, 0)
	yVar := ole.NewVariant(ole.VT_I4, 0)

	oleutil.CallMethod(o.object, "GetCursorPos", &xVar, &yVar)

	// 从 VARIANT.Val 读取值（参考大漠写法）
	xVal := int32(xVar.Val)
	yVal := int32(yVar.Val)

	// 清理 VARIANT
	xVar.Clear()
	yVar.Clear()

	return int(xVal), int(yVal)
}

// ==================== Keypad 键盘操作 ====================

// KeyPress 按下并弹起指定的虚拟键码
func (o *OP) KeyPress(vkCode int) int {
	ret, _ := oleutil.CallMethod(o.object, "KeyPress", vkCode)
	return int(ret.Val)
}

// KeyPressChar 按下并弹起指定的虚拟键码（字符串形式）
func (o *OP) KeyPressChar(keyStr string) int {
	ret, _ := oleutil.CallMethod(o.object, "KeyPressChar", keyStr)
	return int(ret.Val)
}

// KeyDown 按住指定的虚拟键码
func (o *OP) KeyDown(vkCode int) int {
	ret, _ := oleutil.CallMethod(o.object, "KeyDown", vkCode)
	return int(ret.Val)
}

// KeyDownChar 按住指定的虚拟键码（字符串形式）
func (o *OP) KeyDownChar(keyStr string) int {
	ret, _ := oleutil.CallMethod(o.object, "KeyDownChar", keyStr)
	return int(ret.Val)
}

// KeyUp 弹起虚拟键
func (o *OP) KeyUp(vkCode int) int {
	ret, _ := oleutil.CallMethod(o.object, "KeyUp", vkCode)
	return int(ret.Val)
}

// KeyUpChar 弹起虚拟键（字符串形式）
func (o *OP) KeyUpChar(keyStr string) int {
	ret, _ := oleutil.CallMethod(o.object, "KeyUpChar", keyStr)
	return int(ret.Val)
}

// GetKeyState 获取指定的按键状态（前台信息，不是后台）
// 参数:
//   - vkCode: 虚拟按键码
//
// 返回值: 0表示失败，1表示成功
func (o *OP) GetKeyState(vkCode int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetKeyState", vkCode)
	return int(ret.Val)
}

// WaitKey 等待指定的按键按下（前台，不是后台）
// 参数:
//   - vkCode: 虚拟按键码，当此值为0时表示等待任意按键。鼠标左键是1，鼠标右键是2，鼠标中键是4
//   - timeOut: 等待多久，单位毫秒。如果是0，表示一直等待。注意：官方文档中此参数为string类型
//
// 返回值: 0表示失败，1表示成功
func (o *OP) WaitKey(vkCode int, timeOut string) int {
	ret, _ := oleutil.CallMethod(o.object, "WaitKey", vkCode, timeOut)
	return int(ret.Val)
}

// SetKeypadDelay 设置按键时，键盘按下和弹起之间的时间间隔
// 参数:
//   - keypadType: 键盘类型，取值: "normal" | "normal2" | "windows" | "dx"
//   - delay: 指定键盘按下和弹起之间的时间间隔，单位毫秒
//     当取值为"normal"，默认为30ms
//     当取值为"normal2"，默认为30ms
//     当取值为"windows"，默认为10ms
//     当取值为"dx"，默认为50ms
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetKeypadDelay(keypadType string, delay int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetKeypadDelay", keypadType, delay)
	return int(ret.Val)
}

// ==================== ImageProc 图色操作 ====================

// Capture 抓取指定区域的图像，保存为文件
func (o *OP) Capture(x1, y1, x2, y2 int, file string) int {
	ret, _ := oleutil.CallMethod(o.object, "Capture", x1, y1, x2, y2, file)
	return int(ret.Val)
}

// FindPic 查找图片
// 参考大漠写法：使用 oleutil.CallMethod，调用后从 VARIANT.Val 读取值
func (o *OP) FindPic(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) (x, y int, found bool) {
	logf("[FindPic] 开始查找: x1=%d, y1=%d, x2=%d, y2=%d, picName=%s, deltaColor=%s, sim=%f, dir=%d", x1, y1, x2, y2, picName, deltaColor, sim, dir)

	// 创建输出参数 VARIANT
	xVar := ole.NewVariant(ole.VT_I4, 0)
	yVar := ole.NewVariant(ole.VT_I4, 0)

	ret, err := oleutil.CallMethod(o.object, "FindPic", x1, y1, x2, y2, picName, deltaColor, sim, dir, &xVar, &yVar)
	if err != nil {
		logf("[FindPic] 调用失败: %v", err)
		return 0, 0, false
	}

	// 从 VARIANT.Val 读取值（参考大漠写法）
	intX := int32(xVar.Val)
	intY := int32(yVar.Val)

	// 清理 VARIANT
	xVar.Clear()
	yVar.Clear()

	result := int(int32(ret.Val))
	logf("[FindPic] 查找结果: result=%d (原始值=%d), x=%d, y=%d", result, ret.Val, intX, intY)

	return int(intX), int(intY), result >= 0
}

// FindColor 查找指定区域内的颜色
// 参数:
//
//	x1, y1, x2, y2: 查找区域
//	color: 要查找的颜色，格式为"RRGGBB"，如"FF0000"表示红色
//	sim: 相似度，取值范围 0.0-1.0，1.0 表示完全匹配
//	dir: 查找方向，0 表示从左到右从上到下
//
// 返回值:
//
//	x, y: 找到的颜色坐标
//	found: 是否找到
func (o *OP) FindColor(x1, y1, x2, y2 int, color string, sim float64, dir int) (x, y int, found bool) {
	logf("[FindColor] 开始查找: x1=%d, y1=%d, x2=%d, y2=%d, color=%s, sim=%f, dir=%d", x1, y1, x2, y2, color, sim, dir)

	// 创建输出参数 VARIANT
	xVar := ole.NewVariant(ole.VT_I4, 0)
	yVar := ole.NewVariant(ole.VT_I4, 0)

	ret, err := oleutil.CallMethod(o.object, "FindColor", x1, y1, x2, y2, color, sim, dir, &xVar, &yVar)
	if err != nil {
		logf("[FindColor] 调用失败: %v", err)
		return 0, 0, false
	}

	// 从 VARIANT.Val 读取值（参考大漠写法）
	intX := int32(xVar.Val)
	intY := int32(yVar.Val)

	// 清理 VARIANT
	xVar.Clear()
	yVar.Clear()

	result := int(ret.Val)
	logf("[FindColor] 查找结果: result=%d, x=%d, y=%d", result, intX, intY)

	return int(intX), int(intY), result >= 0
}

// FindColorEx 多点找色
// 参考大漠写法：使用 oleutil.CallMethod，调用后从 VARIANT.Val 读取值
// 参数:
//
//	x1, y1, x2, y2: 查找区域
//	color: 颜色描述，格式为"颜色 | 偏移 X|偏移 Y，颜色 | 偏移 X|偏移 Y,..."
//	sim: 相似度，取值范围 0.0-1.0
//	dir: 查找方向
//
// 返回值:
//
//	x, y: 找到的第一个点的坐标
//	found: 是否找到
func (o *OP) FindColorEx(x1, y1, x2, y2 int, color string, sim float64, dir int) (x, y int, found bool) {
	logf("[FindColorEx] 开始查找: x1=%d, y1=%d, x2=%d, y2=%d, color=%s, sim=%f, dir=%d", x1, y1, x2, y2, color, sim, dir)

	// 创建输出参数 VARIANT
	xVar := ole.NewVariant(ole.VT_I4, 0)
	yVar := ole.NewVariant(ole.VT_I4, 0)

	ret, err := oleutil.CallMethod(o.object, "FindColorEx", x1, y1, x2, y2, color, sim, dir, &xVar, &yVar)
	if err != nil {
		logf("[FindColorEx] 调用失败: %v", err)
		return 0, 0, false
	}

	// 从 VARIANT.Val 读取值（参考大漠写法）
	intX := int32(xVar.Val)
	intY := int32(yVar.Val)

	// 清理 VARIANT
	xVar.Clear()
	yVar.Clear()

	result := int(ret.Val)
	logf("[FindColorEx] 查找结果: result=%d, x=%d, y=%d", result, intX, intY)

	return int(intX), int(intY), result >= 0
}

// GetColor 获取指定坐标的颜色
func (o *OP) GetColor(x, y int) string {
	ret, _ := oleutil.CallMethod(o.object, "GetColor", x, y)
	return ret.ToString()
}

// CmpColor 比较指定坐标点的颜色
func (o *OP) CmpColor(x, y int, color string, sim float64) int {
	ret, _ := oleutil.CallMethod(o.object, "CmpColor", x, y, color, sim)
	return int(ret.Val)
}

// FindPicEx 多图查找，返回所有找到的图片信息
// 参数同 FindPic
// 返回值：找到的图片索引，以及所有找到的坐标字符串，格式为"idx,x,y|idx,x,y|..."
func (o *OP) FindPicEx(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) string {
	ret, _ := oleutil.CallMethod(o.object, "FindPicEx", x1, y1, x2, y2, picName, deltaColor, sim, dir)
	return ret.ToString()
}

// FindColorEx2 多点找色，返回所有找到的点
// 参数同 FindColorEx
// 返回值：所有找到的坐标字符串，格式为"x,y|x,y|..."
func (o *OP) FindColorEx2(x1, y1, x2, y2 int, color string, sim float64, dir int) string {
	ret, _ := oleutil.CallMethod(o.object, "FindColorEx2", x1, y1, x2, y2, color, sim, dir)
	return ret.ToString()
}

// FindMultiColor 根据指定的多点颜色查找颜色坐标
// 参数:
//   - x1, y1, x2, y2: 查找区域
//   - firstColor: 第一个点的颜色，格式为"RRGGBB"
//   - offsetColor: 其他点的颜色偏移，格式为"x1|y1|RRGGBB,x2|y2|RRGGBB,..."
//   - sim: 相似度，取值范围 0.0-1.0
//   - dir: 查找方向，0 表示从左到右从上到下
//
// 返回值: x, y 找到的颜色坐标，found 是否找到
func (o *OP) FindMultiColor(x1, y1, x2, y2 int, firstColor, offsetColor string, sim float64, dir int) (x, y int, found bool) {
	ret, _ := oleutil.CallMethod(o.object, "FindMultiColor", x1, y1, x2, y2, firstColor, offsetColor, sim, dir)
	result := int(ret.Val)
	if result == -1 {
		return -1, -1, false
	}
	// 返回值格式为 x*65536 + y
	x = result / 65536
	y = result % 65536
	return x, y, true
}

// FindMultiColorEx 根据指定的多点颜色查找所有颜色坐标
// 参数同 FindMultiColor
// 返回值: 所有找到的坐标字符串，格式为"x,y|x,y|..."
func (o *OP) FindMultiColorEx(x1, y1, x2, y2 int, firstColor, offsetColor string, sim float64, dir int) string {
	ret, _ := oleutil.CallMethod(o.object, "FindMultiColorEx", x1, y1, x2, y2, firstColor, offsetColor, sim, dir)
	return ret.ToString()
}

// FindPicExS 查找多个图片，同 FindPicEx
// 参数同 FindPic
// 返回值: 所有找到的图片信息，格式为"idx,x,y|idx,x,y|..."
func (o *OP) FindPicExS(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) string {
	ret, _ := oleutil.CallMethod(o.object, "FindPicExS", x1, y1, x2, y2, picName, deltaColor, sim, dir)
	return ret.ToString()
}

// FindColorBlock 查找指定区域内的颜色块
// 官方文档: long FindColorBlock(x1, y1, x2, y2, color, sim, count, width, height, intX, intY)
// 参数:
//   - x1, y1, x2, y2: 查找区域
//   - color: 颜色格式串，比如"FFFFFF-000000|CCCCCC-000000"每种颜色用"|"分割
//   - sim: 相似度，取值范围 0.1-1.0
//   - count: 在宽度为width,高度为height的颜色块中，符合color颜色的最小数量
//   - width: 颜色块的宽度
//   - height: 颜色块的高度
//
// 返回值: x, y 找到的颜色块坐标，found 是否找到
func (o *OP) FindColorBlock(x1, y1, x2, y2 int, color string, sim float64, count, width, height int) (x, y int, found bool) {
	// 创建输出参数 VARIANT
	xVar := ole.NewVariant(ole.VT_I4, 0)
	yVar := ole.NewVariant(ole.VT_I4, 0)

	ret, _ := oleutil.CallMethod(o.object, "FindColorBlock", x1, y1, x2, y2, color, sim, count, width, height, &xVar, &yVar)

	// 从 VARIANT.Val 读取值
	xVal := int(int32(xVar.Val))
	yVal := int(int32(yVar.Val))

	// 清理 VARIANT
	xVar.Clear()
	yVar.Clear()

	result := int(ret.Val)
	return xVal, yVal, result >= 0
}

// FindColorBlockEx 查找指定区域内的所有颜色块
// 官方文档: string FindColorBlockEx(x1, y1, x2, y2, color, sim, count, width, height)
// 参数:
//   - x1, y1, x2, y2: 查找区域
//   - color: 颜色格式串，比如"FFFFFF-000000|CCCCCC-000000"每种颜色用"|"分割
//   - sim: 相似度，取值范围 0.1-1.0
//   - count: 在宽度为width,高度为height的颜色块中，符合color颜色的最小数量
//   - width: 颜色块的宽度
//   - height: 颜色块的高度
//
// 返回值: 所有找到的坐标字符串，格式为"x,y|x,y|..."
func (o *OP) FindColorBlockEx(x1, y1, x2, y2 int, color string, sim float64, count, width, height int) string {
	ret, _ := oleutil.CallMethod(o.object, "FindColorBlockEx", x1, y1, x2, y2, color, sim, count, width, height)
	return ret.ToString()
}

// SetDisplayInput 设置图像输入方式
// 官方文档: long SetDisplayInput(mode)
// 参数:
//   - mode: 图色输入模式，取值:
//   - "screen": 默认的模式，表示使用显示器或者后台窗口
//   - "pic:文件名": 指定输入模式为指定的图片
//   - "mem:addr": 指定输入模式为内存中的图片
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetDisplayInput(mode string) int {
	ret, _ := oleutil.CallMethod(o.object, "SetDisplayInput", mode)
	return int(ret.Val)
}

// LoadPic 预加载图片
// 参数:
//   - picName: 图片文件名
//
// 返回值: 0表示失败，1表示成功
func (o *OP) LoadPic(picName string) int {
	ret, _ := oleutil.CallMethod(o.object, "LoadPic", picName)
	return int(ret.Val)
}

// FreePic 释放图片
// 参数:
//   - picName: 图片文件名，空字符串表示释放所有图片
//
// 返回值: 0表示失败，1表示成功
func (o *OP) FreePic(picName string) int {
	ret, _ := oleutil.CallMethod(o.object, "FreePic", picName)
	return int(ret.Val)
}

// GetScreenData 获取指定区域的图像数据
// 官方文档: long GetScreenData(x1,y1,x2,y2)
// 参数:
//   - x1, y1, x2, y2: 区域坐标
//
// 返回值: 返回的是指定区域的二进制图片颜色数据指针，每个颜色是4个字节,表示方式为(00RRGGBB)
func (o *OP) GetScreenData(x1, y1, x2, y2 int) int {
	ret, _ := oleutil.CallMethod(o.object, "GetScreenData", x1, y1, x2, y2)
	return int(ret.Val)
}

// GetScreenDataBmp 获取指定区域的图像（BMP格式）
// 官方文档: long GetScreenDataBmp(x1,y1,x2,y2,data,size)
// 参数:
//   - x1, y1, x2, y2: 区域坐标
//
// 返回值: data 返回图片的数据指针, size 返回图片的数据长度, ret 0表示失败，1表示成功
func (o *OP) GetScreenDataBmp(x1, y1, x2, y2 int) (data, size int, ret int) {
	// 创建输出参数 VARIANT
	dataVar := ole.NewVariant(ole.VT_I4, 0)
	sizeVar := ole.NewVariant(ole.VT_I4, 0)

	result, _ := oleutil.CallMethod(o.object, "GetScreenDataBmp", x1, y1, x2, y2, &dataVar, &sizeVar)

	// 从 VARIANT.Val 读取值
	dataVal := int(int32(dataVar.Val))
	sizeVal := int(int32(sizeVar.Val))

	// 清理 VARIANT
	dataVar.Clear()
	sizeVar.Clear()

	return dataVal, sizeVal, int(result.Val)
}

// MatchPicName 使用通配符并获取文件集合
// 参数:
//   - picName: 图片文件名，支持通配符如 "*.bmp"
//
// 返回值: 匹配的文件名列表，格式为 "file1|file2|file3"
func (o *OP) MatchPicName(picName string) string {
	ret, _ := oleutil.CallMethod(o.object, "MatchPicName", picName)
	return ret.ToString()
}

// LoadMemPic 从内存中加载图片，并将加载结果返回
// 官方文档: long LoadMemPic(file_name, data, size)
// 参数:
//   - fileName: 图片的文件名
//   - data: 图像数据的内存地址（整数指针）
//   - size: 图像数据的大小
//
// 返回值: 0表示失败，1表示成功
func (o *OP) LoadMemPic(fileName string, data, size int) int {
	ret, _ := oleutil.CallMethod(o.object, "LoadMemPic", fileName, data, size)
	return int(ret.Val)
}

// ==================== OCR 文字识别 ====================

// SetDict 设置字库文件
func (o *OP) SetDict(index int, file string) int {
	ret, _ := oleutil.CallMethod(o.object, "SetDict", index, file)
	return int(ret.Val)
}

// UseDict 选择使用哪个字库文件进行识别
func (o *OP) UseDict(index int) int {
	ret, _ := oleutil.CallMethod(o.object, "UseDict", index)
	return int(ret.Val)
}

// Ocr 识别屏幕范围内的字符串
func (o *OP) Ocr(x1, y1, x2, y2 int, colorFormat string, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "Ocr", x1, y1, x2, y2, colorFormat, sim)
	return ret.ToString()
}

// FindStr 在屏幕范围内查找字符串
// 参考大漠写法：使用 oleutil.CallMethod，调用后从 VARIANT.Val 读取值
// 参数:
//
//	x1, y1, x2, y2: 查找区域
//	str: 要查找的字符串，支持多字符串查找，用 | 分隔
//	colorFormat: 颜色格式，格式为"RRGGBB-RRGGBB"
//	sim: 相似度，取值范围 0.0-1.0
//
// 返回值:
//
//	ret: 找到的字符串索引 (从 0 开始)，未找到返回 -1
//	x, y: 找到的字符串左上角坐标
func (o *OP) FindStr(x1, y1, x2, y2 int, str, colorFormat string, sim float64) (ret, x, y int) {
	logf("[FindStr] 开始调用: str=%s, colorFormat=%s, sim=%.2f, region=(%d,%d,%d,%d)", str, colorFormat, sim, x1, y1, x2, y2)

	// 创建输出参数 VARIANT
	xVar := ole.NewVariant(ole.VT_I4, 0)
	yVar := ole.NewVariant(ole.VT_I4, 0)

	result, err := oleutil.CallMethod(o.object, "FindStr", x1, y1, x2, y2, str, colorFormat, sim, &xVar, &yVar)

	if err != nil {
		logf("[FindStr] 调用失败: err=%v", err)
		return -1, 0, 0
	}
	if result == nil {
		logf("[FindStr] 调用失败: result is nil")
		return -1, 0, 0
	}

	// 从 VARIANT.Val 读取值（参考大漠写法）
	ret = int(int32(result.Val))
	x = int(int32(xVar.Val))
	y = int(int32(yVar.Val))

	// 清理 VARIANT
	xVar.Clear()
	yVar.Clear()

	logf("[FindStr] 调用成功: ret=%d, x=%d, y=%d", ret, x, y)
	if ret >= 0 {
		logf("[FindStr] 找到字符串: 索引=%d, 坐标=(%d,%d)", ret, x, y)
	} else {
		logf("[FindStr] 未找到字符串")
	}

	return ret, x, y
}

// FindStrEx 在屏幕范围内查找字符串，返回所有找到的
// 参数同 FindStr
// 返回值：所有找到的字符串信息，格式为"idx,x,y,str|idx,x,y,str|..."
func (o *OP) FindStrEx(x1, y1, x2, y2 int, str, colorFormat string, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "FindStrEx", x1, y1, x2, y2, str, colorFormat, sim)
	return ret.ToString()
}

// SetMemDict 设置内存字库文件
// 官方文档: long SetMemDict(idx,data,size)
// 参数:
//   - index: 字库的序号，范围0-9
//   - data: 字库内容数据的内存地址（整数指针）
//   - size: 字库大小
//
// 返回值: 0表示失败，1表示成功
func (o *OP) SetMemDict(index, data, size int) int {
	ret, _ := oleutil.CallMethod(o.object, "SetMemDict", index, data, size)
	return int(ret.Val)
}

// OcrEx 识别屏幕范围内的字符串，返回识别到的字符串以及每个字符的坐标
// 参数:
//   - x1, y1, x2, y2: 识别区域
//   - colorFormat: 颜色格式串，比如"FFFFFF-000000|CCCCCC-000000"每种颜色用"|"分割
//   - sim: 相似度，取值范围0.1-1.0
//
// 返回值: 返回识别到的字符串以及坐标，格式为"char$x$y|char$x$y|..."
func (o *OP) OcrEx(x1, y1, x2, y2 int, colorFormat string, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "OcrEx", x1, y1, x2, y2, colorFormat, sim)
	return ret.ToString()
}

// OcrAuto 识别屏幕范围内的字符串，自动二值化，无需指定颜色
// 适用于字体颜色和背景相差较大的场合
// 参数:
//   - x1, y1, x2, y2: 识别区域
//   - sim: 相似度，取值范围0.1-1.0
//
// 返回值: 返回识别到的字符串
func (o *OP) OcrAuto(x1, y1, x2, y2 int, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "OcrAuto", x1, y1, x2, y2, sim)
	return ret.ToString()
}

// OcrFromFile 从文件中识别图片
// 参数:
//   - fileName: 文件名
//   - colorFormat: 颜色格式串
//   - sim: 相似度，取值范围0.1-1.0
//
// 返回值: 返回识别到的字符串
func (o *OP) OcrFromFile(fileName, colorFormat string, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "OcrFromFile", fileName, colorFormat, sim)
	return ret.ToString()
}

// OcrAutoFromFile 从文件中识别图片，自动二值化，无需指定颜色
// 参数:
//   - fileName: 文件名
//   - sim: 相似度，取值范围0.1-1.0
//
// 返回值: 返回识别到的字符串
func (o *OP) OcrAutoFromFile(fileName string, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "OcrAutoFromFile", fileName, sim)
	return ret.ToString()
}

// FindLine 在指定的屏幕坐标范围内，查找指定颜色的直线
// 参数:
//   - x1, y1, x2, y2: 查找区域
//   - colorFormat: 颜色格式串
//   - sim: 相似度，取值范围0.1-1.0
//
// 返回值: 返回识别到的结果
func (o *OP) FindLine(x1, y1, x2, y2 int, colorFormat string, sim float64) string {
	ret, _ := oleutil.CallMethod(o.object, "FindLine", x1, y1, x2, y2, colorFormat, sim)
	return ret.ToString()
}

// ==================== Memory 内存操作 ====================

// WriteData 向某进程写入数据
// 官方文档: long WriteData(hwnd,address,data,size)
// 参数:
//   - hwnd: 窗口句柄，用于指定要在哪个窗口内写入数据
//   - address: 写入数据的地址（字符串类型）
//   - data: 写入的数据
//   - size: 写入的数据的大小
//
// 返回值: 0表示失败，1表示成功
func (o *OP) WriteData(hwnd int, address, data string, size int) int {
	ret, _ := oleutil.CallMethod(o.object, "WriteData", hwnd, address, data, size)
	return int(ret.Val)
}

// ReadData 读取数据
// 官方文档: string ReadData(hwnd,address,size)
// 参数:
//   - hwnd: 窗口句柄，用于指定要从哪个窗口内读取数据
//   - address: 表示要读取数据的地址（字符串类型）
//   - size: 要读取的数据的大小
//
// 返回值: 读取到的数值
func (o *OP) ReadData(hwnd int, address string, size int) string {
	ret, _ := oleutil.CallMethod(o.object, "ReadData", hwnd, address, size)
	return ret.ToString()
}

// ==================== Algorithm 算法 ====================

// AStarFindPath A星算法
// 官方文档: string AStarFindPath(mapWidth, mapHeight, disable_points, beginX, beginY, endX, endY)
// 参数:
//   - mapWidth: 地图宽度
//   - mapHeight: 地图高度
//   - disablePoints: 不可通行的坐标，以"|"分割，例如:"10,15|20,30"
//   - beginX, beginY: 源坐标
//   - endX, endY: 目的坐标
//
// 返回值: 找到的路径结果，格式为"x1,y1|x2,y2|..."
func (o *OP) AStarFindPath(mapWidth, mapHeight int, disablePoints string, beginX, beginY, endX, endY int) string {
	ret, _ := oleutil.CallMethod(o.object, "AStarFindPath", mapWidth, mapHeight, disablePoints, beginX, beginY, endX, endY)
	return ret.ToString()
}

// FindNearestPos 查找最近的位置
// 官方文档: void FindNearestPos(all_pos, type, x, y, ret)
// 参数:
//   - allPos: 位置集合，以"|"分割的坐标列表，例如:"10,15|20,30|50,60"
//   - posType: 类型
//   - x, y: 参考坐标
//
// 返回值: 最接近指定坐标 (x, y) 的位置，格式为"x,y"
func (o *OP) FindNearestPos(allPos string, posType, x, y int) string {
	// 创建输出参数 VARIANT
	retVar := ole.NewVariant(ole.VT_BSTR, 0)
	_, _ = oleutil.CallMethod(o.object, "FindNearestPos", allPos, posType, x, y, &retVar)
	result := retVar.ToString()
	retVar.Clear()
	return result
}

// ==================== System 系统命令 ====================

// GetScreenWidth 获取屏幕宽度
func (o *OP) GetScreenWidth() int {
	ret, _ := oleutil.CallMethod(o.object, "GetScreenWidth")
	return int(ret.Val)
}

// GetScreenHeight 获取屏幕高度
func (o *OP) GetScreenHeight() int {
	ret, _ := oleutil.CallMethod(o.object, "GetScreenHeight")
	return int(ret.Val)
}

// GetClipboard 获取剪贴板内容
func (o *OP) GetClipboard() string {
	ret, _ := oleutil.CallMethod(o.object, "GetClipboard")
	return ret.ToString()
}

// SetClipboard 设置剪贴板内容
func (o *OP) SetClipboard(str string) int {
	ret, _ := oleutil.CallMethod(o.object, "SetClipboard", str)
	return int(ret.Val)
}

// RunApp 运行可执行文件，可指定模式
// 参数:
//   - appPath: 指定的可执行程序全路径
//   - mode: 取值0表示普通模式，1表示加强模式
//
// 返回值: 0表示失败，1表示成功
func (o *OP) RunApp(appPath string, mode int) int {
	ret, _ := oleutil.CallMethod(o.object, "RunApp", appPath, mode)
	return int(ret.Val)
}

// WinExec 运行可执行文件，可指定显示模式
// 参数:
//   - cmdLine: 指定的可执行程序全路径
//   - cmdShow: 取值0表示隐藏，1表示用最近的大小和位置显示并激活
//
// 返回值: 0表示失败，1表示成功
func (o *OP) WinExec(cmdLine string, cmdShow int) int {
	ret, _ := oleutil.CallMethod(o.object, "WinExec", cmdLine, cmdShow)
	return int(ret.Val)
}

// GetCmdStr 运行命令行并返回结果
// 参数:
//   - cmdLine: 指定的命令行
//   - milliseconds: 等待的时间（毫秒）
//
// 返回值: cmd输出的字符
func (o *OP) GetCmdStr(cmdLine string, milliseconds int) string {
	ret, _ := oleutil.CallMethod(o.object, "GetCmdStr", cmdLine, milliseconds)
	return ret.ToString()
}

// Delay 实现一个指定毫秒数的延迟，同时确保在此期间不会阻塞用户界面（UI）操作
// 参数:
//   - ms: 指定延迟的时间，单位为毫秒
//
// 返回值: 0表示失败，1表示成功
func (o *OP) Delay(ms int) int {
	ret, _ := oleutil.CallMethod(o.object, "Delay", ms)
	return int(ret.Val)
}

// Delays 实现一个指定毫秒数的延迟，同时确保在此期间不会阻塞用户界面（UI）操作
// 参数:
//   - msMin: 指定延迟时间的最小值，单位为毫秒
//   - msMax: 指定延迟时间的最大值，单位为毫秒
//
// 返回值: 0表示失败，1表示成功
// 函数将随机选择一个介于msMin和msMax之间的延迟时间
func (o *OP) Delays(msMin, msMax int) int {
	ret, _ := oleutil.CallMethod(o.object, "Delays", msMin, msMax)
	return int(ret.Val)
}

// SendMessage 向窗口发送消息
// 参数:
//   - hwnd: 窗口句柄
//   - msg: 消息类型
//   - wParam: 消息参数1
//   - lParam: 消息参数2
//
// 返回值: 消息处理结果
func (o *OP) SendMessage(hwnd, msg, wParam, lParam int) int {
	ret, _ := oleutil.CallMethod(o.object, "SendMessage", hwnd, msg, wParam, lParam)
	return int(ret.Val)
}

// PostMessage 向窗口发送消息（异步）
// 参数:
//   - hwnd: 窗口句柄
//   - msg: 消息类型
//   - wParam: 消息参数1
//   - lParam: 消息参数2
//
// 返回值: 0表示失败，非0表示成功
func (o *OP) PostMessage(hwnd, msg, wParam, lParam int) int {
	ret, _ := oleutil.CallMethod(o.object, "PostMessage", hwnd, msg, wParam, lParam)
	return int(ret.Val)
}
