// 机器人信息输出
// @author MoGuQAQ
// @version 1.0.0

package telebot

import (
	"TeleGPT/config"
	"TeleGPT/global"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
			if update.Message.Chat.ID != adminid && !contains(config.Config.GroupInfo.Whitelist, update.Message.Chat.ID) {
				global.Log.Infof("%d 非白名单账户,主动忽略", update.Message.Chat.ID)
				continue
			}
		}
		global.Log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		mothod := 0
		replyid := 0
		if update.Message.Chat.IsPrivate() {
			mothod = 2
		}
		if update.Message.IsCommand() {
			mothod = 1
		}
		/*
			replyid = update.Message.Chat.
			global.Log.Debugf("回复消息ID: %d", replyid)
			global.Log.Debugf("当前模式: %d", mothod)
			todo 获取回复状态
		*/
		Vector(bot, update.Message.Text, update.Message.MessageID, update.Message.Chat.ID, replyid, mothod)
	}
}

func Vector(bot *tgbotapi.BotAPI, text string, messageid int, chatid int64, replayid int, mothod int) {
	if strings.Contains(text, "/") && mothod == 1 {
		if is_command(text, "start", config.Config.BotInfo.Bot_uname) {
			textreturn := "欢迎使用本机器人"
			msg := tgbotapi.NewMessage(chatid, textreturn)
			bot.Send(msg)
		} else if is_command(text, "help", config.Config.BotInfo.Bot_uname) {
			textreturn := "欢迎使用本机器人\n指令列表:\n/help - 帮助\n/add - 添加白名单\n/del - 删除白名单\n/whitelist - 查看白名单"
			msg := tgbotapi.NewMessage(chatid, textreturn)
			bot.Send(msg)
		} else if is_command(text, "add", config.Config.BotInfo.Bot_uname) {
			if chatid != int64(config.Config.GroupInfo.Admin_id) {
				global.Log.Infof("%d尝试使用管理员命令", chatid)
				msg := tgbotapi.NewMessage(chatid, "你没有权限使用此命令")
				msg.ReplyToMessageID = messageid
				bot.Send(msg)
			} else {
				pattern := fmt.Sprintf(`^\/add(@%s)?\s+(\d+)$`, config.Config.BotInfo.Bot_uname)
				re := regexp.MustCompile(pattern)
				matches := re.FindStringSubmatch(text)
				if len(matches) > 2 {
					addid, err := strconv.Atoi(matches[2])
					if err != nil {
						global.Log.Errorf("数据转换失败: %v", err)
						msg := tgbotapi.NewMessage(chatid, "数据转换失败")
						msg.ReplyToMessageID = messageid
						bot.Send(msg)
						return
					}
					if contains(config.Config.GroupInfo.Whitelist, int64(addid)) {
						textreturn := fmt.Sprintf("%d已在白名单中,当前白名单: %v", addid, config.Config.GroupInfo.Whitelist)
						msg := tgbotapi.NewMessage(chatid, textreturn)
						msg.ReplyToMessageID = messageid
						bot.Send(msg)
						return
					}
					var whitelistPointer *[]int64 = &config.Config.GroupInfo.Whitelist
					*whitelistPointer = append(*whitelistPointer, int64(addid))
					err = config.UpdateYaml(config.Config)
					if err != nil {
						global.Log.Fatalf("更新配置文件失败: %v", err)
						msg := tgbotapi.NewMessage(chatid, "更新配置文件失败")
						msg.ReplyToMessageID = messageid
						bot.Send(msg)
						return
					}
					textreturn := fmt.Sprintf("添加%d成功,当前白名单: %v", addid, config.Config.GroupInfo.Whitelist)
					msg := tgbotapi.NewMessage(chatid, textreturn)
					msg.ReplyToMessageID = messageid
					bot.Send(msg)
				} else {
					textreturn := "请输入要添加的ID"
					msg := tgbotapi.NewMessage(chatid, textreturn)
					msg.ReplyToMessageID = messageid
					bot.Send(msg)
					return
				}
			}
		} else if is_command(text, "del", config.Config.BotInfo.Bot_uname) {
			if chatid != int64(config.Config.GroupInfo.Admin_id) {
				global.Log.Infof("%d尝试使用管理员命令", chatid)
				msg := tgbotapi.NewMessage(chatid, "你没有权限使用此命令")
				msg.ReplyToMessageID = messageid
				bot.Send(msg)
			} else {
				pattern := fmt.Sprintf(`^\/del(@%s)?\s+(\d+)$`, config.Config.BotInfo.Bot_uname)
				re := regexp.MustCompile(pattern)
				matches := re.FindStringSubmatch(text)
				if len(matches) > 2 {
					delid, err := strconv.Atoi(matches[2])
					if err != nil {
						global.Log.Errorf("数据转换失败: %v", err)
						msg := tgbotapi.NewMessage(chatid, "数据转换失败")
						msg.ReplyToMessageID = messageid
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
							msg := tgbotapi.NewMessage(chatid, "更新配置文件失败")
							msg.ReplyToMessageID = messageid
							bot.Send(msg)
							return
						}
						textreturn := fmt.Sprintf("删除%d成功,当前白名单: %v", delid, config.Config.GroupInfo.Whitelist)
						msg := tgbotapi.NewMessage(chatid, textreturn)
						msg.ReplyToMessageID = messageid
						bot.Send(msg)
					} else {
						textreturn := fmt.Sprintf("未在白名单中找到%d,当前白名单: %v", delid, config.Config.GroupInfo.Whitelist)
						msg := tgbotapi.NewMessage(chatid, textreturn)
						msg.ReplyToMessageID = messageid
						bot.Send(msg)
					}
				} else {
					msg := tgbotapi.NewMessage(chatid, "请输入要删除的ID")
					msg.ReplyToMessageID = messageid
					bot.Send(msg)
					return
				}
			}
		} else if is_command(text, "whitelist", config.Config.BotInfo.Bot_uname) {
			if chatid != int64(config.Config.GroupInfo.Admin_id) {
				global.Log.Infof("%d尝试使用管理员命令", chatid)
				msg := tgbotapi.NewMessage(chatid, "你没有权限使用此命令")
				msg.ReplyToMessageID = messageid
				bot.Send(msg)
			} else {
				textreturn := fmt.Sprintf("当前白名单: %v", config.Config.GroupInfo.Whitelist)
				msg := tgbotapi.NewMessage(chatid, textreturn)
				msg.ReplyToMessageID = messageid
				bot.Send(msg)
			}
		} else if is_command(text, "model", config.Config.BotInfo.Bot_uname) {
			if chatid != int64(config.Config.GroupInfo.Admin_id) {
				global.Log.Infof("%d尝试使用管理员命令", chatid)
				msg := tgbotapi.NewMessage(chatid, "你没有权限使用此命令")
				msg.ReplyToMessageID = messageid
				bot.Send(msg)
			} else {
				pattern := fmt.Sprintf(`^\/model(@%s)?\s+(.*)$`, config.Config.BotInfo.Bot_uname)
				re := regexp.MustCompile(pattern)
				matches := re.FindStringSubmatch(text)
				model := config.Config.API.Omodel
				if len(matches) > 2 {
					model = matches[2]
					if len(model) > 25 {
						msg := tgbotapi.NewMessage(chatid, "模型名称过长,请手动修改配置文件")
						msg.ReplyToMessageID = messageid
						bot.Send(msg)
						return
					}
				} else {
					textreturn := fmt.Sprintf("当前模型为: %s", model)
					msg := tgbotapi.NewMessage(chatid, textreturn)
					msg.ReplyToMessageID = messageid
					bot.Send(msg)
					return
				}
				config.Config.API.Omodel = model
				err := config.UpdateYaml(config.Config)
				if err != nil {
					global.Log.Fatalf("更新配置文件失败: %v", err)
					msg := tgbotapi.NewMessage(chatid, "更新配置文件失败")
					msg.ReplyToMessageID = messageid
					bot.Send(msg)
					return
				}
				msg := tgbotapi.NewMessage(chatid, "修改模型成功")
				msg.ReplyToMessageID = messageid
				bot.Send(msg)
			}
		} else {
			global.Log.Infof("非指令,主动忽略")
		}
	} else if strings.Contains(text, "@") || mothod == 2 {
		if mothod == 2 {
			// 私聊模式 无需指令
		} else {
			pattern := fmt.Sprintf(`^@%s(?:\S+)?(?:\s+\S+.*)?$`, config.Config.BotInfo.Bot_uname)
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(text)
			if len(matches) > 2 {
				text := matches[2]
				if text != "" {
					textreturn := text
					msg := tgbotapi.NewMessage(chatid, textreturn)
					msg.ReplyToMessageID = messageid
					bot.Send(msg)
				} else {
					global.Log.Infof("空输入,主动忽略")
					return
				}
			} else {
				global.Log.Infof("空输入,主动忽略")
				return
			}
		}
	} else {
		global.Log.Infof("未匹配的模式,主动忽略")
		return
	}
}

func contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func is_command(str string, command string, uname string) bool {
	pattern := fmt.Sprintf(`^\/%s(?:@%s)?(?:\s+\S.*)?$`, command, uname)
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return true
	} else {
		return false
	}
}
