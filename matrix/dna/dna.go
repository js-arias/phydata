// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package dna stores DNA sequences
// for taxon specimens.
package dna

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A Collection is a collection of taxa
// and their sequences.
type Collection struct {
	specs map[string]*specimen
}

// New creates a new empty collection.
func New() *Collection {
	return &Collection{
		specs: make(map[string]*specimen),
	}
}

// Add adds a new sequence to the collection
// for a given taxon specimen
// and molecule.
// The GenBank accession can be empty.
// If no specimen is given,
// it will generate an specimen ID using the form
// "genbank:<genbank-ID>",
// in this case if no specimen is given,
// it will return an error.
// The sequence can be aligned or unaligned.
func (c *Collection) Add(taxon, spec, gene, genBank, seq string, aligned bool) error {
	taxon = canon(taxon)
	if taxon == "" {
		return nil
	}

	genBank = strings.TrimSpace(genBank)
	spec = specID(spec)
	if spec == "" && genBank == "" {
		return fmt.Errorf("sequence without identifier")
	}
	if spec == "" {
		spec = specID("genbank:" + genBank)
	}
	if genBank == "" {
		genBank = "no-gb:" + spec
	}

	seq = formatSequence(seq)

	gene = strings.TrimSpace(gene)
	if gene == "" {
		return fmt.Errorf("sequence %q without a defined gene-molecule identifier", genBank)
	}
	gene = strings.ToLower(gene)

	sp, ok := c.specs[spec]
	if !ok {
		sp = &specimen{
			taxon: taxon,
			name:  spec,
			genes: make(map[string]map[string]genBankSequence),
		}
		c.specs[spec] = sp
	}

	gb, ok := sp.genes[gene]
	if !ok {
		gb = make(map[string]genBankSequence)
		sp.genes[gene] = gb
	}
	gb[genBank] = genBankSequence{
		seq:     seq,
		aligned: aligned,
	}

	return nil
}

// GenBank returns the GenBank accessions
// for the sequences in a collection.
func (c *Collection) GenBank() []string {
	ids := make(map[string]bool)
	for _, sp := range c.specs {
		for _, g := range sp.genes {
			for gb := range g {
				ids[gb] = true
			}
		}
	}

	gbIDs := make([]string, 0, len(ids))
	for gb := range ids {
		gbIDs = append(gbIDs, gb)
	}
	slices.Sort(gbIDs)
	return gbIDs
}

// Genes returns the genes-molecules with sequences
// in the collection.
func (c *Collection) Genes() []string {
	genNames := make(map[string]bool)
	for _, sp := range c.specs {
		for g := range sp.genes {
			genNames[g] = true
		}
	}

	genes := make([]string, 0, len(genNames))
	for g := range genNames {
		genes = append(genes, g)
	}
	slices.Sort(genes)
	return genes
}

// Sequence returns a sequence for a given specimen,
// gene,
// and genBank accession.
func (c *Collection) Sequence(specimen, gene, genBank string) string {
	specimen = specID(specimen)
	if specimen == "" {
		return ""
	}

	sp, ok := c.specs[specimen]
	if !ok {
		return ""
	}
	gene = strings.TrimSpace(strings.ToLower(gene))
	gb, ok := sp.genes[gene]
	if !ok {
		return ""
	}
	seq, ok := gb[genBank]
	if !ok {
		return ""
	}
	return seq.seq
}

// Specimens returns the specimens in the collection.
func (c *Collection) Specimens() []string {
	specs := make([]string, 0, len(c.specs))
	for _, sp := range c.specs {
		specs = append(specs, sp.name)
	}
	slices.Sort(specs)
	return specs
}

type specimen struct {
	taxon string
	name  string
	genes map[string]map[string]genBankSequence
}

type genBankSequence struct {
	seq     string
	aligned bool
}

// Canon returns a taxon name
// in its canonical form.
func canon(name string) string {
	name = strings.Join(strings.Fields(name), " ")
	if name == "" {
		return ""
	}
	name = strings.ToLower(name)
	r, n := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[n:]
}

func specID(spec string) string {
	spec = strings.Join(strings.Fields(spec), "_")
	if spec == "" {
		return ""
	}
	return strings.ToLower(spec)
}

func formatSequence(seq string) string {
	seq = strings.Join(strings.Fields(seq), "")
	seq = strings.ToLower(seq)
	return seq
}
