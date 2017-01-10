package gitapi

import (
	"os"
	"os/exec"
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("cst")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	//backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.DEBUG, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

func Exist(filename string) bool {
	_,err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//执行命令 返回结果
func ExecCommand(commands string) (string,error) {
    out,err := exec.Command("bash", "-c", commands).Output()
    return string(out),err
}

func CheckErr(rs string,err error) string {
	if err != nil {
		return err.Error()
	}
	return rs
}