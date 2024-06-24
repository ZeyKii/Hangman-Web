package hangman_classic

import (
	"bufio"
	"log"
	"os"
)

func Scan(s string) []string {
	lst_mot := []string{}
	f, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lst_mot = append(lst_mot, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lst_mot
}
