package middlewares

import (
	"blog/database"
	"blog/models"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CheckMatch(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckLogin(basic_user, basic_pass string) bool {
	var user_db models.User
	database.Database.Db.Where("email = ?", basic_user).First(&user_db)

	if basic_user == os.Getenv("ADMIN_EMAIL") && basic_pass == os.Getenv("ADMIN_PASSWORD") { //Super Admin
		return true
	}
	if user_db.Email == "" {
		return false
	}

	match := CheckMatch(basic_pass, user_db.Password)
	return match
}
