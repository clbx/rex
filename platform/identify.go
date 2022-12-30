package platform

import (
	"log"

	"github.com/google/uuid"
)

// Map of supported platforms to identification function.
var identMap = map[string]func(Platform, string) Game{
	"gcn": IdentifyGamecube,
}

//TODO: Add error handling
func IdenfityGameByPlatform(platform Platform, dir string, filename string) (Game, error) {
	loc := dir + "/" + filename

	if val, ok := identMap[platform.Platform]; ok {
		game := val(platform, loc)
		log.Printf("[%s][%s] Game idenfitied as %s", platform.Name, game.Name, game.Name)
		return game, nil
	}
	//Not in supported handlers, so generic handler is used.
	//This is temporary
	game := Game{
		ID:       uuid.New().String(),
		Name:     filename,
		Platform: platform,
		Path:     loc,
	}
	log.Printf("[%s][%s] Platform not supported, falling back to file name.", platform.Name, game.Name)
	return game, nil
}
