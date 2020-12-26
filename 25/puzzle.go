package main

import (
	"fmt"
	"os"
	"strconv"
)

func CryptoRounds(subject, value, rounds int) (int, int) {
	for i := 0; i < rounds; i++ {
		value = value * subject
		value = value % 20201227
	}

	return subject, value
}

func run() error {
	// The order doesn't really matter, but let's assume card then door
	var cardPubKey, doorPubKey int

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}
	cardPubKey = n

	n, err = strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}
	doorPubKey = n

	fmt.Println(cardPubKey, doorPubKey)

	value := 1
	subject := 7

	cardRounds := -1
	doorRounds := -1

	for round := 1; cardRounds < 0 || doorRounds < 0; round++ {
		subject, value = CryptoRounds(subject, value, 1)
		if value == cardPubKey {
			cardRounds = round
		}
		if value == doorPubKey {
			doorRounds = round
		}
	}
	fmt.Println("card rounds:", cardRounds)
	fmt.Println("door rounds:", doorRounds)

	_, cardEncKey := CryptoRounds(doorPubKey, 1, cardRounds)
	_, doorEncKey := CryptoRounds(cardPubKey, 1, doorRounds)

	if cardEncKey != doorEncKey {
		return fmt.Errorf("encryption keys don't match: %v != %v", cardEncKey, doorEncKey)
	}
	fmt.Println("encryption key:", cardEncKey)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
