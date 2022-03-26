package lib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// RunGitCommand 执行任意Git命令的封装
func RunGitCommand(name string, arg ...string) string {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	//msg, err := cmd.CombinedOutput() // 混合输出stdout+stderr
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(err, stderr)
		// 报错时 exit status 1
		os.Exit(1)
	}

	return stdout.String()
}

// ConvertOctonaryUtf8 转换8进制utf-8字符串到中文 eg: `\346\200\241` -> 怡
func ConvertOctonaryUtf8(in string) string {
	s := []byte(in)
	reg := regexp.MustCompile(`\\[0-7]{3}`)

	out := reg.ReplaceAllFunc(
		s,
		func(b []byte) []byte {
			i, _ := strconv.ParseInt(string(b[1:]), 8, 0)
			return []byte{byte(i)}
		},
	)
	return string(out)
}
