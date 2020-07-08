package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uwaifo/lmsvideoapi/infrastructure/auth"
	"github.com/uwaifo/lmsvideoapi/infrastructure/persistence"
	"github.com/uwaifo/lmsvideoapi/interfaces/middleware"
	"github.com/uwaifo/lmsvideoapi/interfaces/upload"
	"github.com/uwaifo/lmsvideoapi/routes"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no enviroment variable found")
	}

}

func main() {

	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Latter add some redis stuff
	//redis details
	//redis://h:p933aa7a328ef19d91da349d6443a5e3dd1dea814802e661b796325771e2169e6@ec2-34-204-117-137.compute-1.amazonaws.com:27129
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbDriver, user, password, port, host, dbname)
	//postgres://nvjuszfxcnmuye:3d836cb2cc5a14103b5addf9048e4a3a21e2ab0e313664a955fc4228b81a2837@ec2-107-22-7-9.compute-1.amazonaws.com:5432/dblgr2metnehkn
	if err != nil {
		panic(err)
	}
	defer services.Close()
	// Run migration if none
	services.AutoMigrate()

	//redis connection
	redisService, err := auth.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	fu := fileupload.NewFileUpload()

	users := routes.NewUsers(services.User, redisService.Auth, tk)
	videos := routes.NewVideo(services.Video, services.User, fu, redisService.Auth, tk)
	//snippet := routes.NewSnippet(services)
	authenticate := routes.NewAuthenticate(services.User, redisService.Auth, tk)

	r := gin.Default()

	//Set max file upload size to 100 Mib
	r.MaxMultipartMemory = 100 << 20
	r.Use(middleware.CORSMiddleware()) //For CORS

	v1 := r.Group("/api/v1")

	{
		// User Routes
		v1.POST("/signup", users.SaveUser)
		v1.GET("/users", users.GetUsers)
		v1.GET("/users/:user_id", users.GetUser)

		// Videoo Routes
		v1.POST("/video", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), videos.SaveVideo)
		v1.PUT("/video/:video_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), videos.UpdateVideo)
		v1.GET("/video/:video_id", videos.GetVideoAndCreator)
		v1.DELETE("/video/:video_id", middleware.AuthMiddleware(), videos.DeleteVideo)
		v1.GET("/video", videos.GetAllVideo)

		//Video Segments
		v1.POST("/snippet", routes.SaveSnippet)

		// Authentication routes
		v1.POST("/login", authenticate.Login)
		v1.POST("/logout", authenticate.Logout)
		v1.POST("/refresh", authenticate.Refresh)

	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8888"

	}
	log.Fatal(r.Run(":" + appPort))

}
