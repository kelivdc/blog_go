package controllers

import (
	"blog/database"
	"blog/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AdminCategoryIndex(c *fiber.Ctx) error {
	categories := []models.Category{}
	var count int64

	sql := "SELECT * FROM categories WHERE deleted_at is NULL"
	if s := c.Query("s"); s != "" {
		sql = fmt.Sprintf("%s AND LOWER(name) LIKE '%%%s%%'", sql, s)
	}
	if sort := c.Query("_sort"); sort != "" {
		sql = fmt.Sprintf("%s ORDER by %s", sql, sort)
	}

	if order := c.Query("_order"); order != "" {
		sql = fmt.Sprintf("%s %s", sql, order)
	}

	database.Database.Db.Raw(sql).Scan(&categories).Count(&count) //Save count before OFFSET and LIMIT

	total := strconv.Itoa(int(count))

	if start := c.Query("_start"); start != "" {
		var end, _ = strconv.Atoi(c.Query("_end"))
		var pstart, _ = strconv.Atoi(start)
		sql = fmt.Sprintf("%s LIMIT %s, %s", sql, start, strconv.Itoa(end-pstart))
	}

	database.Database.Db.Raw(sql).Scan(&categories)

	c.Set("X-Total-Count", total)
	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	return c.Status(fiber.StatusOK).JSON(&categories)
}

func AdminCategoryCreate(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error")
	}

	// errors := controllers.ValidateStruct(*category)
	// if errors != nil {
	// 	return c.Status(400).JSON(errors)
	// }

	database.Database.Db.Create(&category)
	return c.Status(fiber.StatusOK).JSON(category)
}

func AdminCategoryShow(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please enter id only Integer"})
	}
	var category models.Category
	database.Database.Db.First(&category, id)
	if category.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(category)
}

func AdminCategoryUpdate(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please enter id only Integer"})
	}
	var category models.Category
	database.Database.Db.First(&category, id)
	if category.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	updateCategory := new(models.Category)
	if err := c.BodyParser(updateCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error")
	}

	// errors := controllers.ValidateStruct(*updateCategory)
	// if errors != nil {
	// 	return c.Status(400).JSON(errors)
	// }

	category.Name = updateCategory.Name
	category.Slug = updateCategory.Slug
	category.Publish = updateCategory.Publish
	database.Database.Db.Save(&category)

	return c.JSON(category)
}

func AdminCategoryDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please enter id only Integer"})
	}
	var category models.Category
	database.Database.Db.First(&category, id)
	if category.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	database.Database.Db.Delete(&category, id)
	return c.JSON(fiber.Map{"message": "Succesfully deleted"})
}
