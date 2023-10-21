package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	webHandler "bwastartup/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)


	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	// validasi token manual
	// token ,err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMX0.atId9l0QTbiWTZcb0A8SYxciCOft5Rlb-FAqutSvNVY")
	// if err != nil {
	// 	fmt.Println("error")
	// }

	// if token.Valid {
	// 	fmt.Println("Valid")
	// } else {
	// 	fmt.Println("invalid")
	// }

	// fmt.Println(authService.GenerateToken(1001)) manually generated token

	// upload avatar manually
	// userService.SaveAvatar(1, "images/1-profile.png")

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
	
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHanlder(transactionService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionsHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())

	// middleware cookies sessions
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("bwastartup",cookieStore))

	// router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)
	
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.GET("/trasactions/notification", transactionHandler.GetNotification)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users", userWebHandler.Create)
	router.GET("/users/edit/:id",authAdminMiddleware(), userWebHandler.Edit)
	router.POST("/users/update/:id", authAdminMiddleware(), userWebHandler.Update)
	router.GET("users/avatar/:id", authAdminMiddleware(), userWebHandler.NewAvatar)
	router.POST("users/avatar/:id", authAdminMiddleware(), userWebHandler.CreateAvatar)

	router.GET("/campaigns", authAdminMiddleware(), campaignWebHandler.Index)
	router.GET("/campaigns/new", authAdminMiddleware(), campaignWebHandler.New)
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHandler.Create)
	router.GET("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.NewImage)
	router.POST("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.CreateImage)
	router.GET("/campaigns/edit/:id", authAdminMiddleware(), campaignWebHandler.Edit)
	router.POST("/campaigns/update/:id",authAdminMiddleware(), campaignWebHandler.Update)
	router.GET("/campaigns/show/:id", authAdminMiddleware(), campaignWebHandler.Show)

	router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.Create)
	router.GET("/logout", sessionWebHandler.Destroy)
	

	router.Run()

	// input dari user
	// handler, mapping input dari user -> struct input
	// service : melakukan mapping dari struct input ke struct User
	// repository
	// db

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context) {
		// ambil nilai header authorization: Bearer `token`,
		// dari header authorization, kita ambul nilai tokennya
		// kita validasi token
		// kita ambil user_id
		// ambil user dari db berdasarkan user_id lewat service
		// kita set context isinya user
		
			authHeader := c.GetHeader("Authorization")
		
			if !strings.Contains(authHeader, "Bearer") {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return 
			}
			// `Bearer  aoisjkdoiasjdiaosjdsiadosjdai`
			tokenString := ""
			arrayToken := strings.Split(authHeader, " ")
			if len(arrayToken) == 2 {
				tokenString = arrayToken[1]
			}
		
			token, err := authService.ValidateToken(tokenString)
			if err != nil {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

			claim, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

			userID := int(claim["user_id"].(float64))

			user, err := userService.GetUserByID(userID)
			if err != nil {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

			c.Set("currentUser", user)
		}
}

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// keywords login admin = userID
		userIDSession := session.Get("userID")

		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
}


func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
  
	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
	  panic(err.Error())
	}
  
	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
	  panic(err.Error())
	}
  
	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
	  layoutCopy := make([]string, len(layouts))
	  copy(layoutCopy, layouts)
	  files := append(layoutCopy, include)
	  r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
  }


