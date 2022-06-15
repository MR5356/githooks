package utils

import (
	"bufio"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	Script = make(map[string]string, 0)
)

func LoadScripts() {
	Script = make(map[string]string, 0)
	fmt.Println("脚本加载开始...")
	files := GetExtFiles(GetAbsPath(), ".sh")
	for _, file := range files {
		filename := filepath.Base(file)
		Script[filename] = file
		fmt.Println(filename + ": " + file)
	}
	fmt.Println("脚本加载完毕")
}

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

func RunScript(scriptName string, args []string) []string {
	// 执行脚本前加载脚本
	LoadScripts()
	var scriptPath string
	if script, ok := Script[scriptName]; ok {
		scriptPath = script
	} else {
		fmt.Printf("%s 不存在", scriptName)
	}
	scriptPath = scriptPath + " " + strings.Join(args, " ")
	res := make([]string, 0)
	enc := mahonia.NewDecoder("gbk")

	optSys := runtime.GOOS
	cmd := exec.Command("bash", scriptPath)
	if optSys == "windows" {
		panic("仅支持在linux中执行脚本")
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
