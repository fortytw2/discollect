// +build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fortytw2/discollect/countries"
)

func main() {
	goPkg := "countries"

	// Parse template for go source
	t, err := template.New("_").Funcs(template.FuncMap{
		"countrySrc": countrySrc,
	}).Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// Parse csv data
	f, err := os.Open("countries.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.Comma = ';'
	csvr.Comment = '#'
	csvr.FieldsPerRecord = 27

	var allCountries []countries.Country
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		tr := make([]string, len(record))
		for i, rec := range record {
			tr[i] = strings.TrimSpace(rec)
		}

		country := countries.Country{
			ISO2: tr[0],
			ISO3: tr[1],
		}

		allCountries = append(allCountries, country)
	}

	// Prepare template context
	ctx := struct {
		PkgName   string
		Countries []countries.Country
	}{
		PkgName:   goPkg,
		Countries: allCountries,
	}

	var buf bytes.Buffer
	// Exec template
	if err := t.Execute(&buf, ctx); err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		os.Stderr.Write(buf.Bytes())
		log.Fatalf("fmt: %v", err)
	}

	outf, err := os.Create("countries_generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outf.Close()

	if _, err = io.Copy(outf, bytes.NewReader(src)); err != nil {
		log.Fatal(err)
	}
}

func countrySrc(country countries.Country) string {
	src := strings.Replace(fmt.Sprintf("%#v", country), "countries.", "", 1)
	return src
}

const tmpl = `// Code generated by generator.go DO NOT EDIT.
package {{ .PkgName }}

// All countries
var (
{{ range .Countries }}  {{ .ISO2 }} = {{ countrySrc . }}
{{ end }}
)

// Countries defines all countries with they ISO 3166-1 Alpha-2 code as key.
var Countries = map[string]Country{
{{ range .Countries }}"{{ .ISO2 }}": {{ .ISO2 }},
{{ end }}
}
`
