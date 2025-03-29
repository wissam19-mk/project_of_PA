package menu

import (
	"checker-pa/src/checker-modules"
	"checker-pa/src/display"
	"checker-pa/src/manager"
	"checker-pa/src/utils"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

type Menu struct {
	*display.Display
	*manager.Manager
	nav      *tview.List
	launched bool

	redraw func()
}

const (
	labelColor   = tcell.ColorWhite
	labelBgColor = tcell.ColorDefault
	fieldColor   = tcell.ColorWhite
	fieldBgColor = tcell.ColorBlue

	fieldSelColor   = tcell.ColorBlue
	fieldSelBgColor = tcell.ColorWhite

	labelWidth = 20
	fieldWidth = 30
)

func orString(str1, str2 string) string {
	if str1 == "" {
		return str2
	}

	return str1
}

func (m *Menu) displayOptions() {

	m.CurrentContainer().Clear()
	m.CurrentContainer().Title("Options", tview.AlignLeft)

	exPath := tview.NewInputField()
	exPath.SetLabel("Executable Path")
	exPath.SetText(utils.Config.ExecutablePath)
	exPath.SetFieldWidth(fieldWidth)

	iPath := tview.NewInputField()
	iPath.SetLabel("Input Path")
	iPath.SetText(utils.Config.InputPath)
	iPath.SetFieldWidth(fieldWidth)

	sPath := tview.NewInputField()
	sPath.SetLabel("Source Path")
	sPath.SetText(utils.Config.SourcePath)
	sPath.SetFieldWidth(fieldWidth)

	oPath := tview.NewInputField()
	oPath.SetLabel("Output Path")
	oPath.SetText(utils.Config.OutputPath)
	oPath.SetFieldWidth(fieldWidth)

	rPath := tview.NewInputField()
	rPath.SetLabel("Ref Path")
	rPath.SetText(utils.Config.RefPath)
	rPath.SetFieldWidth(fieldWidth)

	fPath := tview.NewInputField()
	fPath.SetLabel("Forward Path")
	fPath.SetText(utils.Config.ForwardPath)
	fPath.SetFieldWidth(fieldWidth)

	rValgrind := tview.NewCheckbox()
	rValgrind.SetDisabled(false)
	rValgrind.SetLabel("Valgrind")
	rValgrind.SetChecked(utils.Config.RunValgrind)

	form := tview.NewFlex()
	form.SetDirection(tview.FlexRow)

	form.AddItem(exPath, 0, 1, false)
	form.AddItem(iPath, 0, 1, false)
	form.AddItem(sPath, 0, 1, false)
	form.AddItem(oPath, 0, 1, false)
	form.AddItem(rPath, 0, 1, false)
	form.AddItem(fPath, 0, 1, false)
	form.AddItem(rValgrind, 0, 1, false)

	items := []*tview.InputField{
		exPath, iPath, sPath, oPath, rPath, fPath,
	}

	formItems := []tview.FormItem{
		exPath, iPath, sPath, oPath, rPath, fPath, rValgrind,
	}

	m.redraw = func() {
		// Avoid saving of config when checker starts because of exec change
		for _, item := range items {
			item.SetText("")
		}
		rValgrind.SetChecked(utils.Config.RunValgrind)
		m.redraw = nil
		// m.displayHome()

	}

	itemIndex := 0
	firstLaunch := true

	for i, item := range items {
		item.SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldColor, fieldBgColor)
		rValgrind.SetCheckedStyle(tcell.StyleDefault.Foreground(fieldColor).Background(fieldBgColor))
		rValgrind.SetUncheckedStyle(tcell.StyleDefault.Foreground(fieldColor).Background(fieldBgColor))
		// Update selected item on click
		item.SetFocusFunc(func() {
			firstLaunch = false
			formItems[itemIndex].SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldColor, fieldBgColor)
			itemIndex = i
			item.SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldSelColor, fieldSelBgColor)
		})
	}

	rValgrind.SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldColor, fieldBgColor)
	rValgrind.SetFocusFunc(func() {
		firstLaunch = false
		formItems[itemIndex].SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldColor, fieldBgColor)
		itemIndex = len(items)
		rValgrind.SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldSelColor, fieldSelBgColor)
		rValgrind.SetCheckedStyle(tcell.StyleDefault.Foreground(fieldSelColor).Background(fieldSelBgColor))
		rValgrind.SetUncheckedStyle(tcell.StyleDefault.Foreground(fieldSelColor).Background(fieldSelBgColor))

	})

	// Stop the input fields from changing cursor the end / beginning
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
			return nil
		}

		return event
	})

	m.CurrentContainer().AddInputCallback(func(event *tcell.EventKey) {

		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
			oldItem := formItems[itemIndex]
			oldItem.SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldColor, fieldBgColor)
			// Dirty hack to change checkbox style
			if itemIndex == len(items) {
				rValgrind.SetCheckedStyle(tcell.StyleDefault.Foreground(fieldColor).Background(fieldBgColor))
				rValgrind.SetUncheckedStyle(tcell.StyleDefault.Foreground(fieldColor).Background(fieldBgColor))
			}
		}

		if !firstLaunch {
			switch event.Key() {
			case tcell.KeyUp:
				itemIndex--
			case tcell.KeyDown:
				itemIndex++
			default:
				return
			}
		} else {
			firstLaunch = false
		}

		if itemIndex < 0 {
			itemIndex = len(formItems) - 1
		} else if itemIndex >= len(formItems) {
			itemIndex = 0
		}

		currItem := formItems[itemIndex]
		m.App.SetFocus(currItem)
		currItem.SetFormAttributes(labelWidth, labelColor, labelBgColor, fieldSelColor, fieldSelBgColor)
		// Dirty hack to change checkbox style
		if itemIndex == len(items) {
			rValgrind.SetCheckedStyle(tcell.StyleDefault.Foreground(fieldColor).Background(fieldBgColor))
			rValgrind.SetUncheckedStyle(tcell.StyleDefault.Foreground(fieldColor).Background(fieldBgColor))
		}
		m.App.ForceDraw()
	})

	// TODO: maybe make a better way to save the fields

	m.CurrentContainer().AddChangeCallback(func() {
		// Save the new values
		utils.Config.ExecutablePath = orString(exPath.GetText(), utils.Config.ExecutablePath)
		utils.Config.InputPath = orString(iPath.GetText(), utils.Config.InputPath)
		utils.Config.SourcePath = orString(sPath.GetText(), utils.Config.SourcePath)
		utils.Config.OutputPath = orString(oPath.GetText(), utils.Config.OutputPath)
		utils.Config.RefPath = orString(rPath.GetText(), utils.Config.RefPath)
		utils.Config.ForwardPath = orString(fPath.GetText(), utils.Config.ForwardPath)
		utils.Config.RunValgrind = rValgrind.IsChecked()

		// TODO: trigger manager

		// Save the current config
		utils.SaveUserConfig()
	})

	m.CurrentContainer().AddPrimitive(form, false, 0, 1)
	m.UpdateDisplay()
}

func (m *Menu) displayRef() {
	m.CurrentContainer().Clear()

	m.redraw = func() {
		// Pop the pages until the nav page
		for m.IsStacked() {
			m.PreviousPage()
		}
		m.displayRef()
	}

	checkermodules.AvailableModules["ref_checker"].Display(m.Display)
}

func (m *Menu) displayStyle() {
	m.CurrentContainer().Clear()
	m.redraw = func() {
		// Pop the pages until the nav page
		for m.IsStacked() {
			m.PreviousPage()
		}
		m.displayStyle()
	}
	checkermodules.AvailableModules["style_checker"].Display(m.Display)
}

func (m *Menu) displayMemory() {
	m.CurrentContainer().Clear()

	m.redraw = func() {
		// Pop the pages until the nav page
		for m.IsStacked() {
			m.PreviousPage()
		}
		m.displayMemory()
	}
	checkermodules.AvailableModules["memory_checker"].Display(m.Display)
}

func (m *Menu) displayCommits() {
	m.CurrentContainer().Clear()
	m.redraw = func() {
		m.displayCommits()
	}
	checkermodules.AvailableModules["commit_checker"].Display(m.Display)
	m.CurrentContainer().WrapInput(m.CurrentContainer().Sections[0])
}

/*
func (m *Menu) displayHome() {
	// Set focus back to nav
	m.App.SetFocus(m.nav)
	m.nav.SetCurrentItem(0)

	m.CurrentContainer().Clear()
	m.CurrentContainer().Title("Home", tview.AlignLeft)
	m.CurrentContainer().SetDirection(tview.FlexRow)
	// TODO: set this to info about the modules, redraw the page when the manager triggers
	//m.CurrentContainer().PrintIndex(0, "Summary", "Placeholder for module summary")

	//m.StatusPing = func(caption string) {
	//	if len(m.CurrentContainer().Sections) > 0 {
	//		m.CurrentContainer().Sections[0].Clear()
	//	}
	//	summary := strings.Builder{}
	//	if caption != "" {
	//		summary.WriteString(caption)
	//	}
	//
	//	summary.WriteString("\n\n\n")
	//
	//	for _, module := range m.Modules {
	//		summary.WriteString(fmt.Sprintf("%-30s - %10s\n", module.GetName(), module.GetStatus().String()))
	//	}
	//
	//	m.CurrentContainer().PrintIndex(0, "Summary", summary.String())
	//
	//	m.App.ForceDraw()
	//}
	//
	//m.StatusPing("")

	m.CurrentContainer().PrintIndex(1, "Tutorial", "TUTORIAL\nTUTORIAL\nTUTORIAL\nTUTORIAL\nTUTORIAL\nTUTORIAL\n")
}
*/

func centerView(view *tview.TextView) *tview.Flex {
	proportion := 4
	container := tview.NewFlex()
	container.SetDirection(tview.FlexColumn)
	container.AddItem(tview.NewBox(), 0, 1, false)
	container.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(view, 0, proportion, false).
		AddItem(tview.NewBox(), 0, 1, false),
		0, proportion, false)
	container.AddItem(tview.NewBox(), 0, 1, false)

	return container
}

func (m *Menu) displayTutorial() {
	tutorialArea := tview.NewTextView()
	tutorialArea.SetTextAlign(tview.AlignCenter)
	tutorialArea.SetDynamicColors(true)
	tutorialArea.SetBorder(true)
	tutorialArea.SetTitle("TUTORIAL")
	tutorialArea.SetText(`
[deeppink]"Credits"[-]: [[greenyellow]"cosmin"[-], [greenyellow]"Ciprian"[-], [greenyellow]"steffe"[-], [greenyellow]"Horicuz"[-]]

Use the [yellow]mouse[-] or the [yellow]arrow keys[-] to navigate around

Press [yellow]TAB[-] to switch between navigation or current screen

Inside a full window press [yellow]ESC[-] to go back

To close the program press [yellow]ESC[-] inside the main page or [yellow]Ctrl+C[-]

To close this tutorial press [yellow]any key[-]

You can enable this back in the config
`)
	tutorialArea.SetBackgroundColor(tcell.ColorDefault)
	centeredContainer := centerView(tutorialArea)
	centeredContainer.SetInputCapture(func(_ *tcell.EventKey) *tcell.EventKey {
		utils.Config.Tutorial = false
		utils.SaveUserConfig()
		m.createMainMenu()
		m.App.SetRoot(m.Root, true)
		m.App.SetFocus(m.nav)
		m.UpdateDisplay()

		return nil
	})

	m.App.SetRoot(centeredContainer, true)
	m.App.SetFocus(tutorialArea)
}

func colorScore(score int) string {
	color := "[red]"
	if score < 35 { //nolint:gocritic
		color = "[red]"
	} else if score < 60 {
		color = "[yellow]"
	} else if score < 90 {
		color = "[green]"
	} else {
		color = "[aqua]"
	}

	return color + strconv.Itoa(score)
}

func (m *Menu) createMainMenu() {

	mainContainer := tview.NewFlex().SetDirection(tview.FlexColumn)

	buttonContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContainer.AddItem(buttonContainer, 0, 1, false)

	infoContainer := m.CurrentContainer().Container
	mainContainer.AddItem(infoContainer, 0, 3, false)

	m.nav = tview.NewList()
	m.nav.SetBorderPadding(1, 1, 1, 1)
	m.nav.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))

	m.nav.SetTitle("Navigation").
		SetTitleAlign(tview.AlignLeft)

	// This prevents double selections
	clicked := false

	m.nav.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			clicked = true
		}

		return action, event
	})

	m.nav.SetSelectedFunc(func(_ int, _ string, _ string, _ rune) {
		// m.CurrentContainer().Title("", 0)
		// m.App.SetFocus(m.CurrentContainer().Container)
		// utils.Log(fmt.Sprintf("Navigation: %s", main))
	})

	currNavItem := 0
	selNavItem := 0

	m.nav.SetChangedFunc(func(index int, _ string, _ string, _ rune) {
		// utils.Log("changed func")
		if clicked {
			clicked = false
			return
		}
		if m.launched {
			selNavItem = index
			m.nav.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
			currNavItem = index
		}
	})

	m.nav.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			if currNavItem == selNavItem {
				m.App.SetFocus(m.CurrentContainer().Container)
				return nil
			}
		}

		return event
	})

	/*
		m.nav.AddItem("home", "", 0, func() {
			m.displayHome()
		})
	*/

	m.nav.AddItem("Refs", "", 0, func() {
		m.displayRef()
	})
	m.nav.AddItem("Style", "", 0, func() {
		m.displayStyle()
	})
	m.nav.AddItem("Memory", "", 0, func() {
		m.displayMemory()
	})
	m.nav.AddItem("Commit", "", 0, func() {
		m.displayCommits()
	})
	m.nav.AddItem("Options", "", 0, func() {
		m.displayOptions()
	})

	m.nav.SetBorder(true)

	buttonContainer.AddItem(m.nav, 0, 4, true)

	/*
		scoreContainer := tview.NewTextView().SetText("80 / 100")
		scoreContainer.SetTitle("Score").
			SetTitleAlign(tview.AlignLeft).
			SetBorder(true)

		buttonContainer.AddItem(scoreContainer, 0, 1, true)
	*/

	infoBox := tview.NewTextView().
		SetDynamicColors(true) // .SetText("mucel Ciprian\nsteffe Horicuz").SetTextAlign(tview.AlignCenter)
	infoBox.SetTitle("Score - " + colorScore(0)).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	infoBox.SetBorder(true)
	infoBox.SetChangedFunc(func() {
		m.App.ForceDraw()
	})

	buttonContainer.AddItem(infoBox, 0, 3, false)

	m.StatusPing = func(caption string) {
		infoBox.SetTitle("Score - " + colorScore(m.TotalScore()))
		infoBox.Clear()
		summary := strings.Builder{}
		if caption != "" {
			summary.WriteString(caption)
		}

		summary.WriteString("\n\n\n")

		for _, module := range m.Modules {
			if module.GetStatus() != checkermodules.Ready {
				summary.WriteString(fmt.Sprintf("%-7s - %-8s\n", module.GetName(), module.GetStatus().String()))
			} else {
				summary.WriteString(fmt.Sprintf("%-7s - %-8s\n", module.GetName(), module.GetResult()))
			}
		}

		fmt.Fprintf(tview.ANSIWriter(infoBox), "%s", summary.String())

		// utils.Log("redrawing")
		if m.redraw != nil {
			m.redraw()
			m.App.ForceDraw()
		}

	}

	m.StatusPing("")

	// m.displayHome()

	m.displayRef()

	m.AddElement(&display.PageElement{Element: mainContainer, Proportion: 1, Focused: false})
	m.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// IDK if this should be called, like save the options and stuff
			// m.CurrentContainer().TriggerChange()
			m.PreviousPage()
			return nil
		}

		if event.Key() == tcell.KeyTab {
			if m.nav.HasFocus() {
				m.App.SetFocus(m.CurrentContainer().Container)
			} else if m.CurrentPageIndex() == 0 {
				m.App.SetFocus(m.nav)
			}

			return nil
		}

		if event.Rune() == '`' {
			err := m.Run()
			if err != nil {
				utils.Err(err.Error())
			}

			return nil
		}

		return event
	})
	m.App.SetFocus(m.nav)
	m.launched = true
}

func (m *Menu) Launch() {
	if utils.Config.Tutorial {
		m.displayTutorial()
	} else {
		m.createMainMenu()
	}
	m.UpdateDisplay()
}
