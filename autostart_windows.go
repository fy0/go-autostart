package autostart

import (
	"os"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

var startupDir string

func init() {
	startupDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
}

func (a *App) path() string {
	return filepath.Join(startupDir, a.Name+".lnk")
}

// IsEnabled checks if app is enabled on startup.
func (a *App) IsEnabled() bool {
	_, err := os.Stat(a.path())
	return err == nil
}

// Enable this app on startup.
func (a *App) Enable() error {
	exePath := a.Exec[0]
	args := a.Exec[1:]

	if err := os.MkdirAll(startupDir, 0777); err != nil {
		return err
	}
	if err := createShortcut(a.path(), exePath, args...); err != nil {
		return err
	}

	return nil
}

// Disable this app on startup.
func (a *App) Disable() error {
	return os.Remove(a.path())
}

func createShortcut(autostartPath, exePath string, params ...string) error {
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		return err
	}
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", autostartPath)
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	_, err = oleutil.PutProperty(idispatch, "TargetPath", exePath)
	if err != nil {
		return err
	}
	var args []interface{}
	for _, param := range params {
		args = append(args, param)
	}
	_, err = oleutil.CallMethod(idispatch, "Save", args...)
	return err
}
