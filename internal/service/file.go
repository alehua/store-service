package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func setDir(path string) error {
	_, err := os.Stat(path)
	// 如果不存在则创建
	if err != nil {
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	return nil
}

func Save(rFile *multipart.FileHeader) error {
	errChan := make(chan error, 1)
	var filePath, fileDir string
	var err error
	go func() {
		fileDir = fmt.Sprintf("databases/%s", rFile.Filename)
		if err := setDir(fileDir); err != nil {
			errChan <- fmt.Errorf("文件路径创建失败: %s", err.Error())
			return
		}
		file, _ := rFile.Open()
		defer file.Close()
		fileName := fmt.Sprintf("%s_create", time.Now().Format("2006-01-02_15:04:05"))
		filePath = fmt.Sprintf("%s/%s", fileDir, fileName)
		out, _ := os.Create(filePath)
		if _, err := io.Copy(out, file); err != nil {
			errChan <- fmt.Errorf("文件保存失败: %s", err.Error())
			return
		}
		absPath, _ := filepath.Abs(filePath)
		var cmd = exec.Command("ln", "-snf", absPath, fmt.Sprintf("%s/last", fileDir))
		if err = cmd.Run(); err != nil {
			errChan <- fmt.Errorf("设置软连接失败: %s", err.Error())
			return
		}
		errChan <- nil
	}()

	err = <-errChan
	return err
}

func DownLoad(fileName string) string {
	targetPath := fmt.Sprintf("./databases/%s/last", fileName) //文件存放目录
	_, err := os.Open(targetPath)
	if err != nil {
		return ""
	}
	return targetPath
}

func Dir() ([]string, error) {
	root := "databases"
	files := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() != root {
			spl := strings.Split(path, "/")
			files = append(files, spl[len(spl)-1])
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, nil
}
