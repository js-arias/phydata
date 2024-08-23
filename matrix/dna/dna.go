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
func (c *Collection) Add(taxon, spec, gene, genBank, seq string) error {
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
			genes: make(map[string]map[string]*genBankSequence),
		}
		c.specs[spec] = sp
	}

	gb, ok := sp.genes[gene]
	if !ok {
		gb = make(map[string]*genBankSequence)
		sp.genes[gene] = gb
	}
	gb[genBank] = &genBankSequence{
		seq: seq,
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

// GeneAccession returns the accession
// for a given gene
// of a given specimen.
func (c *Collection) GeneAccession(specimen, gene string) []string {
	specimen = specID(specimen)
	if specimen == "" {
		return nil
	}

	sp, ok := c.specs[specimen]
	if !ok {
		return nil
	}
	gene = strings.TrimSpace(strings.ToLower(gene))
	gb, ok := sp.genes[gene]
	if !ok {
		return nil
	}

	acc := make([]string, 0, len(gb))
	for a := range gb {
		acc = append(acc, a)
	}
	slices.Sort(acc)
	return acc
}

// Sequence returns a sequence for a given specimen,
// gene,
// and genBank accession.
func (c *Collection) Sequence(specimen, gene, genBank string) string {
	seq := c.sequence(specimen, gene, genBank)
	if seq == nil {
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

// SpecGene return the genes defined for a given specimen.
func (c *Collection) SpecGene(specimen string) []string {
	specimen = specID(specimen)
	sp, ok := c.specs[specimen]
	if !ok {
		return nil
	}

	genes := make([]string, 0, len(sp.genes))
	for g := range sp.genes {
		genes = append(genes, g)
	}
	slices.Sort(genes)

	return genes
}

// MaxLen returns the maximum length
// of a sequence for a given gene.
func (c *Collection) MaxLen(gene string) int {
	gene = strings.ToLower(strings.TrimSpace(gene))
	var max int
	for _, sp := range c.specs {
		gb, ok := sp.genes[gene]
		if !ok {
			continue
		}
		for _, s := range gb {
			ln := len(s.seq)
			if ln > max {
				max = ln
			}
		}
	}

	return max
}

// Taxa returns the taxa defined in the matrix.
func (c *Collection) Taxa() []string {
	taxa := make(map[string]bool)
	for _, sp := range c.specs {
		taxa[sp.taxon] = true
	}

	txLs := make([]string, 0, len(taxa))
	for t := range taxa {
		txLs = append(txLs, t)
	}
	slices.Sort(txLs)

	return txLs
}

// TaxSpec returns the specimens of a given taxon.
func (c *Collection) TaxSpec(name string) []string {
	name = canon(name)
	var specs []string
	for _, sp := range c.specs {
		if sp.taxon != name {
			continue
		}
		specs = append(specs, sp.name)
	}
	slices.Sort(specs)

	return specs
}

// Field is used to define additional information fields
// of a DNA gene.
type Field string

// Additional sequence fields.
const (
	Aligned   Field = "aligned"
	Protein   Field = "protein"
	Organelle Field = "organelle"
	Reference Field = "reference"
	Comments  Field = "comments"
)

// Set sets the value of an additional information
// for a sequence.
func (c *Collection) Set(specimen, gene, genBank, val string, field Field) {
	seq := c.sequence(specimen, gene, genBank)
	if seq == nil {
		return
	}

	val = strings.Join(strings.Fields(val), " ")

	switch field {
	case Aligned:
		seq.aligned = false
		if strings.ToLower(val) == "true" {
			seq.aligned = true
		}
	case Protein:
		seq.protein = false
		if strings.ToLower(val) == "true" {
			seq.protein = true
		}
	case Organelle:
		seq.organelle = strings.ToLower(val)
	case Reference:
		seq.ref = val
	case Comments:
		seq.comment = val
	}
}

func (c *Collection) Val(specimen, gene, genBank string, field Field) string {
	seq := c.sequence(specimen, gene, genBank)
	if seq == nil {
		return ""
	}

	switch field {
	case Aligned:
		if seq.aligned {
			return "true"
		}
		return "false"
	case Protein:
		if seq.protein {
			return "true"
		}
		return "false"
	case Organelle:
		return seq.organelle
	case Reference:
		return seq.ref
	case Comments:
		return seq.comment
	}

	return ""
}

func (c *Collection) sequence(specimen, gene, genBank string) *genBankSequence {
	specimen = specID(specimen)
	if specimen == "" {
		return nil
	}

	sp, ok := c.specs[specimen]
	if !ok {
		return nil
	}
	gene = strings.TrimSpace(strings.ToLower(gene))
	gb, ok := sp.genes[gene]
	if !ok {
		return nil
	}
	seq, ok := gb[genBank]
	if !ok {
		return nil
	}

	return seq
}

type specimen struct {
	taxon string
	name  string
	genes map[string]map[string]*genBankSequence
}

type genBankSequence struct {
	seq       string
	aligned   bool
	protein   bool
	organelle string
	ref       string
	comment   string
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
