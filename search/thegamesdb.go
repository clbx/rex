package search

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/clbx/rex/platform"
)

var originalURL = "https://cdn.thegamesdb.net/images/original/"
var smallURL = "https://cdn.thegamesdb.net/images/small/"
var thumbURL = "https://cdn.thegamesdb.net/images/thumb/"
var boxartFrontPath = "boxart/front/"
var boxartBackPath = "boxart/back/"

type TGDBResponse struct {
	Code                      int                    `json:"code"`
	Status                    string                 `json:"status"`
	Data                      TGDBData               `json:"data"`
	Pages                     TGDBPages              `json:"pages"`
	Include                   map[string]interface{} `json:"include"`
	RemainingMonthlyAllowance int                    `json:"remaining_monthly_allowance"`
	ExtraAllowance            int                    `json:"extra_allowance"`
	AllowanceRefreshTimer     int                    `json:"allowance_refresh_timer"`
}

type TGDBData struct {
	Count int        `json:"count"`
	Games []TGDBGame `json:"games"`
}

type TGDBGame struct {
	ID          int    `json:"id"`
	GameTitle   string `json:"game_title"`
	ReleaseDate string `json:"release_date"`
	PlatformID  int    `json:"platform"`
	RegionID    int    `json:"region_id"`
	CountryID   int    `json:"country_id"`
	Players     int    `json:"players"`
	Overview    string `json:"overview"`
	Rating      string `json:"rating"`
	Developers  []int  `json:"developers"`
}

type TGDBPages struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
}

var PlatformMapping = map[string]int{
	"gcn": 2,
}

// TGDBseachGameByName searches for a game with a game object and returns an integer array of game IDs
func TGDBsearchGameByName(apikey string, game platform.Game) []int {
	req := "https://api.thegamesdb.net/v1/Games/ByGameName?apikey=" + url.QueryEscape(apikey) +
		"&name=" + url.QueryEscape(game.Name) +
		"&fields=" + url.QueryEscape("players, publishers, genres, overview, last_updated, rating, platform, coop, youtube, os, processor, ram, hdd, video, sound, alternates,overview,rating") +
		"&include=" + url.QueryEscape("boxart")
	//fmt.Println(req)

	tgdbResp := &TGDBResponse{}
	res, err := http.Get(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &tgdbResp)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("\n\n TGDB %+v\n", tgdbResp)

	//If there are no games, just return what was put in.
	if tgdbResp.Data.Count == 0 {
		log.Printf("GAME NOT FOUND")
		return []int{-1}
	}

	//If game is returned for multiple platforms, return the one for the current platform and game region.
	// TODO: Implement Region Filtering. Right now just search for US games (Presumably 1 in TGDB??)
	if tgdbResp.Data.Count > 1 {
		log.Printf("MULTIPLE GAMES RECEVIED - NOT SUPPORTED YET")
		return []int{-1}
	}

	return []int{tgdbResp.Data.Games[0].ID}
}

func TGDBsearchGameByID(apikey string, tgdbid int, game platform.Game) platform.Game {
	req := "https://api.thegamesdb.net/v1/Games/ByGameID?apikey=" + url.QueryEscape(apikey) +
		"&id=" + url.QueryEscape(strconv.Itoa(tgdbid)) +
		"&fields=" + url.QueryEscape("players, publishers, genres, overview, last_updated, rating, platform, coop, youtube, os, processor, ram, hdd, video, sound, alternates,overview,rating") +
		"&include=" + url.QueryEscape("boxart")

	tgdbResp := &TGDBResponse{}
	res, err := http.Get(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &tgdbResp)
	if err != nil {
		log.Fatal(err)
	}

	//Get image filename.
	originalSize := tgdbResp.Include["boxart"].(map[string]interface{})["base_url"].(map[string]interface{})["original"]
	front := tgdbResp.Include["boxart"].(map[string]interface{})["data"].(map[string]interface{})[strconv.Itoa(tgdbResp.Data.Games[0].ID)].([]interface{})[0].(map[string]interface{})["filename"]
	//back := tgdbResp.Include["boxart"].(map[string]interface{})["data"].(map[string]interface{})[strconv.Itoa(tgdbResp.Data.Games[0].ID)].([]interface{})[1].(map[string]interface{})["filename"]
	//fmt.Printf("\noriginal: %s\n  front: %s\n  back: %s\n", originalSize, front, back)

	//Get Images
	filepath := fmt.Sprintf("/cache/front-%s.jpg", strconv.Itoa(tgdbResp.Data.Games[0].ID))
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	resp, err := http.Get(fmt.Sprintf("%v%v", originalSize, front))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	game.Name = tgdbResp.Data.Games[0].GameTitle
	game.TGDBID = tgdbResp.Data.Games[0].ID
	game.ReleaseDate = tgdbResp.Data.Games[0].ReleaseDate
	game.Overview = tgdbResp.Data.Games[0].Overview
	game.BoxartFrontPath = fmt.Sprintf("/cache/front-%s.jpg", strconv.Itoa(tgdbResp.Data.Games[0].ID))

	return game

}
