package golmods

import (
    "bitbucket.org/drannoc/golbot"
	"github.com/bwmarrin/discordgo"
    // this is auto generated, import golbot modules from ./pkg/
    #MODS#
)

func healthCheck() string {
    return "OK"
}

func GetCommands(cachePath string, session *discordgo.Session) []golbot.Command {
    return []golbot.Command{
        #ADD_COMMAND#
    }
}
