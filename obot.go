// "Obot" is a simple tool for querying multiple orders using "Bitmex" api.
package main

import (
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"strings"
	"github.com/ragan/bitmex-client-go/bitmex"
	"encoding/json"
)

var logger *log.Logger

// Entries read from input (or file)
var entries []entry

// Commands. Executed in order:
// No Command, Check Balance, Delete All Orders, Create Order
const (
	CmdNoCmd = 1 << iota
	CmdBalance //todo
	CmdDeleteAll
	CmdOrder
)

// Default command. No action will be taken.
var command = CmdNoCmd

// Flags
var (
	// Data file location
	file string
	// True when deleting orders selected
	del bool
	// True when order command selected
	order bool
	// Symbol (e.g. XBTUSD)
	symbol string
	// Order quantity
	orderQty int
	// Price
	price int
)

type cmdFun func(c bitmex.Client) (o *bitmex.Order, e *bitmex.Error)

type entry struct {
	apiKey, apiSecret string
}

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags)

	flag.StringVar(&file, "file", "", "Data file location")

	flag.BoolVar(&del, "d", false, "Deleting all orders if selected")
	flag.BoolVar(&order, "o", false, "Placing orders if selected")
	flag.StringVar(&symbol, "sym", "", "Symbol for placing orders")
	flag.IntVar(&orderQty, "qty", 0, "Selected order quantity")
	flag.IntVar(&price, "price", 0, "Selected order price")

	flag.Parse()

	// Parsing commands
	if del {
		command = command | CmdDeleteAll
	}
	if order {
		command = command | CmdOrder
	}

	logger.Printf("Filename is \"%s\"", file)

	data := ""

	if file != "" {
		logger.Printf("Reading from file %s\n", file)
		d, err := ioutil.ReadFile(file)
		if err != nil {
			exit(err)
		}
		data = string(d)
	} else {
		logger.Printf("Reading from input...")
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			exit(err)
		}
		data = string(in)
	}

	if err := run(data, command); err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func run(val string, command int) error {
	reader := csv.NewReader(strings.NewReader(val))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	entries = make([]entry, len(entries))
	for _, r := range records {
		logger.Printf("Parsing record: %v", r)
		entries = append(entries, entry{r[0], r[1]})
	}
	logger.Printf("Entry count: %d", len(entries))
	return do(entries, command)
}

func do(en []entry, command int) error {
	var clients = make([]bitmex.Client, 0)
	for _, e := range en {
		clients = append(clients, bitmex.NewTestNet(e.apiKey, e.apiSecret))
	}

	if command == CmdNoCmd {
		return fmt.Errorf("no command selected. Program will exit")
	}
	if isCommand(command, CmdDeleteAll) {
		logger.Printf("Executing \"Delete All Orders\"\n")
		exec(clients, func(c bitmex.Client) (o *bitmex.Order, e *bitmex.Error) {
			return bitmex.DeleteOrderAll(c)
		})
	}
	if isCommand(command, CmdOrder) {
		logger.Printf("Executing \"Place Order Command\"")
		if symbol == "" {
			return fmt.Errorf("symbol should have value")
		}
		if orderQty == 0 {
			return fmt.Errorf("qty should not be zero")
		}
		if price == 0 {
			return fmt.Errorf("price should not be zero")
		}
		exec(clients, func(c bitmex.Client) (o *bitmex.Order, e *bitmex.Error) {
			order := bitmex.OrderForm{
				Symbol:   symbol,
				OrderQty: float64(orderQty),
				Price:    float64(price),
			}
			prettyLog("Order:\n%v\n", order)
			return bitmex.PostOrder(c, &order)
		})
	}
	return nil
}

func isCommand(cVal, c int) bool {
	return c == cVal & c
}

func exec(clients []bitmex.Client, fn cmdFun) {
	for i := 0; i < len(clients); i++ {
		prettyLog("Running client:\n%s", clients[i])
		o, err := fn(clients[i])
		if err != nil {
			prettyLog("Error:\n%v\n", *err)
		} else {
			prettyLog("Result:\n%v\n", *o)
		}
	}
}

func prettyLog(format string, v interface{}) {
	s, _ := json.MarshalIndent(v, "", "  ")
	logger.Printf(format, string(s))
}
