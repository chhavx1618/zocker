package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

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

					// Create and configure the process
					process := exec.Command(cmd, args...)
					process.Stdout = os.Stdout
					process.Stdin = os.Stdin
					process.Stderr = os.Stderr

					process.SysProcAttr = &syscall.SysProcAttr{
						Cloneflags: syscall.CLONE_NEWUTS,
					}

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
