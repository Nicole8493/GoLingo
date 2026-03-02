package usecase

import (
	db "github.com/Nicole8493/GoLingo/database"
	"github.com/Nicole8493/GoLingo/models"
	"gorm.io/gorm"
)

type Usecase interface {
	CreateArticle(translations []models.Translation) (id int, err error)
	UpdateTranslations(id int, translations []models.Translation) (err error)
	GetFullArticle(id int) (models.Article, error)
	GetArticle(id int, languages []string) (models.Article, error)
	DeleteTranslations(articleID int, languages []string) (err error)
	DeleteArticle(id int) (err error)
}

type UC struct {
	db *gorm.DB
}

func (u UC) CreateArticle(translations []models.Translation) (id int, err error) {
	article := db.Article{
		ID:           0,
		Translations: make([]db.Translation, len(translations)),
	}
	for i, translation := range translations {
		article.Translations[i] = db.Translation{
			ID:       translation.ID,
			Language: translation.Language,
			Text:     translation.Text,
		}
	}
	err = u.db.Create(&article).Error
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (u UC) UpdateTranslations(id int, translations []models.Translation) (err error) {
	for _, translation := range translations {
		err = u.db.Save(&db.Translation{
			ID:        translation.ID,
			Language:  translation.Language,
			Text:      translation.Text,
			ArticleID: id,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (u UC) GetFullArticle(id int) (models.Article, error) {
	result := db.Article{}
	err := u.db.
		Joins("Translations"). // указываем, что нужно подгрузить все переводы для статьи
		First(&result, id).    // указываем, что нужно подгрузить статью по айди
		Error
	if err != nil {
		return models.Article{}, err
	}

	article := models.Article{
		ID:           result.ID,
		Translations: make([]models.Translation, len(result.Translations)),
	}
	for i, translation := range result.Translations {
		article.Translations[i] = models.Translation{
			Language:  translation.Language,
			Text:      translation.Text,
			ID:        translation.ID,
			ArticleID: id,
		}
	}
	return article, nil
}

func (u UC) GetArticle(id int, languages []string) (models.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (u UC) DeleteTranslations(articleID int, languages []string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u UC) DeleteArticle(id int) (err error) {
	//TODO implement me
	panic("implement me")
}
