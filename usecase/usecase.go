package usecase

import (
	"crypto/ecdsa"
	"errors"
	db "github.com/Nicole8493/GoLingo/database"
	"github.com/Nicole8493/GoLingo/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Usecase interface {
	CreateArticle(translations models.Article) (id int, err error)
	CreateDictionary(dictionary models.Dictionary) (id int, err error)
	CreateGroup(group models.Group) (id int, err error)
	UpdateTranslations(id int, translations []models.Translation) (err error)
	AddGroupArticles(groupID int, articlesID []int) (err error)
	GetFullArticle(id int) (models.Article, error)
	GetArticle(id int, languages []string) (models.Article, error)
	GetArticlesByGroup(groupID int, languages []string, limit int, offset int, order models.Order) ([]models.Article, error)
	GetArticlesByDictionary(dictionaryID int, languages []string, limit int, offset int, order models.Order) ([]models.Article, error)
	Register(email, name string, password []byte) error
	Login(email string, password []byte) (models.User, string, error)
	DeleteTranslations(articleID int, languages []string) (err error)
	DeleteArticle(id int) (err error)
	DeleteGroup(id int) (err error)
	DeleteDictionary(id int) (err error)
	DeleteGroupArticles(groupID int, articlesID []int) (err error)
}

type UC struct {
	db  *gorm.DB
	key *ecdsa.PrivateKey
}

func New(db *gorm.DB, key *ecdsa.PrivateKey) *UC {
	return &UC{db: db, key: key}
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

func (u *UC) GetArticlesByGroup(groupID int, languages []string, limit int, offset int, order models.Order) (articles []models.Article, err error) {
	request := u.db.Joins("ArticleAndGroup").Where("group_id = ?", groupID).Preload("Translations", "Translations.language IN ?", languages).
		Limit(limit).Offset(offset)

	request, err = u.SortArticles(request, order)
	if err != nil {
		return articles, err
	}

	err = request.Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (u *UC) GetArticlesByDictionary(dictionaryID int, languages []string, limit int, offset int, order models.Order) (articles []models.Article, err error) {
	request := u.db.Where("dictionary_id = ?", dictionaryID).
		Preload("Translations", "Translations.language IN ?", languages).
		Limit(limit).Offset(offset)

	request, err = u.SortArticles(request, order)
	if err != nil {
		return articles, err
	}

	err = request.Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (u *UC) SortArticles(request *gorm.DB, order models.Order) (*gorm.DB, error) {
	if order.Type == "" {
		return request, nil // default
	}

	isDesc := false
	switch order.Direction {
	case "asc", "":
	case "desc":
		isDesc = true
	default:
		return request, errors.New("direction is wrong")
	}
	switch order.Type {
	case "language":
		request = request.Joins("JOIN translations t ON t.article_id = articles.id AND t.language = ?", order.Language).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "t.text"}, Desc: isDesc})
	case "date":
		request = request.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: isDesc})
	default:
		return request, errors.New("order is wrong")
	}
	return request, nil
}

func (u *UC) Register(email, name string, password []byte) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10) // уровень сложности хэширования (средний)
	if err != nil {
		return err
	}
	user := models.User{
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
	}
	err = u.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UC) Login(email string, password []byte) (models.User, string, error) {
	user := models.User{}
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, "", err
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, password)
	if err != nil {
		return user, "", err
	}

	var (
		t    *jwt.Token
		sign string // signed token
	)

	t = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"user_id": user.ID,
		})
	sign, err = t.SignedString(u.key)
	if err != nil {
		return user, "", err
	}
	return user, sign, nil
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
