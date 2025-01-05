package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github/clover0/github-issue-agent/cli"
	"github/clover0/github-issue-agent/config"
)

func main() {
	// TODO: input args for issue command

	imageName := "agent-dev"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var passEnvs []string
	for _, env := range os.Environ() {
		envName := strings.Split(env, "=")[0]
		if slices.Contains(cli.EnvNames(), envName) {
			passEnvs = append(passEnvs, env)
		}
	}
	configPath, err := GetConfigPath()
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	var dockerEnvs []string
	for _, env := range passEnvs {
		varName := strings.Split(env, "=")
		if len(varName) == 2 {
			dockerEnvs = append(dockerEnvs, "-e", env)
		}
	}
	containerName := "issue-agent"
	args := []string{
		"run",
		"--rm",
		"--name", containerName,
		"-v", configPath + ":" + config.ConfigFilePath,
	}
	args = append(args, dockerEnvs...)
	args = append(args, imageName)
	args = append(args, "agent")
	args = append(args, os.Args[1:]...)
	//args = append(args, "echo", "hello")

	cmd := exec.CommandContext(ctx, dockerCmd(), args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error running container:", err)
		panic(err)
	}

	go func(containerName string) {
		sig := <-sigChan
		fmt.Println("Received signal")
		if err := cmd.Process.Signal(sig); err != nil {
			fmt.Println("Error sending signal to container:", err)
		}
		stopContainer(containerName)
		cancel()
	}(containerName)

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			fmt.Printf("Process exited with error: %v\n", err)
		}
	}
}

func dockerCmd() string {
	com, ok := os.LookupEnv("_DOCKER_CMD")
	if ok {
		return com
	}
	return "docker"
}

func stopContainer(containerName string) {
	cmd := exec.Command(dockerCmd(), "kill", containerName)
	bytes, _ := cmd.CombinedOutput()
	fmt.Println(string(bytes))
}

func GetConfigPath() (string, error) {
	configStart := len(os.Args)
	for i, arg := range os.Args {
		if strings.HasSuffix(arg, "-config") {
			configStart = i
			break
		}
	}

	if len(os.Args) <= configStart+1 {
		return "", fmt.Errorf("-config option is required")
	}

	path := os.Args[configStart+1]
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
