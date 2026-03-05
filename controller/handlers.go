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
	languages := ctx.Params("languages")
	languagesList := strings.Split(languages, ",")

	// вызов юзкейса
	data, err := c.usecase.GetArticle(idInt, languagesList)
	if err != nil {
		return err
	}
	// обратно кодируем в JSON id для пользователя
	return ctx.JSON(data)
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
