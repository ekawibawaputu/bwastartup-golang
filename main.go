package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)


	// test service
	// input := user.LoginInput{
	// 	Email: "contohemail@test.com",
	// 	Password: "password",
	// }

	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("Terjadi kesalahanan")
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(user.Email)
	// fmt.Println(user.Name)

	// userByEmail, err := userRepository.FindByEmail("contohemail@test.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// if (userByEmail.ID == 0) {
	// 	fmt.Println("User tidak ditemukan...")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	// service ke repository
	// userRepository := user.NewRepository(db)
	// userService := user.NewService(userRepository)

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Test simple user"
	// userInput.Email = "contohemail@test.com"
	// userInput.Occupation = "Designer"
	// userInput.Password = "password"

	// userService.RegisterUser(userInput)
	
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)

	router.Run()

	// input dari user
	// handler, mapping input dari user -> struct input
	// service : melakukan mapping dari struct input ke struct User
	// repository
	// db

}

