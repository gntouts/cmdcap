package cmdcap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mitchellh/go-ps"
	"github.com/sirupsen/logrus"
)

var once sync.Once

type ProcessInfo struct {
	Pid    int
	Pname  string
	Ppid   int
	Ppname string
	Cmd    string
	Cwd    string
}

func GetProcessInfo() *ProcessInfo {
	pid := Pid()
	ppid := Ppid()
	cmd := strings.Join(os.Args[:], " ")
	cwd, err := os.Getwd()
	catch(err)
	cwd, err = filepath.Abs(cwd)
	catch(err)
	return &ProcessInfo{
		Pid:    pid,
		Ppid:   ppid,
		Pname:  Name(pid),
		Ppname: Name(ppid),
		Cmd:    cmd,
		Cwd:    cwd,
	}
}

func (p *ProcessInfo) ToString() string {
	s := fmt.Sprintf("cwd: %s, p: %s (%d), pp: %s (%d), cmd: %s", p.Cwd, p.Pname, p.Pid, p.Ppname, p.Ppid, p.Cmd)
	return s
}

func (p *ProcessInfo) ToFields() logrus.Fields {
	return logrus.Fields{
		"pid":    p.Pid,
		"pname":  p.Pname,
		"ppid":   p.Ppid,
		"ppname": p.Ppname,
		"cmd":    p.Cmd,
		"cwd":    p.Cwd,
	}
}

func Pid() int {
	return os.Getpid()
}

func Ppid() int {
	return os.Getppid()
}

func Name(process int) string {
	exec, err := ps.FindProcess(process)
	if err != nil {
		return ""
	}
	return exec.Executable()
}

func Pname() string {
	return Name(Pid())
}

func PPname() string {
	return Name(Ppid())
}

func Test() func() {
	return func() { fmt.Println("hi") }
}

func Capture(file string) {
	once.Do(func() {
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
			ForceQuote:    true,
		})
		logger.Level = logrus.DebugLevel
		logger.SetReportCaller(true)
		logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		defer logFile.Close()
		logger.SetOutput(logFile)
		info := GetProcessInfo()
		logger.WithFields(info.ToFields()).Debug("cmdcap")
	})
}

func catch(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
