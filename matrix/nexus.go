// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package matrix

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ReadNexus reads a character matrix from a NEXUS file.
// It require an ID for the matrix,
// and a ID for a bibliographic reference.
func (m *Matrix) ReadNexus(r io.Reader, ref string) error {
	nxf := bufio.NewReader(r)
	token := &strings.Builder{}

	// header
	if _, err := readToken(nxf, token); err != nil {
		return fmt.Errorf("expecting '#nexus' header: %v", err)
	}
	if t := strings.ToLower(token.String()); t != "#nexus" {
		return fmt.Errorf("got %q, expecting '#nexus' header", t)
	}

	// ignore all blocks except character block
	for {
		if _, err := readToken(nxf, token); err != nil {
			return fmt.Errorf("expecting 'begin' token: %v", err)
		}
		if t := strings.ToLower(token.String()); t != "begin" {
			return fmt.Errorf("got %q, expecting 'begin' block", t)
		}

		if _, err := readToken(nxf, token); err != nil {
			return fmt.Errorf("expecting block name: %v", err)
		}
		block := strings.ToLower(token.String())
		if block == "characters" {
			break
		}

		if err := skipBlock(nxf, token); err != nil {
			return fmt.Errorf("incomplete block %q: %v", block, err)
		}
	}

	var chars []nexusChar
	for {
		if _, err := readToken(nxf, token); err != nil {
			return fmt.Errorf("incomplete block 'characters': %v", err)
		}
		t := strings.ToLower(token.String())
		if t == "end" || t == "endblock" {
			break
		}
		if t == "charstatelabels" {
			var err error
			chars, err = readNexusCharStateLabels(nxf, token)
			if err != nil {
				return err
			}
			continue
		}
		if t == "charlabels" {
			var err error
			chars, err = readNexusCharLabels(nxf, token)
			if err != nil {
				return err
			}
			continue
		}
		if t == "statelabels" {
			if err := readNexusStateLabels(nxf, token, chars); err != nil {
				return err
			}
			continue
		}
		if t == "matrix" {
			if err := m.readNexusMatrix(nxf, token, ref, chars); err != nil {
				return err
			}
			continue
		}
		if err := skipDefinition(nxf, token); err != nil {
			return fmt.Errorf("incomplete block 'characters', token %q: %v", t, err)
		}
	}

	return nil
}

// Nexus writes an observation matrix as a NEXUS file.
func (m *Matrix) Nexus(w io.Writer) error {
	// header
	fmt.Fprintf(w, "#NEXUS\n")
	fmt.Fprintf(w, "[written %s]\n\n", time.Now().Format(time.RFC3339))

	// taxa block
	taxa := m.Taxa()
	fmt.Fprintf(w, "BEGIN TAXA;\n")
	fmt.Fprintf(w, "\tTITLE Taxa;\n")
	fmt.Fprintf(w, "\tDIMENSIONS NTAX=%d;\n", len(taxa))
	fmt.Fprintf(w, "\tTAXLABELS\n")
	for _, n := range taxa {
		n = strings.Join(strings.Fields(n), "_")
		fmt.Fprintf(w, "\t\t%s\n", n)
	}
	fmt.Fprintf(w, "\t;\n")
	fmt.Fprintf(w, "END;\n\n")

	// character block
	chars := m.Chars()
	fmt.Fprintf(w, "BEGIN CHARACTERS;\n")
	fmt.Fprintf(w, "\tTITLE 'Phylogenetic data matrix';\n")
	fmt.Fprintf(w, "\tDIMENSIONS NCHAR=%d;\n", len(chars))
	fmt.Fprintf(w, "\tFORMAT DATATYPE = STANDARD RESPECTCASE GAP = - MISSING = ? SYMBOLS = \"0 1 2 3 4 5 6 7 8 9 A B C D E F\";\n")
	fmt.Fprintf(w, "\tCHARSTATELABELS\n")
	states := make(map[string][]string, len(chars))
	for i, c := range chars {
		st := m.States(c)
		states[c] = st
		cn := strings.Join(strings.Fields(c), "_")
		fmt.Fprintf(w, "\t\t%d '%s' /", i+1, cn)
		for _, s := range st {
			fmt.Fprintf(w, " '%s'", s)
		}
		if i+1 < len(chars) {
			fmt.Fprintf(w, ",\n")
			continue
		}
		fmt.Fprintf(w, " ;\n")
	}

	// matrix
	fmt.Fprintf(w, "\tMATRIX\n")
	for _, n := range taxa {
		nm := strings.Join(strings.Fields(n), "_")
		fmt.Fprintf(w, "\t%s\t", nm)
		sp := m.TaxSpec(n)
		for _, c := range chars {
			val := "?"
			chSt := make(map[string]bool)
			for _, spec := range sp {
				obs := m.Obs(spec, c)
				for _, o := range obs {
					if o == NotApplicable {
						val = "-"
						continue
					}
					if o == Unknown {
						continue
					}

					chSt[o] = true
				}
			}
			if len(chSt) == 0 {
				fmt.Fprintf(w, "%s", val)
				continue
			}
			val = ""
			for i, s := range states[c] {
				if !chSt[s] {
					continue
				}
				val += strconv.FormatInt(int64(i), 16)
			}
			if len(val) > 1 {
				val = "{" + val + "}"
			}
			fmt.Fprintf(w, "%s", val)
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "\t;\n")
	fmt.Fprintf(w, "END;\n\n")
	return nil
}

type nexusChar struct {
	name   string
	states []string
}

func readNexusCharStateLabels(r *bufio.Reader, token *strings.Builder) ([]nexusChar, error) {
	var chars []nexusChar
	for i := 0; ; i++ {
		// read character number
		if _, err := readToken(r, token); err != nil {
			return nil, fmt.Errorf("while reading char state labels: %v, last character read: %d", err, i)
		}

		id, err := strconv.Atoi(token.String())
		if err != nil {
			return nil, fmt.Errorf("while reading char state labels: char %d [%q]: %v", i+1, token.String(), err)
		}
		if id != i+1 {
			return nil, fmt.Errorf("while reading char state labels: char %d [%q]: expecting %d", i+1, token.String(), i+1)
		}

		// read character name
		delim, err := readToken(r, token)
		if err != nil {
			return nil, fmt.Errorf("while reading char state labels: char %d [%q]: %v", i+1, token.String(), err)
		}
		cName := strings.ReplaceAll(token.String(), "_", " ")
		cName = strings.Join(strings.Fields(cName), " ")

		if delim == ',' || delim == ';' {
			chars = append(chars, nexusChar{
				name: cName,
			})
			if delim == ';' {
				break
			}
			continue
		}
		if delim != '/' {
			return nil, fmt.Errorf("while reading char state labels: char %d [%q]: expecting '/' delimiter", i+1, token.String())
		}

		// read state names
		var states []string
		for {
			delim, err = readToken(r, token)
			if err != nil {
				return nil, fmt.Errorf("while reading char state labels: char %d [%q]: %v", i+1, token.String(), err)
			}
			sName := strings.ReplaceAll(token.String(), "_", " ")
			sName = strings.Join(strings.Fields(sName), " ")
			states = append(states, sName)
			if delim == ',' || delim == ';' {
				break
			}
		}

		chars = append(chars, nexusChar{
			name:   cName,
			states: states,
		})

		if delim == ';' {
			break
		}
	}
	return chars, nil
}

func readNexusCharLabels(r *bufio.Reader, token *strings.Builder) ([]nexusChar, error) {
	var chars []nexusChar
	for i := 0; ; i++ {
		// read character name
		delim, err := readToken(r, token)
		if err != nil {
			return nil, fmt.Errorf("while reading char labels: char %d [%q]: %v", i+1, token.String(), err)
		}
		cName := strings.ReplaceAll(token.String(), "_", " ")
		cName = strings.Join(strings.Fields(cName), " ")

		chars = append(chars, nexusChar{
			name: cName,
		})

		if delim == ';' {
			break
		}
	}
	return chars, nil
}

func readNexusStateLabels(r *bufio.Reader, token *strings.Builder, chars []nexusChar) error {
	for i := 0; ; i++ {
		// read character number
		delim, err := readToken(r, token)
		if err != nil {
			return fmt.Errorf("while reading state labels: %v, last character read: %d", err, i)
		}
		if t := token.String(); t == "" && delim == ';' {
			break
		}

		id, err := strconv.Atoi(token.String())
		if err != nil {
			return fmt.Errorf("while reading state labels: char %d [%q]: %v", i+1, token.String(), err)
		}
		if id != i+1 {
			return fmt.Errorf("while reading state labels: char %d [%q]: expecting %d", i+1, token.String(), i+1)
		}

		// read state names
		var states []string
		for {
			delim, err = readToken(r, token)
			if err != nil {
				return fmt.Errorf("while reading char state labels: char %d [%q]: %v", i+1, token.String(), err)
			}
			sName := strings.ReplaceAll(token.String(), "_", " ")
			sName = strings.Join(strings.Fields(sName), " ")
			states = append(states, sName)
			if delim == ',' || delim == ';' {
				break
			}
		}
		if i < len(chars) {
			chars[i].states = states
		}

		if delim == ';' {
			break
		}
	}
	return nil
}

func (m *Matrix) readNexusMatrix(r *bufio.Reader, token *strings.Builder, ref string, chars []nexusChar) error {
	last := ""
	for {
		// read taxon name
		if _, err := readToken(r, token); err != nil {
			return fmt.Errorf("while reading matrix: %v, last taxon read %q", err, last)
		}
		tax := strings.ReplaceAll(token.String(), "_", " ")
		tax = strings.Join(strings.Fields(tax), " ")
		tax = canon(tax)
		spec := specID(ref + ":" + tax)

		// read characters
		char := 0
		for {
			r1, _, err := r.ReadRune()
			if err != nil {
				return fmt.Errorf("while reading matrix: taxon %q: %v", tax, err)
			}
			cName := fmt.Sprintf("char %d", char+1)
			var c nexusChar
			if char < len(chars) {
				c = chars[char]
				cName = c.name
			}

			if r1 == '\n' || r1 == '\r' {
				break
			}
			if unicode.IsSpace(r1) {
				continue
			}
			char++

			if r1 == '-' {
				m.Add(tax, spec, cName, NotApplicable)
				m.Set(spec, cName, NotApplicable, ref, Reference)
				continue
			}
			if r1 == '?' {
				m.Add(tax, spec, cName, Unknown)
				continue
			}
			if r1 == '(' || r1 == '{' {
				// polymorphic characters
				empty := true
				for {
					r1, _, err := r.ReadRune()
					if err != nil {
						return fmt.Errorf("while reading matrix: taxon %q: char: %d: %v", tax, char, err)
					}
					if r1 == '}' || r1 == ')' {
						break
					}
					if unicode.IsSpace(r1) {
						continue
					}

					s, err := strconv.ParseInt(string(r1), 16, 0)
					if err != nil {
						return fmt.Errorf("while reading matrix: taxon %q: char: %d [%q]: %v", tax, char, string(r1), err)
					}
					sName := fmt.Sprintf("state %d", s)
					if int(s) < len(c.states) {
						sName = c.states[int(s)]
					}
					m.Add(tax, spec, cName, sName)
					m.Set(spec, cName, sName, ref, Reference)
					empty = false
				}
				if empty {
					return fmt.Errorf("while reading matrix: taxon %q: char: %d: empty polymorph", tax, char)
				}
				continue
			}
			s, err := strconv.ParseInt(string(r1), 16, 0)
			if err != nil {
				return fmt.Errorf("while reading matrix: taxon %q: char: %d [%q]: %v", tax, char, string(r1), err)
			}
			sName := fmt.Sprintf("state %d", s)
			if int(s) < len(c.states) {
				sName = c.states[int(s)]
			}
			m.Add(tax, spec, cName, sName)
			m.Set(spec, cName, sName, ref, Reference)
		}
		last = tax

		// check if there is a next taxon
		if err := skipSpaces(r); err != nil {
			return fmt.Errorf("while reading matrix: %v, last taxon read %q", err, last)
		}
		r1, _, err := r.ReadRune()
		if err != nil {
			return fmt.Errorf("while reading matrix: %v, last taxon read %q", err, last)
		}
		if r1 == ';' {
			break
		}
		r.UnreadRune()
	}
	return nil
}

func skipBlock(r *bufio.Reader, token *strings.Builder) error {
	for {
		_, err := readToken(r, token)
		t := strings.ToLower(token.String())
		if t == "end" || t == "endblock" {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func skipDefinition(r *bufio.Reader, token *strings.Builder) error {
	for {
		delim, err := readToken(r, token)
		if delim == ';' {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func readToken(r *bufio.Reader, token *strings.Builder) (delim rune, err error) {
	token.Reset()

	if err := skipSpaces(r); err != nil {
		return 0, err
	}

	r1, _, err := r.ReadRune()
	if err != nil {
		return 0, err
	}
	if r1 == '\'' || r1 == '"' {
		// quoted block
		stop := r1
		for {
			r1, _, err := r.ReadRune()
			if err != nil {
				return 0, err
			}
			if r1 == stop {
				nx, _, err := r.ReadRune()
				if err != nil {
					return 0, err
				}
				if nx != stop {
					r.UnreadRune()
					delim = ' '
					break
				}
				if stop == '\'' {
					continue
				}
			}
			token.WriteRune(r1)
		}
	} else {
		r.UnreadRune()
		for {
			r1, _, err := r.ReadRune()
			if err != nil {
				return 0, err
			}
			if unicode.IsSpace(r1) {
				delim = ' '
				break
			}
			if r1 == ';' || r1 == ',' || r1 == '/' || r1 == '=' {
				delim = r1
				break
			}
			token.WriteRune(r1)
		}
	}

	if unicode.IsSpace(delim) {
		if err := skipSpaces(r); err != nil {
			return 0, err
		}
		r1, _, err := r.ReadRune()
		if err != nil {
			return 0, err
		}
		if r1 == ';' || r1 == ',' || r1 == '/' || r1 == '=' {
			delim = r1
		} else {
			r.UnreadRune()
		}
	}
	return delim, nil
}

func skipSpaces(r *bufio.Reader) error {
	for {
		r1, _, err := r.ReadRune()
		if err != nil {
			return err
		}

		// a comment
		if r1 == '[' {
			if err := skipComment(r); err != nil {
				return err
			}
			continue
		}

		if !unicode.IsSpace(r1) {
			r.UnreadRune()
			return nil
		}
	}
}

func skipComment(r *bufio.Reader) error {
	for {
		r1, _, err := r.ReadRune()
		if err != nil {
			return err
		}

		// a comment
		if r1 == ']' {
			return nil
		}
	}
}
