package gormstorage

import "time"

type UserModel struct {
	ID           string `gorm:"primaryKey;column:id"`
	Email        string `gorm:"column:email;not null"`
	Username     string `gorm:"column:username;not null"`
	PasswordHash string `gorm:"column:password_hash;not null"`
	GlobalRole   string `gorm:"column:global_role;not null"`
	TokenVersion int    `gorm:"column:token_version;not null;default:1"`
	Active       bool   `gorm:"column:active;not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (UserModel) TableName() string {
	return "users"
}
