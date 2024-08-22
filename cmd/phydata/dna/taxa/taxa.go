// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package taxa implements a command to print the taxa
// with DNA sequences in a PhyData project.
package taxa

import (
	"fmt"
	"os"

	"github.com/js-arias/command"
	"github.com/js-arias/phydata/matrix/dna"
	"github.com/js-arias/phydata/project"
)

var Command = &command.Command{
	Usage: "taxa <project-file>",
	Short: "print taxa",
	Long: `
Command taxa reads a PhyData project and print the list of taxa with
DNA sequences in the project.

The argument of the command is the name of the project-file.
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

	df := p.Path(project.DNA)
	if df == "" {
		return fmt.Errorf("undefined DNA file")
	}
	coll := dna.New()
	if err := readDNAFile(df, coll); err != nil {
		return fmt.Errorf("on project %q: %v", args[0], err)
	}

	for _, tx := range coll.Taxa() {
		fmt.Fprintf(c.Stdout(), "%s\n", tx)
	}

	return nil
}

func readDNAFile(name string, c *dna.Collection) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.ReadTSV(f); err != nil {
		return fmt.Errorf("while reading file %q: %v", name, err)
	}
	return nil
}
