package controller

import (
	"github.com/Nicole8493/GoLingo/models"
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) handlerCreateArticle(ctx *fiber.Ctx) error {
	// достаем данные для передачи в юзкейс
	data := new([]models.Translation)
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
func (c *Controller) handlerUpdateTranslations(ctx *fiber.Ctx) error {}
func (c *Controller) handlerGetFullArticle(ctx *fiber.Ctx) error     {}
func (c *Controller) handlerGetArticle(ctx *fiber.Ctx) error         {}
func (c *Controller) handlerDeleteTranslations(ctx *fiber.Ctx) error {}
func (c *Controller) handlerDeleteArticle(ctx *fiber.Ctx) error      {}
