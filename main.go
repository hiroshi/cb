package main

import (
	"encoding/json"
	"flag"
	// "fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Step struct {
	ImageName string `json:"name"`
	Args []string
}

type Config struct {
	Steps []Step
}

func cbMain() (exitCode int) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Take arguments and options
	source := os.Args[1]
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

	// ##  Expand source to workspace
	// create workspace
	// NOTE: On macOS default prefix start with "/var/folders/", but default File sharing option of Docker for Mac don't allow to mount /var, but /tmp
	workspace, err := ioutil.TempDir("/tmp", "cb_workspace")
	if err != nil {
		log.Print(err)
		return 1
	}
	defer func() {
		log.Println("cleanup:", workspace)
		if err := os.RemoveAll(workspace); err != nil {
			log.Print(err)
		}
	}()
	log.Println("workspace:", workspace)
	// get abs path to source
	sourcePath, err := filepath.Abs(source)
	if err != nil {
		log.Print(err)
		return 1
	}
	// cd workspace
	if err:= os.Chdir(workspace); err != nil {
		log.Print(err)
		return 1
	}
	// expand source
	cmd := exec.Command("tar", "xzvf", sourcePath)
	log.Println(cmd.Args)
	out, err := cmd.CombinedOutput()
	log.Printf("%s", out)
	if err != nil {
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
