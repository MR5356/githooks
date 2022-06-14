package utils

import (
	"bufio"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

func RunCommand(command string) []string {
	res := make([]string, 0)
	enc := mahonia.NewDecoder("gbk")

	optSys := runtime.GOOS
	cmd := exec.Command("cmd", "/c", command)
	if optSys != "windows" {
		cmd = exec.Command("bash", "-c", command)
	}
	stdout, _ := cmd.StdoutPipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(enc.ConvertString(readString))
			res = append(res, enc.ConvertString(readString))
		}
	}()

	err := cmd.Start()
	wg.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Process PID: ", cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}

	return res
}
