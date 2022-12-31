package platform

//Master class, all games have one of these.

type Game struct {
	ID              string   `bson:"id,omitempty"`
	Name            string   `bson:"name,omitempty"`
	Platform        Platform `bson:"platform,omitempty"`
	TGDBID          int      `bson:"TGDBID,omitempty"`
	Path            string   `bson:"path,omitempty"`
	ReleaseDate     string   `bson:"releasedate,omitempty"`
	Overview        string   `bson:"overview,omitempty"`
	BoxartFrontPath string   `bson:"boxartback,omitempty"`
	BoxartBackPath  string   `bson:"boxartfront,omitempty"`
}

func CompareGames(game0 Game, game1 Game) bool {
	if game0.Name != game1.Name {
		return false
	}

	if game0.Platform.Name != game1.Platform.Name {
		return false
	}

	if game0.Path != game1.Path {
		return false
	}

	if game0.TGDBID != game1.TGDBID {
		return false
	}

	return true
}
