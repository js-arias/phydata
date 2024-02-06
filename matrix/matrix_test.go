// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package matrix_test

import (
	"reflect"
	"testing"

	"github.com/js-arias/phydata/matrix"
)

func TestMatrix(t *testing.T) {
	m := newMatrix()

	terms := []string{"Ascaphidae", "Bufonidae", "Discoglossidae", "Pipidae", "Ranidae", "Rhinophrynidae"}
	mt := m.Terminals()
	if !reflect.DeepEqual(mt, terms) {
		t.Errorf("terminals: got %v, want %v", mt, terms)
	}

	chars := []string{"pectoral girdle", "ribs, fusion", "scapula, relation to clavical", "tail muscle", "vertebral ossification"}
	c := m.Chars()
	if !reflect.DeepEqual(c, chars) {
		t.Errorf("characters: got %v, want %v", c, chars)
	}

	states := map[string][]string{
		"tail muscle":                   {"absent", "present"},
		"ribs, fusion":                  {"free", "fused", "fused in adults"},
		"vertebral ossification":        {"ectochordal", "holochordal", "stegochordal"},
		"pectoral girdle":               {"arciferal", "finnisternal"},
		"scapula, relation to clavical": {"juxtapose", "overlap"},
	}
	for c, s := range states {
		states := m.States(c)
		if !reflect.DeepEqual(states, s) {
			t.Errorf("character %q states: got %v, want %v", c, states, s)
		}
	}

	tests := map[string]struct {
		term  string
		char  string
		state []string
	}{
		"Ranidae:ossification": {
			term:  "Ranidae",
			char:  "vertebral ossification",
			state: []string{"holochordal"},
		},
		"Ascaphidae:ribs": {
			term:  "Ascaphidae",
			char:  "tail muscle",
			state: []string{"present"},
		},
		"polymorphic": {
			term:  "Pipidae",
			char:  "pectoral girdle",
			state: []string{"arciferal", "finnisternal"},
		},
		"not applicable": {
			term:  "Rhinophrynidae",
			char:  "ribs, fusion",
			state: []string{"<na>"},
		},
		"unknown, unassigned taxon": {
			term:  "Hylidae",
			char:  "ribs, fusion",
			state: []string{"<unknown>"},
		},
		"unknown, unassigned character": {
			term:  "Bufonidae",
			char:  "spiracle",
			state: []string{"<unknown>"},
		},
	}

	for name, test := range tests {
		obs := m.Obs(test.term, test.char)
		if !reflect.DeepEqual(obs, test.state) {
			t.Errorf("observation %s: got %v, want %v", name, obs, test.state)
		}
	}

	// special cases
	m.Add("Discoglossidae", "tail muscle", "<na>")
	obs := m.Obs("Discoglossidae", "tail muscle")
	if !reflect.DeepEqual(obs, []string{"<na>"}) {
		t.Errorf("adding <na>: got %v, want %v", obs, []string{"<na>"})
	}

	m.Add("Discoglossidae", "tail muscle", "absent")
	obs = m.Obs("Discoglossidae", "tail muscle")
	if !reflect.DeepEqual(obs, []string{"absent"}) {
		t.Errorf("remove <na>: got %v, want %v", obs, []string{"absent"})
	}
}

func newMatrix() *matrix.Matrix {
	m := matrix.New()

	m.Add("Ascaphidae", "tail muscle", "present")
	m.Add("Ascaphidae", "ribs, fusion", "free")
	m.Add("Ascaphidae", "vertebral ossification", "ectochordal")
	m.Add("Ascaphidae", "pectoral girdle", "arciferal")
	m.Add("Ascaphidae", "scapula, relation to clavical", "overlap")
	m.Add("Discoglossidae", "tail muscle", "absent")
	m.Add("Discoglossidae", "ribs, fusion", "free")
	m.Add("Discoglossidae", "vertebral ossification", "stegochordal")
	m.Add("Discoglossidae", "pectoral girdle", "arciferal")
	m.Add("Discoglossidae", "scapula, relation to clavical", "overlap")
	m.Add("Pipidae", "tail muscle", "absent")
	m.Add("Pipidae", "ribs, fusion", "fused in adults")
	m.Add("Pipidae", "vertebral ossification", "stegochordal")
	m.Add("Pipidae", "pectoral girdle", "arciferal")
	m.Add("Pipidae", "pectoral girdle", "finnisternal")
	m.Add("Pipidae", "scapula, relation to clavical", "overlap")
	m.Add("Rhinophrynidae", "tail muscle", "absent")
	m.Add("Rhinophrynidae", "ribs, fusion", "<NA>")
	m.Add("Rhinophrynidae", "vertebral ossification", "ectochordal")
	m.Add("Rhinophrynidae", "pectoral girdle", "arciferal")
	m.Add("Rhinophrynidae", "scapula, relation to clavical", "overlap")
	m.Add("Bufonidae", "tail muscle", "absent")
	m.Add("Bufonidae", "ribs, fusion", "fused")
	m.Add("Bufonidae", "vertebral ossification", "holochordal")
	m.Add("Bufonidae", "pectoral girdle", "arciferal")
	m.Add("Bufonidae", "scapula, relation to clavical", "juxtapose")
	m.Add("Ranidae", "tail muscle", "absent")
	m.Add("Ranidae", "ribs, fusion", "fused")
	m.Add("Ranidae", "vertebral ossification", "holochordal")
	m.Add("Ranidae", "pectoral girdle", "finnisternal")
	m.Add("Ranidae", "scapula, relation to clavical", "juxtapose")

	return m
}
