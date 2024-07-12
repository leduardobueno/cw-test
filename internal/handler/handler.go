package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/leduardobueno/cw-test/internal/application/parser"
	"github.com/leduardobueno/cw-test/internal/application/report"
)

type Handler struct {
	parser *parser.Service
	report *report.Service
}

func NewHandler(parser *parser.Service, report *report.Service) *Handler {
	return &Handler{
		parser: parser,
		report: report,
	}
}

func (h *Handler) ParseLogFile() {
	err := h.parser.ParseLogFile()
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) PrintMatchesReport() {
	matches := h.parser.Matches
	matchesReport := h.report.GetMatchesGroupedInfo(matches)
	printAsJson(matchesReport)
}

func (h *Handler) PrintDeathCausesReport() {
	matches := h.parser.Matches
	deathCausesReport := h.report.GetDeathCausesReport(matches)
	printAsJson(deathCausesReport)
}

func printAsJson(report any) {
	reportAsBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("error converting report to JSON: %s", err)
	}
	fmt.Println(string(reportAsBytes))
}
