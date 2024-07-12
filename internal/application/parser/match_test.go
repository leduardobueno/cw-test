package parser

import (
	"testing"
)

func TestMatch_GetPlayerNames(t *testing.T) {
	testCases := []struct {
		name     string
		players  map[string]string
		expected []string
	}{
		{
			name:     "empty",
			players:  make(map[string]string),
			expected: []string{},
		},
		{
			name:     "one player",
			players:  map[string]string{"1": "Player 1"},
			expected: []string{"Player 1"},
		},
		{
			name:     "two players",
			players:  map[string]string{"1": "Player 1", "2": "Player 2"},
			expected: []string{"Player 1", "Player 2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := Match{Players: tc.players}
			actual := match.GetPlayerNames()
			if len(actual) != len(tc.expected) {
				t.Errorf("expected %d player names, got %d", len(tc.expected), len(actual))
			}
			for i, name := range actual {
				if name != tc.expected[i] {
					t.Errorf("expected '%s', got '%s'", tc.expected[i], name)
				}
			}
		})
	}
}

func TestMatch_GetKillsByPlayer(t *testing.T) {
	players := map[string]string{"22": "Player 1", "11": "Player 2"}

	testCases := []struct {
		name      string
		players   map[string]string
		killsById map[string]int
		expected  map[string]int
	}{
		{
			name:      "no kills",
			players:   players,
			killsById: make(map[string]int),
			expected:  make(map[string]int),
		},
		{
			name:      "all players killed",
			players:   players,
			killsById: map[string]int{"22": 10, "11": 5},
			expected:  map[string]int{"Player 1": 10, "Player 2": 5},
		},
		{
			name:      "player 1 killed",
			players:   players,
			killsById: map[string]int{"22": 666},
			expected:  map[string]int{"Player 1": 666},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := Match{Players: tc.players, killsByPlayerId: tc.killsById}
			actual := match.GetKillsByPlayerNames()
			if len(actual) != len(tc.expected) {
				t.Errorf("expected %d players, got %d", len(tc.expected), len(actual))
			}
			for name, kills := range actual {
				if kills != tc.expected[name] {
					t.Errorf("expected '%d' kills for %s, got '%d'", tc.expected[name], name, kills)
				}
			}
		})
	}
}

func TestMatch_GetKillsByMeans(t *testing.T) {
	testCases := []struct {
		name      string
		killsById map[string]int
		expected  map[string]int
	}{
		{
			name:      "no kills",
			killsById: make(map[string]int),
			expected:  make(map[string]int),
		},
		{
			name:      "unmapped means of death",
			killsById: map[string]int{"66": 10, "99": 5, "UNKNOWN": 60},
			expected:  map[string]int{"CAUSE_ID_66": 10, "CAUSE_ID_99": 5, "CAUSE_ID_UNKNOWN": 60},
		},
		{
			name:      "mapped and unmapped means of death",
			killsById: map[string]int{"1": 10, "2": 5, "99": 60},
			expected:  map[string]int{"MOD_SHOTGUN": 10, "MOD_GAUNTLET": 5, "CAUSE_ID_99": 60},
		},
		{
			name:      "mapped means of death",
			killsById: map[string]int{"0": 10, "1": 5, "2": 60},
			expected:  map[string]int{"MOD_UNKNOWN": 10, "MOD_SHOTGUN": 5, "MOD_GAUNTLET": 60},
		},
	}

	meansOfDeath := []string{
		"MOD_UNKNOWN",
		"MOD_SHOTGUN",
		"MOD_GAUNTLET",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := Match{killsByMeansId: tc.killsById}
			actual := match.GetKillsByMeans(meansOfDeath)
			if len(actual) != len(tc.expected) {
				t.Errorf("expected %d means, got %d", len(tc.expected), len(actual))
			}
			for mean, kills := range actual {
				if kills != tc.expected[mean] {
					t.Errorf("expected '%d' kills for %s mean, got '%d'", tc.expected[mean], mean, kills)
				}
			}
		})
	}
}

func TestMatch_upsertPlayer(t *testing.T) {
	testCases := []struct {
		name       string
		playerId   string
		playerName string
	}{
		{
			name:       "insert id: 2, name: Player 1",
			playerId:   "2",
			playerName: "Player 1",
		},
		{
			name:       "insert id: 1, name: Player 2",
			playerId:   "1",
			playerName: "Player 2",
		},
		{
			name:       "update id: 2, name: Player 3",
			playerId:   "1",
			playerName: "Player 3",
		},
	}

	match := Match{Players: make(map[string]string)}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match.upsertPlayer(tc.playerId, tc.playerName)
			if match.Players[tc.playerId] != tc.playerName {
				t.Errorf("expected player name '%s', got '%s'", tc.playerName, match.Players[tc.playerId])
			}
		})
	}
}

func TestMatch_updateKills(t *testing.T) {
	worldId := "1022"
	testCases := []struct {
		name     string
		killerId string
		deadId   string
		causeId  string
		expected *Match
	}{
		{
			name:     "player kill",
			killerId: "2",
			deadId:   "3",
			causeId:  "4",
			expected: &Match{
				Kills:           1,
				killsByMeansId:  map[string]int{"4": 1},
				killsByPlayerId: map[string]int{"2": 1},
			},
		},
		{
			name:     "world kill",
			killerId: worldId,
			deadId:   "2",
			causeId:  "22",
			expected: &Match{
				Kills:           2,
				killsByMeansId:  map[string]int{"22": 1},
				killsByPlayerId: map[string]int{"2": 0},
			},
		},
		{
			name:     "another player kill",
			killerId: "2",
			deadId:   "3",
			causeId:  "4",
			expected: &Match{
				Kills:           3,
				killsByMeansId:  map[string]int{"4": 2},
				killsByPlayerId: map[string]int{"2": 1},
			},
		},
	}

	match := &Match{
		Kills:           0,
		killsByPlayerId: make(map[string]int),
		killsByMeansId:  make(map[string]int),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match.updateKills(tc.killerId, tc.deadId, tc.causeId, worldId)
			if match.Kills != tc.expected.Kills {
				t.Errorf("expected %d kills, got %d", tc.expected.Kills, match.Kills)
			}

			expected := tc.expected.killsByPlayerId[tc.causeId]
			actual := match.killsByPlayerId[tc.causeId]
			if expected != actual {
				t.Errorf("expected %d kills for %s mean ID, got %d", expected, tc.causeId, actual)
			}

			killsId := tc.killerId
			if tc.killerId == worldId {
				killsId = tc.deadId
			}

			expected = tc.expected.killsByPlayerId[killsId]
			actual = match.killsByPlayerId[killsId]
			if expected != actual {
				t.Errorf("expected %d kills for %s ID, got %d", expected, killsId, actual)
			}
		})
	}
}
