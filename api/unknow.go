// 未知命令
// @author MoGuQAQ
// @version 1.0.0

package api

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func Unknown(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	unknownMessage := "未知命令，请发送 /help 查看可用命令。"
	msg := tgbotapi.NewMessage(update.Chat.ID, unknownMessage)
	msg.ReplyToMessageID = update.MessageID
	bot.Send(msg)
}
