package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "zocker",
		Usage: "A simple container runtime",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a new container",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("container name is required")
					}

					containerName := c.Args().First()
					containerDir := fmt.Sprintf("/tmp/zocker/%s", containerName)

					// Check if container directory already exists
					if _, err := os.Stat(containerDir); !os.IsNotExist(err) {
						return fmt.Errorf("container '%s' already exists", containerName)
					}

					// Create the container directory
					err := os.MkdirAll(containerDir, 0755)
					if err != nil {
						return fmt.Errorf("failed to create container directory: %v", err)
					}

					// Create the subvolume (or a placeholder for the subvolume in this case)
					subvolumeDir := filepath.Join(containerDir, "subvolume")
					err = os.Mkdir(subvolumeDir, 0755)
					if err != nil {
						return fmt.Errorf("failed to create subvolume directory: %v", err)
					}

					// Initialize the container with a simple file structure
					initFile := filepath.Join(subvolumeDir, "init.txt")
					f, err := os.Create(initFile)
					if err != nil {
						return fmt.Errorf("failed to create init file: %v", err)
					}
					defer f.Close()

					// Write a simple message in the init file
					_, err = f.WriteString("Container initialized\n")
					if err != nil {
						return fmt.Errorf("failed to write to init file: %v", err)
					}

					fmt.Printf("Container '%s' created at %s\n", containerName, containerDir)
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run a command in a container",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("container name is required")
					}

					// containerName := c.Args().First()
					// containerDir := fmt.Sprintf("/tmp/zocker/%s", containerName)

					// // Check if the container directory exists
					// if _, err := os.Stat(containerDir); os.IsNotExist(err) {
					// 	return fmt.Errorf("container directory does not exist: %s", containerDir)
					// }

					// // Ensure the subvolume exists (container initialization)
					// subvolumeDir := filepath.Join(containerDir, "subvolume")
					// if _, err := os.Stat(subvolumeDir); os.IsNotExist(err) {
					// 	return fmt.Errorf("container subvolume does not exist: %s", subvolumeDir)
					// }

					// Default to "bash" if no command is provided
					// cmd := "/bin/bash"
					// args := []string{}
					// if c.Args().Len() > 1 {
					// 	cmd = c.Args().Get(1)
					// 	args = c.Args().Slice()[2:]
					// }

					cmd := c.Args().First()
					args := c.Args().Tail()

					// // Mount the container root
					// mntFlags := syscall.MS_REC | syscall.MS_PRIVATE
					// if err := syscall.Mount(subvolumeDir, "/mnt/container", "", uintptr(mntFlags), ""); err != nil {
					// 	return fmt.Errorf("failed to mount container root: %v", err)
					// }

					// Create and configure the process to run in the new namespace
					process := exec.Command(cmd, args...)
					process.Stdout = os.Stdout
					process.Stdin = os.Stdin
					process.Stderr = os.Stderr

					process.SysProcAttr = &syscall.SysProcAttr{
						Cloneflags: syscall.CLONE_NEWUTS |
							syscall.CLONE_NEWPID |
							syscall.CLONE_NEWNS |
							syscall.CLONE_NEWIPC,
						Credential: &syscall.Credential{
							Uid: uint32(os.Getuid()),
							Gid: uint32(os.Getgid()),
						},
					}

					// Run the command
					err := process.Run()
					if err != nil {
						return fmt.Errorf("failed to run command '%s': %v", cmd, err)
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
