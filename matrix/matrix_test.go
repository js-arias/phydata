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

	specs := []string{"kluge1969:ascaphidae", "kluge1969:bufonidae", "kluge1969:discoglossidae", "kluge1969:pipidae", "kluge1969:ranidae", "kluge1969:rhinophrynidae"}
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
			spec:  "kluge1969:Ranidae",
			char:  "vertebral ossification",
			state: []string{"holochordal"},
		},
		"Ascaphidae:ribs": {
			spec:  "kluge1969:Ascaphidae",
			char:  "tail muscle",
			state: []string{"present"},
		},
		"polymorphic": {
			spec:  "kluge1969:Pipidae",
			char:  "pectoral girdle",
			state: []string{"arciferal", "finnisternal"},
		},
		"not applicable": {
			spec:  "kluge1969:Rhinophrynidae",
			char:  "ribs, fusion",
			state: []string{"<na>"},
		},
		"unknown, unassigned taxon": {
			spec:  "kluge1969:Hylidae",
			char:  "ribs, fusion",
			state: []string{"<unknown>"},
		},
		"unknown, unassigned character": {
			spec:  "kluge1969:Bufonidae",
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
		"Ascaphidae":     {"kluge1969:ascaphidae"},
		"Discoglossidae": {"kluge1969:discoglossidae"},
		"Pipidae":        {"kluge1969:pipidae"},
		"Rhinophrynidae": {"kluge1969:rhinophrynidae"},
		"Bufonidae":      {"kluge1969:bufonidae"},
		"Ranidae":        {"kluge1969:ranidae"},
	}
	for tn, txSp := range taxSpec {
		if sp := m.TaxSpec(tn); !reflect.DeepEqual(sp, txSp) {
			t.Errorf("specimens of %q: got %v, want %v", tn, sp, txSp)
		}

	}

	// special cases
	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "tail muscle", "<na>")
	obs := m.Obs("kluge1969:Discoglossidae", "tail muscle")
	if !reflect.DeepEqual(obs, []string{"<na>"}) {
		t.Errorf("adding <na>: got %v, want %v", obs, []string{"<na>"})
	}

	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "tail muscle", "absent")
	obs = m.Obs("kluge1969:Discoglossidae", "tail muscle")
	if !reflect.DeepEqual(obs, []string{"absent"}) {
		t.Errorf("remove <na>: got %v, want %v", obs, []string{"absent"})
	}
}

func newMatrix() *matrix.Matrix {
	m := matrix.New()

	m.Add("Ascaphidae", "kluge1969:Ascaphidae", "tail muscle", "present")
	m.Add("Ascaphidae", "kluge1969:Ascaphidae", "ribs, fusion", "free")
	m.Add("Ascaphidae", "kluge1969:Ascaphidae", "vertebral ossification", "ectochordal")
	m.Add("Ascaphidae", "kluge1969:Ascaphidae", "pectoral girdle", "arciferal")
	m.Add("Ascaphidae", "kluge1969:Ascaphidae", "scapula, relation to clavical", "overlap")
	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "tail muscle", "absent")
	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "ribs, fusion", "free")
	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "vertebral ossification", "stegochordal")
	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "pectoral girdle", "arciferal")
	m.Add("Discoglossidae", "kluge1969:Discoglossidae", "scapula, relation to clavical", "overlap")
	m.Add("Pipidae", "kluge1969:Pipidae", "tail muscle", "absent")
	m.Add("Pipidae", "kluge1969:Pipidae", "ribs, fusion", "fused in adults")
	m.Add("Pipidae", "kluge1969:Pipidae", "vertebral ossification", "stegochordal")
	m.Add("Pipidae", "kluge1969:Pipidae", "pectoral girdle", "arciferal")
	m.Add("Pipidae", "kluge1969:Pipidae", "pectoral girdle", "finnisternal")
	m.Add("Pipidae", "kluge1969:Pipidae", "scapula, relation to clavical", "overlap")
	m.Add("Rhinophrynidae", "kluge1969:Rhinophrynidae", "tail muscle", "absent")
	m.Add("Rhinophrynidae", "kluge1969:Rhinophrynidae", "ribs, fusion", "<NA>")
	m.Add("Rhinophrynidae", "kluge1969:Rhinophrynidae", "vertebral ossification", "ectochordal")
	m.Add("Rhinophrynidae", "kluge1969:Rhinophrynidae", "pectoral girdle", "arciferal")
	m.Add("Rhinophrynidae", "kluge1969:Rhinophrynidae", "scapula, relation to clavical", "overlap")
	m.Add("Bufonidae", "kluge1969:Bufonidae", "tail muscle", "absent")
	m.Add("Bufonidae", "kluge1969:Bufonidae", "ribs, fusion", "fused")
	m.Add("Bufonidae", "kluge1969:Bufonidae", "vertebral ossification", "holochordal")
	m.Add("Bufonidae", "kluge1969:Bufonidae", "pectoral girdle", "arciferal")
	m.Add("Bufonidae", "kluge1969:Bufonidae", "scapula, relation to clavical", "juxtapose")
	m.Add("Ranidae", "kluge1969:Ranidae", "tail muscle", "absent")
	m.Add("Ranidae", "kluge1969:Ranidae", "ribs, fusion", "fused")
	m.Add("Ranidae", "kluge1969:Ranidae", "vertebral ossification", "holochordal")
	m.Add("Ranidae", "kluge1969:Ranidae", "pectoral girdle", "finnisternal")
	m.Add("Ranidae", "kluge1969:Ranidae", "scapula, relation to clavical", "juxtapose")

	m.Set("kluge1969:Ascaphidae", "tail muscle", "present", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ascaphidae", "ribs, fusion", "free", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ascaphidae", "vertebral ossification", "ectochordal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ascaphidae", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ascaphidae", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Discoglossidae", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Discoglossidae", "ribs, fusion", "free", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Discoglossidae", "vertebral ossification", "stegochordal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Discoglossidae", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Discoglossidae", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Pipidae", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Pipidae", "ribs, fusion", "fused in adults", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Pipidae", "vertebral ossification", "stegochordal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Pipidae", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Pipidae", "pectoral girdle", "finnisternal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Pipidae", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Rhinophrynidae", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Rhinophrynidae", "ribs, fusion", "<NA>", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Rhinophrynidae", "vertebral ossification", "ectochordal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Rhinophrynidae", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Rhinophrynidae", "scapula, relation to clavical", "overlap", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Bufonidae", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Bufonidae", "ribs, fusion", "fused", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Bufonidae", "vertebral ossification", "holochordal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Bufonidae", "pectoral girdle", "arciferal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Bufonidae", "scapula, relation to clavical", "juxtapose", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ranidae", "tail muscle", "absent", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ranidae", "ribs, fusion", "fused", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ranidae", "vertebral ossification", "holochordal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ranidae", "pectoral girdle", "finnisternal", "kluge1969", matrix.Reference)
	m.Set("kluge1969:Ranidae", "scapula, relation to clavical", "juxtapose", "kluge1969", matrix.Reference)

	return m
}

func newMatrixWithComments() *matrix.Matrix {
	m := newMatrix()
	m.Set("kluge1969:Ascaphidae", "tail muscle", "present", "ascaphus-tail.png", matrix.ImageLink)
	m.Set("kluge1969:Ascaphidae", "tail muscle", "present", "it might be not homologous with tail muscles of salamanders", matrix.Comments)

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
