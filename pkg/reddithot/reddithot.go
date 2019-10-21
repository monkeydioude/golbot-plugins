package reddithot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monkeydioude/lgtR"
	"github.com/turnage/graw/reddit"
)

type redditHot struct {
	// cachedSub
	hot                *lgtR.Hot
	session            *discordgo.Session
	watchList          map[string]*lgtR.Watcher
	subListByChannelID *subList
}

func batchSubscribeToSubs(subsByChannelID map[string][]string) {

}

func AddCommand(cachePath string, session *discordgo.Session) *redditHot {
	r := &redditHot{
		session:            session,
		hot:                lgtR.New(cachePath, (2*time.Minute + 30*time.Second)),
		watchList:          make(map[string]*lgtR.Watcher),
		subListByChannelID: newSubList(cachePath),
	}

	r.subListByChannelID.addSavedSubFromCache(r)
	return r
}

type action func(string, *discordgo.MessageCreate)

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

func watchCallback(channelID, sub string, session *discordgo.Session) func(p *reddit.Post) {
	return func(p *reddit.Post) {
		if p.IsSelf {
			return
		}
		_, err := session.ChannelMessageSendEmbed(channelID, getEmbedMessage(sub, p))

		if err != nil {
			log.Printf("[ERR ] Could not send embed message. Reason: %s", err)
		}
	}
}

func watchSub(channelID, sub string, r *redditHot) error {
	lsub := strings.ToLower(sub)
	watchID := channelID + lsub
	if _, ok := r.watchList[watchID]; ok {
		return fmt.Errorf("Info: Sub '%s' is already being stalked", sub)
	}
	r.watchList[watchID] = r.hot.WatchMe(sub, watchCallback(channelID, sub, r.session))
	return nil
}

func (r *redditHot) addSub(sub string, m *discordgo.MessageCreate) {
	err := watchSub(m.ChannelID, sub, r)
	if err != nil {
		r.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Info: Sub '%s' is already being stalked.", sub))
		return
	}
	r.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Info: Stalking sub '%s' is now part of the keikaku.", sub))
	r.subListByChannelID.addSubToSubList(m.ChannelID, sub)
}

func (r *redditHot) rmSub(sub string, m *discordgo.MessageCreate) {
	lsub := strings.ToLower(sub)
	watchID := m.ChannelID + lsub
	if _, ok := r.watchList[watchID]; !ok {
		r.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Warning: Sub '%s' is not being stalked.", sub))
		return
	}

	r.watchList[watchID].Cancel()
	r.subListByChannelID.removeSubFromSubList(m.ChannelID, lsub)
	delete(r.watchList, watchID)
	// delete(r.subListByChannelID[m.ChannelID], lsub)
	r.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Info: Will not follow the sub '%s' anymore.", sub))
}

func (r *redditHot) Do(m *discordgo.MessageCreate, p []string) error {
	if len(p) < 3 {
		return nil
	}

	funcs := r.getFunctionMap()
	if _, ok := funcs[p[1]]; ok {
		funcs[p[1]](p[2], m)
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
