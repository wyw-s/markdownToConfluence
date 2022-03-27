package lib

import (
	"fmt"
	"os"
	"strings"
)

func GitLogger() ([]string, []string) {

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
		return []string{}, []string{}
	}

	// 这里必须使文件先被跟踪
	RunGitCommand("git", "add", ".")

	// 获取更改的文件列表
	changeFileList := RunGitCommand("git", "diff", "--name-only", "HEAD")

	statusFileList := RunGitCommand("git", "diff", "--name-status", "HEAD")

	formatDeleteFiles := strings.Split(statusFileList, "\n")

	// 获取删除的文件列表
	var deleteFiles []string

	for _, value := range formatDeleteFiles {
		if strings.HasPrefix(value, "D") {
			temp := strings.Split(value, "\t")
			deleteFiles = append(deleteFiles, ConvertOctonaryUtf8(strings.Replace(temp[1], "\"", "", -1)))
		}
	}

	if changeFileList == "" {
		return []string{}, []string{}
	}

	fileList := strings.Split(changeFileList, "\n")

	if len(fileList) != 0 && fileList[len(fileList)-1] == "" {
		fileList = fileList[:len(fileList)-1]
	}

	var uploadFiles []string

	fmt.Printf("待同步的文件列表：(只同步md文件)\n")

	for i, value := range fileList {
		if strings.HasSuffix(value, "\"") && strings.HasPrefix(value, "\"") {
			fileList[i] = ConvertOctonaryUtf8(strings.Replace(value, "\"", "", -1))
		}
		isExist := isValueInList(fileList[i], deleteFiles)
		if !isExist {
			uploadFiles = append(uploadFiles, fileList[i])
		}
		fmt.Printf("%v: %s\n", i, fileList[i])
	}

	fmt.Printf("---------------------------同步至云端--------------------------\n")

	return uploadFiles, deleteFiles
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
