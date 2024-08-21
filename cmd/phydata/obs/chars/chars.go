// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package chars implements a command to print the characters
// defined for a PhyData project.
package chars

import (
	"fmt"
	"os"

	"github.com/js-arias/command"
	"github.com/js-arias/phydata/matrix"
	"github.com/js-arias/phydata/project"
)

var Command = &command.Command{
	Usage: "chars <project-file>",
	Short: "print characters",
	Long: `
Command chars reads a PhyData project and print the character names sued for
the observations stored in a PhyData project.

The argument of the command is the name of the project file.
	`,
	Run: run,
}

func run(c *command.Command, args []string) error {
	if len(args) < 1 {
		return c.UsageError("expecting project file")
	}

	p, err := project.Read(args[0])
	if err != nil {
		return fmt.Errorf("unable ot open project %q: %v", args[0], err)
	}

	mf := p.Path(project.Observations)
	if mf == "" {
		return fmt.Errorf("undefined observations file")
	}
	m := matrix.New()
	if err := readObsFile(mf, m); err != nil {
		return fmt.Errorf("on project %q: %v", args[0], err)
	}

	for _, ch := range m.Chars() {
		fmt.Fprintf(c.Stdout(), "%s\n", ch)
	}

	return nil
}

func readObsFile(name string, m *matrix.Matrix) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := m.ReadTSV(f); err != nil {
		return fmt.Errorf("while reading file %q: %v", name, err)
	}
	return nil
}
