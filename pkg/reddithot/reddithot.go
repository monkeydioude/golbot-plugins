package reddithot

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monkeydioude/lgtR"
	"github.com/turnage/graw/reddit"
)

type redditHot struct {
	hot     *lgtR.Hot
	subList map[string]*lgtR.Watcher
}

func AddCommand(cachePath string) *redditHot {
	return &redditHot{
		hot:     lgtR.New(cachePath, (2*time.Minute + 30*time.Second)),
		subList: make(map[string]*lgtR.Watcher),
	}
}

type action func(string, *discordgo.Session, *discordgo.MessageCreate)

func (r *redditHot) GetRegex() string {
	return "/hot (add|rm) (.+)"
}

func (r *redditHot) getFunctionMap() map[string]action {
	return map[string]action{
		"add": r.addSub,
		"rm":  r.rmSub,
	}
}

func getEmbedMessage(sub string, p *reddit.Post) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: p.URL,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("**[%s]**", sub),
				Value: fmt.Sprintf("[%s](https://www.reddit.com%s)", p.Title, p.Permalink),
			},
		},
	}
}

func (r *redditHot) addSub(sub string, s *discordgo.Session, m *discordgo.MessageCreate) {
	subID := m.ChannelID + sub
	if _, ok := r.subList[subID]; ok {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le sub '%s' est deja stalke !", sub))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Stalker le sub '%s' fait partie du keikaku !", sub))
	r.subList[subID] = r.hot.WatchMe(sub, func(p *reddit.Post) {
		if p.IsSelf {
			return
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, getEmbedMessage(sub, p))

		if err != nil {
			log.Printf("[ERR ] Could not send embed message. Reason: %s", err)
		}
	})
}

func (r *redditHot) rmSub(sub string, s *discordgo.Session, m *discordgo.MessageCreate) {
	subID := m.ChannelID + sub
	if _, ok := r.subList[subID]; !ok {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le sub '%s' n'est pas stalke !", sub))
		return
	}

	r.subList[subID].Cancel()
	delete(r.subList, subID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("JA-JA-J'arrete de stalker le sub '%s' !", sub))
}

func (r *redditHot) Do(s *discordgo.Session, m *discordgo.MessageCreate, p []string) error {
	if len(p) < 3 {
		return nil
	}

	funcs := r.getFunctionMap()
	if _, ok := funcs[p[1]]; ok {
		funcs[p[1]](p[2], s, m)
		return nil
	}

	return nil
}

func (r *redditHot) GetHelp() string {
	return "/hot [add|rm] allows to mirrors/remove hot post of a subbreddit (ex: /hot add shitpostcrusaders)"
}

func (r *redditHot) GetName() string {
	return "redditHot"
}
