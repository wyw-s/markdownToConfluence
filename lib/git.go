package lib

import (
	"fmt"
	"os"
	"strings"
)

type MarkdownFileFromGit struct {
	status string
	path   string
}

func GetMarkdownFile(m *Markdown2Confluence) []MarkdownFileFromGit {
	// 获取当前工作目录
	workspaceDir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("当前工作目录路径：%s\n", workspaceDir)
	fmt.Println()

	// 检查返回值 workspace 是否为空 为空则表明工作区没有未提交的变更
	workspaceFiles := RunGitCommand("git", "status", "--short")

	if workspaceFiles == "" {
		return nil
	}

	// 这里必须使文件先被跟踪
	RunGitCommand("git", "add", ".")

	statusFileList := strings.Split(workspaceFiles, "\n")

	var uploadFileList []MarkdownFileFromGit

	for _, value := range statusFileList {
		if value != "" {
			uploadFileList = append(uploadFileList, ProcessGitFilePath(value, m)...)
		}
	}

	if len(uploadFileList) == 0 {
		fmt.Println("暂无同步的markdown文件")
		return nil
	}

	fmt.Printf("---------------------------正在同步...--------------------------\n")

	return uploadFileList
}

func ProcessGitFilePath(path string, m *Markdown2Confluence) []MarkdownFileFromGit {

	path = strings.TrimSpace(path)

	// M: 修改 R：重命名 D：删除 A: 新增 ?: 未跟踪的
	var statusList = [5]string{"M", "R", "D", "A", "?"}
	var markdownFiles []MarkdownFileFromGit

	for _, code := range statusList {
		if !strings.HasPrefix(path, code) {
			continue
		}
		if strings.Contains(path, "->") {

			var list = strings.Split(path, "->")

			var st = MarkdownFileFromGit{
				status: code,
				path:   path,
			}

			for i, v := range list {

				var str = strings.Trim(v, " ")

				st.path = str

				if i == 0 {
					st.status = "D"
				} else {
					st.status = "A"
				}

				if strings.HasSuffix(str, "\"") {

					firstIndex := strings.Index(str, "\"")
					lastIndex := strings.LastIndex(str, "\"")

					str = str[firstIndex+1 : lastIndex]

					st.path = ConvertOctonaryUtf8(str)
				} else if i == 0 {
					st.path = v[strings.LastIndex(v, " "):len(v)]
				}
				if m.GitSyncDir != "" {
					if strings.HasPrefix(st.path, m.GitSyncDir) {
						markdownFiles = append(markdownFiles, st)
					}
				} else {
					markdownFiles = append(markdownFiles, st)
				}
			}
			return markdownFiles
		} else {
			var st = MarkdownFileFromGit{
				status: code,
				path:   path,
			}
			if strings.HasSuffix(path, "\"") {

				firstIndex := strings.Index(path, "\"")
				lastIndex := strings.LastIndex(path, "\"")

				st.path = ConvertOctonaryUtf8(path[firstIndex+1 : lastIndex])
			} else {
				st.path = path[strings.LastIndex(path, " "):len(path)]
			}
			if m.GitSyncDir != "" {
				if strings.HasPrefix(st.path, m.GitSyncDir) {
					return append(markdownFiles, st)
				}
			} else {
				return append(markdownFiles, st)
			}
		}
	}
	return nil
}
