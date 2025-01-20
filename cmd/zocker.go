package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bocker",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run a command in a container",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no command to run")
					}

					cmd := c.Args().First()
					args := c.Args().Slice()[1:] // Define the args variable properly

					process := exec.Command(cmd, args...)
					process.Stdout = os.Stdout
					process.Stdin = os.Stdin
					process.Stderr = os.Stderr

					err := process.Run()
					if err != nil {
						return fmt.Errorf("failed: %v", err) // Corrected error message format
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err) // Corrected this as well
		os.Exit(1)
	}
}
