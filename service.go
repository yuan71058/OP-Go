package op

import (
	"errors"
	"log"
	"sync"
)

var ErrNotInitialized = errors.New("op: not initialized")

type Service struct {
	mu      sync.RWMutex
	op      *OP
	dllPath string
	isReady bool
}

func NewService(dllPath string) *Service {
	return &Service{
		dllPath: dllPath,
	}
}

func (s *Service) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isReady {
		return nil
	}

	op, err := NewOP(s.dllPath)
	if err != nil {
		return err
	}

	s.op = op
	s.isReady = true
	return nil
}

func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.op != nil {
		s.op.Release()
		s.op = nil
	}
	s.isReady = false
	return nil
}

func (s *Service) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isReady
}

func (s *Service) GetOP() (*OP, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.isReady || s.op == nil {
		return nil, ErrNotInitialized
	}
	return s.op, nil
}

func (s *Service) GetVersion() (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.Ver(), nil
}

func (s *Service) SetDLLPath(dllPath string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.isReady {
		s.dllPath = dllPath
	}
}

type Status struct {
	IsReady bool   `json:"is_ready"`
	Version string `json:"version"`
	DllPath string `json:"dll_path"`
}

func (s *Service) GetStatus() (*Status, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := &Status{
		IsReady: s.isReady,
		DllPath: s.dllPath,
	}

	if s.isReady && s.op != nil {
		status.Version = s.op.Ver()
	}

	return status, nil
}

func (s *Service) GetScreenSize() (int, int, error) {
	op, err := s.GetOP()
	if err != nil {
		return 0, 0, err
	}
	return op.GetScreenWidth(), op.GetScreenHeight(), nil
}

type WindowInfo struct {
	Hwnd      int    `json:"hwnd"`
	Title     string `json:"title"`
	ClassName string `json:"class_name"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	IsVisible bool   `json:"is_visible"`
}

func (s *Service) FindWindow(className, titleName string) (int, error) {
	op, err := s.GetOP()
	if err != nil {
		return 0, err
	}
	return op.FindWindow(className, titleName), nil
}

func (s *Service) GetWindowText(hwnd int) (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.GetWindowTitle(hwnd), nil
}

func (s *Service) GetWindowRect(hwnd int) (x1, y1, x2, y2 int, err error) {
	op, e := s.GetOP()
	if e != nil {
		return 0, 0, 0, 0, e
	}
	x1, y1, x2, y2 = op.GetWindowRect(hwnd)
	return x1, y1, x2, y2, nil
}

func (s *Service) GetClientSize(hwnd int) (width, height int, err error) {
	log.Printf("[Service] GetClientSize called: hwnd=%d", hwnd)
	op, e := s.GetOP()
	if e != nil {
		log.Printf("[Service] GetClientSize error: %v", e)
		return 0, 0, e
	}
	width, height = op.GetClientSize(hwnd)
	log.Printf("[Service] GetClientSize success: width=%d, height=%d", width, height)
	return width, height, nil
}

func (s *Service) BindWindow(hwnd int, display, mouse, keypad string, mode int) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	ret := op.BindWindow(hwnd, display, mouse, keypad, mode)
	if ret == 0 {
		return errors.New("绑定窗口失败")
	}
	return nil
}

func (s *Service) UnBindWindow() error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.UnBindWindow()
	return nil
}

func (s *Service) MoveTo(x, y int) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.MoveTo(x, y)
	return nil
}

func (s *Service) LeftClick() error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.LeftClick()
	return nil
}

func (s *Service) RightClick() error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.RightClick()
	return nil
}

func (s *Service) KeyPress(keyStr string) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.KeyPressChar(keyStr)
	return nil
}

func (s *Service) Capture(x1, y1, x2, y2 int, file string) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	ret := op.Capture(x1, y1, x2, y2, file)
	if ret == 0 {
		return errors.New("截图失败")
	}
	return nil
}

func (s *Service) GetColor(x, y int) (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.GetColor(x, y), nil
}

func (s *Service) FindColor(x1, y1, x2, y2 int, color string, sim float64, dir int) (int, int, bool, error) {
	op, err := s.GetOP()
	if err != nil {
		return 0, 0, false, err
	}
	x, y, found := op.FindColor(x1, y1, x2, y2, color, sim, dir)
	return x, y, found, nil
}

func (s *Service) FindPic(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) (int, int, bool, error) {
	op, err := s.GetOP()
	if err != nil {
		return 0, 0, false, err
	}
	x, y, found := op.FindPic(x1, y1, x2, y2, picName, deltaColor, sim, dir)
	return x, y, found, nil
}

func (s *Service) FindPicEx(x1, y1, x2, y2 int, picName, deltaColor string, sim float64, dir int) (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.FindPicEx(x1, y1, x2, y2, picName, deltaColor, sim, dir), nil
}

func (s *Service) Ocr(x1, y1, x2, y2 int, colorFormat string, sim float64) (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.Ocr(x1, y1, x2, y2, colorFormat, sim), nil
}

func (s *Service) GetCursorPos() (int, int, error) {
	op, err := s.GetOP()
	if err != nil {
		return 0, 0, err
	}
	x, y := op.GetCursorPos()
	return x, y, nil
}

func (s *Service) SendString(hwnd int, str string) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.SendString(hwnd, str)
	return nil
}

func (s *Service) SetPath(path string) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.SetPath(path)
	return nil
}

func (s *Service) GetPath() (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.GetPath(), nil
}

func (s *Service) GetBasePath() (string, error) {
	op, err := s.GetOP()
	if err != nil {
		return "", err
	}
	return op.GetBasePath(), nil
}

func (s *Service) SetShowErrorMsg(show int) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.SetShowErrorMsg(show)
	return nil
}

func (s *Service) EnablePicCache(enable int) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.EnablePicCache(enable)
	return nil
}

func (s *Service) SetDict(index int, file string) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.SetDict(index, file)
	return nil
}

func (s *Service) UseDict(index int) error {
	op, err := s.GetOP()
	if err != nil {
		return err
	}
	op.UseDict(index)
	return nil
}
