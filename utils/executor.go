package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	Script      = make(map[string]string, 0)
	ScriptTasks = make(map[string]ScriptTask, 0)
)

type ScriptTask struct {
	Name       string `json:"name"`
	StartTime  int64  `json:"startTime"`
	FinishTime int64  `json:"finishTime"`
	Command    string `json:"command"`
}

func LoadScripts() {
	Script = make(map[string]string, 0)
	log.Println("脚本加载开始...")
	files := GetExtFiles(GetAbsPath(), ".sh")
	for _, file := range files {
		filename := filepath.Base(file)
		Script[filename] = file
		RunCommand(fmt.Sprintf("chmod +x %s", file))
		log.Println(filename + ": " + file)
	}
	log.Println("脚本加载完毕")
}

func RunCommand(command string) []string {
	res := make([]string, 0)

	// 支持中文编码
	enc := mahonia.NewDecoder("gbk")

	// 判断运行系统
	optSys := runtime.GOOS
	cmd := exec.Command("cmd", "/c", command)
	if optSys != "windows" {
		cmd = exec.Command("bash", "-c", command)
	}

	// 标准输出与错误输出
	var stderr = bytes.Buffer{}
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = &stderr

	// 日志处理
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
			log.Print(enc.ConvertString(readString))
			res = append(res, enc.ConvertString(readString))
		}
	}()

	// 流程处理
	err := cmd.Start()
	log.Printf("PID: %d command: %s", cmd.Process.Pid, command)

	wg.Wait()
	if err != nil {
		log.Println(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(enc.ConvertString(stderr.String()))
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
		log.Printf("%s 不存在", scriptName)
	}

	scriptPath = scriptPath + " " + strings.Join(args, " ")

	res := make([]string, 0)
	enc := mahonia.NewDecoder("gbk")

	optSys := runtime.GOOS
	cmd := exec.Command("bash", "-c", scriptPath)
	if optSys == "windows" {
		panic("仅支持在linux中执行脚本")
	}

	ScriptTasks[args[0]+":"+args[2]] = ScriptTask{args[0], time.Now().Unix(), 0, scriptPath}

	stdout, _ := cmd.StdoutPipe()
	var stderr = bytes.Buffer{}
	cmd.Stderr = &stderr

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
			log.Print(enc.ConvertString(readString))
			res = append(res, enc.ConvertString(readString))
		}
	}()

	err := cmd.Start()
	log.Printf("PID: %d command: %s", cmd.Process.Pid, scriptPath)
	wg.Wait()
	if err != nil {
		log.Println(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(enc.ConvertString(stderr.String()))
	}
	delete(ScriptTasks, args[0]+":"+args[2])
	return res
}
