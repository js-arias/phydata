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
	taxa := []string{"Ascaphidae", "Bufonidae", "Discoglossidae", "Pipidae", "Ranidae", "Rhinophrynidae"}
	if tx := m.Taxa(); !reflect.DeepEqual(tx, taxa) {
		t.Errorf("taxa: got %v, want %v", tx, taxa)
	}
	taxSpec := map[string][]string{
		"Ascaphidae":     {"ascaphidae:kluge69"},
		"Discoglossidae": {"discoglossidae:kluge69"},
		"Pipidae":        {"pipidae:kluge69"},
		"Rhinophrynidae": {"rhinophrynidae:kluge69"},
		"Bufonidae":      {"bufonidae:kluge69"},
		"Ranidae":        {"ranidae:kluge69"},
	}
	for tn, txSp := range taxSpec {
		if sp := m.TaxSpec(tn); !reflect.DeepEqual(sp, txSp) {
			t.Errorf("specimens of %q: got %v, want %v", tn, sp, txSp)
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

	m.Set("Ascaphidae:kluge69", "tail muscle", "present", "kluge1969", matrix.Reference)
	m.Set("Ascaphidae:kluge69", "ribs, fusion", "free", "kluge1969", matrix.Reference)
	m.Set("Ascaphidae:kluge69", "vertebral ossification", "ectochordal", "kluge1969", matrix.Reference)
	m.Set("Ascaphidae:kluge69", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("Ascaphidae:kluge69", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("Discoglossidae:kluge69", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("Discoglossidae:kluge69", "ribs, fusion", "free", "kluge1969", matrix.Reference)
	m.Set("Discoglossidae:kluge69", "vertebral ossification", "stegochordal", "kluge1969", matrix.Reference)
	m.Set("Discoglossidae:kluge69", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("Discoglossidae:kluge69", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("Pipidae:kluge69", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("Pipidae:kluge69", "ribs, fusion", "fused in adults", "kluge1969", matrix.Reference)
	m.Set("Pipidae:kluge69", "vertebral ossification", "stegochordal", "kluge1969", matrix.Reference)
	m.Set("Pipidae:kluge69", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("Pipidae:kluge69", "pectoral girdle", "finnisternal", "kluge1969", matrix.Reference)
	m.Set("Pipidae:kluge69", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("Rhinophrynidae:kluge69", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("Rhinophrynidae:kluge69", "ribs, fusion", "<NA>", "kluge1969", matrix.Reference)
	m.Set("Rhinophrynidae:kluge69", "vertebral ossification", "ectochordal", "kluge1969", matrix.Reference)
	m.Set("Rhinophrynidae:kluge69", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("Rhinophrynidae:kluge69", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("Bufonidae:kluge69", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("Bufonidae:kluge69", "ribs, fusion", "fused", "kluge1969", matrix.Reference)
	m.Set("Bufonidae:kluge69", "vertebral ossification", "holochordal", "kluge1969", matrix.Reference)
	m.Set("Bufonidae:kluge69", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("Bufonidae:kluge69", "scapula, relation to clavical", "juxtapose", "kluge1969", matrix.Reference)
	m.Set("Ranidae:kluge69", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("Ranidae:kluge69", "ribs, fusion", "fused", "kluge1969", matrix.Reference)
	m.Set("Ranidae:kluge69", "vertebral ossification", "holochordal", "kluge1969", matrix.Reference)
	m.Set("Ranidae:kluge69", "pectoral girdle", "finnisternal", "kluge1969", matrix.Reference)
	m.Set("Ranidae:kluge69", "scapula, relation to clavical", "juxtapose", "kluge1969", matrix.Reference)

	return m
}

func newMatrixWithComments() *matrix.Matrix {
	m := newMatrix()
	m.Set("Ascaphidae:kluge69", "tail muscle", "present", "ascaphus-tail.png", matrix.ImageLink)
	m.Set("Ascaphidae:kluge69", "tail muscle", "present", "it might be not homologous with tail muscles of salamanders", matrix.Comments)

	return m
}

func cmpMatrix(t testing.TB, got, want *matrix.Matrix) {
	t.Helper()

	specs := want.Specimens()
	sp := got.Specimens()
	if !reflect.DeepEqual(sp, specs) {
		t.Errorf("specimens: got %v, want %v", sp, specs)
	}

	chars := want.Chars()
	c := got.Chars()
	if !reflect.DeepEqual(c, chars) {
		t.Errorf("characters: got %v, want %v", c, chars)
	}

	for _, cn := range chars {
		states := want.States(cn)
		s := got.States(cn)
		if !reflect.DeepEqual(s, states) {
			t.Errorf("character %q states: got %v, want %v", cn, s, states)
		}
	}

	fields := []matrix.Field{matrix.Reference, matrix.ImageLink, matrix.Comments}

	for _, sn := range specs {
		for _, cn := range chars {
			obs := want.Obs(sn, cn)
			o := got.Obs(sn, cn)
			if !reflect.DeepEqual(o, obs) {
				t.Errorf("observation %s-%s: got %v, want %v", sn, cn, o, obs)
			}

			for _, s := range obs {
				for _, f := range fields {
					val := want.Val(sn, cn, s, f)
					v := got.Val(sn, cn, s, f)
					if v != val {
						t.Errorf("value %s-%s-%s [%q]: got %q, want %q", sn, cn, s, f, v, val)
					}
				}
			}
		}
	}

	taxa := want.Taxa()
	if tx := got.Taxa(); !reflect.DeepEqual(tx, taxa) {
		t.Errorf("taxa: got %v, want %v", tx, taxa)
	}

	for _, tax := range taxa {
		txSp := got.TaxSpec(tax)
		if sp := got.TaxSpec(tax); !reflect.DeepEqual(sp, txSp) {
			t.Errorf("specimens of %q: got %v, want %v", tax, sp, txSp)
		}
	}
}
