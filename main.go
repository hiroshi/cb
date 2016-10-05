package main

import (
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

	// ##  Expand source to workspace
	// create workspace
	// NOTE: On macOS default prefix start with "/var/folders/", but default File sharing option of Docker for Mac don't allow to mount /var, but /tmp
	// workspace, err := ioutil.TempDir("/tmp", "cb_workspace")
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// defer func() {
	// 	log.Println("cleanup:", workspace)
	// 	if err := os.RemoveAll(workspace); err != nil {
	// 		log.Print(err)
	// 	}
	// }()
	// log.Println("workspace:", workspace)
	// // get abs path to source
	// sourcePath, err := filepath.Abs(source)
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// // cd workspace
	// if err:= os.Chdir(workspace); err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// // expand source
	// cmd := exec.Command("tar", "xzvf", sourcePath)
	// log.Println(cmd.Args)
	// out, err := cmd.CombinedOutput()
	// log.Printf("%s", out)
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }

	// Create workspace volume
	workspace, err := DockerOutput("volume", "create")
	if err != nil {
		log.Print(err)
		return 1
	}
	log.Println("workspace:", workspace)
	defer func() {
		_, err := DockerOutput("volume", "rm", workspace)
		if err != nil {
			log.Print(err)
		}
	}()
	// Create a container to be destination of `docker cp`
	cmd := exec.Command("docker", "create", "--volume", workspace + ":/workspace", "busybox")
	log.Println(cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		log.Print(err)
		return 1
	}
	// lines := bytes.Split(out, []byte("\n"))
	// log.Println(lines)
	// log.Println(string(out))
	log.Println(strings.Split(string(out), "\n"))

	
	// buf := bytes.NewBuffer(out)

	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// if err := cmd.Start(); err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// scanner := bufio.NewScanner(stdout)
	// ch := make(chan string)
	// go func() {
	// 	for scanner.Scan() {}
	// 	ch <- scanner.Text()
	// }()
	// if err := scanner.Err(); err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// container := <- ch
	// log.Println("container:", container)

	// out, err := cmd.Output()
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	// buf := bytes.NewBuffer(out)

	return 1

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
