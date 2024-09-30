// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package matrix implements a command to build a phylogenetic matrix from the
// data stored in a PhyData project.
package matrix

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/js-arias/command"
	"github.com/js-arias/phydata/matrix"
	"github.com/js-arias/phydata/matrix/dna"
	"github.com/js-arias/phydata/project"
)

var Command = &command.Command{
	Usage: `matrix
	[-f|--format <format>]
	[-o|--output <file>]
	[--taxa <file>] [--chars <file>]
	<project> <data-type>...`,
	Short: "build a phylogenetic data matrix",
	Long: `
Command matrix reads a PhyData project and builds a phylogenetic data matrix
with the data stored in the project.

The first argument is the name of the project file.

The second and following arguments, are the types of data that will be
included in the data matrix. Valid values are:

	obs	used for morphological characters
	dna	used for DNA sequences

By default, the matrix will be printed in the standard output. To define an
output file use the flag --output, or -o to define the file name.

by default, the matrix format is the TNT format. Use the flag -f or --format
to define a format. Valid formats are:

	tnt   used for tnt output (default)
	nexus used for nexus output

By default, all taxa in the project will be used to build the matrix. If the
flag --taxa is defined with a file, the taxa in that file will be used as the
terminals of the matrix, using the order given in the file. In the file each
line will be read as a taxon name. Blank lines and lines starting with '#'
will be ignored.

By default, when making a matrix with observations, all characters will be
used to build the matrix. If the flag --chars is defined with a file, the
characters in the file will be used in the given order. In the file each line
will be interpreted as a character. Blank lines and lines starting with '#'
will be ignored.
	`,
	SetFlags: setFlags,
	Run:      run,
}

var output string
var format string
var txLsFile string
var charFile string

func setFlags(c *command.Command) {
	c.Flags().StringVar(&output, "output", "", "")
	c.Flags().StringVar(&output, "o", "", "")
	c.Flags().StringVar(&txLsFile, "taxa", "", "")
	c.Flags().StringVar(&charFile, "chars", "", "")
	c.Flags().StringVar(&format, "format", "tnt", "")
	c.Flags().StringVar(&format, "f", "tnt", "")
}

func run(c *command.Command, args []string) (err error) {
	if len(args) < 1 {
		return c.UsageError("expecting project file")
	}
	if len(args) < 2 {
		return c.UsageError("expecting data type definitions")
	}

	p, err := project.Read(args[0])
	if err != nil {
		return fmt.Errorf("unable ot open project %q: %v", args[0], err)
	}

	var m *matrix.Matrix
	var coll *dna.Collection
	withData := false
	for _, a := range args[1:] {
		switch strings.ToLower(a) {
		case "obs":
			mf := p.Path(project.Observations)
			if mf == "" {
				return fmt.Errorf("undefined observations file")
			}
			m = matrix.New()
			if err := readObsFile(mf, m); err != nil {
				return fmt.Errorf("on project %q: %v", args[0], err)
			}
			withData = true
		case "dna":
			df := p.Path(project.DNA)
			if df == "" {
				return fmt.Errorf("undefined DNA file")
			}
			coll = dna.New()
			if err := readDNAFile(df, coll); err != nil {
				return fmt.Errorf("on project %q: %v", args[0], err)
			}
			withData = true
		}
	}
	if !withData {
		return fmt.Errorf("data types %v not defined in the project", args[1:])
	}

	out := c.Stdout()
	if output != "" {
		var f *os.File
		f, err = os.Create(output)
		if err != nil {
			return err
		}
		defer func() {
			e := f.Close()
			if e != nil && err == nil {
				err = e
			}
		}()
		out = f
	}

	switch strings.ToLower(format) {
	case "tnt":
		if err := printTNTMatrix(out, m, coll); err != nil {
			return err
		}
	case "nexus":
		if err := printNexusMatrix(out, m, coll); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format %q", format)
	}

	return nil
}

func readObsFile(name string, m *matrix.Matrix) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := m.ReadTSV(f); err != nil {
		return fmt.Errorf("while reading file %q: %v", name, err)
	}
	return nil
}

func readDNAFile(name string, c *dna.Collection) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.ReadTSV(f); err != nil {
		return fmt.Errorf("while reading file %q: %v", name, err)
	}
	return nil
}

type taxaer interface {
	Taxa() []string
}

func getNumTaxa(d ...taxaer) int {
	tn := make(map[string]bool)
	for _, v := range d {
		if reflect.ValueOf(v).IsNil() {
			continue
		}
		for _, tx := range v.Taxa() {
			tn[tx] = true
		}
	}

	return len(tn)
}

func getTaxaList(d ...taxaer) []string {
	tn := make(map[string]bool)
	for _, v := range d {
		if reflect.ValueOf(v).IsNil() {
			continue
		}
		for _, tx := range v.Taxa() {
			tn[tx] = true
		}
	}

	ls := make([]string, 0, len(tn))
	for n := range tn {
		ls = append(ls, n)
	}

	return ls
}

func validTaxNames(ls []string) map[string]string {
	m := make(map[string]string, len(ls))
	for _, n := range ls {
		v := n
		if strings.ContainsRune(v, '&') {
			v = strings.ReplaceAll(v, "&", "+")
		}
		if strings.ContainsRune(v, '"') {
			v = strings.ReplaceAll(v, `"`, "")
		}

		v = strings.Join(strings.Fields(v), "_")
		m[n] = v
	}
	return m
}

func getNumChars(chLs []string, m *matrix.Matrix, coll *dna.Collection) int {
	var nc int
	if m != nil {
		nc = len(m.Chars())
		if len(chLs) > 0 {
			nc = len(chLs)
		}
	}

	if coll != nil {
		for _, gene := range coll.Genes() {
			nc += coll.MaxLen(gene)
		}
	}

	return nc
}

func printTNTMatrix(w io.Writer, m *matrix.Matrix, coll *dna.Collection) error {
	var txLs []string
	if txLsFile != "" {
		var err error
		txLs, err = readTaxa(txLsFile)
		if err != nil {
			return err
		}
	}

	var chLs []string
	if charFile != "" {
		var err error
		chLs, err = readFileList(charFile)
		if err != nil {
			return err
		}
	}

	bw := bufio.NewWriter(w)

	nt := getNumTaxa(m, coll)
	if len(txLs) > 0 {
		nt = len(txLs)
	}
	nc := getNumChars(chLs, m, coll)

	fmt.Fprintf(bw, "mxram 250 ;\ntaxname +255 ;\nxread %d %d\n\n", nc, nt)
	if m != nil {
		fmt.Fprintf(bw, "&[num]\n")

		states := make(map[string]map[int]string)
		chars := m.Chars()
		if len(chLs) > 0 {
			chars = chLs
		}
		for _, c := range chars {
			st := m.States(c)
			stID := make(map[int]string, len(st))
			for i, s := range st {
				if i > 9 {
					break
				}
				stID[i] = s
			}
			states[c] = stID
		}

		ls := m.Taxa()
		if len(txLs) > 0 {
			ls = txLs
		}

		for _, tx := range ls {
			ntx := strings.Join(strings.Fields(tx), "_")
			fmt.Fprintf(bw, "%s\t", ntx)
			txSp := m.TaxSpec(tx)
			for _, c := range chars {
				na := false
				st := make(map[string]bool, len(states[c]))
				for _, sp := range txSp {
					obs := m.Obs(sp, c)
					if len(obs) == 0 {
						continue
					}
					if obs[0] == matrix.NotApplicable {
						na = true
						continue
					}
					if obs[0] == matrix.Unknown {
						continue
					}
					for _, o := range obs {
						st[o] = true
					}
				}
				if len(st) == 0 {
					if na {
						fmt.Fprintf(bw, "-")
						continue
					}
					fmt.Fprintf(bw, "?")
					continue
				}
				obSt := states[c]
				if len(st) > 1 {
					fmt.Fprintf(bw, "[")
					for i := 0; i < len(obSt); i++ {
						v := obSt[i]
						if !st[v] {
							continue
						}
						fmt.Fprintf(bw, "%d", i)
					}
					fmt.Fprintf(bw, "]")
					continue
				}
				for i := 0; i < len(obSt); i++ {
					v := obSt[i]
					if st[v] {
						fmt.Fprintf(bw, "%d", i)
						break
					}
				}
			}
			fmt.Fprintf(bw, "\n")
		}
		fmt.Fprintf(bw, "\n")
	}

	if coll != nil {
		for _, gene := range coll.Genes() {
			fmt.Fprintf(bw, "&[dna nogaps]\n")

			ls := coll.Taxa()
			if len(txLs) > 0 {
				ls = txLs
			}
			for _, tx := range ls {
				var seq string
				for _, spec := range coll.TaxSpec(tx) {
					for _, acc := range coll.GeneAccession(spec, gene) {
						s := coll.Sequence(spec, gene, acc)
						if countNucleotides(s) > countNucleotides(seq) {
							seq = s
						}
					}
				}
				if len(seq) == 0 {
					continue
				}
				ntx := strings.Join(strings.Fields(tx), "_")
				fmt.Fprintf(bw, "%s\t%s\n", ntx, seq)
			}
			fmt.Fprintf(bw, "\n")
		}
	}

	fmt.Fprintf(bw, ";\n\ncc - . ;\n\nproc /; \n")
	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}

func printNexusMatrix(w io.Writer, m *matrix.Matrix, coll *dna.Collection) error {
	var txLs []string
	if txLsFile != "" {
		var err error
		txLs, err = readTaxa(txLsFile)
		if err != nil {
			return err
		}
	}

	var chLs []string
	if charFile != "" {
		var err error
		chLs, err = readFileList(charFile)
		if err != nil {
			return err
		}
	}

	bw := bufio.NewWriter(w)

	fmt.Fprintf(bw, "#NEXUS\n\n")

	nt := getNumTaxa(m, coll)
	if len(txLs) > 0 {
		nt = len(txLs)
	}
	nc := getNumChars(chLs, m, coll)

	nMorf := getNumChars(chLs, m, nil)
	nDNA := getNumChars(nil, nil, coll)

	fmt.Fprintf(bw, "Begin data;\n")
	fmt.Fprintf(bw, "\tDimensions ntax=%d nchar=%d;\n", nt, nc)
	if nMorf > 0 && nDNA > 0 {
		fmt.Fprintf(bw, "\tFormat datatype=mixed(standard:1-%d,DNA:%d-%d) interleave=yes gap=- missing=?;\n\n", nMorf, nMorf+1, nc)
	} else if nMorf > 0 {
		fmt.Fprintf(bw, "\tFormat datatype=standard missing=?;\n\n")
	} else {
		fmt.Fprintf(bw, "\tFormat datatype=DNA interleave=yes gap=- missing=?;\n\n")
	}

	if len(txLs) == 0 {
		txLs = getTaxaList(m, coll)
	}
	names := validTaxNames(txLs)

	fmt.Fprintf(bw, "\tMatrix\n\n")

	if m != nil {
		fmt.Fprintf(bw, "[Morphology]\n")

		states := make(map[string]map[int]string)
		chars := m.Chars()
		if len(chLs) > 0 {
			chars = chLs
		}
		for _, c := range chars {
			st := m.States(c)
			stID := make(map[int]string, len(st))
			for i, s := range st {
				if i > 9 {
					break
				}
				stID[i] = s
			}
			states[c] = stID
		}

		for _, tx := range txLs {
			ntx := names[tx]
			fmt.Fprintf(bw, "%s\t", ntx)
			txSp := m.TaxSpec(tx)
			for _, c := range chars {
				na := false
				st := make(map[string]bool, len(states[c]))
				for _, sp := range txSp {
					obs := m.Obs(sp, c)
					if len(obs) == 0 {
						continue
					}
					if obs[0] == matrix.NotApplicable {
						na = true
						continue
					}
					if obs[0] == matrix.Unknown {
						continue
					}
					for _, o := range obs {
						st[o] = true
					}
				}
				if len(st) == 0 {
					if na {
						fmt.Fprintf(bw, "-")
						continue
					}
					fmt.Fprintf(bw, "?")
					continue
				}
				obSt := states[c]
				if len(st) > 1 {
					fmt.Fprintf(bw, "{")
					for i := 0; i < len(obSt); i++ {
						v := obSt[i]
						if !st[v] {
							continue
						}
						fmt.Fprintf(bw, "%d", i)
					}
					fmt.Fprintf(bw, "}")
					continue
				}
				for i := 0; i < len(obSt); i++ {
					v := obSt[i]
					if st[v] {
						fmt.Fprintf(bw, "%d", i)
						break
					}
				}
			}
			fmt.Fprintf(bw, "\n")
		}
		fmt.Fprintf(bw, "\n")
	}
	if coll != nil {
		for _, gene := range coll.Genes() {
			fmt.Fprintf(bw, "[%s]\n", gene)
			ns := coll.MaxLen(gene)

			for _, tx := range txLs {
				var seq string
				for _, spec := range coll.TaxSpec(tx) {
					for _, acc := range coll.GeneAccession(spec, gene) {
						s := coll.Sequence(spec, gene, acc)
						if countNucleotides(s) > countNucleotides(seq) {
							seq = s
						}
					}
				}
				ntx := names[tx]
				if len(seq) == 0 {
					fmt.Fprintf(bw, "%s\t", ntx)
					for i := 0; i < ns; i++ {
						fmt.Fprintf(bw, "?")
					}
					fmt.Fprintf(bw, "\n")
					continue
				}
				fmt.Fprintf(bw, "%s\t%s\n", ntx, seq)
			}
			fmt.Fprintf(bw, "\n")
		}
	}

	fmt.Fprintf(bw, "\t;\n\n")
	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}

func countNucleotides(seq string) float64 {
	num := 0.0
	for _, p := range seq {
		switch p {
		case 'a', 'c', 'g', 't', 'u':
			num += 1
		case 'm', 'r', 'w', 's', 'y', 'k':
			num += 0.5
		case 'v', 'h', 'd', 'b':
			num += 0.25
		}
	}
	return num
}

func readTaxa(name string) ([]string, error) {
	ls, err := readFileList(name)
	if err != nil {
		return nil, err
	}

	for i, n := range ls {
		n = canon(n)
		ls[i] = n
	}

	return ls, nil
}

func readFileList(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var ls []string
	for i := 1; ; i++ {
		ln, err := r.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("on file %q: line %d: %v", name, i, err)
		}

		n := strings.Join(strings.Fields(ln), " ")
		if n == "" {
			continue
		}
		if n[0] == '#' {
			continue
		}
		ls = append(ls, strings.ToLower(n))
	}

	return ls, nil
}

// Canon returns a taxon name
// in its canonical form.
func canon(name string) string {
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.Join(strings.Fields(name), " ")
	if name == "" {
		return ""
	}
	name = strings.ToLower(name)
	r, n := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[n:]
}
