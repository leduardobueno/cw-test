package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/leduardobueno/cw-test/internal/application/parser"
	"github.com/leduardobueno/cw-test/internal/application/report"
	"github.com/leduardobueno/cw-test/internal/config"
	"github.com/leduardobueno/cw-test/internal/handler"
)

func main() {
	cfg := config.FromConfigFile()
	prs := parser.New(cfg.LogFile)
	rep := report.New(cfg.MeansOfDeath)

	inputHandler := handler.NewHandler(prs, rep)
	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) > 1 {
		runFromCommand(inputHandler)
		return
	}

	runCLIReader(inputHandler, reader)
}

func runFromCommand(inputHandler *handler.Handler) {
	switch os.Args[1] {
	case "matches", "m":
		inputHandler.ParseLogFile()
		inputHandler.PrintMatchesReport()
		return
	case "deaths", "d":
		inputHandler.ParseLogFile()
		inputHandler.PrintDeathCausesReport()
		return
	default:
		fmt.Println("Invalid argument. Available options:")
		fmt.Println("  'matches' (m) to print matches report")
		fmt.Println("  'deaths'  (d) to print death causes report")
		return
	}
}

func runCLIReader(inputHandler *handler.Handler, reader *bufio.Reader) {
	for {
		fmt.Println("Enter one of the following options:")
		fmt.Println("  'parse'   (p) to parse the Quake 3 Arena log file")
		fmt.Println("  'matches' (m) to print matches report")
		fmt.Println("  'deaths'  (d) to print death causes report")
		fmt.Println("  'quit'    (q) to exit")
		fmt.Print(":")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "parse", "p":
			started := time.Now()
			inputHandler.ParseLogFile()
			fmt.Printf("Log file parsed successfully in %0.6f seconds\n\n", time.Since(started).Seconds())
		case "matches", "m":
			inputHandler.PrintMatchesReport()
		case "deaths", "d":
			inputHandler.PrintDeathCausesReport()
		case "quit", "q":
			fmt.Println("Exiting...")
			os.Exit(0)
			return
		default:
			fmt.Println("Invalid input. Please try again.")
		}
	}
}
