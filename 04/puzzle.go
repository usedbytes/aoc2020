package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ScanPassport(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "\n\n"); i >= 0 {
		// Data contains a blank line
		// Note: my input.txt uses \n newlines, this won't work for
		// Windows-style line endings
		return i + 1, bytes.TrimSpace(data[0 : i+1]), nil
	}

	if atEOF {
		// Return what we've got
		return len(data), bytes.TrimSpace(data), nil
	}

	// Request more data.
	return 0, nil, nil
}

type Passport struct {
	Fields map[string]string
}

func (p *Passport) String() string {
	return fmt.Sprintf("byr: %s, iyr: %s, eyr: %s hgt: %s, hcl: %s, ecl: %s, pid: %s",
		p.Fields["byr"],
		p.Fields["iyr"],
		p.Fields["eyr"],
		p.Fields["hgt"],
		p.Fields["hcl"],
		p.Fields["ecl"],
		p.Fields["pid"])
}

type fieldValidFunc func(val string) bool
type policy map[string]fieldValidFunc

func (p *Passport) Valid(rules policy) bool {
	requiredFields := []string{
		"byr",
		"iyr",
		"eyr",
		"hgt",
		"hcl",
		"ecl",
		"pid",
	}

	for _, f := range requiredFields {
		val, ok := p.Fields[f]
		if !ok {
			return false
		}

		if vf, ok := rules[f]; ok {
			if !vf(val) {
				return false
			}
		}
	}

	return true
}

func (p *Passport) UnmarshalText(text []byte) error {

	if p.Fields == nil {
		p.Fields = make(map[string]string)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(text))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text := scanner.Text()

		kv := strings.Split(text, ":")
		if len(kv) != 2 {
			return fmt.Errorf("couldn't parse as key:value: %s", text)
		}

		p.Fields[kv[0]] = kv[1]
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var yearRE *regexp.Regexp = regexp.MustCompile("[0-9]{4}")

func rangeYearField(val string, min, max int) bool {
	if !yearRE.MatchString(val) {
		return false
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return false
	}

	return i >= min && i <= max
}

func heightValid(val string) bool {
	var num int
	var unit string

	n, err := fmt.Sscanf(val, "%d%s", &num, &unit)
	if n != 2 {
		return false
	} else if err != nil {
		return false
	}

	switch unit {
	case "cm":
		return num >= 150 && num <= 193
	case "in":
		return num >= 59 && num <= 76
	default:
		return false
	}

	return false
}

var colorRE *regexp.Regexp = regexp.MustCompile("^#[0-9a-f]{6}$")

func colorValid(val string) bool {
	return colorRE.MatchString(val)
}

var pidRE *regexp.Regexp = regexp.MustCompile("^[0-9]{9}$")

func pidValid(val string) bool {
	return pidRE.MatchString(val)
}

func listValid(val string, vals []string) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}

func alwaysValid(val string) bool {
	return true
}

var strictRules policy = map[string]fieldValidFunc{
	"byr": func(val string) bool { return rangeYearField(val, 1920, 2002) },
	"iyr": func(val string) bool { return rangeYearField(val, 2010, 2020) },
	"eyr": func(val string) bool { return rangeYearField(val, 2020, 2030) },
	"hgt": heightValid,
	"hcl": colorValid,
	"ecl": func(val string) bool {
		return listValid(val, []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"})
	},
	"pid": pidValid,
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("Usage: %s RULES INPUT", os.Args[0])
	}

	policies := map[string]policy{
		"relaxed": nil,
		"strict":  strictRules,
	}

	policy, ok := policies[os.Args[1]]
	if !ok {
		return fmt.Errorf("unknown rules: %s", os.Args[1])
	}

	f, err := os.Open(os.Args[2])
	if err != nil {
		return err
	}
	defer f.Close()

	valid := 0

	scanner := bufio.NewScanner(f)
	scanner.Split(ScanPassport)
	for scanner.Scan() {
		var ppt Passport

		err := (&ppt).UnmarshalText(scanner.Bytes())
		if err != nil {
			return err
		}

		if ppt.Valid(policy) {
			valid++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println(valid)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
