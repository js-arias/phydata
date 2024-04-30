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
Ascaphidae	kluge1969:ascaphidae	tail muscle	present	kluge1969	ascaphus-tail.png	it might be not homologous with tail muscles of salamanders
Ascaphidae	kluge1969:ascaphidae	ribs, fusion	free	kluge1969		
Ascaphidae	kluge1969:ascaphidae	vertebral ossification	ectochordal	kluge1969		
Ascaphidae	kluge1969:ascaphidae	pectoral girdle	arciferal	kluge1969		
Ascaphidae	kluge1969:ascaphidae	scapula, relation to clavical	overlap	kluge1969		
Discoglossidae	kluge1969:discoglossidae	tail muscle	absent	kluge1969		
Discoglossidae	kluge1969:discoglossidae	ribs, fusion	free	kluge1969		
Discoglossidae	kluge1969:discoglossidae	vertebral ossification	stegochordal	kluge1969		
Discoglossidae	kluge1969:discoglossidae	pectoral girdle	arciferal	kluge1969		
Discoglossidae	kluge1969:discoglossidae	scapula, relation to clavical	overlap	kluge1969		
Pipidae	kluge1969:pipidae	tail muscle	absent	kluge1969		
Pipidae	kluge1969:pipidae	ribs, fusion	fused in adults	kluge1969		
Pipidae	kluge1969:pipidae	vertebral ossification	stegochordal	kluge1969		
Pipidae	kluge1969:pipidae	pectoral girdle	arciferal	kluge1969		
Pipidae	kluge1969:pipidae	pectoral girdle	finnisternal	kluge1969		
Pipidae	kluge1969:pipidae	scapula, relation to clavical	overlap	kluge1969		
Rhinophrynidae	kluge1969:rhinophrynidae	tail muscle	absent	kluge1969		
Rhinophrynidae	kluge1969:rhinophrynidae	ribs, fusion	<na>	kluge1969		
Rhinophrynidae	kluge1969:rhinophrynidae	vertebral ossification	ectochordal	kluge1969		
Rhinophrynidae	kluge1969:rhinophrynidae	pectoral girdle	arciferal	kluge1969		
Rhinophrynidae	kluge1969:rhinophrynidae	scapula, relation to clavical	overlap	kluge1969		
Bufonidae	kluge1969:bufonidae	tail muscle	absent	kluge1969		
Bufonidae	kluge1969:bufonidae	ribs, fusion	fused	kluge1969		
Bufonidae	kluge1969:bufonidae	vertebral ossification	holochordal	kluge1969		
Bufonidae	kluge1969:bufonidae	pectoral girdle	arciferal	kluge1969		
Bufonidae	kluge1969:bufonidae	scapula, relation to clavical	juxtapose	kluge1969		
Ranidae	kluge1969:ranidae	tail muscle	absent	kluge1969		
Ranidae	kluge1969:ranidae	ribs, fusion	fused	kluge1969		
Ranidae	kluge1969:ranidae	vertebral ossification	holochordal	kluge1969		
Ranidae	kluge1969:ranidae	pectoral girdle	finnisternal	kluge1969		
Ranidae	kluge1969:ranidae	scapula, relation to clavical	juxtapose	kluge1969		
`

func TestReadTSV(t *testing.T) {
	m := matrix.New()
	if err := m.ReadTSV(strings.NewReader(obsText)); err != nil {
		t.Fatalf("unable to read TSV data: %v", err)
	}

	want := newMatrixWithComments()
	cmpMatrix(t, m, want)
}

func TestWriteTSV(t *testing.T) {
	m := newMatrixWithComments()
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
