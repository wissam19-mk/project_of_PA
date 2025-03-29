package checkermodules

import (
	"checker-pa/src/display"
	"checker-pa/src/utils"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	definitelyLeaked = "Leak_DefinitelyLost"
)

// ValgrindOutput represents a simplified version of Valgrind XML output focused on errors
type ValgrindOutput struct {
	Errors []Error `xml:"error"`
}

// Error represents a single error detected by Valgrind
type Error struct {
	Kind    string `xml:"kind"`
	What    string `xml:"what,omitempty"`  // Regular error description
	XWhat   XWhat  `xml:"xwhat,omitempty"` // Extended error description (for leaks)
	Stack   Stack  `xml:"stack"`
	AuxWhat string `xml:"auxwhat,omitempty"` // Additional error information
}

func (err *Error) isUserGenerated() bool {
	return strings.Contains(err.Stack.Frames[0].Obj, "/home/")
}

// XWhat contains simplified extended error information for memory leaks
type XWhat struct {
	Text        string `xml:"text"`
	LeakedBytes int    `xml:"leakedbytes"`
}

// Stack represents a call stack for an error
type Stack struct {
	Frames []Frame `xml:"frame"`
}

// Frame represents a simplified single stack frame
type Frame struct {
	Fn   string `xml:"fn"`
	File string `xml:"file,omitempty"`
	Line int    `xml:"line,omitempty"`
	Obj  string `xml:"obj,omitempty"`
}

type memoryCheckerIssue struct {
	message  string
	function string
	file     string
	line     int
}

func (mci *memoryCheckerIssue) String() string {
	str := mci.file + ":" + strconv.Itoa(mci.line) + " inside " + mci.function + " "
	str += mci.message

	return str
}

type TestMemoryResult struct {
	testName    string
	criticalMsg string
	issues      []memoryCheckerIssue
	warnings    []memoryCheckerIssue
}

type TestStatus int

const (
	OK TestStatus = iota
	WARNING
	ISSUE
	CRITICAL
)

func (tmr *TestMemoryResult) GetStatus() TestStatus {
	if tmr.criticalMsg != "" {
		return CRITICAL
	}

	if len(tmr.issues) > 0 {
		return ISSUE
	} else if len(tmr.warnings) > 0 {
		return WARNING
	}

	return OK
}

func (tmr *TestMemoryResult) String() string {

	if tmr.GetStatus() == OK {
		return fmt.Sprintf("%s - OK", tmr.testName)
	}

	str := strings.Builder{}

	if tmr.GetStatus() == CRITICAL {
		str.WriteString(fmt.Sprintf("%s - CRITICAL ERROR\n\n", tmr.testName))
		str.WriteString(tmr.criticalMsg + "\n")
		return str.String()
	}

	if len(tmr.issues) > 0 {
		str.WriteString(fmt.Sprintf("%s - Issues\n\n", tmr.testName))

		for _, issue := range tmr.issues {
			str.WriteString(issue.String() + "\n\n")
		}

		if len(tmr.warnings) > 0 {
			str.WriteString(strings.Repeat("-", 20) + "\n\n")
		}

	}

	if len(tmr.warnings) > 0 {
		str.WriteString(fmt.Sprintf("%s - Warnings\n\n", tmr.testName))

		for _, warning := range tmr.warnings {
			str.WriteString(warning.String() + "\n\n")
		}

	}

	return str.String()
}

type MemoryChecker struct {
	score  int
	tests  []TestMemoryResult
	status ModuleStatus
}

func (*MemoryChecker) GetName() string {
	return "MEMORY"
}

func (*MemoryChecker) IsOutputDependent() bool {
	return utils.Config.MemoryChecker.OutputDependent
}

func (*MemoryChecker) GetDependencies() []string { return utils.Config.MemoryChecker.Dependencies }

func (mc *MemoryChecker) Disable(fail bool) {
	if fail {
		mc.status = DependencyFail
	} else {
		mc.status = Disabled
	}
}

func (mc *MemoryChecker) Enable() {
	mc.status = Queued
}

func (mc *MemoryChecker) GetStatus() ModuleStatus {
	return mc.status
}

func (mc *MemoryChecker) GetResult() string {
	return fmt.Sprintf("%d leaks", mc.getTotalIssues())
}

func (mc *MemoryChecker) Reset() {
	if mc.status == Disabled || mc.status == DependencyFail {
		return
	}
	mc.tests = nil
	mc.score = 0
	mc.status = Queued
}

func (mc *MemoryChecker) Score() int {
	return int(float32(mc.score) * utils.Config.MemoryChecker.Grade)
}

func (mc *MemoryChecker) Panic() {
	mc.status = Panic
}

func (mc *MemoryChecker) getTotalIssues() int {
	totalIssues := 0
	for _, test := range mc.tests {
		totalIssues += len(test.issues)
	}

	return totalIssues
}

func (mc *MemoryChecker) getStatus() TestStatus {

	currentGravity := OK

	for _, test := range mc.tests {
		if test.GetStatus() > currentGravity {
			currentGravity = test.GetStatus()
		}
	}

	return currentGravity
}

func (mc *MemoryChecker) Display(d *display.Display) {
	d.CurrentContainer().Title("Memory checker - "+strconv.Itoa(int(mc.Score())), tview.AlignLeft)

	if statusStr := StatusStr(mc); statusStr != "" {
		d.PrintPage(0, "$nb", statusStr)
		return
	}

	if mc.getStatus() == OK {
		// Disable border
		d.PrintPage(0, "$nb", "")
		d.Println(
			fmt.Sprintf("No issues found! Great job you got %d/%d :)!",
				mc.Score(), mc.Score()))
		return
	}

	fileTable := tview.NewTable()

	fileTable.SetInputCapture(utils.TableSelector(len(mc.tests), fileTable))

	currentRow := 0
	currentCol := 0

	MaxRow, MaxCol := utils.ComputeBestArea(len(mc.tests))

	for _, test := range mc.tests {
		if currentRow >= MaxRow && currentCol < MaxCol {
			currentRow = 0
			currentCol++
		}

		cell := tview.NewTableCell(test.testName)

		color := "[white]"

		switch test.GetStatus() {
		case OK:
			cell.SetTextColor(tcell.ColorGreen)
			color = "[green]"
		case WARNING:
			cell.SetTextColor(tcell.ColorYellow)
			color = "[yellow]"
		case ISSUE:
			cell.SetTextColor(tcell.ColorRed)
			color = "[red]"
		case CRITICAL:
			cell.SetTextColor(tcell.ColorDarkRed)
			color = "[red]"
		}

		cell.SetSelectable(true)
		cell.SetClickedFunc(func() bool {

			// TODO: don't add new page, replace container with the view and pop it on enter

			d.NewPage(color+test.testName, true)
			d.CurrentContainer().SetDirection(tview.FlexColumn)
			d.CurrentContainer().SyncSections(true)
			d.AddWritableContainer(d.CurrentContainer(), 0, 1)

			d.PrintPage(0, "$nb", "")

			// Disable border
			d.Println(test.String())

			d.App.SetFocus(d.CurrentContainer().Container)
			d.CurrentContainer().WrapInput(d.CurrentContainer().Sections[0])

			return false
		})
		fileTable.SetCell(currentRow, currentCol, cell)

		currentRow++
	}

	firstCell := fileTable.GetCell(0, 0)

	textColor, _, _ := firstCell.Style.Decompose()

	// Create reverse style
	firstCell.SetBackgroundColor(textColor)
	firstCell.SetTextColor(tcell.ColorWhite)

	d.CurrentContainer().AddPrimitive(fileTable, true, 0, 1)

}

func (mc *MemoryChecker) Dump() {
	fmt.Printf("===== %s - %d =====\n\n", "Memory checker", mc.Score())

	if mc.status != Ready {
		fmt.Println("This module is disabled.")
		return
	}

	if mc.getStatus() == OK {
		fmt.Println("No issues found! Great job :)!")
	}

	for i, test := range mc.tests {
		fmt.Println(test.String())

		if i < len(mc.tests)-1 {
			fmt.Println(strings.Repeat("=", 20) + "\n")
		}

	}

	fmt.Println()
}

func (mc *MemoryChecker) Run() {
	mc.status = Running
	defer func() { mc.status = Ready }()

	mc.score = 100

	// Preallocate to keep order and avoid conflicts in the goroutines
	mc.tests = make([]TestMemoryResult, len(utils.Config.Tests))

	// WaitGroup for goroutines
	wg := sync.WaitGroup{}

	for i, test := range utils.Config.Tests {

		wg.Add(1)

		go func() {
			defer wg.Done()

			absTempPath, err := filepath.Abs(utils.Config.TempPath)
			if err != nil {
				utils.Err("Failed to get absolute temp")
				return
			}

			data, err := os.ReadFile(fmt.Sprintf("%s/%s.xml", absTempPath, test.File))
			if err != nil {
				utils.Err(fmt.Sprintf("Failed to read file: %s.xml", test.File))
				return
			}

			testResult := TestMemoryResult{testName: test.DisplayName}

			var output ValgrindOutput
			err = xml.Unmarshal(data, &output)
			if err != nil {
				testResult.criticalMsg = err.Error()
				mc.tests = append(mc.tests, testResult)
				return
			}

			idx := len(output.Errors) - 1
			for idx > -1 && output.Errors[idx].Kind == definitelyLeaked {
				mci := memoryCheckerIssue{message: output.Errors[idx].XWhat.Text}
				mci.file = output.Errors[idx].Stack.Frames[1].File
				mci.function = output.Errors[idx].Stack.Frames[1].Fn
				mci.line = output.Errors[idx].Stack.Frames[1].Line

				testResult.issues = append(testResult.issues, mci)

				idx--
			}

			for idx > -1 {
				if output.Errors[idx].isUserGenerated() {
					w := memoryCheckerIssue{message: output.Errors[idx].What}
					w.file = output.Errors[idx].Stack.Frames[1].File
					w.function = output.Errors[idx].Stack.Frames[1].Fn
					w.line = output.Errors[idx].Stack.Frames[1].Line

					testResult.warnings = append(testResult.warnings, w)
				}

				idx--
			}

			mc.tests[i] = testResult

		}()

	}

	wg.Wait()

	deduction := 100 / utils.Config.MemoryChecker.MaxWarning
	if mc.getTotalIssues() >= utils.Config.MemoryChecker.MaxWarning {
		mc.score = 0
		return
	}
	if mc.getTotalIssues() == 0 {
		return
	}

	mc.score -= (mc.getTotalIssues() * deduction)
}
