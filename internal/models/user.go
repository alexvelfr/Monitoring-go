package models

import (
	"fmt"
)

//User model
type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	Phone        string `db:"phone"`
	TelegramID   int64  `db:"telegram_id"`
	SendMessages bool   `db:"send_messages"`
	SendReports  bool   `db:"send_reports"`
	IsAdmin      bool   `db:"is_admin"`
}

//SelectByTeledgramID - return query string for select one user by telegram_id
func (u *User) SelectByTeledgramID(id int64) string {
	return fmt.Sprintf(`SELECT * FROM users WHERE telegram_id=%d LIMIT 1`, id)
}

//UpdateQuery - return query string for update user by telegram_id
func (u *User) UpdateQuery() string {
	return `UPDATE users SET phone=:phone, send_messages=:send_messages, send_reports=:send_reports, is_admin=:is_admin, name=:name WHERE telegram_id=:telegram_id`
}

//GetAllRecipients - return all recipients
func GetAllRecipients() string {
	return `SELECT * FROM users WHERE send_messages=true`
}

//GetAllReportRecipients - return all report recipients
func GetAllReportRecipients() string {
	return `SELECT * FROM users WHERE send_messages=true and send_reports=true`
}
