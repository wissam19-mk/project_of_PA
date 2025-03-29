package manager

import (
	"bytes"
	"checker-pa/src/checker-modules"
	"checker-pa/src/utils"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Manager struct {
	Modules []checkermodules.CheckerModule

	capabilities map[string]bool

	StatusPing func(caption string)
}

func (m *Manager) BasicSummary(caption string) {
	summary := strings.Builder{}

	summary.WriteString("\n===== Summary =====\n")
	if caption != "" {
		summary.WriteString("\n" + caption + "\n\n")
	}

	for _, module := range m.Modules {
		if module.GetStatus() != checkermodules.Ready {
			summary.WriteString(fmt.Sprintf("%-7s - %-8s\n", module.GetName(), module.GetStatus().String()))
		} else {
			summary.WriteString(fmt.Sprintf("%-7s - %-8s\n", module.GetName(), module.GetResult()))
		}
	}

	summary.WriteString(fmt.Sprintf("\nScore: %d\n", m.TotalScore()))

	fmt.Println(summary.String())
}

func (m *Manager) checkCapabilities() {

	// Check for Valgrind

	utils.Log("Checking capabilities...")

	for _, module := range m.Modules {
		for _, dependency := range module.GetDependencies() {

			if _, err := exec.LookPath(dependency); err != nil {
				utils.Log("[ERR] " + dependency)
				module.Disable(true)
				continue
				// return errors.New("couldn't find valgrind on your system")
			}
			if module.GetStatus() != checkermodules.DependencyFail && dependency == "valgrind" && !utils.Config.RunValgrind {
				utils.Log("[Disabled] " + dependency)
				module.Disable(false)
			} else {
				module.Enable()
			}

			utils.Log("[OK] " + dependency)
			m.capabilities[dependency] = true
		}
	}

	/*
		if _, err := exec.LookPath("valgrind"); err != nil {
			utils.Log("[ERR] valgrind")
			//return errors.New("couldn't find valgrind on your system")
		}

		utils.Log("[OK] valgrind")
		m.capabilities["valgrind"] = true

		// Check for cppcheck
		if _, err := exec.LookPath("cppcheck"); err != nil {
			utils.Log("[ERR] cppcheck")
			//return errors.New("couldn't find cppcheck on your system")
		}

		utils.Log("[OK] cppcheck")
		m.capabilities["cppcheck"] = true

		return nil
	*/
}

func abs(rel string) string {
	absPath, _ := filepath.Abs(rel)

	return absPath
}

func NewManager() (*Manager, error) {
	var m Manager

	m.capabilities = make(map[string]bool)

	err := m.registerModules()
	if err != nil {
		return nil, err
	}

	for _, module := range m.Modules {
		module.Reset()
	}

	m.checkCapabilities()

	err = m.RetrieveConfig()
	if err != nil {
		return nil, err
	}

	go func() {

		prevPath := abs(utils.Config.ExecutablePath)
		prevStat, err := os.Stat(prevPath)
		for {
			currentPath := abs(utils.Config.ExecutablePath)
			currentStat, err2 := os.Stat(currentPath)
			if (err != nil && err2 == nil) || (err == nil && err2 == nil && (prevPath != currentPath || prevStat.ModTime().Before(currentStat.ModTime()))) {
				utils.Log("file change detected!")
				err := m.Run()
				if err != nil {
					utils.Err("failed manager run with error: " + err.Error())
				}
			}
			/*
				if err != nil && err2 == nil {
					// Something happened here, maybe the file got created or idk

				} else if err == nil && err2 == nil {
					if prevPath != currentPath {

					} else if prevStat.ModTime().Before(currentStat.ModTime()) {
						// prev file is modified before current, but how can we know if they're the same
					}
				}
			*/

			prevPath = currentPath
			prevStat, err = currentStat, err2

			time.Sleep(2 * time.Second)
		}
	}()

	return &m, nil
}

func (m *Manager) register(module checkermodules.CheckerModule) {
	m.Modules = append(m.Modules, module)
}

func (m *Manager) registerModules() error {
	if utils.Config.ModuleConfig.RefChecker != nil && checkermodules.AvailableModules["ref_checker"] == nil {
		return errors.New("ref_checker not available")
	}

	if utils.Config.ModuleConfig.MemoryChecker != nil && checkermodules.AvailableModules["memory_checker"] == nil {
		return errors.New("memory_checker not available")
	}

	if utils.Config.ModuleConfig.StyleChecker != nil && checkermodules.AvailableModules["style_checker"] == nil {
		return errors.New("style_checker not available")
	}

	if utils.Config.ModuleConfig.CommitChecker != nil && checkermodules.AvailableModules["commit_checker"] == nil {
		return errors.New("commit_checker not available")
	}

	for _, module := range checkermodules.AvailableModules {
		m.register(module)
	}

	return nil
}

func updateMacros() {
	// Output path
	outPath, err := filepath.Abs(utils.Config.OutputPath)
	if err == nil {
		utils.ConfigMacros["OUT_DIR"] = outPath
	}

	// Make sure input path exists
	inPath, err := filepath.Abs(utils.Config.InputPath)
	if err == nil {
		if _, err := os.Stat(inPath); err == nil {
			utils.ConfigMacros["IN_DIR"] = inPath
		}
	}

	srcPath, err := filepath.Abs(utils.Config.SourcePath)
	if err == nil {
		if _, err := os.Stat(srcPath); err == nil {
			utils.ConfigMacros["SRC_DIR"] = srcPath
		}
	}

	// Load module config macros
	for k, v := range utils.Config.Macros {
		utils.ConfigMacros[k] = v
	}

}

func (m *Manager) RetrieveConfig() error {
	defer updateMacros()

	if _, err := os.Stat(utils.UserConfigPath); err == nil {
		// Read the config from there
		data, err := os.ReadFile(utils.UserConfigPath)
		if err != nil {
			return err
		}

		// Bug: if fields are not present, they get changed to ""
		utils.Config.UserConfig, err = utils.NewUserConfig(string(data))
		if err != nil {
			return err
		}
	} else {
		f, err := os.Create(utils.UserConfigPath)
		if err != nil {
			return err
		}

		defer f.Close()

		_, err = f.WriteString(utils.Config.DefaultUserConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func forwardBytes(bytes bytes.Buffer, filename string) error {

	absForward, err := filepath.Abs(utils.Config.ForwardPath)
	if err != nil {
		return err
	}
	if err := os.Mkdir(absForward, 0777); err != nil {
		if !errors.Is(err, os.ErrExist) {
			utils.Err(fmt.Sprintf("failed creating forward directory: %s", absForward))
			return err
		}
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", absForward, filename))
	if err != nil {
		utils.Err(fmt.Sprintf("failed creating forward file: %s", filename))
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes.Bytes())
	if err != nil {
		utils.Err(fmt.Sprintf("failed writing to forward file: %s", filename))
		return err
	}

	return nil
}

// TODO: what to do when a new run is triggered but the current one isn't finished?
// TODO: should it be queued? idk how to cancel de current running routines

func (m *Manager) IsRunning() bool {
	// TODO: bug where the state is Queued instead of running, launching the run anyway
	// quick fix: check if at least one module has already finished
	oneReady := false
	oneQueued := false
	for _, module := range m.Modules {
		utils.Log("status: " + strconv.Itoa(int(module.GetStatus())))
		switch module.GetStatus() {
		case checkermodules.Ready:
			oneReady = true
		case checkermodules.Running:
			return true
		case checkermodules.Queued:
			oneQueued = true
		default:
		}
	}

	// utils.Log(fmt.Sprintf("%t", oneReady && oneQueued))

	return oneReady && oneQueued
}

func (m *Manager) Run() error {

	if m.IsRunning() {
		return errors.New("already running")
	}

	utils.Log("launched new run")

	m.checkCapabilities()

	if _, err := exec.LookPath(utils.Config.ExecutablePath); err != nil {
		for _, module := range m.Modules {
			module.Panic()
		}
		// Return nil only when in interactive mode

		if m.StatusPing != nil {
			m.StatusPing("[ERR] " + utils.Config.ExecutablePath + " not found")
			// This one is recoverable, don't crash
			return nil
		}
		
		m.BasicSummary("[ERR] " + utils.Config.ExecutablePath + " not found")
		return errors.New("executable not found: " + utils.Config.ExecutablePath)
	}

	start := time.Now()

	// Make sure temp path exists
	tempPath, err := filepath.Abs(utils.Config.TempPath)
	if err != nil {
		return err
	}

	if err := os.Mkdir(tempPath, 0777); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	// Make sure output path exists
	if err := os.Mkdir(utils.ConfigMacros["OUT_DIR"], 0777); err != nil {
		if !errors.Is(err, os.ErrExist) {
			for _, module := range m.Modules {
				module.Panic()
			}

			// Return nil only when in interactive mode

			if m.StatusPing != nil {
				m.StatusPing("[ERR] could not create directory: " + utils.Config.OutputPath)
				// This one is recoverable, don't crash
				return nil
			}

			m.BasicSummary("[ERR] could not create directory: " + utils.Config.OutputPath)
			return errors.New("could not create directory: " + utils.Config.OutputPath)
		}
	}

	wg := sync.WaitGroup{}

	for _, module := range m.Modules {
		module.Reset()
		if !module.IsOutputDependent() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer utils.Log(module.GetName() + " done!")
				if module.GetStatus() == checkermodules.Queued {
					module.Run()
				}
			}()
		}
	}

	var ranTests int32

	for i, test := range utils.Config.Tests {
		wg.Add(1)
		go func() {
			defer func() { wg.Done(); atomic.AddInt32(&ranTests, 1) }()

			// Create Context macros
			contextMacros := map[string]string{
				"FILE": test.File,
				"IN":   fmt.Sprintf("%s/%s.in", utils.ConfigMacros["IN_DIR"], test.File),
				"OUT":  fmt.Sprintf("%s/%s.out", utils.ConfigMacros["OUT_DIR"], test.File),
				"N":    strconv.Itoa(i),
			}

			var processedArgs []string

			// Process args
			for _, arg := range test.Args {
				processedArgs = append(processedArgs, utils.ExpandMacros(arg, contextMacros))
			}

			var cmd *exec.Cmd

			if m.capabilities["valgrind"] && utils.Config.RunValgrind {

				xmlPath := filepath.Join(tempPath, fmt.Sprintf("%s.xml", test.File))

				execPath, err := filepath.Abs(utils.Config.ExecutablePath)
				if err != nil {
					utils.Err(fmt.Sprintf("failed getting executable path: %s", xmlPath))
					return // err
				}

				valgrindArgs := []string{
					"--leak-check=yes",
					"--xml=yes",
					fmt.Sprintf("--xml-file=%s", xmlPath),
				}

				cmd = exec.Command("valgrind", append(append(valgrindArgs, execPath), processedArgs...)...) //nolint:gosec
				// fmt.Println("running: valgrind " + strings.Join(append(append(valgrindArgs, execPath), processedArgs...), " "))
			} else {
				cmd = exec.Command(utils.Config.ExecutablePath, processedArgs...) //nolint:gosec
			}

			// fmt.Printf("%d: %s %s\n\n", i+1, utils.Config.ExecutablePath, strings.Join(processedArgs, " "))

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			start = time.Now()

			if err := cmd.Run(); err != nil {
				utils.Err("Error running " + test.File)
			}

			// Forward stdout
			if err := forwardBytes(stdout, fmt.Sprintf("%s.stdout", test.File)); err != nil {
				utils.Err(fmt.Sprintf("failed forwarding stdout %s", test.File))
				return // err
			}

			// Forward stderr
			if err := forwardBytes(stderr, fmt.Sprintf("%s.stderr", test.File)); err != nil {
				utils.Err(fmt.Sprintf("failed forwarding stderr %s", test.File))
				return // err
			}

			utils.Log(fmt.Sprintf("[%s] %s", time.Since(start).String(), test.File))

		}()
	}

	updateDisplay := true

	const (
		barLength = 20
	)

	go func() {
		// Func available only in interactive mode
		if m.StatusPing == nil {
			return
		}
		for updateDisplay {
			builder := strings.Builder{}
			builder.WriteString("[")
			filled := int(math.Ceil(float64(ranTests) / float64(len(utils.Config.Tests)) * barLength))
			for i := 0; i < filled; i++ {
				builder.WriteString("#")
			}
			for i := filled; i < barLength; i++ {
				builder.WriteString(".")
			}
			builder.WriteString("]")

			m.StatusPing(builder.String())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
	m.Check()
	updateDisplay = false
	utils.Log("finished updating display")
	if m.StatusPing != nil {
		m.StatusPing("")
	}
	return nil
}

func (m *Manager) Check() {
	wg := sync.WaitGroup{}

	for _, module := range m.Modules {
		if module.IsOutputDependent() {
			wg.Add(1)
			utils.Log("Running " + module.GetName())
			go func() {
				defer wg.Done()
				if module.GetStatus() == checkermodules.Queued {
					module.Run()
				}
			}()
		}
	}

	wg.Wait()
}

func (m *Manager) TotalScore() int {
	var total int
	for _, module := range m.Modules {
		total += module.Score()
	}

	return total
}
