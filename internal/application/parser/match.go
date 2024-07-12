package parser

import "strconv"

type Match struct {
	Kills           int
	Players         map[string]string
	killsByPlayerId map[string]int
	killsByMeansId  map[string]int
}

func (m *Match) GetPlayerNames() []string {
	players := make([]string, 0, len(m.Players))
	for _, name := range m.Players {
		players = append(players, name)
	}

	return players
}

func (m *Match) GetKillsByPlayerNames() map[string]int {
	killsByPlayer := make(map[string]int)
	for playerId, kills := range m.killsByPlayerId {
		playerName := m.Players[playerId]
		killsByPlayer[playerName] = kills
	}

	return killsByPlayer
}

func (m *Match) GetKillsByMeans(meansOfDeath []string) map[string]int {
	maxMeanId := len(meansOfDeath) - 1
	killsByMeans := make(map[string]int)
	for meanId, kills := range m.killsByMeansId {
		meanIdAsInt, err := strconv.Atoi(meanId)
		if err != nil || meanIdAsInt < 0 || meanIdAsInt > maxMeanId {
			killsByMeans["CAUSE_ID_"+meanId] = kills
			continue
		}

		meanName := meansOfDeath[meanIdAsInt]
		killsByMeans[meanName] = kills
	}

	return killsByMeans
}

func (m *Match) upsertPlayer(id, name string) {
	m.Players[id] = name
}

func (m *Match) updateKills(killerId, deadId, causeId, worldId string) {
	m.Kills++
	m.killsByMeansId[causeId]++
	if killerId != worldId {
		m.killsByPlayerId[killerId]++
		return
	}

	m.killsByPlayerId[deadId]--
}
