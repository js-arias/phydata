// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package project_test

import (
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/js-arias/phydata/project"
)

type setPath struct {
	set  project.Dataset
	path string
}

func TestProject(t *testing.T) {
	p := project.New()

	sets := []setPath{
		{project.Observations, "observations.tab"},
		{project.Homologues, "homologues.tab"},
		{project.DNA, "dna.tab"},
	}

	for _, s := range sets {
		p.Add(s.set, s.path)
	}

	name := "tmp-project-for-test.tab"
	defer os.Remove(name)

	if err := p.Write(name); err != nil {
		t.Fatalf("error when writing data: %v", err)
	}

	np, err := project.Read(name)
	if err != nil {
		t.Fatalf("error when reading data: %v", err)
	}
	testProject(t, np, sets)
}

func testProject(t testing.TB, p *project.Project, sets []setPath) {
	t.Helper()

	for _, s := range sets {
		if path := p.Path(s.set); path != s.path {
			t.Errorf("set %s: got path %q, want %q", s.set, path, s.path)
		}
	}
	datasets := make([]project.Dataset, 0, len(sets))
	for _, v := range sets {
		datasets = append(datasets, v.set)
	}
	slices.Sort(datasets)

	if ls := p.Sets(); !reflect.DeepEqual(ls, datasets) {
		t.Errorf("sets: got %v, want %v", ls, datasets)
	}
}
