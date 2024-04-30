// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package matrix

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var headerFields = []string{
	"taxon",
	"specimen",
	"character",
	"state",
}

var valFields = []Field{
	Reference,
	ImageLink,
	Comments,
}

// ReadTSV reads a set of specimen observations
// from a TSV file.
//
// The TSV file must contains the following fields:
//
//   - taxon, the taxonomic name of the taxon
//   - specimen, the ID of the particular specimen observed
//   - character, the name of the observed character
//   - state, the observed character state
//
// Additional fields are:
//
//   - reference, an ID of a bibliographic reference
//   - image, a path to an image of the observation
//   - comments, simple comments about the observation
//
// Here is an example file:
//
//	# character observations
//	taxon	specimen	character	state	reference	image	comments
//	Ascaphus truei	kluge1969:ascaphus_truei	tail muscle	present	kluge1969	ascaphus-tail.png	it might be not homologous with tail muscles of salamanders
//	Ascaphus truei	kluge1969:ascaphus_truei	ribs, fusion	free	kluge1969
//	Discoglossidae	kluge1969:discoglossidae	tail muscle	absent	kluge1969
//	Discoglossidae	kluge1969:discoglossidae	ribs, fusion	free	kluge1969
//	Pipidae	kluge1969:pipidae	tail muscle	absent	kluge1969
//	Pipidae	kluge1969:pipidae	ribs, fusion	fused in adults	kluge1969
func (m *Matrix) ReadTSV(r io.Reader) error {
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

		f = "character"
		char := row[fields[f]]
		if char == "" {
			continue
		}

		f = "state"
		state := row[fields[f]]
		if state == "" {
			continue
		}

		m.Add(tax, spec, char, state)

		for _, ff := range valFields {
			f = string(ff)
			i, ok := fields[f]
			if !ok {
				continue
			}

			v := row[i]
			m.Set(spec, char, state, v, ff)
		}
	}

	return nil
}

// TSV writes an observation matrix as a TSV file.
func (m *Matrix) TSV(w io.Writer) error {
	tab := csv.NewWriter(w)
	tab.Comma = '\t'
	tab.UseCRLF = true

	// header
	header := []string{"taxon", "specimen", "character", "state", "reference", "image", "comments"}
	if err := tab.Write(header); err != nil {
		return fmt.Errorf("unable to write header: %v", err)
	}

	tax := make(map[string][]string)
	var tn []string
	for _, sp := range m.specs {
		t, ok := tax[sp.taxon]
		if !ok {
			tn = append(tn, sp.taxon)
		}
		t = append(t, sp.name)
		tax[sp.taxon] = t
	}
	slices.Sort(tn)

	chars := m.Chars()

	for _, tt := range tn {
		t := tax[tt]
		slices.Sort(t)
		for _, spv := range t {
			sp := m.specs[spv]

			for _, c := range chars {
				obs, ok := sp.obs[c]
				if !ok {
					continue
				}

				// special case: not aplicable
				if o, ok := obs[NotApplicable]; ok {
					row := []string{
						sp.taxon,
						sp.name,
						c,
						NotApplicable,
						o.ref,
						o.img,
						o.comment,
					}
					if err := tab.Write(row); err != nil {
						return fmt.Errorf("while writing data: %v", err)
					}
					continue
				}

				sts := make([]string, 0, len(obs))
				for _, o := range obs {
					sts = append(sts, o.name)
				}
				slices.Sort(sts)

				for _, s := range sts {
					o := obs[s]
					row := []string{
						sp.taxon,
						sp.name,
						c,
						o.name,
						o.ref,
						o.img,
						o.comment,
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
