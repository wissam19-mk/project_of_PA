package checkermodules

import (
	"checker-pa/src/display"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

type ModuleIssue struct {
	File        string
	Line        int
	Col         int
	Message     string
	ShowLineCol bool
	Critical    bool
}

type ModuleError struct {
	ErrorMessage string
	Issues       []ModuleIssue
}

type ModuleOutput struct {
	Score int32
	ModuleError
}

type ModuleStatus int

const (
	Ready ModuleStatus = iota
	Running
	Queued
	Disabled
	DependencyFail
	Panic
)

func (ms ModuleStatus) String() string {
	switch ms {
	case Ready:
		return color.New(color.FgGreen).Sprint("READY")
	case Queued:
		return color.New(color.FgYellow).Sprint("QUEUED")
	case Running:
		return color.New(color.FgYellow).Sprint("RUNNING")
	case Disabled:
		return color.New(color.FgBlack).Add(color.BgWhite).Sprint("DISABLED")
	case DependencyFail:
		return color.New(color.FgRed).Sprint("ERR")
	case Panic:
		return color.New(color.FgHiRed).Sprint("PANIC!")
	default:
		return "UNKNOWN"
	}
}

type CheckerModule interface {
	GetName() string
	IsOutputDependent() bool
	GetDependencies() []string
	Run()
	Display(d *display.Display)
	Dump()
	Reset()
	Score() int
	GetResult() string
	Disable(fail bool)
	Enable()
	GetStatus() ModuleStatus
	Panic()
}

func (err *ModuleError) String() string {
	message := err.ErrorMessage + "\n"

	for _, issue := range err.Issues {
		message += "\n"
		if issue.ShowLineCol {
			message += strconv.Itoa(issue.Line) + ":" + strconv.Itoa(issue.Col) + " "
		}
		message += issue.Message + "\n"
	}

	return message
}

func (err *ModuleError) groupIssues(groupBy func(issue *ModuleIssue) string) map[string][]ModuleIssue {

	group := make(map[string][]ModuleIssue)

	for _, issue := range err.Issues {
		group[groupBy(&issue)] = append(group[groupBy(&issue)], issue)
	}

	return group
}

func StatusStr(cm CheckerModule) string {
	switch cm.GetStatus() {
	case Disabled:
		return "This module is disabled."
	case DependencyFail:
		msg := strings.Builder{}
		msg.WriteString("One or more dependencies have failed.\nCheck if you have the following installed:")
		for _, dependency := range cm.GetDependencies() {
			msg.WriteString(dependency)
		}
		return msg.String()
	case Queued:
		fallthrough
	case Running:
		return "This module is currently running. Please wait"
	case Panic:
		return "The checker went into panic. Check the config and run again"
	default:
	}

	return ""
}
