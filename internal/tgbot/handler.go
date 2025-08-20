package tgbot

import (
	"context"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SetBotMenu(ctx context.Context, bot *telego.Bot) error {
	commands := []telego.BotCommand{
		{Command: "list", Description: "список"},
		{Command: "mygroups", Description: "мои группы"},
		{Command: "groupsettings", Description: "настройки группы"},
		{Command: "deleteprofile", Description: "удалить профиль"},
	}
	return bot.SetMyCommands(ctx, &telego.SetMyCommandsParams{Commands: commands})
}

func handleMyGroups(bot *telego.Bot, ctx context.Context, msg telego.Message) {
	_, _ = bot.SendMessage(ctx,
		tu.Message(tu.ID(msg.Chat.ID), "Ваши группы:").WithReplyMarkup(MyGroupsKeyboard()))
}

func handleGroupSettings(bot *telego.Bot, ctx context.Context, msg telego.Message) {
	_, _ = bot.SendMessage(ctx,
		tu.Message(tu.ID(msg.Chat.ID), "Натройки группы:").WithReplyMarkup(GroupSettingKeyboard()))
}

func handleListInteractions(bot *telego.Bot, ctx context.Context, msg telego.Message) {
	_, _ = bot.SendMessage(ctx,
		tu.Message(tu.ID(msg.Chat.ID), "Управление списком:").WithReplyMarkup(ListInteractions()))
}

func handleDeleteProfile(bot *telego.Bot, ctx context.Context, msg telego.Message) {
	_, _ = bot.SendMessage(ctx,
		tu.Message(tu.ID(msg.Chat.ID), "Вы уверены, что ходтите удалить профиль?").WithReplyMarkup(DeleteOrNot()))
}

func handleCallback(bot *telego.Bot, ctx context.Context, cb telego.CallbackQuery) {
	_ = bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
		Text:            "Вы нажали: " + cb.Data,
		ShowAlert:       false,
	})
}
