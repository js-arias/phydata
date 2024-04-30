// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package project

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"
)

// Dataset is a keyword to identify
// the type of dataset file in a project.
type Dataset string

// Valid dataset types
const (
	// File with an hierarchy of homologues.
	Homologues Dataset = "homologues"

	// File for specimen character observations.
	Observations Dataset = "observations"
)

// A Project represents a collection of paths
// for particular datasets.
type Project struct {
	paths map[Dataset]string
}

// New creates a new empty project.
func New() *Project {
	return &Project{
		paths: make(map[Dataset]string),
	}
}

var header = []string{
	"dataset",
	"path",
}

// Read reads a project file from a TSV file.
//
// The TSV must contain the following fields:
//
//   - dataset, for the kind of file
//   - path, for the path of the file
//
// Here is an example file:
//
//	# phydata project files
//	dataset	path
//	homologues	homologues.tab
//	observations	observations.tab
func Read(name string) (*Project, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	tsv := csv.NewReader(f)
	tsv.Comma = '\t'
	tsv.Comment = '#'

	head, err := tsv.Read()
	if err != nil {
		return nil, fmt.Errorf("on file %q: header: %v", name, err)
	}
	fields := make(map[string]int, len(head))
	for i, h := range head {
		h = strings.ToLower(h)
		fields[h] = i
	}
	for _, h := range header {
		if _, ok := fields[h]; !ok {
			return nil, fmt.Errorf("on file %q: expecting field %q", name, h)
		}
	}

	p := New()
	for {
		row, err := tsv.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		ln, _ := tsv.FieldPos(0)
		if err != nil {
			return nil, fmt.Errorf("on file %q: on row %d: %v", name, ln, err)
		}

		f := "dataset"
		s := Dataset(row[fields[f]])

		f = "path"
		path := row[fields[f]]
		p.paths[s] = path
	}

	return p, nil
}

// Add adds a filepath to a dataset to a given project.
// It returns the previous value
// for the dataset.
func (p *Project) Add(set Dataset, path string) string {
	prev := p.paths[set]
	if path == "" {
		delete(p.paths, set)
		return prev
	}

	p.paths[set] = path
	return prev
}

// Path returns the path of the given dataset.
func (p *Project) Path(set Dataset) string {
	return p.paths[set]
}

// Sets returns the datasets defined on a project.
func (p *Project) Sets() []Dataset {
	var sets []Dataset
	for s := range p.paths {
		sets = append(sets, s)
	}
	slices.Sort(sets)
	return sets
}

// Write writes a project into a file
// with the indicated name.
func (p *Project) Write(name string) (err error) {
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

	bw := bufio.NewWriter(f)
	fmt.Fprintf(bw, "# phydata project files\n")
	fmt.Fprintf(bw, "# data save on: %s\n", time.Now().Format(time.RFC3339))
	tsv := csv.NewWriter(bw)
	tsv.Comma = '\t'
	tsv.UseCRLF = true

	if err := tsv.Write(header); err != nil {
		return fmt.Errorf("on file %q: while writing header: %v", name, err)
	}

	sets := p.Sets()
	for _, s := range sets {
		row := []string{
			string(s),
			p.paths[s],
		}
		if err := tsv.Write(row); err != nil {
			return fmt.Errorf("on file %q: %v", name, err)
		}
	}

	tsv.Flush()
	if err := tsv.Error(); err != nil {
		return fmt.Errorf("on file %q: while writing data: %v", name, err)
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("on file %q: while writing data: %v", name, err)
	}
	return nil
}
