package routes

import (
	"blog/controllers"
	admin_mod "blog/controllers/admin"
	"blog/database"
	"blog/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	api := app.Group("/api")

	v1 := api.Group("/v1")

	// v1.Get("/categories", controllers.Index)
	// v1.Get("/categories/:id", controllers.Show)
	// v1.Post("/categories", controllers.Create)
	// v1.Put("/categories/:id", controllers.Update)
	// v1.Delete("/categories/:id", controllers.Delete)

	admin := v1.Group("/admin")

	// Login
	admin.Post("/login", controllers.Login)
	admin.Post("/auth", controllers.Auth)
	admin.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			token := claims["token"].(string)
			var user_db models.User
			database.Database.Db.Where("token = ? and active = 1", token).First(&user_db)
			if user_db.Email == "" {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			return c.Next()
		},
	}))
	// Category
	admin.Get("/categories", admin_mod.AdminCategoryIndex)
	admin.Get("/categories/:id", admin_mod.AdminCategoryShow)
	admin.Post("/categories", admin_mod.AdminCategoryCreate)
	admin.Put("/categories/:id", admin_mod.AdminCategoryUpdate)
	admin.Delete("/categories/:id", admin_mod.AdminCategoryDelete)

	// Post
	admin.Get("/posts", admin_mod.AdminPostIndex)
	admin.Post("/posts", admin_mod.AdminPostCreate)
	admin.Get("/posts/:id", admin_mod.AdminPostShow)

	// User
	admin.Post("/users", controllers.CreateUser)
	admin.Get("/users/:id", controllers.ShowUser)
	admin.Put("/users/:id", controllers.UpdateUser)
	admin.Delete("/users/:id", controllers.DeleteUser)

}
