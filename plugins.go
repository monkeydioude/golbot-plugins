package golmods

import (
    "bitbucket.org/drannoc/golbot"
    // this is auto generated, import golbot modules from ./pkg/
    "github.com/monkeydioude/golmods/pkg/alert"
	"github.com/monkeydioude/golmods/pkg/giphy"
	
)

func healthCheck() string {
    return "OK"
}


func GetCommands(g *golbot.Golbot) []golbot.Command {
    return []golbot.Command{
        alert.AddCommand(g),
		giphy.AddCommand(g),
		
    }
}
