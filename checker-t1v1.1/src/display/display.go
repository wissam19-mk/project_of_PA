package display

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

type PageStack []Page

func (ps *PageStack) Push() {

	*ps = append(*ps, Page{})
}

func (ps *PageStack) Pop() {
	*ps = (*ps)[:len(*ps)-1]
}

func (ps *PageStack) Page() *Page {
	return &(*ps)[len(*ps)-1]
}

func (ps *PageStack) Len() int {
	return len(*ps)
}

type PageElement struct {
	Element    tview.Primitive
	Fixed      int
	Proportion int
	Hidden     bool
	Focused    bool
}

type Page struct {
	Title    string
	Elements []*PageElement
	*WritableContainer
}

type Display struct {
	pageStack *PageStack
	Root      *tview.Flex
	App       *tview.Application
}

func (d *Display) currentPage() *Page {
	return d.pageStack.Page()
}

func NewDisplay() *Display {
	// Create the display
	display := &Display{}

	// Initialize the page stack
	display.pageStack = &PageStack{}
	display.NewPage("[yellow]Interactive Display", true)

	// Configure the Root element
	display.Root = tview.NewFlex()
	display.Root.SetBorder(true)
	display.Root.SetDirection(tview.FlexRow)

	// Configure App
	display.App = tview.NewApplication().SetRoot(display.Root, true)
	display.App.EnableMouse(true)

	display.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// PreviousPage page
		if event.Key() == tcell.KeyEscape {
			display.PreviousPage()
		}
		// Exit App
		if event.Key() == tcell.KeyCtrlC {
			display.App.Stop()
		}
		return event
	})

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	return display
}

func (d *Display) UpdateDisplay() {
	// Clear the screen
	d.Root.Clear()

	if d.currentPage().Title == "" {
		d.Root.SetTitle("[yellow]Interactive Display")
	} else {
		d.Root.SetTitle(d.currentPage().Title)
	}

	// Add back the elements in the Page
	for _, element := range d.currentPage().Elements {
		if element.Hidden {
			continue
		}
		d.Root.AddItem(element.Element, element.Fixed, element.Proportion, element.Focused)
	}
	if d.App != nil {
		d.App.ForceDraw()
	}
}

func (d *Display) Enable() {

	if err := d.App.Run(); err != nil {
		panic(err)
	}

}

func (d *Display) NewPage(title string, includeWContainer bool) {
	d.pageStack.Push()
	d.currentPage().Title = title

	if includeWContainer {
		d.currentPage().WritableContainer = NewWritableContainer(d)
	}

	if d.Root != nil {
		d.UpdateDisplay()
	}
}

func (d *Display) PreviousPage() {

	// Stop the App if the current page is the main menu
	if d.pageStack.Len() == 1 {
		d.Stop()
	}

	// Pop the current page and reload the last
	d.pageStack.Pop()
	d.App.SetFocus(d.CurrentContainer().Container)
	d.UpdateDisplay()

}

func (d *Display) IsStacked() bool {
	return d.pageStack.Len() > 1
}

func (d *Display) AddElement(element *PageElement) {
	d.currentPage().Elements = append(d.currentPage().Elements, element)
}

func (d *Display) AddWritableContainer(container *WritableContainer, fixed int, proportion int) {
	element := &PageElement{}
	element.Element = container.Container
	element.Fixed = fixed
	element.Proportion = proportion

	d.AddElement(element)
	d.currentPage().WritableContainer = container
}

func (d *Display) CurrentPageIndex() int {
	return d.pageStack.Len() - 1
}

func (d *Display) CurrentContainer() *WritableContainer {
	return d.currentPage().WritableContainer
}

func (d *Display) Stop() {
	d.App.Stop()
	os.Exit(0)
}

func (d *Display) PrintPage(index int, title string, buffer string) {
	d.CurrentContainer().PrintIndex(index, title, buffer)
	d.UpdateDisplay()
}

func (d *Display) Print(buffer string) {
	d.PrintPage(0, "", buffer)
}

func (d *Display) Println(buffer string) {
	d.Print(buffer)
	d.Print("\n")
}
