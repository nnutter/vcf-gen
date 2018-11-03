package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func withFile(name string, flag int, perm os.FileMode, do func(f *os.File) error) error {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}
	if err := do(f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

const contactSize = 500

func genSuffixSets(areaCode, prefix int16) [][]string {
	xss := make([][]string, 0, 10)
	xs := make([]string, 0, contactSize)
	for i := 0; i <= 9999; i++ {
		xs = append(xs, fmt.Sprintf("(%03d) %03d-%04d", areaCode, prefix, i))
		if len(xs) == cap(xs) {
			xss = append(xss, xs)
			xs = make([]string, 0, contactSize)
		}
	}
	return xss
}

func main() {
	areaCode := int16(618)
	prefix := int16(420)

	tmpl, err := template.New("vcf").Parse(vcfTemplate)
	if err != nil {
		log.Fatal(err)
	}

	suffixSets := genSuffixSets(areaCode, prefix)
	for _, suffixes := range suffixSets {
		fname := fmt.Sprintf("%s.vcf", suffixes[0])
		err := withFile(fname, os.O_RDWR|os.O_CREATE, 0755, func(f *os.File) error {
			return tmpl.Execute(f, suffixes)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

const vcfTemplate = `BEGIN:VCARD
VERSION:3.0
PRODID:-//Block6//NutterInk//EN
N:;;;;
FN:{{ index . 0 }}
ORG:{{ index . 0 }};
{{range .}}{{printf "TEL;type=OTHER;type=VOICE:%s\n" .}}{{end -}}
X-ABShowAs:COMPANY
END:VCARD`
