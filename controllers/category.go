package controllers

import (
	"blog/database"
	"blog/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func CreateResponseCategory(data models.Category) models.Category {
	return models.Category{ID: data.ID, Name: data.Name, Slug: data.Slug, Publish: data.Publish, CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt}
}

func ValidateStruct(category models.Category) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(category)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func Index(c *fiber.Ctx) error {
	categories := []models.Category{}
	keyword := c.Query("s")
	if keyword != "" {
		database.Database.Db.Where("LOWER(name) LIKE ?", "%"+keyword+"%").Find(&categories)
	} else {
		database.Database.Db.Find(&categories)
	}
	responseCategories := []models.Category{}

	for _, category := range categories {
		responseCategory := CreateResponseCategory(category)
		responseCategories = append(responseCategories, responseCategory)
	}
	return c.Status(fiber.StatusOK).JSON(responseCategories)
}

func Show(c *fiber.Ctx) error {
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

func Create(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error")
	}

	errors := ValidateStruct(*category)
	if errors != nil {
		return c.Status(400).JSON(errors)
	}

	database.Database.Db.Create(&category)
	return c.Status(fiber.StatusOK).JSON(category)
}

func Update(c *fiber.Ctx) error {
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

	errors := ValidateStruct(*updateCategory)
	if errors != nil {
		return c.Status(400).JSON(errors)
	}

	category.Name = updateCategory.Name
	category.Slug = updateCategory.Slug
	category.Publish = updateCategory.Publish
	database.Database.Db.Save(&category)

	return c.JSON(category)
}

func Delete(c *fiber.Ctx) error {
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
