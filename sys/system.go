package sys

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// open file
// 打开文件
func open(uri string) error {
	var commands = map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("%s platform ？？？", runtime.GOOS)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start ", uri)
		//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	} else {
		cmd = exec.Command(run, uri)
	}
	return cmd.Start()
}

// Create a folder under the execution path
// 在执行路径下面创建文件夹
func mkdir(path string) (string, error) {
	delimiter := "/"
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filePtah := dir + delimiter + path + delimiter
	fmt.Println("exec path----->", filePtah)
	err := os.MkdirAll(filePtah, 0777)
	if err != nil {
		return "", err
	}
	return filePtah, nil
}
