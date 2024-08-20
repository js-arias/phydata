// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package dna

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

var headerFields = []string{
	"taxon",
	"specimen",
	"gene",
	"genbank",
	"bases",
}

var valFields = []Field{
	Protein,
	Organelle,
	Aligned,
	Reference,
	Comments,
}

// ReadTSV reads a set of DNA sequences
// from a TSV file.
//
// The TSV file must contains the following fields:
//
//   - taxon, the taxonomic name of the source taxon
//   - specimen, the ID of the particular source specimen
//   - gene, an identifier for the sequenced region
//   - genbank, the accession in the GenBank database.
//   - bases, the DNA sequence
//
// Additional fields are:
//
//   - protein, if "true" the molecule product is a protein
//   - organelle, the celular organelle that contains the sequence
//   - aligned, if "true" the sequence has been previously aligned
//   - reference, an ID of a bibliographic reference
//   - comments, simple additional comments about the sequence
//
// Here is an example file:
//
//	# DNA sequences
//	taxon	specimen	gene	genbank	protein	organelle	aligned	reference	comments	bases
//	Loxodonta africana	sp-01	cytb	MN148748	true	mitochondrion	true			ccatccaacatctcagcatgatgaaatttc
//	Loxodonta africana	sp-01	eef1a1	XM_064288029	true	nucleus	true			ggtaaactgggaagtgctggcgtgtgctgg
//	Orycteropus afer	sp-02	cytb	OR167429	true	mitochondrion	true			??gaccaacattcgtaaaacccaccctctt
//	Panthera tigris	fmnh_un_2485	cytb	MH290773	true	mitochondrion	true			gactcagacaaa---ccattccacccatac
//	Papio anubis	genbank:ku871221	cytb	KU871221	true	mitochondrion	true			atgaccccaatacgcaaatctaatcctatc
//	Papio anubis	genbank:xm_003897809	eef1a1	XM_003897809	true	nucleus	true			gcagtgagccgagatcgcgccactgcaccc
func (c *Collection) ReadTSV(r io.Reader) error {
	tab := csv.NewReader(r)
	tab.Comma = '\t'
	tab.Comment = '#'

	head, err := tab.Read()
	if err != nil {
		return fmt.Errorf("while reading header: %v", err)
	}
	fields := make(map[string]int, len(head))
	for i, h := range head {
		h = strings.ToLower(h)
		fields[h] = i
	}
	for _, h := range headerFields {
		if _, ok := fields[h]; !ok {
			return fmt.Errorf("expecting field %q", h)
		}
	}

	for {
		row, err := tab.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		ln, _ := tab.FieldPos(0)
		if err != nil {
			return fmt.Errorf("on row %d: %v", ln, err)
		}

		f := "taxon"
		tax := row[fields[f]]
		if tax == "" {
			continue
		}

		f = "specimen"
		spec := row[fields[f]]
		if spec == "" {
			continue
		}

		f = "gene"
		gene := row[fields[f]]
		if gene == "" {
			continue
		}

		f = "genbank"
		gb := row[fields[f]]
		if gb == "" {
			continue
		}

		f = "bases"
		seq := row[fields[f]]
		if seq == "" {
			continue
		}
		c.Add(tax, spec, gene, gb, seq)

		// additional fields
		for _, ff := range valFields {
			f = string(ff)
			i, ok := fields[f]
			if !ok {
				continue
			}

			v := row[i]
			c.Set(spec, gene, gb, v, ff)
		}
	}

	return nil
}

// TSV writes a DNA sequence collection as a TSV file.
func (c *Collection) TSV(w io.Writer) error {
	tab := csv.NewWriter(w)
	tab.Comma = '\t'
	tab.UseCRLF = true

	//header
	header := []string{"taxon", "specimen", "gene", "genbank", "protein", "organelle", "aligned", "reference", "comments", "bases"}
	if err := tab.Write(header); err != nil {
		return fmt.Errorf("unable to write header: %v", err)
	}

	tax := make(map[string][]string)
	var tn []string
	for _, sp := range c.specs {
		t, ok := tax[sp.taxon]
		if !ok {
			tn = append(tn, sp.taxon)
		}
		t = append(t, sp.name)
		tax[sp.taxon] = t
	}
	slices.Sort(tn)

	genes := c.Genes()
	for _, tt := range tn {
		t := tax[tt]
		slices.Sort(t)
		for _, spv := range t {
			sp := c.specs[spv]

			for _, gn := range genes {
				g := sp.genes[gn]
				if len(g) == 0 {
					continue
				}
				acc := make([]string, 0, len(g))
				for gb := range g {
					acc = append(acc, gb)
				}
				slices.Sort(acc)

				for _, a := range acc {
					seq := g[a]
					row := []string{
						sp.taxon,
						sp.name,
						gn,
						a,
						strconv.FormatBool(seq.protein),
						seq.organelle,
						strconv.FormatBool(seq.aligned),
						seq.ref,
						seq.comment,
						seq.seq,
					}
					if err := tab.Write(row); err != nil {
						return fmt.Errorf("while writing data: %v", err)
					}
				}
			}
		}
	}

	tab.Flush()
	if err := tab.Error(); err != nil {
		return fmt.Errorf("while writing data: %v", err)
	}

	return nil
}
