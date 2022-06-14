package utils

import (
	"os"
	"os/exec"
	"runtime"
)

func Clear() {
	optSys := runtime.GOOS
	if optSys == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
