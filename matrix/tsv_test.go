// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package matrix_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/js-arias/phydata/matrix"
)

var obsText = `# character observations
taxon	specimen	character	state	reference	image	comments
Ascaphidae	ascaphidae:kluge69	tail muscle	present	kluge1969	ascaphus-tail.png	it might be not homologous with tail muscles of salamanders
Ascaphidae	ascaphidae:kluge69	ribs, fusion	free	kluge1969		
Ascaphidae	ascaphidae:kluge69	vertebral ossification	ectochordal	kluge1969		
Ascaphidae	ascaphidae:kluge69	pectoral girdle	arciferal	kluge1969		
Ascaphidae	ascaphidae:kluge69	scapula, relation to clavical	overlap	kluge1969		
Discoglossidae	discoglossidae:kluge69	tail muscle	absent	kluge1969		
Discoglossidae	discoglossidae:kluge69	ribs, fusion	free	kluge1969		
Discoglossidae	discoglossidae:kluge69	vertebral ossification	stegochordal	kluge1969		
Discoglossidae	discoglossidae:kluge69	pectoral girdle	arciferal	kluge1969		
Discoglossidae	discoglossidae:kluge69	scapula, relation to clavical	overlap	kluge1969		
Pipidae	pipidae:kluge69	tail muscle	absent	kluge1969		
Pipidae	pipidae:kluge69	ribs, fusion	fused in adults	kluge1969		
Pipidae	pipidae:kluge69	vertebral ossification	stegochordal	kluge1969		
Pipidae	pipidae:kluge69	pectoral girdle	arciferal	kluge1969		
Pipidae	pipidae:kluge69	pectoral girdle	finnisternal	kluge1969		
Pipidae	pipidae:kluge69	scapula, relation to clavical	overlap	kluge1969		
Rhinophrynidae	rhinophrynidae:kluge69	tail muscle	absent	kluge1969		
Rhinophrynidae	rhinophrynidae:kluge69	ribs, fusion	<na>	kluge1969		
Rhinophrynidae	rhinophrynidae:kluge69	vertebral ossification	ectochordal	kluge1969		
Rhinophrynidae	rhinophrynidae:kluge69	pectoral girdle	arciferal	kluge1969		
Rhinophrynidae	rhinophrynidae:kluge69	scapula, relation to clavical	overlap	kluge1969		
Bufonidae	bufonidae:kluge69	tail muscle	absent	kluge1969		
Bufonidae	bufonidae:kluge69	ribs, fusion	fused	kluge1969		
Bufonidae	bufonidae:kluge69	vertebral ossification	holochordal	kluge1969		
Bufonidae	bufonidae:kluge69	pectoral girdle	arciferal	kluge1969		
Bufonidae	bufonidae:kluge69	scapula, relation to clavical	juxtapose	kluge1969		
Ranidae	ranidae:kluge69	tail muscle	absent	kluge1969		
Ranidae	ranidae:kluge69	ribs, fusion	fused	kluge1969		
Ranidae	ranidae:kluge69	vertebral ossification	holochordal	kluge1969		
Ranidae	ranidae:kluge69	pectoral girdle	finnisternal	kluge1969		
Ranidae	ranidae:kluge69	scapula, relation to clavical	juxtapose	kluge1969		
`

func TestReadTSV(t *testing.T) {
	m := matrix.New()
	if err := m.ReadTSV(strings.NewReader(obsText)); err != nil {
		t.Fatalf("unable to read TSV data: %v", err)
	}

	want := newMatrix()
	cmpMatrix(t, m, want)
}

func TestWriteTSV(t *testing.T) {
	m := newMatrix()
	var w bytes.Buffer
	if err := m.TSV(&w); err != nil {
		t.Fatalf("unable to write TSV data: %v", err)
	}
	t.Logf("output:\n%s\n", w.String())

	got := matrix.New()
	if err := got.ReadTSV(&w); err != nil {
		t.Fatalf("unable to read TSV data: %v", err)
	}

	cmpMatrix(t, got, m)
}
