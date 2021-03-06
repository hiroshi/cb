package main

import (
	"compress/gzip"
	// "bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Step struct {
	ImageName string `json:"name" yaml:"name"`
	Args []string
	Env map[string]string
	Dir string
}

type Config struct {
	Steps []Step
}

func ReadConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return config, err
	}
	defer file.Close()
	var b []byte
	b, err = ioutil.ReadAll(file)
	if err != nil {
		log.Print(err)
		return config, err
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		if err = json.Unmarshal(b, &config); err != nil {
			log.Print("JSON Error:", err)
			return config, err
		}
	case ".yaml", ".yml":
		if err = yaml.Unmarshal(b, &config); err != nil {
			log.Print("YAML Error:", err)
			return config, err
		}
	default:
		err := errors.New("Unknown configuration file extension: " + ext)
		log.Print(err)
		return config, err
	}
	return config, nil
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
	flags := flag.NewFlagSet("", flag.ExitOnError)
	configPath := flags.String("config", "", "The .yaml or .json file to use for build configuration.")
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s SOURCE -config CONFIG.(yml|json)\n", os.Args[0])
		flags.PrintDefaults()
		return 1
	}
	source := os.Args[1]
	log.Println("source:", source)
	flags.Parse(os.Args[2:])

	// Read the config file
	config, err := ReadConfig(*configPath)
	if err != nil {
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

	for i, step := range config.Steps {
		log.Println("Step:", i)
		dir := "/workspace"
		if step.Dir != "" {
			dir += "/" + step.Dir
		}
		args := []string{
			"run",
			"--rm",
			"--volume", "/var/run/docker.sock:/var/run/docker.sock",
			// --volume /root/.docker:/root/.docker
			"--volume", workspace + ":/workspace",
			"--workdir", dir,
		}
		for key, value := range step.Env {
			args = append(args, "--env", key + "=" + value)
		}
		args = append(args, step.ImageName)
		args = append(args, step.Args...)
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
	// NOTE: Workaround for os.Exit() skips defer.
	os.Exit(cbMain())
}
