package main

import (
	"crowdfunding-minpro-alterra/database"
	"crowdfunding-minpro-alterra/handler"
	"crowdfunding-minpro-alterra/modules/campaign"
	"crowdfunding-minpro-alterra/modules/chat"
	"crowdfunding-minpro-alterra/modules/donation"
	"crowdfunding-minpro-alterra/modules/payment"
	"crowdfunding-minpro-alterra/modules/user"
	"crowdfunding-minpro-alterra/utils/auth"
	"crowdfunding-minpro-alterra/utils/helper"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	donationRepository := donation.NewRepository(db)
	chatRepository := chat.NewChatRepository()

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	donationService := donation.NewService(donationRepository, campaignRepository, paymentService)
	chatUC := chat.NewChatUseCase(chatRepository)

	cloudinary, err := initCloudinary()
	if err != nil {
		fmt.Println("Failed to initialize Cloudinary:", err)
		return
	}

	userHandler := handler.NewUserHandler(userService, authService, cloudinary)
	campaignHandler := handler.NewCampaignHandler(campaignService, cloudinary)
	donationHandler := handler.NewDonationHandler(donationService)
	chatHandler := handler.NewChatHandler(chatUC)

	router := gin.Default()
	router.Use(cors.Default())
	
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	api.POST("/chatbot", chatHandler.HandleChat)

	api.GET("/admin/users", authMiddleware(authService, userService), userHandler.GetAllUsers)
	api.DELETE("/admin/users/:id", authMiddleware(authService, userService), userHandler.DeleteUser)
	api.GET("/admin/campaigns", campaignHandler.GetCampaigns)
	api.POST("/admin/sessions", userHandler.Login)

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

	api.GET("/campaigns/:id/donations", authMiddleware(authService, userService), donationHandler.GetCampaignDonations)
	api.GET("/donations", authMiddleware(authService, userService), donationHandler.GetUserDonations)
	api.POST("/donations", authMiddleware(authService, userService), donationHandler.CreateDonation)
	api.POST("/donations/notification", donationHandler.GetNotification)

	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	router.Run(":8080")
}

func initCloudinary() (*cloudinary.Cloudinary, error) {
	cloudinaryURL := "cloudinary://753269553777558:-ErnNG6kEVtqQrik-OWxFX4y7JI@dikzx7xyv"
	cloudinary, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return cloudinary, nil
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

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
