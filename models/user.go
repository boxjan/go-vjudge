package models

import "time"

const UserTableName = "users"

type User struct {
	Id uint64 `orm:"pk;auto"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
	DeletedAt time.Time `orm:"null"`
	Username string `orm:"size(64);unique"`
	Nickname string `orm:"size(64)"`
	Password string `orm:"size(255)"`
	Role string `orm:"size(16);default(user)"`
	Email string `orm:"size(255);unique"`
	EmailVerifiedAt time.Time `orm:"null"`
	School string `orm:"size(128)"`
	RememberToken string `orm:"size(128)"`
	RegisterIp string `orm:"size(40);"`
}

func (u *User) TableName() string {
	return UserTableName
}
