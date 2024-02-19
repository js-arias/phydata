// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package matrix provides a matrix of taxon specimens
// and character observations.
package matrix

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Character states without data.
const NotApplicable = "<na>"
const Unknown = "<unknown>"

// A Matrix is a phylogenetic data matrix,
// a collection of taxa
// and their character states.
type Matrix struct {
	chars map[string]*character
	specs map[string]*specimen
}

// New creates a new empty matrix.
func New() *Matrix {
	return &Matrix{
		chars: make(map[string]*character),
		specs: make(map[string]*specimen),
	}
}

// Add adds a new observation
// (i.e., a character state) to the matrix
// for a given taxon specimen,
// and character.
func (m *Matrix) Add(taxon, spec, char, state string) {
	taxon = canon(taxon)
	if taxon == "" {
		return
	}

	spec = strings.Join(strings.Fields(spec), " ")
	if spec == "" {
		return
	}
	spec = strings.ToLower(spec)

	char = strings.Join(strings.Fields(char), " ")
	if char == "" {
		return
	}
	char = strings.ToLower(char)

	state = strings.Join(strings.Fields(state), " ")
	if state == "" {
		return
	}
	state = strings.ToLower(state)

	c, ok := m.chars[char]
	if !ok {
		c = &character{
			name:   char,
			states: make(map[string]bool),
		}
		m.chars[char] = c
	}
	c.states[state] = true

	sp, ok := m.specs[spec]
	if !ok {
		sp = &specimen{
			taxon: taxon,
			name:  spec,
			obs:   make(map[string]map[string]*observation),
		}
		m.specs[spec] = sp
	}
	if sp.taxon != taxon {
		return
	}

	obs, ok := sp.obs[char]
	if !ok {
		obs = make(map[string]*observation)
	}

	if state == NotApplicable {
		obs = make(map[string]*observation)
	} else if state == Unknown {
		delete(sp.obs, char)
		return
	} else if isNoObservation(obs) {
		obs = make(map[string]*observation)
	}

	obs[state] = &observation{name: state}
	sp.obs[char] = obs
}

// Chars returns the characters in the matrix.
func (m *Matrix) Chars() []string {
	chars := make([]string, 0, len(m.chars))
	for _, c := range m.chars {
		chars = append(chars, c.name)
	}
	slices.Sort(chars)
	return chars
}

// Obs returns the states assigned for character
// in a specimen.
func (m *Matrix) Obs(spec, char string) []string {
	spec = strings.Join(strings.Fields(spec), " ")
	if spec == "" {
		return nil
	}
	spec = strings.ToLower(spec)

	sp, ok := m.specs[spec]
	if !ok {
		return []string{Unknown}
	}

	char = strings.Join(strings.Fields(char), " ")
	if char == "" {
		return []string{Unknown}
	}
	char = strings.ToLower(char)

	obs, ok := sp.obs[char]
	if !ok {
		return []string{Unknown}
	}

	states := make([]string, 0, len(obs))
	for _, s := range obs {
		states = append(states, s.name)
	}
	slices.Sort(states)
	return states
}

// States returns the states of a character in the matrix.
func (m *Matrix) States(char string) []string {
	char = strings.Join(strings.Fields(char), " ")
	if char == "" {
		return nil
	}
	char = strings.ToLower(char)
	c, ok := m.chars[char]
	if !ok {
		return nil
	}

	states := make([]string, 0, len(c.states))
	for s := range c.states {
		if s == NotApplicable {
			continue
		}
		states = append(states, s)
	}
	slices.Sort(states)
	return states
}

// Specimens returns the specimens in the matrix.
func (m *Matrix) Specimens() []string {
	specs := make([]string, 0, len(m.specs))
	for _, t := range m.specs {
		specs = append(specs, t.name)
	}
	slices.Sort(specs)
	return specs
}

// Field is used to define additional information fields
// of an observation.
type Field string

// Additional observation fields.
const (
	Reference Field = "reference"
	ImageLink Field = "image"
	Comments  Field = "comments"
)

// Set sets the value of an addition information
// for an observation.
func (m *Matrix) Set(spec, char, state, val string, field Field) {
	spec = strings.Join(strings.Fields(spec), " ")
	if spec == "" {
		return
	}
	spec = strings.ToLower(spec)

	sp, ok := m.specs[spec]
	if !ok {
		return
	}

	char = strings.Join(strings.Fields(char), " ")
	if char == "" {
		return
	}
	char = strings.ToLower(char)

	obsMap, ok := sp.obs[char]
	if !ok {
		return
	}

	state = strings.Join(strings.Fields(state), " ")
	if state == "" {
		return
	}
	state = strings.ToLower(state)

	obs, ok := obsMap[state]
	if !ok {
		return
	}

	val = strings.Join(strings.Fields(val), " ")

	switch field {
	case Reference:
		obs.ref = val
	case ImageLink:
		obs.img = val
	case Comments:
		obs.comment = val
	}
}

// Val returns the value of additional fields
// for an observation.
func (m *Matrix) Val(spec, char, state string, field Field) string {
	spec = strings.Join(strings.Fields(spec), " ")
	if spec == "" {
		return ""
	}
	spec = strings.ToLower(spec)

	sp, ok := m.specs[spec]
	if !ok {
		return ""
	}

	char = strings.Join(strings.Fields(char), " ")
	if char == "" {
		return ""
	}
	char = strings.ToLower(char)

	obsMap, ok := sp.obs[char]
	if !ok {
		return ""
	}

	state = strings.Join(strings.Fields(state), " ")
	if state == "" {
		return ""
	}
	state = strings.ToLower(state)

	obs, ok := obsMap[state]
	if !ok {
		return ""
	}

	switch field {
	case Reference:
		return obs.ref
	case ImageLink:
		return obs.img
	case Comments:
		return obs.comment
	}
	return ""
}

type character struct {
	name   string
	states map[string]bool
}

type specimen struct {
	taxon string
	name  string
	obs   map[string]map[string]*observation
}

type observation struct {
	name    string
	ref     string // bibliographic reference
	img     string // a link to an image
	comment string // a commentary of the observation
}

func isNoObservation(obs map[string]*observation) bool {
	if _, ok := obs[NotApplicable]; ok {
		return true
	}
	if _, ok := obs[Unknown]; ok {
		return true
	}
	return false
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
