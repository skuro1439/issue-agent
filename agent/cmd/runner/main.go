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

const defaultConfigPath = "./issue_agent.yml"

// Use the docker command to start a container and execute the agent binary
func main() {
	// TODO: input args for issue command

	imageName := "agent-dev"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configPath, err := getConfigPathOrDefault()
	if err != nil {
		panic(err)
	}

	dockerEnvs := passEnvs()
	containerName := "issue-agent"
	args := []string{
		"run",
		"--rm",
		"--name", containerName,
		"-v", configPath + ":" + config.ConfigFilePath,
	}
	args = append(args, dockerEnvs...)
	args = append(args, imageName)
	args = append(args, "agent") // agent binary is built by agent/main.go
	args = append(args, os.Args[1:]...)
	for _, a := range os.Args[1:] {
		if strings.HasSuffix(a, "-config") {
			break
		}
		args = append(args, "-config", configPath)
	}

	cmd := exec.CommandContext(ctx, dockerCmd(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error running container:", err)
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

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

func getConfigPathOrDefault() (string, error) {
	configStart := len(os.Args)
	foundConfig := false
	for i, arg := range os.Args {
		if strings.HasSuffix(arg, "-config") {
			configStart = i
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return filepath.Abs(defaultConfigPath)
	}

	if len(os.Args) <= configStart+1 {
		return "", fmt.Errorf("-config option value is required")
	}

	path := os.Args[configStart+1]
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

// Pass only the environment variables that are required by the agent.
// This is to avoid passing sensitive information to the container.
func passEnvs() []string {
	var passEnvs []string
	for _, env := range os.Environ() {
		envName := strings.Split(env, "=")[0]
		if slices.Contains(cli.EnvNames(), envName) {
			passEnvs = append(passEnvs, env)
		}
	}

	var dockerEnvs []string
	for _, env := range passEnvs {
		varName := strings.Split(env, "=")
		if len(varName) == 2 {
			dockerEnvs = append(dockerEnvs, "-e", env)
		}
	}

	return dockerEnvs
}
