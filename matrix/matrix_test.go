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

	specs := []string{"ascaphidae:kluge69", "bufonidae:kluge69", "discoglossidae:kluge69", "pipidae:kluge69", "ranidae:kluge69", "rhinophrynidae:kluge69"}
	msp := m.Specimens()
	if !reflect.DeepEqual(msp, specs) {
		t.Errorf("specimens: got %v, want %v", msp, specs)
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
		spec  string
		char  string
		state []string
	}{
		"Ranidae:ossification": {
			spec:  "Ranidae:kluge69",
			char:  "vertebral ossification",
			state: []string{"holochordal"},
		},
		"Ascaphidae:ribs": {
			spec:  "Ascaphidae:kluge69",
			char:  "tail muscle",
			state: []string{"present"},
		},
		"polymorphic": {
			spec:  "Pipidae:kluge69",
			char:  "pectoral girdle",
			state: []string{"arciferal", "finnisternal"},
		},
		"not applicable": {
			spec:  "Rhinophrynidae:kluge69",
			char:  "ribs, fusion",
			state: []string{"<na>"},
		},
		"unknown, unassigned taxon": {
			spec:  "Hylidae:kluge69",
			char:  "ribs, fusion",
			state: []string{"<unknown>"},
		},
		"unknown, unassigned character": {
			spec:  "Bufonidae:kluge69",
			char:  "spiracle",
			state: []string{"<unknown>"},
		},
	}

	for name, test := range tests {
		obs := m.Obs(test.spec, test.char)
		if !reflect.DeepEqual(obs, test.state) {
			t.Errorf("observation %s: got %v, want %v", name, obs, test.state)
		}
	}

	// special cases
	m.Add("Discoglossidae", "Discoglossidae:kluge69", "tail muscle", "<na>")
	obs := m.Obs("Discoglossidae:kluge69", "tail muscle")
	if !reflect.DeepEqual(obs, []string{"<na>"}) {
		t.Errorf("adding <na>: got %v, want %v", obs, []string{"<na>"})
	}

	m.Add("Discoglossidae", "Discoglossidae:kluge69", "tail muscle", "absent")
	obs = m.Obs("Discoglossidae:kluge69", "tail muscle")
	if !reflect.DeepEqual(obs, []string{"absent"}) {
		t.Errorf("remove <na>: got %v, want %v", obs, []string{"absent"})
	}
}

func newMatrix() *matrix.Matrix {
	m := matrix.New()

	m.Add("Ascaphidae", "Ascaphidae:kluge69", "tail muscle", "present")
	m.Add("Ascaphidae", "Ascaphidae:kluge69", "ribs, fusion", "free")
	m.Add("Ascaphidae", "Ascaphidae:kluge69", "vertebral ossification", "ectochordal")
	m.Add("Ascaphidae", "Ascaphidae:kluge69", "pectoral girdle", "arciferal")
	m.Add("Ascaphidae", "Ascaphidae:kluge69", "scapula, relation to clavical", "overlap")
	m.Add("Discoglossidae", "Discoglossidae:kluge69", "tail muscle", "absent")
	m.Add("Discoglossidae", "Discoglossidae:kluge69", "ribs, fusion", "free")
	m.Add("Discoglossidae", "Discoglossidae:kluge69", "vertebral ossification", "stegochordal")
	m.Add("Discoglossidae", "Discoglossidae:kluge69", "pectoral girdle", "arciferal")
	m.Add("Discoglossidae", "Discoglossidae:kluge69", "scapula, relation to clavical", "overlap")
	m.Add("Pipidae", "Pipidae:kluge69", "tail muscle", "absent")
	m.Add("Pipidae", "Pipidae:kluge69", "ribs, fusion", "fused in adults")
	m.Add("Pipidae", "Pipidae:kluge69", "vertebral ossification", "stegochordal")
	m.Add("Pipidae", "Pipidae:kluge69", "pectoral girdle", "arciferal")
	m.Add("Pipidae", "Pipidae:kluge69", "pectoral girdle", "finnisternal")
	m.Add("Pipidae", "Pipidae:kluge69", "scapula, relation to clavical", "overlap")
	m.Add("Rhinophrynidae", "Rhinophrynidae:kluge69", "tail muscle", "absent")
	m.Add("Rhinophrynidae", "Rhinophrynidae:kluge69", "ribs, fusion", "<NA>")
	m.Add("Rhinophrynidae", "Rhinophrynidae:kluge69", "vertebral ossification", "ectochordal")
	m.Add("Rhinophrynidae", "Rhinophrynidae:kluge69", "pectoral girdle", "arciferal")
	m.Add("Rhinophrynidae", "Rhinophrynidae:kluge69", "scapula, relation to clavical", "overlap")
	m.Add("Bufonidae", "Bufonidae:kluge69", "tail muscle", "absent")
	m.Add("Bufonidae", "Bufonidae:kluge69", "ribs, fusion", "fused")
	m.Add("Bufonidae", "Bufonidae:kluge69", "vertebral ossification", "holochordal")
	m.Add("Bufonidae", "Bufonidae:kluge69", "pectoral girdle", "arciferal")
	m.Add("Bufonidae", "Bufonidae:kluge69", "scapula, relation to clavical", "juxtapose")
	m.Add("Ranidae", "Ranidae:kluge69", "tail muscle", "absent")
	m.Add("Ranidae", "Ranidae:kluge69", "ribs, fusion", "fused")
	m.Add("Ranidae", "Ranidae:kluge69", "vertebral ossification", "holochordal")
	m.Add("Ranidae", "Ranidae:kluge69", "pectoral girdle", "finnisternal")
	m.Add("Ranidae", "Ranidae:kluge69", "scapula, relation to clavical", "juxtapose")

	return m
}
