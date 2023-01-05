package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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

// @title       Rex
// @description Self Hostable Game Library
// @host        localhost:8080
// @BasePath    /
// @schemes     http
func main() {

	r := gin.Default()
	r.Use(CORSMiddleware())

	startup()
	findGames()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/v1/ping", ping)
	r.GET("/v1/games", getGames)
	r.GET("/v1/games/byId", getGamesById)
	r.GET("/v1/platforms", getPlatforms)
	r.GET("/cache/:filename", func(c *gin.Context) {
		fileName := c.Param("filename")
		targetPath := filepath.Join("/cache/", fileName)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/octet-stream")
		c.File(targetPath)
	})
	r.POST("/v1/games/setGameById", setGameById)

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
	// Iterate through every platform
	for _, p := range platforms {
		// Iterate through every directory that is a part of every platform
		for _, dir := range p.Directories {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Fatal(err)
			}
			// Iterate through every file in directory
			for _, file := range files {

				log.Printf("[%s] Found file %s", p.Name, dir+"/"+file.Name())

				// Attempt to Identify the Game by Platform Specific Information
				game, err := platform.IdenfityGameByPlatform(p, dir, file.Name())

				// Search the game in TGDB
				//TODO: Filter for multiple IDs
				ids := search.TGDBsearchGameByName(apikey, game)
				fmt.Printf("ID: %d", ids[0])
				if ids[0] != -1 {
					game = search.TGDBsearchGameByID(apikey, ids[0], game)
				}
				//spew.Dump(game)
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
// @Summary     Ping!
// @Description Pong!
// @Accept      */*
// @Produce     json
// @Success     200 {object} string
// @Router      /v1/ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// getGames godoc
// @Summary     Get All Games
// @Description Get a list of all of the games that Rex can find
// @Produce     json
// @Success     200
// @Router      /v1/games [get]
func getGames(c *gin.Context) {
	allGames, err := db.GetAllGames(gamedb, ctx)
	if err != nil {
		fmt.Printf("HERE1.5\n")
		log.Fatal(err)
	}
	fmt.Printf("HERE2\n")
	c.JSON(http.StatusOK, allGames)
}

// getGames godoc
// @Summary     Get game by UUID
// @Description Get a game by UUID
// @Produce     json
// @Success     200
// @Router      /v1/games/byId [get]
// @Param		id path string true "ID of the game to search for"
func getGamesById(c *gin.Context) {
	gameId := c.Query("id")
	if gameId == "" {
		c.JSON(http.StatusBadRequest, "No UUID provided with game")
		return
	}
	game, err := db.GetGameByID(gamedb, ctx, gameId)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: check if more than one, there never should be, need to figure out mongo better
	c.JSON(http.StatusOK, game[0])
}

// getGames godoc
// @Summary     Get game by UUID
// @Description Get a game by UUID
// @Produce     json
// @Success     200
// @Router      /v1/games/setGameById [post]
func setGameById(c *gin.Context) {
	gameToSet := c.Query("id")
	tgdbIDstr := c.Query("tgdbid")
	if gameToSet == "" || tgdbIDstr == "" {
		c.JSON(http.StatusBadRequest, "No UUID or IGDB ID Provided")
		return
	}

	tgdbID, _ := strconv.Atoi(tgdbIDstr)

	game, err := db.GetGameByID(gamedb, ctx, gameToSet)
	if err != nil {
		log.Fatal(err)
	}
	game[0] = search.TGDBsearchGameByID(apikey, tgdbID, game[0])
	db.AddGame(gamedb, ctx, game[0])
	c.JSON(200, gin.H{"status": "ID Assigned"})
}

// getPlatforms godoc
// @Summary     Get platforms with games
// @Description Returns a list of platforms with games in the library
// @Produce     json
// @Sucess      200
// @Router      /v1/platforms [get]
func getPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, platforms)
}
