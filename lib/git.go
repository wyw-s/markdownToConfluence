package lib

import (
	"fmt"
	"os"
	"strings"
)

func GitLogger() []string {

	// 获取当前工作目录
	workspaceDir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("当前工作目录路径：%s", workspaceDir)
	fmt.Println()

	// 检查返回值 workspace 是否为空 为空则表明工作区没有未提交的变更
	workspace := RunGitCommand("git", "status", "--short")

	if workspace == "" {
		return []string{}
	}

	// 这里必须使文件先被跟踪
	RunGitCommand("git", "add", ".")

	// 获取更改的文件列表
	changeFileList := RunGitCommand("git", "diff", "--name-only", "HEAD")

	if changeFileList == "" {
		return []string{}
	}

	fileList := strings.Split(changeFileList, "\n")

	if len(fileList) != 0 && fileList[len(fileList)-1] == "" {
		fileList = fileList[:len(fileList)-1]
	}

	fmt.Printf("变更的文件列表：\n")

	for i, value := range fileList {
		if strings.HasSuffix(value, "\"") && strings.HasPrefix(value, "\"") {
			fileList[i] = ConvertOctonaryUtf8(strings.Replace(value, "\"", "", -1))
		}
		fmt.Printf("%v: %s\n", i, fileList[i])
	}

	fmt.Printf("已上传的文件列表：(只同步md文件)\n")

	return fileList
}
