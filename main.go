package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"log"
)

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flagP := flag.Int("p", 0, "Amount percent")
	flag.Parse()

	if *flagP == 0 {
		return fmt.Errorf("\"p\" value is %d", *flagP)
	}

	rdr := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		do(scanner.Text())
	}
	println(*flagP)
	return nil
}

func do(s string) {
	logger.Println(s)
}
