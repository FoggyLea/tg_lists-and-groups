package app

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

func GetOrCreateUserID(telegramID string) (int, error) {
	var userID int
	err := DB.QueryRow(`SELECT user_id FROM users WHERE telegram_id = $1`, telegramID).Scan(&userID)
	if err == sql.ErrNoRows {
		err = DB.QueryRow(`INSERT INTO users (telegram_id) VALUES ($1) RETURNING user_id`, telegramID).Scan(&userID)
		if err != nil {
			return 0, err
		}
		if _, err := CreateGroup(userID, "default_group"); err != nil {
			return 0, fmt.Errorf("failed to create default_group of user \"%d\": %v", userID, err)
		}
	}
	return userID, err
}

func CreateGroup(userID int, groupName string) (int, error) {
	var groupID int
	err := DB.QueryRow(`INSERT INTO groups (name) VALUES ($1) RETURNING group_id`, groupName).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %v", err)
	}
	_, err = DB.Exec(`INSERT INTO group_members (user_id, group_id, is_admin)
	VALUES ($1, $2, TRUE)`, userID, groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to add creator to group \"%d\": %v", groupID, err)
	}
	//	msg := fmt.Sprintf("Группа \"%s\" создана", groupName)
	return groupID, nil
}

func DeleteGroup(userID, groupID int) ReturnMessage {
	var isAdmin bool
	err := DB.QueryRow(`SELECT is_admin FROM group_members WHERE user_id = $1 AND group_id = $2`, userID, groupID).Scan(&isAdmin)
	if err == sql.ErrNoRows {
		return ReturnMessage{
			Text: defaultFailMsg,
			Err:  fmt.Errorf("user \"%d\" is not the member of the group %d", userID, groupID),
		}
	}
	if err != nil {
		return ReturnMessage{
			Text: defaultFailMsg,
			Err:  fmt.Errorf("error checking admin rights: %v", err),
		}
	}
	if !isAdmin {
		return ReturnMessage{
			Text: "Только админ может удалить группу",
			Err:  fmt.Errorf("user \"%d\" does not have admin rights in the group \"%d\"", userID, groupID),
		}
	}
	_, err = DB.Exec(`DELETE FROM groups WHERE group_id = $1`, groupID)
	if err != nil {
		return ReturnMessage{
			Text: "Не удалось удалить группу",
			Err:  fmt.Errorf("error deleting group \"%d\": %v", groupID, err),
		}
	}
	return ReturnMessage{Text: "Группа удалена"}
}

func AddUserToGroup(userID, groupID int) ReturnMessage {
	_, err := DB.Exec(`INSERT INTO group_members (user_id, group_id, is_admin)
	VALUES ($1, $2, FALSE)`, userID, groupID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return ReturnMessage{Text: "Пользователь уже состоит в группе"}
			case "23503":
				return ReturnMessage{Text: "Пользователь не зарегистирован"}
			}
		}
		return ReturnMessage{Err: fmt.Errorf("failed to add user \"%d\" to group \"%d\": %v", userID, groupID, err)}
	}
	return ReturnMessage{Text: "Пользователь добавлен"}
}

func UserToAdmin(userID, groupID int) ReturnMessage {
	_, err := DB.Exec(`UPDATE group_members SET is_admin = TRUE WHERE group_id = $1 AND user_id = $2`, groupID, userID)
	if err != nil {
		return ReturnMessage{
			Text: "Не удалось назначить пользователя админом",
			Err:  fmt.Errorf("failed to assign user \"%d\" as admin of group \"%d\": %v", userID, groupID, err),
		}
	}
	return ReturnMessage{Text: "Пользователь назначен админом"}
}

func AdminToUser(userID, groupID int) ReturnMessage {
	_, err := DB.Exec(`UPDATE group_members SET is_admin = FALSE WHERE group_id = $1 AND user_id = $2`, groupID, userID)
	if err != nil {
		return ReturnMessage{
			Text: "Не удалось понизить админа",
			Err:  fmt.Errorf("failed to demote user \"%d\" in group \"%d\": %v", userID, groupID, err),
		}
	}
	return ReturnMessage{Text: "Админ понижен до пользователя"}
}

func DeleteUserFromGroup(userID, groupID int) ReturnMessage {
	_, err := DB.Exec(`DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`, groupID, userID)
	if err != nil {
		return ReturnMessage{
			Text: "Не удалось удалить пользователя",
			Err:  fmt.Errorf("failed to delete user \"%d\" from group \"%d\": %v", userID, groupID, err),
		}
	}
	return ReturnMessage{Text: "Пользователь удален"}
}

func LeaveGroup(userID, groupID int) ReturnMessage {
	err := DeleteUserFromGroup(userID, groupID)
	if err.Err != nil {
		return ReturnMessage{Text: "Не удалось покинуть группу", Err: err.Err}
	}
	return ReturnMessage{Text: "Вы покинули группу"}
}

func DeleteProfile(userID int) ReturnMessage {
	groupIDs, err := getAdminGroupIDs(userID)
	if err != nil {
		return ReturnMessage{
			Text: defaultFailMsg,
			Err:  fmt.Errorf("failed to get group IDs where user \"%d\" was admin: %v", userID, err),
		}
	}
	for _, groupID := range groupIDs {
		only, err := IsOnlyAdmin(groupID)
		if err != nil {
			return ReturnMessage{
				Text: defaultFailMsg,
				Err:  fmt.Errorf("failed to check number of admins of group \"%d\": %v", groupID, err),
			}
		}
		if only {
			msg := DeleteGroup(userID, groupID)
			if msg.Err != nil {
				return ReturnMessage{
					Text: "Не удалось удалить одну из групп",
					Err:  msg.Err,
				}
			}
		} else {
			msg := LeaveGroup(userID, groupID)
			if msg.Err != nil {
				return ReturnMessage{
					Text: "Не удалось выйти из одной из групп",
					Err:  msg.Err,
				}
			}
		}
	}
	_, err = DB.Exec(`DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return ReturnMessage{
			Text: "Не удалось удалить профиль",
			Err:  fmt.Errorf("failed to delete user \"%d\": %v", userID, err),
		}
	}
	return ReturnMessage{Text: "Профиль удален"}
}

func getAdminGroupIDs(userID int) ([]int, error) {
	var groupIDs []int
	rows, err := DB.Query(`SELECT group_id FROM group_members WHERE user_id = $1 AND is_admin = TRUE`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var groupID int
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)
	}
	return groupIDs, nil
}

func IsOnlyAdmin(groupID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND is_admin = TRUE`, groupID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func GetGroupIDs(userID int) ([]int, error) {
	rows, err := DB.Query(`SELECT group_id FROM group_members WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group list of user \"%d\": %v", userID, err)
	}
	defer rows.Close()
	var groupIDs []int
	for rows.Next() {
		var groupID int
		if err := rows.Scan(&groupID); err != nil {
			return nil, fmt.Errorf("failed to get groupID of user \"%d\": %v", userID, err)
		}
		groupIDs = append(groupIDs, groupID)
	}
	return groupIDs, nil
}
