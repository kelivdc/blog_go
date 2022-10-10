package controllers

import (
	"blog/database"
	"blog/models"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorPost struct {
	FailedField string
	Tag         string
	Value       string
}

type PostFields struct {
}

var validate = validator.New()

func ValidatePost(post models.Post) []*ErrorPost {
	var errors []*ErrorPost
	err := validate.Struct(post)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorPost
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func AdminPostIndex(c *fiber.Ctx) error {
	posts := []models.Post{}
	database.Database.Db.Joins("Category").Find(&posts)

	// if result.Error != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(result.Error)
	// }

	// fmt.

	return c.JSON(&posts)
}

func AdminPostIndexOld(c *fiber.Ctx) error {
	posts := []models.Post{}
	var count int64

	sql := "SELECT * from posts WHERE deleted_at is NULL"
	if s := c.Query("s"); s != "" {
		sql = fmt.Sprintf("%s AND LOWER(title) LIKE '%%%s%%'", sql, s)
	}

	if sort := c.Query("_sort"); sort != "" {
		sql = fmt.Sprintf("%s ORDER by %s", sql, sort)
	}

	if order := c.Query("_order"); order != "" {
		sql = fmt.Sprintf("%s %s", sql, order)
	}

	database.Database.Db.Raw(sql).Scan(&posts).Count(&count)
	total := strconv.Itoa(int(count))

	if start := c.Query("_start"); start != "" {
		var end, _ = strconv.Atoi(c.Query("_end"))
		var pstart, _ = strconv.Atoi(c.Query("_end"))
		sql = fmt.Sprintf("%s LIMIT %s %s", sql, start, strconv.Itoa(end-pstart))
	}

	database.Database.Db.Raw(sql).Scan(&posts)
	c.Set("X-Total-Count", total)
	c.Set("Access-Control-Expose-Headers", "X-Total-Count")

	return c.Status(fiber.StatusOK).JSON(&posts)
}

func AdminPostCreate(c *fiber.Ctx) error {
	var post models.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	data := models.Post{
		Title:       post.Title,
		CategoryID:  post.CategoryID,
		Body:        post.Body,
		Description: post.Description,
		ShortDesc:   post.ShortDesc,
		Keyword:     post.Keyword,
		Slug:        post.Slug,
		Image:       post.Image,
		ImageNote:   post.ImageNote,
		Publish:     post.Publish,
	}

	result := database.Database.Db.Joins("Category").Create(&data)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(result.Error)
	}
	database.Database.Db.Joins("Category").First(&post, data.ID)
	return c.Status(fiber.StatusOK).JSON(&post)
}

func AdminPostShow(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please select Id"})
	}
	var post models.Post
	database.Database.Db.First(&post, id)
	if post.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(post)
}
