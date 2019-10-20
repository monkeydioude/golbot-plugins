package golmods

import (
	"bitbucket.org/drannoc/golbot"
	// this is auto generated, import golbot modules from ./pkg/
	"github.com/monkeydioude/golmods/pkg/giphy"
	"github.com/monkeydioude/golmods/pkg/reddithot"
)

func healthCheck() string {
	return "OK"
}

func GetCommands(cachePath string) []golbot.Command {
	return []golbot.Command{
		giphy.AddCommand(cachePath + "giphy/"),
		reddithot.AddCommand(cachePath + "reddithot/"),
	}
}
