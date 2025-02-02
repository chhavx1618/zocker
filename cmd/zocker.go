package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	//"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "zocker",
		Commands: []*cli.Command{
			{
				//create command
				Name:  "create",
				Usage: "Create a new container",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("container name is required")
					}

					containerName := c.Args().First()

					containerDir := fmt.Sprintf("/tmp/%s", containerName)
					err := os.Mkdir(containerDir, 0755)
					if err != nil {
						return fmt.Errorf("failed to create container: %v", err)
					}

					fmt.Printf("Container '%s' created at %s\n", containerName, containerDir)
					return nil
				},
			},
			//run command
			{
				Name:  "run",
				Usage: "Run a command in a container",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no command to run")
					}

					// Capture the command and its arguments
					cmd := c.Args().First()
					args := c.Args().Slice()[1:]

					containerDir := "/tmp/mycontainer"
					if _, err := os.Stat(containerDir); os.IsNotExist(err) {
						return fmt.Errorf("container dir does not exist %s", containerDir)
					}

					// Create and configure the process
					process := exec.Command(cmd, args...)
					process.Stdout = os.Stdout
					process.Stdin = os.Stdin
					process.Stderr = os.Stderr

					process.SysProcAttr = &syscall.SysProcAttr{
						Cloneflags: syscall.CLONE_NEWUTS |
						syscall.CLONE_NEWPID |
						syscall.CLONE_NEWNS |
						syscall.CLONE_NEWIPC,
					}

					//container env

					process.SysProcAttr.Credential = &syscall.Credential{
						Uid: uint32(os.Getuid()),
                        Gid: uint32(os.Getgid()),
					}

					//changing the root filesystem
					process.SysProcAttr.Chroot = containerDir
					process.Dir = "/"

					// Run the command
					err := process.Run()
					if err != nil {
						return fmt.Errorf("failed to run command: %v", err)
					}

					return nil
				},
			},
		},
	}

	// Run the CLI app
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
