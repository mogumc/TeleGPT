// 配置文件读取
// @author MoGuQAQ
// @version 1.0.0

package config

import (
	"TeleGPT/global"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type config struct {
	Logger    logger    `yaml:"logger"`
	API       apis      `yaml:"api"`
	BotInfo   tginfo    `yaml:"bot"`
	GroupInfo groupinfo `yaml:"group"`
}

type logger struct {
	Level        string `yaml:"level"`
	Prefix       string `yaml:"prefix"`
	Showline     bool   `yaml:"show_line"`
	LogInConsole bool   `yaml:"log_in_console"`
}

type apis struct {
	Oapi     string `yaml:"api_link"`
	Otoken   string `yaml:"api_key"`
	Omodel   string `yaml:"api_model"`
	Oprompt  string `yaml:"api_prompt"`
	Otimeout int    `yaml:"api_timeout"`
	Otemp    int    `yaml:"api_temperature"`
	Otopk    int    `yaml:"api_top_k"`
	Otopp    int    `yaml:"api_top_p"`
	Oproxy   string `yaml:"proxy"`
}

type tginfo struct {
	Bot_token   string `yaml:"bottoken"`
	Bot_uname   string `yaml:"bot_user_name"`
	Privacymode bool   `yaml:"privacy"`
}

type groupinfo struct {
	ClearWord string  `yaml:"clear_word"`
	Admin_id  int     `yaml:"admin_id"`
	Whitelist []int64 `yaml:"Whitelist"`
}

var Config *config

func init() {
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		global.Log.Errorf("解析文件失败: %v", err)
	}
	yaml.Unmarshal(yamlFile, &Config)
}

func UpdateYaml(config *config) error {
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		global.Log.Fatalf("发生错误!更新配置失败!")
		return err
	}
	err = ioutil.WriteFile("./config.yaml", yamlData, 0644)
	if err != nil {
		global.Log.Errorf("写入文件失败: %v", err)
		return err
	}
	return nil
}
