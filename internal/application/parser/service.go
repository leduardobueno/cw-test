package parser

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/leduardobueno/cw-test/internal/config"
)

type Service struct {
	Matches []*Match
	conf    *config.LogFileConfig
}

func New(logConfig *config.LogFileConfig) *Service {
	return &Service{conf: logConfig}
}

func (s *Service) ParseLogFile() error {
	logFile, openErr := os.Open(s.conf.FilePath)
	if openErr != nil {
		return errors.New("error opening file: " + openErr.Error())
	}
	defer logFile.Close()

	s.Matches = make([]*Match, 0)

	reader := bufio.NewReader(logFile)
	for {
		recordAsString, err := reader.ReadString(s.conf.Record.Delimiter)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		record := newRecord(recordAsString, &s.conf.Record)
		if record.Type == "" || record.Data == "" {
			continue
		}
		s.collectRecordData(record)
	}

	return nil
}

func (s *Service) collectRecordData(record *Record) {
	switch record.Type {
	case s.conf.RecordTypes.NewMatch:
		s.addNewMatch()
		break
	case s.conf.RecordTypes.PlayerInfo:
		id, name := record.parsePlayerInfo(&s.conf.PlayerRecord)
		s.getCurrentMatch().upsertPlayer(id, name)
		break
	case s.conf.RecordTypes.Kill:
		killerId, deadId, causeId := record.parseKill(&s.conf.KillRecord)
		s.getCurrentMatch().updateKills(killerId, deadId, causeId, s.conf.KillRecord.WorldId)
	}
}

func (s *Service) addNewMatch() {
	s.Matches = append(s.Matches, &Match{
		Kills:           0,
		Players:         make(map[string]string),
		killsByPlayerId: make(map[string]int),
		killsByMeansId:  make(map[string]int),
	})
}

func (s *Service) getCurrentMatch() *Match {
	matchesCount := len(s.Matches)
	if matchesCount == 0 {
		s.addNewMatch()
	}

	return s.Matches[matchesCount-1]
}
