// start命令
// @author MoGuQAQ
// @version 1.0.0

package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	welcomeMessage := "欢迎使用本机器人！\n\n" +
		"请发送 /help 查看可用命令。"
	msg := tgbotapi.NewMessage(update.Chat.ID, welcomeMessage)
	msg.ReplyToMessageID = update.MessageID
	bot.Send(msg)
}
