package giphy

import (
	"encoding/json"
	"net/url"

	"bitbucket.org/drannoc/golbot"
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
	gophy *gophy.Gophy
}

func AddCommand() *giphy {
	return &giphy{
		gophy: gophy.NewGophy(apiKey),
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
func (g *giphy) Do(s *discordgo.Session, m *discordgo.MessageCreate, p []string) golbot.KeepLooking {
	if len(p) < 2 {
		return false
	}

	res, err := g.gophy.Request(&request.Search{
		Query: url.PathEscape(p[1]),
		Limit: searchLimit,
	})
	if err != nil {
		return false
	}

	var entity entity.Gifs

	err = json.Unmarshal(res, &entity)
	if err != nil {
		return false
	}

	n := int64(len(entity.Data))

	if n == 0 {
		s.ChannelMessageSend(m.ChannelID, notFoundMsg)
		return false
	}

	s.ChannelMessageSend(m.ChannelID, entity.Data[tools.RandUnixNano(n)].EmbedURL)
	return false
}

func (g *giphy) GetHelp() string {
	return "/gif, allows to search a gif through Giphy API"
}
