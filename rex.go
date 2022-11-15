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

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ctx = context.Background()
var gamedb *db.DB
var platforms []string

// @title Rex
// @description Self Hostable Game Library
// @host localhost:8080
// @BasePath /
// @schemes http

func main() {

	r := gin.Default()

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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	gamedb, err = db.InitMongoDB(ctx, "mongodb://mongodb:27017")

}

func findGames() {
	// //Nintendo Wii
	// wiiDirs := viper.GetStringSlice("gameDirectories.wii")
	// if wiiDirs != nil {
	// 	platforms = append(platforms, "wii")
	// }

	//Nintendo GameCube
	gcnDirs := viper.GetStringSlice("gameDirectories.gcn")
	if gcnDirs != nil {
		platforms = append(platforms, "gcn")
	}
	log.Printf("Finding Gamecube Games")
	for _, dir := range gcnDirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			log.Printf("Found file %s", dir+"/"+file.Name())
			game := platform.WrapGamecube(platform.IdentifyGamecube(dir + "/" + file.Name()))
			log.Printf("Identified %s as %s", game.Path, game.Name)
			err = db.AddGame(gamedb, ctx, game)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	// //Nintendo 64
	// n64Dirs := viper.GetStringSlice("gameDirectories.n64")
	// if n64Dirs != nil {
	// 	platforms = append(platforms, "n64")
	// }
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