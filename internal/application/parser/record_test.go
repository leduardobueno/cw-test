package parser

import (
	"testing"

	"github.com/leduardobueno/cw-test/internal/config"
)

func TestRecord_newRecord(t *testing.T) {
	conf := &config.LogFileConfig{
		Record: config.LogRecord{
			Delimiter: '\n',
			Offset:    7,
			Separator: ":",
		},
		RecordTypes: config.LogRecordTypes{
			NewMatch:   "InitGame",
			PlayerInfo: "ClientUserinfoChanged",
			Kill:       "Kill",
		},
	}

	testCases := []struct {
		name     string
		record   string
		expected *Record
	}{
		{
			name:     "empty record",
			record:   "",
			expected: &Record{},
		},
		{
			name:     "invalid record",
			record:   "invalid record",
			expected: &Record{},
		},
		{
			name:     "invalid record length",
			record:   "x",
			expected: &Record{},
		},
		{
			name:   "valid record",
			record: " 22:22 MockedType: MockedData...",
			expected: &Record{
				Type: "MockedType",
				Data: " MockedData...",
			},
		},
		{
			name:   "new game record",
			record: " 22:22 InitGame: MockedData...",
			expected: &Record{
				Type: conf.RecordTypes.NewMatch,
				Data: " MockedData...",
			},
		},
		{
			name:   "player info record",
			record: " 22:22 ClientUserinfoChanged: 2 n\\MockedName\\t...",
			expected: &Record{
				Type: conf.RecordTypes.PlayerInfo,
				Data: " 2 n\\MockedName\\t...",
			},
		},
		{
			name:   "kill record",
			record: " 22:22 Kill: 1 2 3: MockedData...",
			expected: &Record{
				Type: conf.RecordTypes.Kill,
				Data: " 1 2 3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := newRecord(tc.record, &conf.Record)
			if actual.Type != tc.expected.Type || actual.Data != tc.expected.Data {
				t.Errorf("expected '%s', got '%s'", tc.expected, actual)
			}
		})
	}
}

func TestRecord_parsePlayerInfo(t *testing.T) {
	type output struct {
		id   string
		name string
	}

	conf := &config.PlayerLogRecord{
		IdDelimiter:   "n\\",
		NameDelimiter: "\\t",
	}

	testCases := []struct {
		name     string
		data     string
		expected output
	}{
		{
			name:     "empty player data",
			data:     "",
			expected: output{},
		},
		{
			name:     "invalid player data",
			data:     "invalid data",
			expected: output{},
		},
		{
			name:     "missing player name",
			data:     " 1 n\\",
			expected: output{},
		},
		{
			name:     "missing player id",
			data:     "n\\MockedName\\t...",
			expected: output{},
		},
		{
			name: "valid player data",
			data: " 2 n\\MockedName\\t...",
			expected: output{
				id:   "2",
				name: "MockedName",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := Record{Data: tc.data}
			actual := output{}
			actual.id, actual.name = rec.parsePlayerInfo(conf)
			if actual.id != tc.expected.id || actual.name != tc.expected.name {
				t.Errorf("expected '%s', got '%s'", tc.expected, actual)
			}
		})
	}
}

func TestRecord_parseKill(t *testing.T) {
	type output struct {
		killerId string
		deadId   string
		causeId  string
	}

	conf := &config.KillLogRecord{
		Separator: " ",
		KillerIdx: 0,
		DeadIdx:   1,
		CauseIdx:  2,
	}

	testCases := []struct {
		name     string
		data     string
		expected output
	}{
		{
			name:     "empty kill data",
			data:     "",
			expected: output{},
		},
		{
			name:     "invalid kill data",
			data:     "invalid data",
			expected: output{},
		},
		{
			name:     "without separator",
			data:     "123",
			expected: output{},
		},
		{
			name: "valid kill data",
			data: " 1 2 3",
			expected: output{
				killerId: "1",
				deadId:   "2",
				causeId:  "3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := Record{Data: tc.data}
			actual := output{}
			actual.killerId, actual.deadId, actual.causeId = rec.parseKill(conf)
			if actual.killerId != tc.expected.killerId ||
				actual.deadId != tc.expected.deadId ||
				actual.causeId != tc.expected.causeId {
				t.Errorf("expected '%s', got '%s'", tc.expected, actual)
			}
		})
	}
}
