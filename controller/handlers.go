package controller

import (
	"github.com/Nicole8493/GoLingo/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func (c *Controller) handlerCreateArticle(ctx *fiber.Ctx) error {
	// достаем данные для передачи в юзкейс
	data := new(models.Article)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	// вызов юзкейса
	id, err := c.usecase.CreateArticle(*data)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON(id)
}

func (c *Controller) handlerCreateDictionary(ctx *fiber.Ctx) error {
	data := new(models.Dictionary)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	id, err := c.usecase.CreateDictionary(*data)
	if err != nil {
		return err
	}
	return ctx.JSON(id)
}

func (c *Controller) handlerCreateGroup(ctx *fiber.Ctx) error {
	data := new(models.Group)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	id, err := c.usecase.CreateGroup(*data)
	if err != nil {
		return err
	}
	return ctx.JSON(id)
}

func (c *Controller) handlerUpdateTranslations(ctx *fiber.Ctx) error {
	data := new([]models.Translation)
	if err := ctx.BodyParser(data); err != nil {
		return err
	}

	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	// вызов юзкейса
	err = c.usecase.UpdateTranslations(idInt, *data)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON("ok")
}

func (c *Controller) handlerAddGroupArticles(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	groupIDInt, err := strconv.Atoi(id)
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

	err = c.usecase.AddGroupArticles(groupIDInt, articlesIDsInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

func (c *Controller) handlerGetFullArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
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

func (c *Controller) handlerLogin(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	passwordBytes := []byte(password)
	user, sign, err := c.usecase.Login(email, passwordBytes)
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]interface{}{"user": user, "token": sign, "ok": true})
}

func (c *Controller) handlerDeleteTranslations(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	languages := ctx.Params("languages")
	languagesList := strings.Split(languages, ",")

	// вызов юзкейса
	err = c.usecase.DeleteTranslations(idInt, languagesList)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}
func (c *Controller) handlerDeleteArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	err = c.usecase.DeleteArticle(idInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

func (c *Controller) handlerDeleteGroup(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	err = c.usecase.DeleteGroup(idInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

func (c *Controller) handlerDeleteDictionary(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	err = c.usecase.DeleteDictionary(idInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}

func (c *Controller) handlerDeleteGroupArticles(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	groupIDInt, err := strconv.Atoi(id)
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

	err = c.usecase.DeleteGroupArticles(groupIDInt, articlesIDsInt)
	if err != nil {
		return err
	}
	return ctx.JSON("ok")
}
