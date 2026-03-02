package models

type Article struct {
	ID           int           `json:"id"`
	Translations []Translation `json:"translations"`
}

type Translation struct {
	ID        int     `json:"id"`
	Language  string  `json:"language"`
	Text      string  `json:"text"`
	ArticleID int     `json:"article_id"` // связ с ID Article foreign key
	Article   Article `json:"article"`    // чтобы призапросе к бд удобнее было подгружать связи (тк в бд нет поля Translations)
}
