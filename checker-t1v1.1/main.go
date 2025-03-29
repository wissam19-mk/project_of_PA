package main

import (
	"checker-pa/src/display"
	"checker-pa/src/manager"
	"checker-pa/src/menu"
	"checker-pa/src/utils"
	_ "embed"
	"flag"
)

//go:embed res/config/module_config.json
var moduleConfigStr string

//go:embed res/config/user_config.json
var defaultUserConfigStr string

var useInteractive bool

func init() {
	flag.BoolVar(&useInteractive, "i", false, "Interactive mode")
}

func main() {
	flag.Parse()

	err := utils.InitConfig(defaultUserConfigStr, moduleConfigStr)
	if err != nil {
		utils.Fatal("FATAL ERROR DETECTED! " + err.Error() + "\n ABORTING!")
	}

	m, err := manager.NewManager()

	if err != nil {
		utils.Fatal("FATAL ERROR DETECTED! " + err.Error() + "\n ABORTING!")
	}

	// TODO?: make this return an error as well?
	// m.Check()

	if useInteractive {

		utils.Log("Interactive Display")
		d := display.NewDisplay()

		go func() {
			err := m.Run()
			if err != nil {
				d.App.Stop()
				utils.Fatal("FATAL ERROR DETECTED! " + err.Error() + "\n ABORTING!")
			}
		}()

		mn := menu.Menu{Display: d, Manager: m}

		mn.Launch()

		d.Enable()

	} else {
		err = m.Run()
		if err != nil {
			utils.Fatal("FATAL ERROR DETECTED! " + err.Error() + "\n ABORTING!")
		}

		utils.Log("Basic Display")

		for _, module := range m.Modules {
			module.Dump()
		}

		m.BasicSummary("")

	}
}
