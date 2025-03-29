package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const (
	UserConfigPath = "./config.json"
)

var logFile *os.File
var logger *slog.Logger

var Config struct {
	*ModuleConfig
	*UserConfig

	DefaultUserConfig string
}

func InitConfig(defaultUserConfigStr string, moduleConfigStr string) error {
	var err error

	logFile, err = os.Create("./checker_log.txt")
	if err != nil {
		return err
	}

	logger = slog.New(slog.NewTextHandler(logFile, nil))

	Config.UserConfig, err = NewUserConfig(defaultUserConfigStr)
	if err != nil {
		return err
	}

	Config.ModuleConfig, err = newModuleConfig(moduleConfigStr)
	if err != nil {
		return err
	}

	Config.DefaultUserConfig = defaultUserConfigStr

	return nil
}

func SaveUserConfig() {
	f, err := os.Create(UserConfigPath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	newData, err := json.MarshalIndent(Config.UserConfig, "", "	")
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(string(newData))
	if err != nil {
		panic(err)
	}
}

func Log(str string) {
	logger.Info(str)
}

func Err(str string) {
	logger.Error(str)
}

func Fatal(str string) {
	fmt.Println(str)
	logger.Error(str)
	os.Exit(1)
}

var ConfigMacros = make(map[string]string)

func convertMacros(srcStr string, contextMacros map[string]string) string {
	// Replace implicit macros
	for k, v := range ConfigMacros {
		srcStr = strings.ReplaceAll(srcStr, fmt.Sprintf("$%s", k), v)
	}

	// Replace context macros
	for k, v := range contextMacros {
		srcStr = strings.ReplaceAll(srcStr, fmt.Sprintf("$%s", k), v)
	}

	return srcStr
}

// ExpandMacros No cyclic macros please!
func ExpandMacros(str string, contextMacros map[string]string) string {
	lastStr := ""

	for lastStr != str {
		lastStr = str
		str = convertMacros(str, contextMacros)
	}

	return str
}
