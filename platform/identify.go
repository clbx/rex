package platform

import (
	"log"
	"strings"
)

// Map of supported platforms to identification function.
var identMap = map[string]func(string) Game{
	"gcn": IdentifyGamecube,
}

//TODO: Add error handling
func IdenfityGameByPlatform(platform Platform, loc string) (Game, error) {
	split := strings.Split(loc, "\\")
	filename := split[len(split)-1]

	if val, ok := identMap[platform.Platform]; ok {
		game := val(loc)
		log.Printf("%s game idenfitied as %s", platform.Name, game.Name)
		return game, nil
	}
	//Not in supported handlers, so generic handler is used.
	//This is temporary
	game := Game{
		Name:     filename,
		Platform: platform.Platform,
		Path:     loc,
	}
	log.Printf("%s platform not supported, falling back to file name.", platform.Name)
	return game, nil
}
