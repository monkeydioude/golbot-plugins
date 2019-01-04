package golmods

import (
    "bitbucket.org/drannoc/golbot"
    // this is auto generated, import golbot modules from ./pkg/
    #MODS#
)

func healthCheck() string {
    return "OK"
}

func GetCommands(g *golbot.Golbot, cachePath string) []golbot.Command {
    return []golbot.Command{
        #ADD_COMMAND#
    }
}
