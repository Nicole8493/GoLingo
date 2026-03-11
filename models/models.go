package models

type Article struct {
	ID           int           `json:"id"`
	Translations []Translation `json:"translations"`
	DictionaryID int
	// Pin          int // для закрепления статей
}

type Translation struct {
	ID        int     `json:"id"`
	Language  string  `json:"language"`
	Text      string  `json:"text"`
	ArticleID int     `json:"article_id"` // связ с ID Article foreign key
	Article   Article `json:"article"`    // чтобы призапросе к бд удобнее было подгружать связи (тк в бд нет поля Translations)
}

type Dictionary struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	UserID           int    `json:"user_id"`
	BaseDictionaryID int    `json:"baseDictionaryID"` // айди моего словаря, кот будет отражаться у дргуих пользователей
	User             User   `json:"user"`
}

type Group struct {
	ID     int    `json:"id"`
	Type   string `json:"type"` //доп (например, объед. по слову/смыслу...)
	Name   string `json:"name"` //например, подгруппы по теме
	Color  string `json:"color"`
	UserID int    `json:"user_id"`
}

type ArticleAndGroup struct {
	ArticleID int `json:"article_id"`
	GroupID   int `json:"group_id"`
	//Order     int // порядок, ручная сорировка
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
