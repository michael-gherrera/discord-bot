package util

import (
	dg "github.com/bwmarrin/discordgo"
)

type Logger struct {
	Session   *dg.Session
	ChannelID string
}

// Use diff as the codeblock highlighter. Ghetto way of getting text colors.
// # BLUE
// + YELLOW-GREEN
// - RED

func (l *Logger) Info(s string) {
	fmtStr := "```md\n"
	fmtStr += "# INFO\n\n"
	fmtStr += "# " + s
	fmtStr += "```\n"
	l.Session.ChannelMessageSend(l.ChannelID, fmtStr)
}

func (l *Logger) Trace(s string) {
	fmtStr := "```diff\n"
	fmtStr += "+ TRACE\n\n"
	fmtStr += "+ " + s
	fmtStr += "```\n"
	l.Session.ChannelMessageSend(l.ChannelID, fmtStr)
}

func (l *Logger) Warn(s string) {
	fmtStr := "```diff\n"
	fmtStr += "- WARN\n\n"
	fmtStr += "- " + s
	fmtStr += "```\n"
	l.Session.ChannelMessageSend(l.ChannelID, fmtStr)
}
