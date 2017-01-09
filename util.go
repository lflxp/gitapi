package gitapi

import (
	"os/exec"
	"fmt"
)

//执行命令 返回结果
func ExecCommand(commands string) string {
    out,err := exec.Command("bash", "-c", commands).Output()
    CheckErr(err)
    return string(out)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}