// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package matrix provides a representation
// of a phylogenetic data matrix.
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
	terms map[string]*terminal
}

// New creates a new empty matrix.
func New() *Matrix {
	return &Matrix{
		chars: make(map[string]*character),
		terms: make(map[string]*terminal),
	}
}

// Add adds a new observation
// (i.e., a character state) to the matrix
// for a given terminal,
// and character.
func (m *Matrix) Add(term, char, state string) {
	term = canon(term)
	if term == "" {
		return
	}

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

	t, ok := m.terms[term]
	if !ok {
		t = &terminal{
			name: term,
			obs:  make(map[string]map[string]*observation),
		}
		m.terms[term] = t
	}

	obs, ok := t.obs[char]
	if !ok {
		obs = make(map[string]*observation)
	}

	if state == NotApplicable {
		obs = make(map[string]*observation)
	} else if state == Unknown {
		delete(t.obs, char)
		return
	} else if isNoObservation(obs) {
		obs = make(map[string]*observation)
	}

	obs[state] = &observation{name: state}
	t.obs[char] = obs
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
// in a taxon.
func (m *Matrix) Obs(term, char string) []string {
	term = canon(term)
	if term == "" {
		return []string{Unknown}
	}
	t, ok := m.terms[term]
	if !ok {
		return []string{Unknown}
	}

	char = strings.Join(strings.Fields(char), " ")
	if char == "" {
		return []string{Unknown}
	}
	char = strings.ToLower(char)

	obs, ok := t.obs[char]
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

// Terminals returns the terminals in the matrix.
func (m *Matrix) Terminals() []string {
	terms := make([]string, 0, len(m.terms))
	for _, t := range m.terms {
		terms = append(terms, t.name)
	}
	slices.Sort(terms)
	return terms
}

type character struct {
	name   string
	states map[string]bool
}

type terminal struct {
	name string
	obs  map[string]map[string]*observation
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
