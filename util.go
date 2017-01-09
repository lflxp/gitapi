package gitapi

import (
	"os/exec"
)

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