// 白名单命令
// @author MoGuQAQ
// @version 1.0.0

package api

import (
	"TeleGPT/config"
	"TeleGPT/global"
	"TeleGPT/utils"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Add(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	if update.From.ID != config.Config.GroupInfo.Admin_id {
		global.Log.Infof("%d尝试使用管理员命令", update.From.ID)
		msg := tgbotapi.NewMessage(update.Chat.ID, "你没有权限使用此命令")
		msg.ReplyToMessageID = update.MessageID
		bot.Send(msg)
	} else {
		gettext := utils.Get_command(update.Text, "add", config.Config.BotInfo.Bot_uname)
		if gettext != "" {
			addid, err := strconv.Atoi(gettext)
			if err != nil {
				global.Log.Errorf("数据转换失败: %v", err)
				msg := tgbotapi.NewMessage(update.Chat.ID, "数据转换失败")
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
				return
			}
			if utils.Contains(config.Config.GroupInfo.Whitelist, int64(addid)) {
				textreturn := fmt.Sprintf("%d已在白名单中,当前白名单: %v", addid, config.Config.GroupInfo.Whitelist)
				msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
				return
			}
			var whitelistPointer *[]int64 = &config.Config.GroupInfo.Whitelist
			*whitelistPointer = append(*whitelistPointer, int64(addid))
			err = config.UpdateYaml(config.Config)
			if err != nil {
				global.Log.Fatalf("更新配置文件失败: %v", err)
				msg := tgbotapi.NewMessage(update.Chat.ID, "更新配置文件失败")
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
				return
			}
			textreturn := fmt.Sprintf("添加%d成功,当前白名单: %v", addid, config.Config.GroupInfo.Whitelist)
			msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
			msg.ReplyToMessageID = update.MessageID
			bot.Send(msg)
		} else {
			textreturn := "请输入要添加的ID"
			msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
			msg.ReplyToMessageID = update.MessageID
			bot.Send(msg)
			return
		}
	}
}

func Del(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	if update.From.ID != config.Config.GroupInfo.Admin_id {
		global.Log.Infof("%d尝试使用管理员命令", update.From.ID)
		msg := tgbotapi.NewMessage(update.Chat.ID, "你没有权限使用此命令")
		msg.ReplyToMessageID = update.MessageID
		bot.Send(msg)
	} else {
		gettext := utils.Get_command(update.Text, "del", config.Config.BotInfo.Bot_uname)
		if gettext != "" {
			delid, err := strconv.Atoi(gettext)
			if err != nil {
				global.Log.Errorf("数据转换失败: %v", err)
				msg := tgbotapi.NewMessage(update.Chat.ID, "数据转换失败")
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
				return
			}
			result := []int64{}
			flag := false
			for _, v := range config.Config.GroupInfo.Whitelist {
				if v != int64(delid) {
					result = append(result, v)
				} else {
					flag = true
				}
			}
			if flag {
				config.Config.GroupInfo.Whitelist = result
				err = config.UpdateYaml(config.Config)
				if err != nil {
					global.Log.Fatalf("更新配置文件失败: %v", err)
					msg := tgbotapi.NewMessage(update.Chat.ID, "更新配置文件失败")
					msg.ReplyToMessageID = update.MessageID
					bot.Send(msg)
					return
				}
				textreturn := fmt.Sprintf("删除%d成功,当前白名单: %v", delid, config.Config.GroupInfo.Whitelist)
				msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
			} else {
				textreturn := fmt.Sprintf("未在白名单中找到%d,当前白名单: %v", delid, config.Config.GroupInfo.Whitelist)
				msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
				msg.ReplyToMessageID = update.MessageID
				bot.Send(msg)
			}
		} else {
			textreturn := "请输入要删除的ID"
			msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
			msg.ReplyToMessageID = update.MessageID
			bot.Send(msg)
			return
		}
	}
}

func WhiteList(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	if update.From.ID != config.Config.GroupInfo.Admin_id {
		global.Log.Infof("%d尝试使用管理员命令", update.Chat.ID)
		msg := tgbotapi.NewMessage(update.Chat.ID, "你没有权限使用此命令")
		msg.ReplyToMessageID = update.MessageID
		bot.Send(msg)
	} else {
		textreturn := fmt.Sprintf("当前白名单: %v", config.Config.GroupInfo.Whitelist)
		msg := tgbotapi.NewMessage(update.Chat.ID, textreturn)
		msg.ReplyToMessageID = update.MessageID
		bot.Send(msg)
	}
}
