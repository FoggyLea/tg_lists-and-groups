package tgbot

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func InitRouter(bh *th.BotHandler, bot *telego.Bot) {

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		handleListInteractions(bot, ctx.Context(), *update.Message)
		return nil
	}, th.CommandEqual("list"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		handleMyGroups(bot, ctx.Context(), *update.Message)
		return nil
	}, th.CommandEqual("mygroups"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		handleGroupSettings(bot, ctx.Context(), *update.Message)
		return nil
	}, th.CommandEqual("groupsettings"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		handleDeleteProfile(bot, ctx.Context(), *update.Message)
		return nil
	}, th.CommandEqual("deleteprofile"))

	// callback заглушка
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		handleCallback(bot, ctx.Context(), *update.CallbackQuery)
		return nil
	}, th.AnyCallbackQuery())
}
