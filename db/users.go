package db

import (
	"crypto/sha512"
	"fmt"
)

type User struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Email string `gorm:"type:varchar(255);not null"`
	Pwd   string `gorm:"type:varchar(255);not null"`
}

const pwd_salt = "3JLk7aNr4GzpQ2u8w1eTMY6v9EKxpBuW"

func (*User) TableName() string {
	return "users.users"
}

func GetUser(username string, pwd string) *User {
	var user User

	pwd_hash := fmt.Sprintf("%x", sha512.Sum512([]byte(pwd+pwd_salt)))
	res := DB.Where("email = ? and pwd = ?", username, pwd_hash).First(&user)

	if res.Error == nil && res.RowsAffected > 0 {
		return &user
	} else {
		return nil
	}
}
