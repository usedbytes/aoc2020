package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type policy func(i, j int, c byte, password string) bool

func oldPolicy(i, j int, c byte, password string) bool {
	count := strings.Count(password, string(c))
	if count >= i && count <= j {
		return true
	}

	return false
}

func xor(a, b bool) bool {
	return (a && !b) || (b && !a)
}

func newPolicy(i, j int, c byte, password string) bool {
	return xor(password[i-1] == c, password[j-1] == c)
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("Usage: %s POLICY INPUT", os.Args[0])
	}

	policies := map[string]policy{
		"old": oldPolicy,
		"new": newPolicy,
	}

	policyFunc, ok := policies[os.Args[1]]
	if !ok {
		return fmt.Errorf("unrecognised policy: %s", os.Args[1])
	}

	f, err := os.Open(os.Args[2])
	if err != nil {
		return err
	}
	defer f.Close()

	numValid := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		var i, j int
		var r byte
		var password string

		n, err := fmt.Sscanf(line, "%d-%d %c: %s", &i, &j, &r, &password)
		if n != 4 {
			return fmt.Errorf("couldn't parse line: %s", line)
		} else if err != nil {
			return nil
		}

		if policyFunc(i, j, r, password) {
			numValid++
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println(numValid)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
