// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package dna_test

import (
	"reflect"
	"testing"

	"github.com/js-arias/phydata/matrix/dna"
)

func TestCollection(t *testing.T) {
	c := newCollection()

	specs := []string{
		"fmnh_un_2485",
		"genbank:ku871221",
		"genbank:xm_003897809",
		"sp-01",
		"sp-02",
	}
	csp := c.Specimens()
	if !reflect.DeepEqual(csp, specs) {
		t.Errorf("specimens: got %v, want %v", csp, specs)
	}

	genes := []string{"cytb", "eef1a1"}
	cgn := c.Genes()
	if !reflect.DeepEqual(cgn, genes) {
		t.Errorf("genes: got %v, want %v", cgn, genes)
	}

	genbank := []string{"KU871221", "MH290773", "MN148748", "OR167429", "XM_003897809", "XM_064288029"}
	cgb := c.GenBank()
	if !reflect.DeepEqual(cgb, genbank) {
		t.Errorf("genbank: got %v, want %v", cgb, genbank)
	}

	seqs := map[string]struct {
		specimen string
		genBank  string
		gene     string
		seq      string
	}{
		"MN148748": {
			specimen: "sp-01",
			genBank:  "MN148748",
			gene:     "cytb",
			seq:      "ccatccaacatctcagcatgatgaaatttc",
		},
		"XM_064288029": {
			specimen: "sp-01",
			genBank:  "XM_064288029",
			gene:     "eef1a1",
			seq:      "ggtaaactgggaagtgctggcgtgtgctgg",
		},
		"OR167429": {
			specimen: "sp-02",
			genBank:  "OR167429",
			gene:     "cytb",
			seq:      "??gaccaacattcgtaaaacccaccctctt",
		},
		"MH290773": {
			specimen: "fmnh_un_2485",
			genBank:  "MH290773",
			gene:     "cytb",
			seq:      "gactcagacaaa---ccattccacccatac",
		},
		"KU871221": {
			specimen: "genbank:KU871221",
			genBank:  "KU871221",
			gene:     "cytb",
			seq:      "atgaccccaatacgcaaatctaatcctatc",
		},
		"XM_003897809": {
			specimen: "genbank:XM_003897809",
			genBank:  "XM_003897809",
			gene:     "eef1a1",
			seq:      "gcagtgagccgagatcgcgccactgcaccc",
		},
	}

	for name, seq := range seqs {
		s := c.Sequence(seq.specimen, seq.gene, seq.genBank)
		if s != seq.seq {
			t.Errorf("sequence %q: specimen %q, gene %q, accession %q: got %q, want %q", name, seq.specimen, seq.gene, seq.genBank, s, seq.seq)
		}
	}
}

func newCollection() *dna.Collection {
	c := dna.New()
	c.Add("Loxodonta africana", "sp-01", "cytb", "MN148748", "ccatccaaca tctcagcatg atgaaatttc", true)
	c.Add("Loxodonta africana", "sp-01", "eef1a1", "XM_064288029", "ggtaaactgg gaagtgctgg cgtgtgctgg", true)
	c.Add("Orycteropus afer", "sp-02", "cytb", "OR167429", "??gaccaaca ttcgtaaaac ccaccctctt", true)
	c.Add("Panthera tigris", "FMNH_UN_2485", "cytb", "MH290773 ", "gactcagaca aa---ccatt ccacccatac", true)
	c.Add("Papio anubis", "", "cytb", "KU871221 ", "atgaccccaa tacgcaaatc taatcctatc", true)
	c.Add("Papio anubis", "", "eef1a1", "XM_003897809", "gcagtgagcc gagatcgcgc cactgcaccc", true)
	return c
}
