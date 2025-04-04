// help命令
// @author MoGuQAQ
// @version 1.0.0

package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Help(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	helpMessage := "欢迎使用本机器人\n指令列表:\n/help - 帮助\n/add - 添加白名单\n/del - 删除白名单\n/whitelist - 查看白名单"
	msg := tgbotapi.NewMessage(update.Chat.ID, helpMessage)
	msg.ReplyToMessageID = update.MessageID
	bot.Send(msg)
}
