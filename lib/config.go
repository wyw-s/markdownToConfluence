package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfluenceConfig struct {
	Username   string
	Password   string
	Endpoint   string
	Space      string
	Parent     string
	GitSyncDir string
	Model      string
}

func (conf *ConfluenceConfig) LoadConfig() error {

	// 获取当前工作目录
	workspaceDir, _ := os.Getwd()

	// 读取配置文件
	buf, err := ioutil.ReadFile(workspaceDir + `/.confluence.json`)

	if err != nil {
		fmt.Println("read file err: ", err)
		return err
	}

	err = json.Unmarshal(buf, conf)

	if err != nil {
		fmt.Println("decode config file failed", err)
		return err
	}
	return nil
}

func (conf *ConfluenceConfig) SetConfig(m Markdown2Confluence) {
	m.Username = conf.Username
	m.Password = conf.Password
	m.Endpoint = conf.Endpoint
	m.Space = conf.Space
	m.Parent = conf.Parent
	m.GitSyncDir = conf.GitSyncDir
	m.Model = conf.Model
}
