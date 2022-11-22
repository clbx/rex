package search

import "github.com/clbx/rex/platform"

type Search interface {
	searchGameByName(Search, string) platform.Game
}
