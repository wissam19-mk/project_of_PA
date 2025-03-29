package checkermodules

import (
	"checker-pa/src/display"
	"checker-pa/src/utils"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

type CommitChecker struct {
	ModuleOutput
	commits []string
	score   int
	status  ModuleStatus
}

var ErrNotFound = errors.New(" The checker couldn't find git on your system. Are you sure it's installed?")

func (*CommitChecker) GetName() string {
	return "COMMIT"
}

func (*CommitChecker) IsOutputDependent() bool {
	return utils.Config.CommitChecker.OutputDependent
}

func (*CommitChecker) GetDependencies() []string { return utils.Config.CommitChecker.Dependencies }

func (cc *CommitChecker) Disable(fail bool) {
	if fail {
		cc.status = DependencyFail
	} else {
		cc.status = Disabled
	}
}

func (cc *CommitChecker) Enable() {
	cc.status = Queued
}

func (cc *CommitChecker) GetStatus() ModuleStatus {
	return cc.status
}

func (cc *CommitChecker) GetResult() string {
	// TODO: add colors?
	return strconv.Itoa(cc.score)
}

func (cc *CommitChecker) Panic() {
	cc.status = Panic
}

func (cc *CommitChecker) Reset() {
	if cc.status == Disabled || cc.status == DependencyFail {
		return
	}
	cc.commits = nil
	cc.Issues = nil
	cc.score = 0
	cc.status = Queued
}

func (cc *CommitChecker) Score() int {
	return int(float32(cc.score) * utils.Config.CommitChecker.Grade)
}

func (cc *CommitChecker) Display(d *display.Display) {

	d.CurrentContainer().Title("Commit checker - "+strconv.Itoa(cc.Score()), tview.AlignLeft)

	// Disable border
	d.PrintPage(0, "$nb", "")

	if statusStr := StatusStr(cc); statusStr != "" {
		d.PrintPage(0, "$nb", statusStr)
		return
	}

	if len(cc.Issues) == 0 {
		d.Println("No issues found!")
		d.Println(fmt.Sprintf("Got %d/%d points! Congrats :)", cc.Score(), cc.Score()))
		return
	}

	// means we have a internal error
	if len(cc.Issues) == 1 {
		// means we have internal error or .git doesn't exit
		if cc.Issues[0].Critical {
			d.Println("Got an error!")
			d.Println(cc.Issues[0].Message)
			return
		}
	}

	d.Println("Detected some issues!")
	for _, issue := range cc.Issues {
		d.Println(issue.Message)
	}

	d.Println(fmt.Sprintf("The final score is %d/%d.", cc.Score(), int(utils.Config.CommitChecker.Grade*100)))
}

func (cc *CommitChecker) Dump() {
	fmt.Printf("===== Commit checker - %d =====\n\n", cc.Score())

	if cc.status != Ready {
		fmt.Println("This module is disabled.")
		return
	}

	fmt.Println(cc.ModuleError.String())
	fmt.Println()
}

// receives the commit line without the commit hash
func checkCommits(line string) error {
	if !strings.Contains(line, ":") {
		return errors.New("invalid format! Hint, the format is: <type of commit>: <message>)")
	}

	typeAndMessage := strings.SplitN(line, ":", 2)
	message := strings.Trim(typeAndMessage[1], " ")

	//TODO: replace 10 later
	if len(message) < 10 {
		return errors.New("the message is too short")
	}

	return nil
}

func (cc *CommitChecker) Run() {
	cc.status = Running
	defer func() { cc.status = Ready }()

	args := []string{"log", "--oneline", "--all"}
	cmd := exec.Command("git", args...)

	output, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			issue := ModuleIssue{Message: ErrNotFound.Error(), Critical: true}
			cc.Issues = append(cc.Issues, issue)
			return
		}
		// if the student didn't "git init" before, this will give an ambiguous error
		_, newErr := os.Stat(".git")
		if errors.Is(newErr, os.ErrNotExist) {
			errMsg := "Couldn't find any commits, are you sure you ran 'git init' first?"
			issue := ModuleIssue{Message: errMsg, Critical: true}
			cc.Issues = append(cc.Issues, issue)
			return
		}

		issue := ModuleIssue{Message: "CRITICAL ERROR! " + err.Error(), Critical: true}
		cc.Issues = append(cc.Issues, issue)
		return
	}

	if len(output) == 0 {
		errMsg := "The checker couldn't find any commits!"
		issue := ModuleIssue{Message: errMsg}
		cc.Issues = append(cc.Issues, issue)
		return
	}

	lines := strings.Split(string(output), "\n")

	// sanity check, it shouldn't happen ... i hope
	if len(lines) == 0 {
		// maybe put something more ... non screaming
		errMsg := "CRITICAL ERROR IN COMMIT CHECKER! Please make contact the team that made the checker. #1"
		issue := ModuleIssue{Message: errMsg, Critical: true}
		cc.Issues = append(cc.Issues, issue)
		return
	}

	lines = lines[0 : len(lines)-1]
	cc.commits = make([]string, len(lines))

	for i, line := range lines {
		splitLine := strings.SplitN(line, " ", 2)

		// this shouldn't happen but ... you never know
		if len(splitLine) == 0 {
			errMsg := "CRITICAL ERROR IN COMMIT CHECKER! Please make contact the team that made the checker. #2"
			issue := ModuleIssue{Message: errMsg, Critical: true}
			cc.Issues = append(cc.Issues, issue)
			return
		}
		if utils.Config.CommitChecker.UseFormat {
			err := checkCommits(splitLine[1])
			if err != nil {
				errMsg := "Bad commit detected: " + err.Error() + " the commit was \"" + splitLine[1] + "\"\n"
				issue := ModuleIssue{Message: errMsg}
				cc.Issues = append(cc.Issues, issue)
				continue
			}
		}

		cc.commits[i] = splitLine[1]
	}

	minCommits := utils.Config.CommitChecker.MinCommits
	cc.score = 100

	if minCommits > len(cc.commits) {
		pointsToDeduct := 1
		cc.score -= pointsToDeduct
		issueMsg := "Not enough commits have been made."
		cc.Issues = append(cc.Issues, ModuleIssue{Message: issueMsg})
	}

	deduction := 100 / utils.Config.CommitChecker.MinCommits

	if len(cc.Issues) == 0 {
		return
	}

	if len(cc.Issues) > 3 {
		cc.score = 0
		return
	}

	cc.score -= len(cc.Issues) * deduction
}
