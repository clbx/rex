package platform

//Master class, all games have one of these.

type Game struct {
	Name     string `bson:"name,omitempty"`
	Platform string `bson:"platform,omitempty"`
	Path     string `bson:"path,omitempty"`
}

func CompareGames(game0 Game, game1 Game) bool {
	if game0.Name != game1.Name {
		return false
	}

	if game0.Platform != game1.Platform {
		return false
	}

	if game0.Path != game1.Path {
		return false
	}

	return true
}
