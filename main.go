package main

import (
	"log"
	"tender_management/config"
	"tender_management/controllers"
	"tender_management/pkg/db"
	"tender_management/pkg/middleware"
	"tender_management/pkg/redise"

	_ "tender_management/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Tender Management REST API
// @version 1.0
// @description Tender Management Golang REST API
// @contact.name Muxtorxon Gofurov
// @contact.url https://github.com/muhtorkhon
// @contact.email muhtorhongofurov@gmail.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token
// @type apiKey
func main() {

	enforcer, err := casbin.NewEnforcer("auth_model.conf", "policy.csv")
	if err != nil {
		log.Println("Error configuring casbin", err)
	}

	r := gin.Default()

	cfg := config.LoadConfig()
	gin.SetMode(gin.ReleaseMode)

	conn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	rd := db.ConnectRedis()
	redisDb := redise.NewRedis(rd)
	if err := redisDb.Ping(); err != nil {
		log.Fatalf("Redis ulanish xatosi: %v", err)
	}
	
	authSt := controllers.NewAuthController(conn, redisDb)
	tenderSt := controllers.NewTenderController(conn)
	offerSt := controllers.NewOfferController(conn)
	notifSt := controllers.NewNotifController(conn)

	public := r.Group("")

	public.POST("/auth/register", authSt.CreateUser)
	public.POST("/auth/verify", authSt.VerifyCode)
	public.POST("/auth/login", authSt.LoginUser)

	r.Use(middleware.AutoMiddleware(enforcer))

	r.POST("/tenders", tenderSt.CreateTender)
	r.GET("/tenders", tenderSt.GetAllTenders)
	r.GET("/tenders/:client_id", tenderSt.GetTenders)
	r.PUT("/tenders/:id", tenderSt.UpdateTender)
	r.DELETE("/tenders/:id", tenderSt.DeleteTender)
	r.PATCH("/tenders/restore/:id", tenderSt.RestoreTender)

	r.POST("/offers", offerSt.CreateOffer)
	r.GET("/offers", offerSt.GetAllOffers)
	r.GET("/offers/sorted", offerSt.GetFilterSort)
	r.GET("/offers/filter", offerSt.GetMaxMinFilter)
	r.GET("/offers/:contractor_id", offerSt.GetOffer)	
	r.PUT("/offers/:id", offerSt.UpdateOffer)
	r.DELETE("/offers/:id", offerSt.DeleteOffer)
	r.PATCH("/offers/restore/:id", offerSt.RestoreOffer)

	r.POST("/notifs", notifSt.CreateNotif)
	r.GET("/notifs/:user_id/:relation_id", notifSt.GetNotifsUser)

	public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server on port: 8080")
	}
}

// 	client := r.Group("/client")
// 	client.Use(middleware.AutoMiddleware(enforcer))
// 	{
// 		client.GET("/dashboard", func(c *gin.Context) {
// 			c.JSON(200, "Welcome to the client dashboard")
// 		})

// 		client.POST("/tenders", tenderSt.CreateTender)
// 		client.GET("/tenders", tenderSt.GetAllTenders)
// 		client.GET("/tenders/:client_id", tenderSt.GetTenders)
// 		client.PUT("/tenders/:id", tenderSt.UpdateTender)
// 		client.DELETE("/tenders/:id", tenderSt.DeleteTender)
// 		client.PATCH("/tenders/restore/:id", tenderSt.RestoreTender)

// 		client.GET("/offers", offerSt.GetAllOffers)
// 		client.GET("/offers/sorted", offerSt.GetFilterSort)
// 		client.GET("/offers/filter", offerSt.GetMaxMinFilter)

// 		client.POST("/notifs", notifSt.CreateNotif)
// 		client.GET("/notifs/:user_id/:relation_id", notifSt.GetNotifClient)
// 	}

// 	contractor := r.Group("/contractor")
// 	contractor.Use(middleware.AutoMiddleware(enforcer))
// 	{
// 		contractor.GET("/profile", func(c *gin.Context) {
// 			email, _ := c.Get("email")
// 			userEmail := fmt.Sprintf("Welcome too your profile %s", email)
// 			c.JSON(200, userEmail)
// 		})

// 		contractor.POST("/offers", offerSt.CreateOffer)
// 		contractor.GET("/offers", offerSt.GetAllOffers)
// 		contractor.GET("/offers/:contractor_id", offerSt.GetOffer)
// 		contractor.PUT("/offers/:id", offerSt.UpdateOffer)
// 		contractor.DELETE("/offers/:id", offerSt.DeleteOffer)
// 		contractor.PATCH("/offers/restore/:id", offerSt.RestoreOffer)

// 		contractor.GET("/tenders", tenderSt.GetAllTenders)

// 		contractor.POST("/notifs", notifSt.CreateNotif)
// 		contractor.GET("/notifs/:user_id/:relation_id", notifSt.GetNotifContractor)
// 	}

// 	r.GET("/Tender-management/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("Failed to run server on port: 8080")
// 	}
// }
