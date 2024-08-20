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
		specimen  string
		genBank   string
		gene      string
		seq       string
		aligned   bool
		protein   bool
		organelle string
	}{
		"MN148748": {
			specimen:  "sp-01",
			genBank:   "MN148748",
			gene:      "cytb",
			seq:       "ccatccaacatctcagcatgatgaaatttc",
			aligned:   true,
			protein:   true,
			organelle: "mitochondrion",
		},
		"XM_064288029": {
			specimen:  "sp-01",
			genBank:   "XM_064288029",
			gene:      "eef1a1",
			seq:       "ggtaaactgggaagtgctggcgtgtgctgg",
			aligned:   true,
			protein:   true,
			organelle: "nucleus",
		},
		"OR167429": {
			specimen:  "sp-02",
			genBank:   "OR167429",
			gene:      "cytb",
			seq:       "??gaccaacattcgtaaaacccaccctctt",
			aligned:   true,
			protein:   true,
			organelle: "mitochondrion",
		},
		"MH290773": {
			specimen:  "fmnh_un_2485",
			genBank:   "MH290773",
			gene:      "cytb",
			seq:       "gactcagacaaa---ccattccacccatac",
			aligned:   true,
			protein:   true,
			organelle: "mitochondrion",
		},
		"KU871221": {
			specimen:  "genbank:ku871221",
			genBank:   "KU871221",
			gene:      "cytb",
			seq:       "atgaccccaatacgcaaatctaatcctatc",
			aligned:   true,
			protein:   true,
			organelle: "mitochondrion",
		},
		"XM_003897809": {
			specimen:  "genbank:XM_003897809",
			genBank:   "XM_003897809",
			gene:      "eef1a1",
			seq:       "gcagtgagccgagatcgcgccactgcaccc",
			aligned:   true,
			protein:   true,
			organelle: "nucleus",
		},
	}

	for name, seq := range seqs {
		s := c.Sequence(seq.specimen, seq.gene, seq.genBank)
		if s != seq.seq {
			t.Errorf("sequence %q: specimen %q, gene %q, accession %q: got %q, want %q", name, seq.specimen, seq.gene, seq.genBank, s, seq.seq)
		}

		aligned := c.Val(seq.specimen, seq.gene, seq.genBank, dna.Aligned)
		if aligned == "true" {
			if !seq.aligned {
				t.Errorf("sequence %q: specimen %q, gene %q, accession %q: aligned is %q", name, seq.specimen, seq.gene, seq.genBank, aligned)
			}
		} else {
			if seq.aligned {
				t.Errorf("sequence %q: specimen %q, gene %q, accession %q: aligned is %q", name, seq.specimen, seq.gene, seq.genBank, aligned)
			}
		}

		protein := c.Val(seq.specimen, seq.gene, seq.genBank, dna.Protein)
		if protein == "true" {
			if !seq.protein {
				t.Errorf("sequence %q: specimen %q, gene %q, accession %q: protein is %q", name, seq.specimen, seq.gene, seq.genBank, protein)
			}
		} else {
			if seq.protein {
				t.Errorf("sequence %q: specimen %q, gene %q, accession %q: protein is %q", name, seq.specimen, seq.gene, seq.genBank, protein)
			}
		}

		organelle := c.Val(seq.specimen, seq.gene, seq.genBank, dna.Organelle)
		if organelle != seq.organelle {
			t.Errorf("sequence %q: specimen %q, gene %q, accession %q: got organelle %q, want %q", name, seq.specimen, seq.gene, seq.genBank, organelle, seq.organelle)
		}
	}
}

func newCollection() *dna.Collection {
	c := dna.New()
	c.Add("Loxodonta africana", "sp-01", "cytb", "MN148748", "ccatccaaca tctcagcatg atgaaatttc")
	c.Add("Loxodonta africana", "sp-01", "eef1a1", "XM_064288029", "ggtaaactgg gaagtgctgg cgtgtgctgg")
	c.Add("Orycteropus afer", "sp-02", "cytb", "OR167429", "??gaccaaca ttcgtaaaac ccaccctctt")
	c.Add("Panthera tigris", "FMNH_UN_2485", "cytb", "MH290773 ", "gactcagaca aa---ccatt ccacccatac")
	c.Add("Papio anubis", "", "cytb", "KU871221 ", "atgaccccaa tacgcaaatc taatcctatc")
	c.Add("Papio anubis", "", "eef1a1", "XM_003897809", "gcagtgagcc gagatcgcgc cactgcaccc")

	c.Set("sp-01", "cytb", "MN148748", "true", dna.Aligned)
	c.Set("sp-01", "cytb", "MN148748", "true", dna.Protein)
	c.Set("sp-01", "cytb", "MN148748", "mitochondrion", dna.Organelle)
	c.Set("sp-01", "eef1a1", "XM_064288029", "true", dna.Aligned)
	c.Set("sp-01", "eef1a1", "XM_064288029", "true", dna.Protein)
	c.Set("sp-01", "eef1a1", "XM_064288029", "nucleus", dna.Organelle)
	c.Set("sp-02", "cytb", "OR167429", "true", dna.Aligned)
	c.Set("sp-02", "cytb", "OR167429", "true", dna.Protein)
	c.Set("sp-02", "cytb", "OR167429", "mitochondrion", dna.Organelle)
	c.Set("fmnh_un_2485", "cytb", "MH290773", "true", dna.Aligned)
	c.Set("fmnh_un_2485", "cytb", "MH290773", "true", dna.Protein)
	c.Set("fmnh_un_2485", "cytb", "MH290773", "mitochondrion", dna.Organelle)
	c.Set("genbank:KU871221", "cytb", "KU871221", "true", dna.Aligned)
	c.Set("genbank:KU871221", "cytb", "KU871221", "true", dna.Protein)
	c.Set("genbank:KU871221", "cytb", "KU871221", "mitochondrion", dna.Organelle)
	c.Set("genbank:xm_003897809", "eef1a1", "XM_003897809", "true", dna.Aligned)
	c.Set("genbank:xm_003897809", "eef1a1", "XM_003897809", "true", dna.Protein)
	c.Set("genbank:xm_003897809", "eef1a1", "XM_003897809", "nucleus", dna.Organelle)
	return c
}

func cmpCollection(t testing.TB, got, want *dna.Collection) {
	t.Helper()

	specs := want.Specimens()
	csp := got.Specimens()
	if !reflect.DeepEqual(csp, specs) {
		t.Errorf("specimens: got %v, want %v", csp, specs)
	}

	genes := want.Genes()
	cgn := got.Genes()
	if !reflect.DeepEqual(cgn, genes) {
		t.Errorf("genes: got %v, want %v", cgn, genes)
	}

	genbank := want.GenBank()
	cgb := got.GenBank()
	if !reflect.DeepEqual(cgb, genbank) {
		t.Errorf("genbank: got %v, want %v", cgb, genbank)
	}

	for _, tax := range want.Taxa() {
		for _, spec := range want.TaxSpec(tax) {
			for _, gene := range want.SpecGene(spec) {
				for _, acc := range want.GeneAccession(spec, gene) {
					seq := want.Sequence(spec, gene, acc)
					s := got.Sequence(spec, gene, acc)
					if s != seq {
						t.Errorf("sequence %q: specimen %q, gene %q, accession %q: got %q, want %q", tax, spec, gene, acc, s, seq)
					}

					alg := want.Val(spec, gene, acc, dna.Aligned)
					aligned := got.Val(spec, gene, acc, dna.Aligned)
					if aligned != alg {
						t.Errorf("sequence %q: specimen %q, gene %q, accession %q: aligned: got %q, want %q", tax, spec, gene, acc, aligned, alg)
					}

					prt := want.Val(spec, gene, acc, dna.Protein)
					protein := got.Val(spec, gene, acc, dna.Aligned)
					if protein != prt {
						t.Errorf("sequence %q: specimen %q, gene %q, accession %q: protein: got %q, want %q", tax, spec, gene, acc, protein, prt)
					}

					org := want.Val(spec, gene, acc, dna.Organelle)
					organelle := got.Val(spec, gene, acc, dna.Organelle)
					if organelle != org {
						t.Errorf("sequence %q: specimen %q, gene %q, accession %q: organelle: got %q, want %q", tax, spec, gene, acc, organelle, org)
					}
				}
			}

		}
	}

}
