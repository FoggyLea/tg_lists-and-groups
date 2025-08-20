package tgbot

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func GroupSettingKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Добавить участника").WithCallbackData("addmember")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Повысить до админа").WithCallbackData("promote")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Понизить админа").WithCallbackData("remote")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Удалить группу").WithCallbackData("deletegroup")),
	)
}

func MyGroupsKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Группа 1").WithCallbackData("group1")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Создать группу").WithCallbackData("newgroup")),
	)
}

func ListInteractions() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Просмотреть список").WithCallbackData("showlist")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Добавить в список").WithCallbackData("additem")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Удалить из списка").WithCallbackData("removeitem")),
	)
}

func DeleteOrNot() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Да").WithCallbackData("deleteid")),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Нет").WithCallbackData("notdeleteid")),
	)
}
