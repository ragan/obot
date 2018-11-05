package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/ragan/bitmex-api-go"
	"golang.org/x/net/context"
)

var logger *log.Logger

var entries []entry

type entry struct {
	apiKey, apiSecret string
}

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
	reader := csv.NewReader(rdr)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Errorf("%v", err.Error())
	}
	entries = make([]entry, len(entries))
	for _, r := range records {
		entries = append(entries, entry{r[0], r[1]})
	}
	logger.Printf("Parsed entries: %v", entries)
	err = checkBalance(entries)
	if err != nil {
		return err
	}
	return nil
}

func checkBalance(en []entry) error {
	for _, e := range en {
		ctx := context.Background()
		config := swagger.NewConfiguration()
		config.Host = "https://testnet.bitmex.com"
		config.BasePath = config.Host + "/api/v1"
		config.AddDefaultHeader()

		c := swagger.NewAPIClient(config)

		c.UserApi.UserGetWallet(nil)

	}
	return nil
}
