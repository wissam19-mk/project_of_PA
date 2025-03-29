package checkermodules

var AvailableModules = map[string]CheckerModule{
	"ref_checker":    NewDiffModule(),
	"memory_checker": &MemoryChecker{},
	"style_checker":  &StyleChecker{},
	"commit_checker": &CommitChecker{},
}
