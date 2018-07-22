package log

import (
	"fmt"
	"time"
)

type Enabler struct {
	Until  time.Time // effective until this time
	Since  time.Time // effective before this time
	Count1 int       // minus 1 for each hit; effective when = 1; 0 disabled
	Count2 int       // enabled when Count1 is 1; minus 1 for each hit; effective when = 1; 0 disabled
	//CountMode int
}

type Probe struct {
	Package   string
	Module    string
	Line      int
	Level     int
	Tags      []string
	ToFile    bool
	Actions   []func(msg string)
	HitCount  uint64
	Condition Enabler
	Sample    string // a sample string illustrate what to be logged
	Msg       string // original log.Logln(...)
}

const (
	//LvlFatal   = 0
	//LvlAction  = -1
	LvlAlways  = 1
	LvlError   = 100
	LvlWarn    = 200
	LvlInfo    = 300
	LvlVerbose = 400
	LvlTrace   = 500
	LvlOff     = 999
)

const (
	ModeDefault  = 0
	ModeFiltered = 1
)

var allProbes = make(map[string]*Probe)

var globalLevel = LvlWarn
var globalMode = ModeDefault

func GetLevel() int {
	return globalLevel
}

func SetLevel(lvl int) {
	globalLevel = lvl
}

func GetMode() int {
	return globalMode
}

func SetMode(mode int) {
	globalMode = mode
}

// name = package/module:line
func RegProbe(name string, probe *Probe) {
	allProbes[name] = probe
}

func GetProbe(name string) *Probe {
	return allProbes[name]
}

func SetProbeLvl(probe *Probe, lvl int) {
	probe.Level = lvl
}

func SetProbeLvlByName(name string, lvl int) {
	if p := allProbes[name]; p != nil {
		p.Level = lvl
	}
}

func IsProbeEnabled(probe *Probe) bool {
	// TODO: ...
	return false
}

func Logln(args ...interface{}) {
	fmt.Println(args...)
}
