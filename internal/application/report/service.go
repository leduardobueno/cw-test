package report

import (
	"fmt"
	"sort"

	"github.com/leduardobueno/cw-test/internal/application/parser"
)

type Service struct {
	meansOfDeath []string
}

func New(meansOfDeath []string) *Service {
	return &Service{meansOfDeath: meansOfDeath}
}

func (s *Service) GetMatchesGroupedInfo(matches []*parser.Match) []map[string]matchReport {
	matchesReport := make([]map[string]matchReport, len(matches))
	for i, currentMatch := range matches {
		killsByPlayer := currentMatch.GetKillsByPlayerNames()
		killsRanking := s.sortMapStringByIntValue(killsByPlayer)

		gameCounter := fmt.Sprintf("game_%d", i+1)
		matchesReport[i] = map[string]matchReport{
			gameCounter: {
				TotalKills: currentMatch.Kills,
				Players:    currentMatch.GetPlayerNames(),
				Kills:      killsRanking,
			},
		}
	}

	return matchesReport
}

func (s *Service) GetDeathCausesReport(matches []*parser.Match) []map[string]matchDeathCausesReport {
	matchesReport := make([]map[string]matchDeathCausesReport, len(matches))
	for i, currentMatch := range matches {
		killsByMeans := currentMatch.GetKillsByMeans(s.meansOfDeath)
		killsRanking := s.sortMapStringByIntValue(killsByMeans)

		gameCounter := fmt.Sprintf("game_%d", i+1)
		matchesReport[i] = map[string]matchDeathCausesReport{
			gameCounter: {KillsByMeans: killsRanking},
		}
	}

	return matchesReport
}

func (s *Service) sortMapStringByIntValue(toBeSorted map[string]int) sortedMapStringInt {
	sortableSlice := make(sortedMapStringInt, 0, len(toBeSorted))
	for mapKey, valueToBeSorted := range toBeSorted {
		sortableSlice = append(
			sortableSlice,
			mapStringInt{key: mapKey, value: valueToBeSorted},
		)
	}

	sort.Slice(sortableSlice, func(i, j int) bool {
		return sortableSlice[i].value > sortableSlice[j].value
	})

	return sortableSlice
}
