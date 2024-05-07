// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package matrix_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/js-arias/phydata/matrix"
)

var nexusMatrix = `#NEXUS

BEGIN TAXA;
 	TITLE Taxa;
	DIMENSIONS NTAX=6;
	TAXLABELS
		Ascaphus_truei
		Bufonidae
		Discoglossidae
		Pipidae
		Ranidae
		Rhinophrynidae
	;
ENDBLOCK;

BEGIN CHARACTERS;
	TITLE 'Phylogenetic data matrix';
	DIMENSIONS NCHAR=5;
	FORMAT DATATYPE = STANDARD RESPECTCASE GAP = - MISSING = ? SYMBOLS = "0 1 2 3 4 5 6 7 8 9 A B C D E F";
	CHARSTATELABELS
		1 'pectoral_girdle' / 'arciferal' 'finnisternal',
		2 'ribs,_fusion' / 'free' 'fused' 'fused_in_adults',
		3 'scapula, relation to clavical' / 'juxtapose' 'overlap',
		4 'tail_muscle' / 'absent' 'present',
		5 'vertebral_ossification' / 'ectochordal' 'holochordal' 'stegochordal' ;
	MATRIX
	Ascaphus_truei	00110
	Bufonidae	01001
	Discoglossidae	00102
	Pipidae	{01}2102
	Ranidae	11001
	Rhinophrynidae	0-100
	;
END;
`

func TestReadNexus(t *testing.T) {
	m := matrix.New()
	if err := m.ReadNexus(strings.NewReader(nexusMatrix), "kluge1969"); err != nil {
		t.Fatalf("unable to read NEXUS data: %v", err)
	}

	want := newMatrix()
	cmpMatrix(t, m, want)
}

func TestWriteNexus(t *testing.T) {
	m := newMatrix()
	var w bytes.Buffer
	if err := m.Nexus(&w); err != nil {
		t.Fatalf("unable to write NEXUS data: %v", err)
	}
	t.Logf("output:\n%s\n", w.String())

	got := matrix.New()
	if err := got.ReadNexus(&w, "kluge1969"); err != nil {
		t.Fatalf("unable to read NEXUS data: %v", err)
	}

	cmpMatrix(t, got, m)
}

var nexusMatrixNoStates = `#NEXUS

BEGIN TAXA;
 	TITLE Taxa;
	DIMENSIONS NTAX=6;
	TAXLABELS
		Ascaphus_truei
		Bufonidae
		Discoglossidae
		Pipidae
		Ranidae
		Rhinophrynidae
	;
END;

BEGIN CHARACTERS;
	TITLE 'Phylogenetic data matrix';
	DIMENSIONS NCHAR=5;
	FORMAT DATATYPE = STANDARD RESPECTCASE GAP = - MISSING = ? SYMBOLS = "0 1 2 3 4 5 6 7 8 9 A B C D E F";
	CHARSTATELABELS
		1 'pectoral_girdle',
		2 'ribs,_fusion',
		3 'scapula, relation to clavical',
		4 'tail_muscle',
		5 'vertebral_ossification';
	MATRIX
	Ascaphus_truei	00110
	Bufonidae	01001
	Discoglossidae	00102
	Pipidae	{01}2102
	Ranidae	11001
	Rhinophrynidae	0-100
	;
END;
`

func TestReadNexusNoStateLabels(t *testing.T) {
	m := matrix.New()
	if err := m.ReadNexus(strings.NewReader(nexusMatrixNoStates), "kluge1969"); err != nil {
		t.Fatalf("unable to read NEXUS data: %v", err)
	}

	chars := []string{"pectoral girdle", "ribs, fusion", "scapula, relation to clavical", "tail muscle", "vertebral ossification"}
	c := m.Chars()
	if !reflect.DeepEqual(c, chars) {
		t.Errorf("characters: got %v, want %v", c, chars)
	}

	states := map[string][]string{
		"tail muscle":                   {"state 0", "state 1"},
		"ribs, fusion":                  {"state 0", "state 1", "state 2"},
		"vertebral ossification":        {"state 0", "state 1", "state 2"},
		"pectoral girdle":               {"state 0", "state 1"},
		"scapula, relation to clavical": {"state 0", "state 1"},
	}
	for c, s := range states {
		states := m.States(c)
		if !reflect.DeepEqual(states, s) {
			t.Errorf("character %q states: got %v, want %v", c, states, s)
		}
	}
}
