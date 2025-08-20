package app

import (
	"errors"
	"fmt"
	"strings"
)

func parseItems(raw string) []string {
	split := strings.Split(raw, ",")
	items := []string{}
	for _, item := range split {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			items = append(items, strings.ToLower(trimmed))
		}
	}
	return items
}

func ShowList(groupID int) ReturnMessage {
	rows, err := DB.Query(`SELECT item FROM shopping_items WHERE group_id = $1`, groupID)
	if err != nil {
		return ReturnMessage{
			Text: "Ошибка при получении списка",
			Err:  fmt.Errorf("failed to get list for group \"%d\": %v", groupID, err),
		}
	}
	defer rows.Close()
	i, hasRows, msg := 1, false, ""
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			msg = "Ошибка чтения строки"
			return ReturnMessage{
				Text: msg,
				Err:  fmt.Errorf("failed to get next list row for group \"%d\": %v", groupID, err),
			}
		}
		if !hasRows {
			msg = "Ваш список:\n"
			hasRows = true
		}
		msg += fmt.Sprintf("%d. %s\n", i, item)
		i++
	}
	if !hasRows {
		msg = "Список пуст"
	}
	return ReturnMessage{Text: msg}
}

func AddItem(groupID int, itemsRaw string) ReturnMessage {
	items := parseItems(itemsRaw)
	if len(items) == 0 {
		return ReturnMessage{Text: "Пустой ввод"}
	}

	var (
		updated                      bool
		duplicates, failed, errorMsg []string
		msg                          string
	)
	for _, item := range items {
		var exists bool
		err := DB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM shopping_items WHERE group_id = $1 AND item = $2)`, groupID, item).Scan(&exists)
		if err != nil {
			errorMsg = append(errorMsg, fmt.Sprintf("failed to verify existence of item \"%s\" in group %d: %v", item, groupID, err))
			continue
		}
		if exists {
			duplicates = append(duplicates, item)
			continue
		}
		_, err = DB.Exec(`INSERT INTO shopping_items (group_id, item) VALUES ($1, $2)`, groupID, item)
		if err != nil {
			errorMsg = append(errorMsg, fmt.Sprintf("failed to add item \"%s\" to group %d: %v", item, groupID, err))
			failed = append(failed, item)
		} else {
			updated = true
		}
	}
	if len(failed) > 0 {
		msg += fmt.Sprintf("Не удалось добавить: %s\n", strings.Join(failed, ", "))
	}
	if len(duplicates) > 0 {
		msg += fmt.Sprintf("Уже в списке: %s\n", strings.Join(duplicates, ", "))
	}
	if updated {
		msg += "Список обновлен"
	}
	var err error
	if len(errorMsg) > 0 {
		err = errors.New(strings.Join(errorMsg, "\n"))
	}
	return ReturnMessage{
		Text: msg,
		Err:  err,
	}
}

func RemoveItem(groupID int, itemsRaw string) ReturnMessage {
	items := parseItems(itemsRaw)
	if len(items) == 0 {
		return ReturnMessage{Text: "Предметы не выбраны"}
	}
	var (
		removed  int
		errorMsg []string
		msg      string
	)
	for _, item := range items {
		res, err := DB.Exec(`DELETE FROM shopping_items WHERE group_id = $1 AND item = $2`, groupID, item)
		if err != nil {
			errorMsg = append(errorMsg, fmt.Sprintf("Ошибка при удалении \"%s\": %v", item, err))
			continue
		}
		count, _ := res.RowsAffected()
		if count > 0 {
			removed++
		}
	}
	if removed > 0 {
		msg = "Удаление выполнено"
	} else {
		msg = "Ничего не было удалено"
	}

	var err error
	if len(errorMsg) > 0 {
		err = errors.New(strings.Join(errorMsg, "\n"))
	}

	return ReturnMessage{
		Text: msg,
		Err:  err,
	}
}
