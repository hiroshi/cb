package main

import (
	"compress/gzip"
	// "bytes"
	"encoding/json"
	"flag"
	// "fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	// "path/filepath"
	"strings"
)

type Step struct {
	ImageName string `json:"name"`
	Args []string
}

type Config struct {
	Steps []Step
}

func DockerOutput(arg ...string) (string, error) {
	cmd := exec.Command("docker", arg...)
	log.Println(cmd.Args)
	out, err := cmd.Output()
	if err != nil {
		log.Print(err)
		return "", err
	}
	return strings.TrimSpace(string(out[:])), nil
}

func cbMain() (exitCode int) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Take arguments and options
	source := os.Args[1]
	log.Println("source:", source)
	flags := flag.NewFlagSet("", flag.ExitOnError)
	configPath := flags.String("config", "", "The .yaml or .json file to use for build configuration.")
	flags.Parse(os.Args[2:])

	// Read the config file
	file, err := os.Open(*configPath)
	if err != nil {
		log.Print(err)
		return 1
	}
	defer file.Close()
	var b []byte
	b, err = ioutil.ReadAll(file)
	if err != nil {
		log.Print(err)
		return 1
	}

	// Create workspace volume
	workspace, err := DockerOutput("volume", "create")
	if err != nil {
		return 1
	}
	log.Println("workspace:", workspace)
	defer DockerOutput("volume", "rm", workspace)
	// Create a container to be destination of `docker cp`
	container, err := DockerOutput("create", "--volume", workspace + ":/workspace", "busybox")
	if err != nil {
		return 1
	}
	log.Println("container:", container)
	defer DockerOutput("rm", container)
	// open source tar stream
	sourceReader, err := os.Open(source)
	if err != nil {
		log.Print(err)
		return 1
	}
	defer sourceReader.Close()
	tarReader, err := gzip.NewReader(sourceReader)
	if err != nil {
		log.Print(err)
		return 1
	}
	defer tarReader.Close()
	// Expand source to workspace volume through container
	cmd := exec.Command("docker", "cp", "-", container + ":/workspace")
	log.Println(cmd.Args)
	cmd.Stdin = tarReader
	if err := cmd.Start(); err != nil {
		log.Print(err)
		return 1
	}
	if err := cmd.Wait(); err != nil {
		log.Print(err)
		return 1
	}

	var config Config
	err = json.Unmarshal(b, &config)
	for _, step := range config.Steps {
		args := append([]string{
			"run",
			"--rm",
			"--volume", "/var/run/docker.sock:/var/run/docker.sock",
			// --volume /root/.docker:/root/.docker
			"--volume", workspace + ":/workspace",
			"--workdir", "/workspace",
			// --env <KEY1=val1>
			// --env <KEY2=val2>
			step.ImageName},
			step.Args...)
		cmd := exec.Command("docker", args...)
		log.Println(cmd.Args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Print(err)
			return 1
		}
		if err := cmd.Wait(); err != nil {
			log.Print(err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(cbMain())
}
