package types

type Mode string

const (
	ModeRelease Mode = "release"
	ModeDebug   Mode = "debug"
)

type BuildContext struct {
	Mode     Mode
	Compiler string
	Sources  []string
	Includes []string
	Defines  []string
	CFlags   []string
	LdFlags  []string
	Output   string
	WorkDir  string
	Verbose  bool
}

type BuildRequest struct {
	Mode       Mode
	Sources    []string
	TargetName string
	Verbose    bool
	WorkDir    string
}

type CompileCommand struct {
	Directory string
	Command   []string
	File      string
	Output    string
}
