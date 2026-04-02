package db

import "time"

type Article struct {
	ID           int `gorm:"primaryKey"`
	Translations []Translation
	DictionaryID int
	Dictionary   Dictionary
	CreatedAt    time.Time
	// Pin          int // для закрепления статей
}

type Translation struct {
	ID        int    `gorm:"primaryKey"`
	Language  string `gorm:"index:articleLanguage"`
	Text      string
	ArticleID int `gorm:"index:articleLanguage"` // внеш ключ
	Article   Article
}

type Dictionary struct {
	ID               int `gorm:"primaryKey"`
	Name             string
	UserID           int
	BaseDictionaryID int         `gorm:"default:null"` // айди моего словаря, кот будет отражаться у дргуих пользователей?, ставим нил по умолч
	BaseDictionary   *Dictionary // встраиваем Dictionary в самого себя для создания рекурсивной ссылки
	User             User
}

type Group struct {
	ID     int
	Type   string //доп (например, объед. по слову/смыслу...)
	Name   string //например, подгруппы по теме
	Color  string
	UserID int
	User   User
}

type ArticleAndGroup struct {
	ArticleID int
	GroupID   int
	Group     Group
	Article   Article
	//Order     int // порядок, ручная сорировка
}

type User struct {
	ID           int
	Name         string
	Email        string `gorm:"uniqueIndex"` // не будет 2 юзеров с одинаковым мэйлом
	PasswordHash []byte
}
