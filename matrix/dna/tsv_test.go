// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package dna_test

import (
	"bytes"
	"testing"

	"github.com/js-arias/phydata/matrix/dna"
)

func TestTSV(t *testing.T) {
	c := newCollection()
	var w bytes.Buffer
	if err := c.TSV(&w); err != nil {
		t.Fatalf("unable to write TSV data: %v", err)
	}
	t.Logf("output:\n%s\n", w.String())

	got := dna.New()
	if err := got.ReadTSV(&w); err != nil {
		t.Fatalf("unable to read TSV data: %v", err)
	}

	cmpCollection(t, got, c)
}
