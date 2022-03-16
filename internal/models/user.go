package models

import "gorm.io/gorm"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Nickname string `json:"nickname"`
	Cover    string `json:"cover"`
	Avatar   string `json:"avatar"`
	Disabled bool   `json:"disabled" gorm:"default:false"`
	gorm.Model
}

// custom table name
func (User) TableName() string {
	return "t_users"
}
