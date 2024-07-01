package cli

import (
	"strings"

	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

const (
	colorReset      = "\033[m"
	colorError      = "\033[31m\033[1m"
	colorDarkYellow = "\033[33m\033[1m"
	colorInfo       = "\033[32m"
	color2          = "\033[34m"
	colorDebug      = "\033[36m"
)

func colored(msg, color string) string {
	return color + msg + colorReset
}

func E(msg string) string {
	return colorError + msg + colorReset
}

func I(msg string) string {
	return colorInfo + msg + colorReset
}

func D(msg string) string {
	return colorDebug + msg + colorReset
}

func C(msg string) string {
	return color2 + msg + colorReset
}

func ColorMsg(msg string) string {
	words := strings.Split(msg, " ")

	ln := len(words)
	if ln < shortMsg {
		return msg
	}

	if !util.ContainNumber(words[ln-1]) {
		return colorError + msg + colorReset
	}

	words[ln-1] = colorDebug + words[ln-1] + colorReset
	words[ln-3] = colorInfo + words[ln-3] + colorReset
	words[ln-5] = color2 + words[ln-5] + colorReset

	return strings.Join(words, " ")
}

func ColorService(ser string) string {
	return colorDarkYellow + ser + colorReset
}
