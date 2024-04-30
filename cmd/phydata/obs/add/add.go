// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package add implements a command to add character observations
// to a PhyData project.
package add

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/js-arias/command"
	"github.com/js-arias/phydata/matrix"
	"github.com/js-arias/phydata/project"
)

var Command = &command.Command{
	Usage: `add [-f|--file <obs-file>]
	[--nexus <ref-id>] <project-file> <obs-file>`,
	Short: "add characters observations to a PhyData project",
	Long: `
Command add, read a character observation file, and add the observations to a
PhyData project.

The first argument of the command is the name of the project file. If no
project file exists, a new project will be created.
	
The second argument of the command is the name of the file that contains the
character observations that will be added to the project.
	
By default, the input is expected to be in the form of a tab-delimited
observations file. To import a nexus matrix, use the flag --nexus with an ID
for the reference of the data matrix that will be used as a prefix for
specimen identifiers.
	
By default, the observations will be stored in the observations file currently
defined for the project. If the project does not have an observations file, a
new one will be created with the name 'observations.tab'. A different
observations file name can be defined using the flag --file or -f. If this
file is used and there is an observations file already defined, then a new
file will be created and used as the observations file for the project
(previously defined observations will be preserved).
	`,
	SetFlags: setFlags,
	Run:      run,
}

var obsFile string
var nexusRef string

func setFlags(c *command.Command) {
	c.Flags().StringVar(&obsFile, "file", "", "")
	c.Flags().StringVar(&obsFile, "f", "", "")
	c.Flags().StringVar(&nexusRef, "nexus", "", "")
}

func run(c *command.Command, args []string) error {
	if len(args) < 1 {
		return c.UsageError("expecting project file")
	}
	if len(args) < 2 {
		return c.UsageError("expecting observations file")
	}

	pFile := args[0]
	p, err := openProject(pFile)
	if err != nil {
		return err
	}

	m := matrix.New()
	if mf := p.Path(project.Observations); mf != "" {
		if err := readObsFile(mf, m); err != nil {
			return fmt.Errorf("on project %q: %v", pFile, err)
		}
	}

	in := args[1]
	if nexusRef != "" {
		if err := readNexusFile(in, m, nexusRef); err != nil {
			return err
		}
	} else {
		if err := readObsFile(in, m); err != nil {
			return err
		}
	}

	if obsFile == "" {
		obsFile = p.Path(project.Observations)
		if obsFile == "" {
			obsFile = "observations.tab"
		}
	}
	if err := writeObs(obsFile, m); err != nil {
		return err
	}

	p.Add(project.Observations, obsFile)
	if err := p.Write(pFile); err != nil {
		return err
	}

	return nil
}

func openProject(name string) (*project.Project, error) {
	p, err := project.Read(name)
	if errors.Is(err, os.ErrNotExist) {
		return project.New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable ot open project %q: %v", name, err)
	}
	return p, nil
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

func readNexusFile(name string, m *matrix.Matrix, ref string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := m.ReadNexus(f, ref); err != nil {
		return fmt.Errorf("while reading file %q: %v", name, err)
	}
	return nil
}

func writeObs(name string, m *matrix.Matrix) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	fmt.Fprintf(f, "# phydata: character observations\n")
	fmt.Fprintf(f, "# data saved on: %s\n", time.Now().Format(time.RFC3339))
	if err := m.TSV(f); err != nil {
		return fmt.Errorf("while writing to %q: %v", name, err)
	}
	return nil
}
