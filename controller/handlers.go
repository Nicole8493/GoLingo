package controller

import (
	"github.com/Nicole8493/GoLingo/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
)

func (c *Controller) getUserID(ctx *fiber.Ctx) int {
	// мидлваря проверяет токен и сохраняет его в Locals (после хэндлер достает токен из Locals)
	user := ctx.Locals("user").(*jwt.Token)    // получаем токен
	claims := user.Claims.(jwt.MapClaims)      // достаем полезную нагрузку (айди юзера)
	userID := int(claims["user_id"].(float64)) // claims["user_id"] получаем пустой интерфейс, .(float64) конвертируем интер-с в флоат, а затем в инт
	return userID
}

// handlerCreateArticle
// @Summary      Create a new article
// @Description  Create a new article in given dictionary (user must be dictionary owner)
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        article   body   models.Article  true   "article with translations"
// @Success      200  {object}  int
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /article [post]
func (c *Controller) handlerCreateArticle(ctx *fiber.Ctx) error {
	// достаем данные для передачи в юзкейс
	data := new(models.Article)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	// вызов юзкейса
	articleID, err := c.usecase.CreateArticle(c.getUserID(ctx), *data)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON(articleID)
}

// handlerCreateDictionary
// @Summary      Create a new dictionary
// @Description  Create a new dictionary
// @Tags         dictionaries
// @Accept       json
// @Produce      json
// @Param        dictionary   body  models.Dictionary  true   "new dictionary"
// @Success      200  {object}  int
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /dictionary [post]
func (c *Controller) handlerCreateDictionary(ctx *fiber.Ctx) error {

	data := new(models.Dictionary)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	id, err := c.usecase.CreateDictionary(c.getUserID(ctx), *data)
	if err != nil {
		return err
	}
	return ctx.JSON(id)
}

// handlerCreateGroup
// @Summary      Create a new group
// @Description  Create a new group
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        group   body   models.Group  true   "new group"
// @Success      200  {object}  int
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /group [post]
func (c *Controller) handlerCreateGroup(ctx *fiber.Ctx) error {
	data := new(models.Group)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	id, err := c.usecase.CreateGroup(c.getUserID(ctx), *data)
	if err != nil {
		return err
	}
	return ctx.JSON(id)
}

// handlerUpdateTranslations
// @Summary      Update translations
// @Description  Update translations in specific article (user must be dictionary owner)
// @Tags         translations
// @Accept       json
// @Produce      json
// @Param        translations   body   []models.Translation  true   "translations to update"
// @Param        articleID      path   int   true   "articleID"
// @Success      200  {object}  int
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /translations/{id} [post]
func (c *Controller) handlerUpdateTranslations(ctx *fiber.Ctx) error {
	data := new([]models.Translation)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}

	articleID := ctx.Params("id")
	idInt, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	// вызов юзкейса
	err = c.usecase.UpdateTranslations(c.getUserID(ctx), idInt, *data)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON("ok")
}

// handlerAddGroupArticles
// @Summary      Add articles to group
// @Description  Add articles to group (user must be group owner)
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        groupID   path   int  true  "groupID to save articles"
// @Param        articles  query   []int   true  "articleIDs comma separated"
// @Success      200  {object}  int
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /group/articles/{id} [post]
func (c *Controller) handlerAddGroupArticles(ctx *fiber.Ctx) error {
	groupID := ctx.Params("id")
	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		return err
	}

	articles := ctx.Query("articles")
	articlesList := strings.Split(articles, ",")        // articlesList список строк
	articlesIDsInt := make([]int, 0, len(articlesList)) // длина 0, чтобы не ограничиваться и при апенде не наращивать слайс
	for _, id := range articlesList {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		articlesIDsInt = append(articlesIDsInt, idInt)
	}

	err = c.usecase.AddGroupArticles(c.getUserID(ctx), groupIDInt, articlesIDsInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

// handlerGetFullArticle
// @Summary      Get full list of translations in article
// @Description  Get full list of translations in article
// @Tags         articles
// @Produce      json
// @Param        articleID   path   int  true  "articleID to get all translations"
// @Success      200  {object}  models.Article
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /article/full/{id} [get]
func (c *Controller) handlerGetFullArticle(ctx *fiber.Ctx) error {
	articleID := ctx.Params("id")
	idInt, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}

	// вызов юзкейса
	data, err := c.usecase.GetFullArticle(idInt)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON(data)
}

// handlerGetArticle
// @Summary      Get list of translations in article
// @Description  Get list of translations in article
// @Tags         articles
// @Produce      json
// @Param        articleID   path   int  true  "articleID to get translations"
// @Param        languages  query  []string  true  "languages comma separated"
// @Success      200  {object}  models.Article
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /article/{id} [get]
func (c *Controller) handlerGetArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	languages := ctx.Query("languages")
	languagesList := strings.Split(languages, ",")

	// вызов юзкейса
	data, err := c.usecase.GetArticle(idInt, languagesList)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON(data)
}

// handlerGetArticlesByGroup
// @Summary      Get articles by group
// @Description  Get articles by group
// @Tags         groups
// @Produce      json
// @Param        groupID   path  int  true  "groupID to fet articles"
// @Param        languages  query  []string  true  "languages comma separated"
// @Param        limit  query  int  true  "limit"
// @Param        offset  query  int  true  "offset"
// @Param        order  query  string  true  "order"
// @Param        orderDirection  query  string  true  "orderDirection"
// @Param        orderLanguage  query  string  true  "orderLanguage"
// @Success      200  {object}  []models.Article
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /articles/group/{id} [get]
func (c *Controller) handlerGetArticlesByGroup(ctx *fiber.Ctx) error {
	groupID := ctx.Params("id")
	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		return err
	}
	languages := ctx.Query("languages")
	languagesList := strings.Split(languages, ",")
	limitString := ctx.Query("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return err
	}
	offsetString := ctx.Query("offset")
	offsetInt, err := strconv.Atoi(offsetString)
	if err != nil {
		return err
	}
	order := ctx.Query("order")
	orderDirection := ctx.Query("orderDirection")
	orderLanguage := ctx.Query("orderLanguage")
	modelOrder := models.Order{
		Type:      order,
		Direction: orderDirection,
		Language:  orderLanguage,
	}

	// вызов юзкейса
	articles, err := c.usecase.GetArticlesByGroup(groupIDInt, languagesList, limit, offsetInt, modelOrder)
	if err != nil {
		return err
	}
	return ctx.JSON(articles)
}

// handlerGetArticlesByDictionary
// @Summary      Get articles by dictionary
// @Description  Get articles by dictionary
// @Tags         dictionaries
// @Produce      json
// @Param        dictionaryID   path   int  true  "dictionaryID to get articles"
// @Param        languages  query  []string  true  "languages comma separated"
// @Param        limit  query  int  true  "limit"
// @Param        offset  query  int  true  "offset"
// @Param        order  query  string  true  "order"
// @Param        orderDirection  query  string  true  "orderDirection"
// @Param        orderLanguage  query  string  true  "orderLanguage"
// @Success      200  {object}  []models.Article
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /articles/dictionary/{id} [get]
func (c *Controller) handlerGetArticlesByDictionary(ctx *fiber.Ctx) error {
	dictionaryID := ctx.Params("id")
	dictionaryIDInt, err := strconv.Atoi(dictionaryID)
	if err != nil {
		return err
	}
	languages := ctx.Query("languages")
	languagesList := strings.Split(languages, ",")
	limitString := ctx.Query("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return err
	}
	offsetString := ctx.Query("offset")
	offsetInt, err := strconv.Atoi(offsetString)
	if err != nil {
		return err
	}
	order := ctx.Query("order")
	orderDirection := ctx.Query("orderDirection")
	orderLanguage := ctx.Query("orderLanguage")
	modelOrder := models.Order{
		Type:      order,
		Direction: orderDirection,
		Language:  orderLanguage,
	}
	articles, err := c.usecase.GetArticlesByDictionary(dictionaryIDInt, languagesList, limit, offsetInt, modelOrder)
	if err != nil {
		return err
	}
	return ctx.JSON(articles)
}

// handlerRegister
// @Summary      Register user
// @Description  Register user
// @Tags         users
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email  formData  string  true  "email"
// @Param        name  formData  string  true  "name of user"
// @Param        password  formData  string  true  "password"
// @Success      200  {object}  string
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /register [post]
func (c *Controller) handlerRegister(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	name := ctx.FormValue("name")
	password := ctx.FormValue("password")
	passwordBytes := []byte(password)
	err := c.usecase.Register(email, name, passwordBytes)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

// handlerLogin
// @Summary      User authorisation
// @Description  User authorisation
// @Tags         users
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email  formData  string  true  "email"
// @Param        password  formData  string  true  "password"
// @Success      200  {object}  controller.Controller.handlerLogin.response
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /login [post]
func (c *Controller) handlerLogin(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	passwordBytes := []byte(password)
	user, sign, err := c.usecase.Login(email, passwordBytes)
	if err != nil {
		return err
	}
	type response struct { // тк свагер не имеет типа мапа для Success, преобразуем в структуру
		User  models.User `json:"user"` // для преобразования полей в json поля с большой буквы
		Token string      `json:"token"`
		Ok    bool        `json:"ok"`
	}
	return ctx.JSON(response{
		User:  user,
		Token: sign,
		Ok:    true,
	})
}

// handlerDeleteTranslations
// @Summary      Delete translations
// @Description  Delete translations
// @Tags         translations
// @Produce      json
// @Param        articleID  path   int  true  "articleID"
// @Param        languages   query  []string  true  "languages comma separated"
// @Success      200  {object}  string
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /translations/{id} [delete]
func (c *Controller) handlerDeleteTranslations(ctx *fiber.Ctx) error {
	articleID := ctx.Params("id")
	idInt, err := strconv.Atoi(articleID)
	if err != nil {
		return err
	}
	languages := ctx.Query("languages")
	languagesList := strings.Split(languages, ",")

	// вызов юзкейса
	err = c.usecase.DeleteTranslations(c.getUserID(ctx), idInt, languagesList)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

// handlerDeleteArticle
// @Summary      Delete article
// @Description  Delete article from dictionary (user must be article owner)
// @Tags         articles
// @Produce      json
// @Param        articleID  path   int  true  "articleID"
// @Success      200  {object}  string
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /article/{id} [delete]
func (c *Controller) handlerDeleteArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	err = c.usecase.DeleteArticle(c.getUserID(ctx), idInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

// handlerDeleteGroup
// @Summary      Delete group
// @Description  Delete group (user must be group owner)
// @Tags         groups
// @Produce      json
// @Param        groupID  path   int  true  "groupID"
// @Success      200  {object}  string
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /group/{id} [delete]
func (c *Controller) handlerDeleteGroup(ctx *fiber.Ctx) error {
	groupID := ctx.Params("id")
	idInt, err := strconv.Atoi(groupID)
	if err != nil {
		return err
	}
	err = c.usecase.DeleteGroup(c.getUserID(ctx), idInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

// handlerDeleteDictionary
// @Summary      Delete dictionary
// @Description  Delete dictionary (user must be dictionary owner)
// @Tags         dictionaries
// @Produce      json
// @Param        dictionaryID  path   int  true  "dictionaryID"
// @Success      200  {object}  string
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /dictionary/{id} [delete]
func (c *Controller) handlerDeleteDictionary(ctx *fiber.Ctx) error {
	dictionaryID := ctx.Params("id")
	idInt, err := strconv.Atoi(dictionaryID)
	if err != nil {
		return err
	}
	err = c.usecase.DeleteDictionary(c.getUserID(ctx), idInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

// handlerDeleteGroupArticles
// @Summary      Delete articles fron group
// @Description  Delete articles fron group (user must be group owner)
// @Tags         groups
// @Produce      json
// @Param        groupID  path  int  true  "groupID"
// @Param        articles  query  []int  true  "articles comma separated"
// @Success      200  {object}  string
// @Failure      400  {object}  fiber.Error
// @Failure      404  {object}  fiber.Error
// @Failure      500  {object}  fiber.Error
// @Router       /group/articles/{id} [delete]
func (c *Controller) handlerDeleteGroupArticles(ctx *fiber.Ctx) error {
	groupID := ctx.Params("id")
	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		return err
	}

	articles := ctx.Query("articles")
	articlesList := strings.Split(articles, ",")        // articlesList список строк
	articlesIDsInt := make([]int, 0, len(articlesList)) // длина 0, чтобы не ограничиваться и при апенде не наращивать слайс
	for _, id := range articlesList {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		articlesIDsInt = append(articlesIDsInt, idInt)
	}

	err = c.usecase.DeleteGroupArticles(c.getUserID(ctx), groupIDInt, articlesIDsInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}
