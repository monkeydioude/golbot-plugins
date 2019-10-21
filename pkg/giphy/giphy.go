package giphy

import (
	"encoding/json"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/monkeydioude/gophy"
	"github.com/monkeydioude/gophy/pkg/entity"
	"github.com/monkeydioude/gophy/pkg/request"
	"github.com/monkeydioude/tools"
)

const (
	apiKey      = "NaYGr3kaEMDdR8eTjsbZOUrAxUPhE7Iq"
	searchLimit = 3
	notFoundMsg = ":middle_finger::joy:"
)

type giphy struct {
	session *discordgo.Session
	gophy   *gophy.Gophy
}

func AddCommand(cachePath string, session *discordgo.Session) *giphy {
	return &giphy{
		session: session,
		gophy:   gophy.NewGophy(apiKey),
	}
}

// GetRegex implements golbot.Command interface
func (g *giphy) GetRegex() string {
	return "/gif (.+)"
}

func (g *giphy) GetName() string {
	return "gif"
}

// Do implements golbot.Command interface
func (g *giphy) Do(m *discordgo.MessageCreate, p []string) error {
	if len(p) < 2 {
		return nil
	}

	res, err := g.gophy.Request(&request.Search{
		Query: url.PathEscape(p[1]),
		Limit: searchLimit,
	})
	if err != nil {
		return nil
	}

	var entity entity.Gifs

	err = json.Unmarshal(res, &entity)
	if err != nil {
		return nil
	}

	n := int64(len(entity.Data))

	if n == 0 {
		g.session.ChannelMessageSend(m.ChannelID, notFoundMsg)
		return nil
	}

	g.session.ChannelMessageSend(m.ChannelID, entity.Data[tools.RandUnixNano(n)].EmbedURL)
	return nil
}

func (g *giphy) GetHelp() string {
	return "/gif [term], allows to search a gif through Giphy API"
}
