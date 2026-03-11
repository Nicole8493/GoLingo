package db

type Article struct {
	ID           int `gorm:"primaryKey"`
	Translations []Translation
	DictionaryID int
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
	ID               int
	Name             string
	UserID           int
	BaseDictionaryID int // айди моего словаря, кот будет отражаться у дргуих пользователей
	User             User
}

type Group struct {
	ID     int
	Type   string //доп (например, объед. по слову/смыслу...)
	Name   string //например, подгруппы по теме
	Color  string
	UserID int
}

type ArticleAndGroup struct {
	ArticleID int
	GroupID   int
	//Order     int // порядок, ручная сорировка
}

type User struct {
	ID    int
	Name  string
	Email string
}
