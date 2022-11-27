package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/clbx/rex/db"
	_ "github.com/clbx/rex/docs"
	"github.com/clbx/rex/platform"
	"github.com/clbx/rex/search"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var platforms []platform.Platform

var ctx = context.Background()
var gamedb *db.DB
var apikey = "f8e8aae3d8fdcd3d4d29c1e2a65d899410001610d398afd6df7964fa1b527e1a"

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// @title Rex
// @description Self Hostable Game Library
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	r := gin.Default()
	r.Use(CORSMiddleware())

	startup()
	findGames()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/v1/ping", ping)
	r.GET("/v1/games", getGames)
	r.GET("/v1/platforms", getPlatforms)

	r.Run()
}

func startup() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.UnmarshalKey("platforms", &platforms)
	if err != nil {
		log.Fatal(err)
	}

	//platforms := []platform.Platform(viper.Get("platforms"))
	//fmt.Printf("%+v\n", platforms)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	gamedb, err = db.InitMongoDB(ctx, "mongodb://mongodb:27017")

}

func findGames() {
	for _, p := range platforms {
		for _, dir := range p.Directories {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {

				log.Printf("[%s] Found file %s", p.Name, dir+"/"+file.Name())

				game, err := platform.IdenfityGameByPlatform(p, dir, file.Name())
				game = search.TGDBsearchGameByName(apikey, game)
				if err != nil {
					log.Fatal(err)
				}
				err = db.AddGame(gamedb, ctx, game)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// ping godoc
// @Summary Ping!
// @Description Pong!
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /v1/ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// getGames godoc
// @Summary Get All Games
// @Description Get a list of all of the games that Rex can find
// @Produce json
// @Success 200
// @Router /v1/games [get]
func getGames(c *gin.Context) {
	allGames, err := db.GetAllGames(gamedb, ctx)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allGames)
}

// getPlatforms godoc
// @Summary Get platforms with games
// @Description Returns a list of platforms with games in the library
// @Produce json
// @Sucess 200
// @Router /v1/platforms [get]
func getPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, platforms)
}
