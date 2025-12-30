package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RoleType string

const (
	USER      RoleType = "user"
	STAFF     RoleType = "staff"
	SUPERUSER RoleType = "superuser"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
}

type Group struct {
	gorm.Model
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:group_permissions;"`
}

type User struct {
	gorm.Model
	Name     string   `gorm:"not null"`
	Email    string   `gorm:"unique;not null;<-:create"`
	Password string   `gorm:"not null"`
	Role     RoleType `gorm:"type:varchar(20);default:'user'"`
	Groups   []Group  `gorm:"many2many:user_groups;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	u.Password = string(bytes)

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	var oldUser User
	if err := tx.First(&oldUser, u.ID).Error; err != nil {
		log.Println("Error: ", err)
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(oldUser.Password), []byte(u.Password))
	if err != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)

		return nil
	}

	u.Password = ""

	return nil
}
