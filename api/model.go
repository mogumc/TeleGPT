// model命令
// @author MoGuQAQ
// @version 1.0.0

package api

import (
	"TeleGPT/config"
	"TeleGPT/global"
	"TeleGPT/utils"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Model(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	if update.Chat.ID != int64(config.Config.GroupInfo.Admin_id) {
		global.Log.Infof("%d尝试使用管理员命令", update.Chat.ID)
		msg := tgbotapi.NewMessage(update.Chat.ID, "你没有权限使用此命令")
		msg.ReplyToMessageID = update.MessageID
		bot.Send(msg)
	} else {
		model := utils.Get_command(update.Text, "model", config.Config.BotInfo.Bot_uname)
		if model != "" {
			if len(model) > 25 {
				msg := tgbotapi.NewMessage(update.Chat.ID, "模型名称过长,请手动修改配置文件")
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
				return
			}
			config.Config.API.Omodel = model
			err := config.UpdateYaml(config.Config)
			if err != nil {
				global.Log.Fatalf("更新配置文件失败: %v", err)
				msg := tgbotapi.NewMessage(update.Chat.ID, "更新配置文件失败")
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Chat.ID, "修改模型成功")
			msg.ReplyToMessageID = update.MessageID
			bot.Send(msg)
		} else {
			textreturn := fmt.Sprintf("当前模型为: %s", model)
			msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
			msg.ReplyToMessageID = update.MessageID
			bot.Send(msg)
			return
		}
	}
}
