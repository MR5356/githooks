package runner

import (
	"fmt"
	"githooks/utils"
	"log"
)

const (
	buildPath  = "/tmp"
	buildFile  = "release/docker/build.sh"
	dockerfile = "release/docker/Dockerfile"
)

var (
	BuildTasks = BuildTaskQueue{make([]Build, 0)}
)

type BuildTaskQueue struct {
	Queue []Build `json:"queue"`
}

func (q *BuildTaskQueue) Enqueue(b Build) {
	q.Queue = append(q.Queue, b)
	for {
		if q.Size() < 100 || (q.Queue[0].Success && q.Queue[0].StepCurrent != q.Queue[0].StepTotal) {
			break
		}
		q.Queue = q.Queue[1:]
	}
}

func (q *BuildTaskQueue) Size() int {
	return len(q.Queue)
}

type Build struct {
	Name        string `json:"name"`
	From        string `json:"from"`
	SshUrl      string `json:"sshUrl"`
	HttpUrl     string `json:"httpUrl"`
	CommitId    string `json:"commitId"`
	UserName    string `json:"userName"`
	StepCurrent int    `json:"stepCurrent"`
	StepTotal   int    `json:"stepTotal"`
	Success     bool   `json:"success"`
}

func NewDefaultBuild() Build {
	return Build{
		StepTotal:   11,
		StepCurrent: 0,
		Success:     true,
	}
}

func (b *Build) stepPrint(message string) {
	b.StepCurrent = b.StepCurrent + 1
	log.Printf("build Runner Step %d/%d: %s", b.StepCurrent, b.StepTotal, message)
}

func (b *Build) failedPrint() {
	b.Success = false
	log.Printf("build Runner Step %d/%d: Build %s failed", b.StepCurrent, b.StepTotal, b.Name)
	log.Printf("build Runner Step %d/%d Build info: %+v", b.StepCurrent, b.StepTotal, b)
}

func (b Build) Run() {
	BuildTasks.Enqueue(b)
	var (
		gitClone = fmt.Sprintf("cd %s && git clone %s", buildPath, b.SshUrl)
		buildF   = fmt.Sprintf("%s/%s/%s", buildPath, b.Name, buildFile)
		dockerF  = fmt.Sprintf("%s/%s/%s", buildPath, b.Name, dockerfile)

		buildCommand        = fmt.Sprintf("cd %s/%s && /bin/bash %s %s", buildPath, b.Name, buildFile, b.CommitId)
		defaultBuildCommand = fmt.Sprintf("cd %s/%s && docker build -t %s:%s -f %s .", buildPath, b.Name, b.Name, b.CommitId, dockerfile)

		cleanNoneTagImagesCommand = fmt.Sprintf("docker rmi `docker images|grep none|awk '{print $3 }'|xargs` || echo none")
		cleanBuildPathCommand     = fmt.Sprintf("rm -rf %s/%s", buildPath, b.Name)
	)

	//执行预清理
	b.stepPrint(fmt.Sprintf("clean build path"))
	if !utils.RunCommand(cleanBuildPathCommand) {
		b.failedPrint()
		return
	}

	//克隆仓库
	b.stepPrint(fmt.Sprintf("Clone %s", b.SshUrl))
	if !utils.RunCommand(gitClone) {
		b.failedPrint()
		return
	}

	//执行构建
	b.stepPrint(fmt.Sprintf("%s IsExist: ???", dockerF))
	if utils.IsExist(dockerF) {
		b.stepPrint(fmt.Sprintf("%s IsExist: true", dockerF))
		b.stepPrint(fmt.Sprintf("%s IsExist: ???", buildF))
		if utils.IsExist(buildF) {
			b.stepPrint(fmt.Sprintf("%s IsExist: true", buildF))
			b.stepPrint(fmt.Sprintf("Build %s start", b.Name))
			if !utils.RunCommand(buildCommand) {
				b.failedPrint()
				return
			}
			b.stepPrint(fmt.Sprintf("Build %s successed", b.Name))
		} else {
			b.stepPrint(fmt.Sprintf("%s IsExist: false", buildF))
			b.stepPrint(fmt.Sprintf("%s Build file isn't exist, use default build action", b.Name))
			if !utils.RunCommand(defaultBuildCommand) {
				b.failedPrint()
				return
			}
			b.stepPrint(fmt.Sprintf("Build %s successed", b.Name))
		}
	} else {
		b.stepPrint(fmt.Sprintf("%s IsExist: false", dockerF))
		b.stepPrint(fmt.Sprintf("%s Dockerfile isn't exist, build skip", b.Name))
		b.stepPrint(fmt.Sprintf("%s IsExsit: skip", buildF))
		b.stepPrint(fmt.Sprintf("Build %s: skip", b.Name))
		b.stepPrint(fmt.Sprintf("Default build action: skip"))
	}

	//执行清理
	b.stepPrint(fmt.Sprintf("clean none tag images"))
	if !utils.RunCommand(cleanNoneTagImagesCommand) {
		b.failedPrint()
		return
	}

	b.stepPrint(fmt.Sprintf("clean build path"))
	if !utils.RunCommand(cleanBuildPathCommand) {
		b.failedPrint()
		return
	}

	//构建完成
	b.stepPrint(fmt.Sprintf("Build %s finished", b.Name))
}
