package controllers

import (
	"blog/config"
	"blog/database"
	"blog/models"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type ErrorUser struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateUser(user models.User) []*ErrorUser {
	var errors []*ErrorUser
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorUser
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type AuthStruct struct {
	Token string
}

type UserLogin struct {
	Username string
	Password string
}

func HashPassword(password string) (string, error) {
	fmt.Println(config.BCRYPT_COST)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomStr() string {
	rand.Seed(time.Now().Unix())
	str := config.RANDOM_STR
	shuff := []rune(str)
	rand.Shuffle(len(shuff), func(i, j int) {
		shuff[i], shuff[j] = shuff[j], shuff[i]
	})
	return string(shuff)
}

func Auth(c *fiber.Ctx) error {
	authBody := new(AuthStruct)
	if err := c.BodyParser(authBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var user_db models.User
	database.Database.Db.Where("token = ?", authBody.Token).First(&user_db)
	if user_db.Email != "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
	}
}

func Login(c *fiber.Ctx) error {
	user := new(UserLogin)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user_db models.User
	database.Database.Db.Where("email = ?", user.Username).First(&user_db)
	if user_db.Email == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true, "message": "Email or Password not match"})
	}
	match := CheckPasswordHash(user.Password, user_db.Password)
	if !match {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true, "message": "Email or Password not match"})
	}

	acak := RandomStr()
	user_db.Token = acak
	database.Database.Db.Save(&user_db)
	// Create the Claims
	claims := jwt.MapClaims{
		"name":  user_db.Name,
		"token": acak,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": user_db.ID, "fullName": user_db.Name, "avatar": "https://ssl.gstatic.com/ui/v1/icons/mail/rfr/logo_gmail_lockup_default_1x_r5.png", "token": t})
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateUser(*user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	hash, _ := HashPassword(user.Password)
	data := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hash),
		Active:   user.Active,
	}
	if err := database.Database.Db.Create(&data).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"email": user.Email})
}

func ShowUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please enter id only Integer"})
	}
	var user models.User
	database.Database.Db.First(&user, id)
	if user.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"active": user.Active,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please enter id only Integer"})
	}
	var user models.User
	database.Database.Db.First(&user, id)
	if user.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	updateUser := new(models.User)
	if err := c.BodyParser(updateUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	errors := ValidateUser(*updateUser)
	if errors != nil {
		return c.Status(400).JSON(errors)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 4)

	fmt.Println(err)

	user.Name = updateUser.Name
	user.Email = updateUser.Email
	user.Password = string(bytes)
	user.Active = updateUser.Active
	database.Database.Db.Save(&user)

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please enter id only Integer"})
	}
	var user models.User
	database.Database.Db.First(&user, id)
	if user.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	database.Database.Db.Delete(&user, id)
	return c.JSON(fiber.Map{"message": "Succesfully deleted"})
}
