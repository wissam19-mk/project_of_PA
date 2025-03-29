package checkermodules

import (
	"checker-pa/src/display"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

type DummyModule struct {
	totalScore int
	uniqueName string
	randGen    *big.Int
	ModuleOutput
}

func NewDummyModule() *DummyModule {

	newDummy := &DummyModule{}
	newDummy.totalScore = 0
	var err error
	newDummy.randGen, err = rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		panic(err)
	}
	newDummy.uniqueName = "dummy-" + fmt.Sprintf("%x", newDummy.randGen.Int64()%255)

	return newDummy
}

func (dummy *DummyModule) GetName() string {
	return dummy.uniqueName
}

func (dummy *DummyModule) IsOutputDependent() bool { return false }

func (dummy *DummyModule) GetDependencies() []string { return nil }

func (dummy *DummyModule) Disable(_ bool) {}
func (dummy *DummyModule) Enable()        {}

func (dummy *DummyModule) GetStatus() ModuleStatus { return Ready }

func (dummy *DummyModule) Run() {
	const issueCount = 25

	for i := 0; i < issueCount; i++ {
		dummy.Issues = append(
			dummy.Issues,
			ModuleIssue{
				Message: "Lorem ipsum dolor sit amet",
				Line:    int(dummy.randGen.Int64() % 255),
				Col:     int(dummy.randGen.Int64() % 100),
			})
	}

	dummy.totalScore = int(dummy.randGen.Int64() % 70)
}

func (dummy *DummyModule) Display(d *display.Display) {

	// Set the page title
	d.PrintPage(0, dummy.uniqueName, "")

	d.Println("\nTotal module score: " + strconv.Itoa(int(dummy.totalScore)))

	if len(dummy.Issues) > 0 {

		d.PrintPage(1, dummy.uniqueName+" errors", dummy.ModuleError.String())

	}
}

func (dummy *DummyModule) Dump() {
	fmt.Printf("===== %s - %d =====\n\n", dummy.uniqueName, dummy.totalScore)
	fmt.Println(dummy.ModuleError.String())
	fmt.Println()

}

func (dummy *DummyModule) Reset() {
	dummy.totalScore = 0
	dummy.Issues = nil
}

func (dummy *DummyModule) Score() int {
	return dummy.totalScore
}
