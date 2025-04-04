// 机器人信息输出
// @author MoGuQAQ
// @version 1.0.0

package telebot

import (
	"TeleGPT/api"
	"TeleGPT/config"
	"TeleGPT/global"
	"TeleGPT/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func InitTelegramBot() {
	bottoken := config.Config.BotInfo.Bot_token
	bot, err := tgbotapi.NewBotAPI(bottoken)
	if err != nil {
		global.Log.Fatalf("发生致命错误 %s", err)
	}

	bot.Debug = false

	global.Log.Infof("认证账户到 %s", bot.Self.UserName)
	config.Config.BotInfo.Bot_uname = bot.Self.UserName
	err = config.UpdateYaml(config.Config)
	if err != nil {
		global.Log.Fatalf("更新配置文件失败: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		global.Log.Warnf("发生致命错误 %s", err)
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		privacymode := config.Config.BotInfo.Privacymode
		adminid := int64(config.Config.GroupInfo.Admin_id)
		global.Log.Infof("%v", config.Config.GroupInfo.Whitelist)
		if privacymode {
			if update.Message.Chat.ID != adminid && !utils.Contains(config.Config.GroupInfo.Whitelist, update.Message.Chat.ID) {
				global.Log.Infof("%d 非白名单账户,主动忽略", update.Message.Chat.ID)
				continue
			}
		}
		global.Log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		Vector(bot, update.Message)
	}
}

func Vector(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	if update.IsCommand() {
		if utils.Is_command(update.Text, "start", config.Config.BotInfo.Bot_uname) {
			api.Start(bot, update)
		} else if utils.Is_command(update.Text, "help", config.Config.BotInfo.Bot_uname) {
			api.Help(bot, update)
		} else if utils.Is_command(update.Text, "add", config.Config.BotInfo.Bot_uname) {
			api.Add(bot, update)
		} else if utils.Is_command(update.Text, "del", config.Config.BotInfo.Bot_uname) {
			api.Del(bot, update)
		} else if utils.Is_command(update.Text, "whitelist", config.Config.BotInfo.Bot_uname) {
			api.WhiteList(bot, update)
		} else if utils.Is_command(update.Text, "model", config.Config.BotInfo.Bot_uname) {
			api.Model(bot, update)
		} else {
			api.Unknown(bot, update)
		}
	} else if update.Chat.IsPrivate() || utils.Is_at(update.Text, config.Config.BotInfo.Bot_uname) {
		api.Chat(bot, update)
	} else {
		return
	}
}
