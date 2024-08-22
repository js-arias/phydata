// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package add implements a command to add DNA sequences
// to a PhyData project.
package add

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/js-arias/command"
	"github.com/js-arias/phydata/matrix/dna"
	"github.com/js-arias/phydata/project"
)

var Command = &command.Command{
	Usage: `add [-f|--file <dna-file>]
	<project-file> <dna-data-file>`,
	Short: "add DNA sequences to a project",
	Long: `
Command add, read a DNA sequence file, and add the sequences to a PhyData
project.

The first argument of the command is the name of the project file. If no
project file exists, a new project will be created.

The second arguments is the name of the file that contains the DNA sequences
that will be added to the project. The input file must be DNA sequence file.

By default, the DNA data will be stored in the DNA file currently defined for
the project. If the project does not have a DNA file, a ew one will be created
with the name 'dna.tab'. A different DNA file name can be defined using the
flag --file or -f. If this flag is given and there is a DNA file already
defined, then a new file will be created and used as the DNA file for the
project (previously defined DNA sequences will be preserved).
	`,
	SetFlags: setFlags,
	Run:      run,
}

var dnaFile string

func setFlags(c *command.Command) {
	c.Flags().StringVar(&dnaFile, "file", "", "")
	c.Flags().StringVar(&dnaFile, "f", "", "")
}

func run(c *command.Command, args []string) error {
	if len(args) < 1 {
		return c.UsageError("expecting project file")
	}
	if len(args) < 2 {
		return c.UsageError("expecting DNA file")
	}

	pFile := args[0]
	p, err := openProject(pFile)
	if err != nil {
		return err
	}

	coll := dna.New()
	if df := p.Path(project.DNA); df != "" {
		if err := readDNAFile(df, coll); err != nil {
			return fmt.Errorf("on project %q: %v", pFile, err)
		}
	}

	in := args[1]
	nd := dna.New()
	if err := readDNAFile(in, nd); err != nil {
		return err
	}

	for _, tax := range nd.Taxa() {
		for _, spec := range nd.TaxSpec(tax) {
			for _, gene := range nd.SpecGene(spec) {
				for _, acc := range nd.GeneAccession(spec, gene) {
					seq := nd.Sequence(spec, gene, acc)
					if err := coll.Add(tax, spec, gene, acc, seq); err != nil {
						return fmt.Errorf("when adding %q (%s, %s): %v", acc, gene, tax, err)
					}

					alg := nd.Val(spec, gene, acc, dna.Aligned)
					coll.Set(spec, gene, acc, alg, dna.Aligned)
					prt := nd.Val(spec, gene, acc, dna.Protein)
					coll.Set(spec, gene, acc, prt, dna.Protein)
					org := nd.Val(spec, gene, acc, dna.Organelle)
					coll.Set(spec, gene, acc, org, dna.Organelle)
					ref := nd.Val(spec, gene, acc, dna.Reference)
					coll.Set(spec, gene, acc, ref, dna.Reference)
					com := nd.Val(spec, gene, acc, dna.Comments)
					coll.Set(spec, gene, acc, com, dna.Comments)
				}
			}
		}
	}

	if dnaFile == "" {
		dnaFile = p.Path(project.DNA)
		if dnaFile == "" {
			dnaFile = "dna.tab"
		}
	}
	if err := writeDNA(dnaFile, coll); err != nil {
		return err
	}

	p.Add(project.DNA, dnaFile)
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

func writeDNA(name string, c *dna.Collection) (err error) {
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

	fmt.Fprintf(f, "# phydata: DNA sequences\n")
	fmt.Fprintf(f, "# data saved on: %s\n", time.Now().Format(time.RFC3339))
	if err := c.TSV(f); err != nil {
		return fmt.Errorf("while writing to %q: %v", name, err)
	}
	return nil
}
