package report

import "fmt"

type matchReport struct {
	TotalKills int                `json:"total_kills"`
	Players    []string           `json:"players"`
	Kills      sortedMapStringInt `json:"kills"`
}

type matchDeathCausesReport struct {
	KillsByMeans sortedMapStringInt `json:"kills_by_means"`
}

type mapStringInt struct {
	key   string
	value int
}

type sortedMapStringInt []mapStringInt

func (s sortedMapStringInt) MarshalJSON() ([]byte, error) {
	var result []byte
	result = append(result, '{')

	mapLen := len(s)
	for i := 0; i < mapLen; i++ {
		jsonByte := []byte(fmt.Sprintf(`"%s":%d`, s[i].key, s[i].value))
		result = append(result, jsonByte...)

		if i < mapLen-1 {
			result = append(result, ',')
		}
	}

	result = append(result, '}')
	return result, nil
}
