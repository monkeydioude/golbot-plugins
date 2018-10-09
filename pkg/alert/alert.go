package alert

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/drannoc/golbot"
	"github.com/bwmarrin/discordgo"
)

type alert struct {
	session  *discordgo.Session
	emitter  *discordgo.MessageCreate
	duration time.Duration
	reason   string
	timer    *time.Timer
}

const (
	ackMessage        = "Ok %s, j'te dirais !"
	tellMessage       = "Wesh %s"
	tellMessageReason = ", %s"
)

func (a *alert) sleep() {
	fmt.Println(a.duration.Seconds())
	a.timer = time.NewTimer(a.duration)

	a.session.ChannelMessageSend(a.emitter.ChannelID, fmt.Sprintf(ackMessage, a.emitter.Author.Mention()))
	<-a.timer.C
	heyListen(a.session, a.emitter, a.reason)
}

func init() {
	golbot.AddCommand(&cmd{
		alerts: make(map[string]map[string]*alert),
	})
}

func getDuration(t string) time.Duration {
	tParts := strings.Split(t, "h")

	if len(tParts) < 2 {
		return 0
	}

	now := time.Now()

	h, err := strconv.Atoi(tParts[0])
	if err != nil {
		return 0
	}

	d := time.Duration(0)
	if now.Hour() != h {
		d = d + time.Duration((24-now.Hour()+h)*int(time.Hour))
	}

	m, err := strconv.Atoi(tParts[1])
	if err != nil {
		return 0
	}

	return d + time.Duration((m-now.Minute())*int(time.Minute))
}

func errorMsg(s *discordgo.Session, m *discordgo.MessageCreate, msg string) golbot.KeepLooking {
	s.ChannelMessageSend(m.ChannelID, msg)
	return false
}

func heyListen(s *discordgo.Session, m *discordgo.MessageCreate, reason string) {
	msg := fmt.Sprintf(tellMessage, m.Author.Mention())

	if reason != "" {
		msg = fmt.Sprintf("%s%s", msg, fmt.Sprintf(tellMessageReason, reason))
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s !", msg))
}

type cmd struct {
	alerts map[string]map[string]*alert
}

func (c *cmd) GetRegex() string {
	return "/alert ([0-9]{1,2}h[0-9]{1,2})( .+)?"
}

func (c *cmd) cancelAlert(ID, time string) bool {
	if _, ok := c.alerts[ID][time]; !ok {
		return false
	}

	c.alerts[ID][time].timer.Stop()
	delete(c.alerts[ID], time)
	return true
}

func (c *cmd) Do(s *discordgo.Session, m *discordgo.MessageCreate, p []string) golbot.KeepLooking {
	if len(p) < 2 {
		return false
	}
	var reason string

	if len(p) == 3 {
		reason = strings.Trim(p[2], " ")

		if reason == "cancel" {
			c.cancelAlert(m.Author.ID, p[1])
			return errorMsg(s, m, ":middle_finger::joy:")
		}
		reason = p[2]
	}

	duration := getDuration(p[1])

	if duration == 0 {
		heyListen(s, m, reason)
		return false
	}

	a := &alert{
		session:  s,
		emitter:  m,
		duration: duration,
		reason:   reason,
	}

	go a.sleep()

	if _, ok := c.alerts[m.Author.ID]; !ok {
		c.alerts[m.Author.ID] = make(map[string]*alert)
	}
	c.alerts[m.Author.ID][p[1]] = a

	return false
}

func (c *cmd) GetHelp() string {
	return `Alert creates a reminding at the next specified time. Usage:
		* */alert 16:30 snack time* <= create a reminder to eat a snack
		* */alert 16:30 cancel* <= cancel the previously snack time reminder`
}
