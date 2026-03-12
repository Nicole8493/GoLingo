package usecase

import (
	db "github.com/Nicole8493/GoLingo/database"
	"github.com/Nicole8493/GoLingo/models"
	"gorm.io/gorm"
)

type Usecase interface {
	CreateArticle(translations models.Article) (id int, err error)
	CreateDictionary(dictionary models.Dictionary) (id int, err error)
	CreateGroup(group models.Group) (id int, err error)
	UpdateTranslations(id int, translations []models.Translation) (err error)
	AddGroupArticles(groupID int, articlesID []int) (err error)
	GetFullArticle(id int) (models.Article, error)
	GetArticle(id int, languages []string) (models.Article, error)
	GetArticlesByGroup(groupID int, languages []string, limit int, offset int) ([]models.Article, error)
	GetArticlesByDictionary(dictionaryID int, languages []string, limit int, offset int) ([]models.Article, error)
	DeleteTranslations(articleID int, languages []string) (err error)
	DeleteArticle(id int) (err error)
	DeleteGroup(id int) (err error)
	DeleteDictionary(id int) (err error)
	DeleteGroupArticles(groupID int, articlesID []int) (err error)
}

type UC struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UC {
	return &UC{db: db}
}

func (u *UC) CreateArticle(article models.Article) (id int, err error) {
	articleDB := db.Article{
		ID:           0,
		Translations: make([]db.Translation, len(article.Translations)),
	}
	for i, translation := range article.Translations {
		articleDB.Translations[i] = db.Translation{
			ID:       translation.ID,
			Language: translation.Language,
			Text:     translation.Text,
		}
	}
	err = u.db.Create(&articleDB).Error
	if err != nil {
		return 0, err
	}
	return articleDB.ID, nil
}

func (u *UC) CreateDictionary(dictionary models.Dictionary) (id int, err error) {
	dictionaryDB := db.Dictionary{
		ID:               0,
		Name:             dictionary.Name,
		UserID:           dictionary.UserID,
		BaseDictionaryID: dictionary.BaseDictionaryID,
	}
	err = u.db.Create(&dictionaryDB).Error
	if err != nil {
		return 0, err
	}
	return dictionaryDB.ID, nil
}

func (u *UC) CreateGroup(group models.Group) (id int, err error) {
	groupDB := db.Group{
		ID:     0,
		Type:   group.Type,
		Name:   group.Name,
		UserID: group.UserID,
		Color:  group.Color,
	}
	err = u.db.Create(&groupDB).Error
	if err != nil {
		return 0, err
	}
	return groupDB.ID, nil
}

func (u *UC) UpdateTranslations(id int, translations []models.Translation) (err error) {
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

func (u *UC) AddGroupArticles(groupID int, articlesID []int) (err error) {
	var group = make([]db.ArticleAndGroup, 0, len(articlesID))
	for _, id := range articlesID {
		group = append(group, db.ArticleAndGroup{
			ArticleID: id,
			GroupID:   groupID,
		})
	}
	err = u.db.Create(&group).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UC) GetFullArticle(id int) (models.Article, error) {
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

func (u *UC) GetArticle(id int, languages []string) (models.Article, error) {
	result := db.Article{}
	err := u.db.
		Joins("JOIN translations ON translations.language IN ? AND article_id = articles.id", languages). // выбранные языки
		First(&result, id).                                                                               // указываем, что нужно подгрузить статью по айди
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

func (u *UC) GetArticlesByGroup(groupID int, languages []string, limit int, offset int) (articles []models.Article, err error) {
	err = u.db.Joins("ArticleAndGroup").Where("group_id = ?", groupID).Preload("Translations", "Translations.language IN ?", languages).
		Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (u *UC) GetArticlesByDictionary(dictionaryID int, languages []string, limit int, offset int) (articles []models.Article, err error) {
	err = u.db.Where("dictionary_id = ?", dictionaryID).
		Preload("Translations", "Translations.language IN ?", languages).
		Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (u *UC) DeleteTranslations(articleID int, languages []string) (err error) {
	return u.db.Delete(&db.Translation{}, "article_id = ? AND translations.language IN ?", articleID, languages).Error
}

func (u *UC) DeleteArticle(id int) (err error) {
	u.db.Delete(&db.Translation{}, "article_id = ?", id)
	return u.db.Delete(&db.Article{}, "id = ?", id).Error
}

func (u *UC) DeleteGroup(id int) (err error) {
	return u.db.Delete(&db.Group{}, "id = ?", id).Error
}

func (u *UC) DeleteDictionary(id int) (err error) {
	return u.db.Delete(&db.Dictionary{}, "id = ?", id).Error
}

func (u *UC) DeleteGroupArticles(groupID int, articlesID []int) (err error) {
	return u.db.Delete(&db.ArticleAndGroup{}, "article_id IN ? AND group_id = ?", articlesID, groupID).Error
}
