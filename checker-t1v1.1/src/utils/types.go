package utils

import "encoding/xml"

type Test struct {
	DisplayName string   `json:"displayName"`
	File        string   `json:"file"`
	Args        []string `json:"args"`
	Ordered     bool     `json:"ordered"`
	WhiteSpace  bool     `json:"whitespace"`
	Score       int      `json:"score"`
}

type RefChecker struct {
	OutputDependent bool    `json:"output_dependent"`
	Grade           float32 `json:"grade"`
}

type CommitChecker struct {
	Dependencies    []string `json:"dependencies"`
	OutputDependent bool     `json:"output_dependent"`
	MinCommits      int      `json:"minCommits"`
	UseFormat       bool     `json:"useFormat"`
	Grade           float32  `json:"grade"`
}

type MemoryChecker struct {
	Dependencies    []string `json:"dependencies"`
	OutputDependent bool     `json:"output_dependent"`
	MaxWarning      int      `json:"maxWarnings"`
	MaxLeak         int      `json:"maxLeak"`
	Grade           float32  `json:"grade"`
}

type StyleThreshold struct {
	Under int `json:"under"`
	Score int `json:"score"`
}
type StyleChecker struct {
	Dependencies    []string         `json:"dependencies"`
	OutputDependent bool             `json:"output_dependent"`
	ScoreThreshold  int              `json:"score_threshold"`
	Grade           float32          `json:"grade"`
	Thresholds      []StyleThreshold `json:"thresholds"`
}

type ModuleConfig struct {
	TempPath string            `json:"temp_path"`
	Macros   map[string]string `json:"macros"`
	Tests    []Test            `json:"tests"`

	*RefChecker    `json:"ref_checker"`
	*CommitChecker `json:"commit_checker"`
	*MemoryChecker `json:"memory_checker"`
	*StyleChecker  `json:"style_checker"`
}

type UserConfig struct {
	SourcePath     string `json:"source_path"`
	ExecutablePath string `json:"executable_path"`
	InputPath      string `json:"input_path"`
	OutputPath     string `json:"output_path"`
	RefPath        string `json:"ref_path"`
	ForwardPath    string `json:"forward_path"`
	RunValgrind    bool   `json:"run_valgrind"`
	Tutorial       bool   `json:"tutorial"`
}

type CppcheckResults struct {
	XMLName xml.Name   `xml:"results"`
	Version string     `xml:"version,attr"`
	Errors  []CppError `xml:"errors>error"`
}

type CppError struct {
	ID        string        `xml:"id,attr"`
	Severity  string        `xml:"severity,attr"`
	Msg       string        `xml:"msg,attr"`
	Verbose   string        `xml:"verbose,attr"`
	Locations []CppLocation `xml:"location"`
}

type CppLocation struct {
	File   string `xml:"file,attr"`
	Line   int    `xml:"line,attr"`
	Column int    `xml:"column,attr"`
	Info   string `xml:"info,attr"`
}
