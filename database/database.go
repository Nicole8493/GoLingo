package db

type Article struct {
	ID           int `gorm:"primaryKey"`
	Translations []Translation
}

type Translation struct {
	ID        int    `gorm:"primaryKey"`
	Language  string `gorm:"index:articleLanguage"`
	Text      string
	ArticleID int `gorm:"index:articleLanguage"`
	Article   Article
}
