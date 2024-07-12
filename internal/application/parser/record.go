package parser

import (
	"strings"

	"github.com/leduardobueno/cw-test/internal/config"
)

type Record struct {
	Type string
	Data string
}

func newRecord(record string, conf *config.LogRecord) *Record {
	if len(record) < conf.Offset {
		return &Record{}
	}

	parsed := strings.Split(record[conf.Offset:], conf.Separator)
	if len(parsed) < 2 {
		return &Record{}
	}

	return &Record{
		Type: parsed[0],
		Data: parsed[1],
	}
}

func (r *Record) parsePlayerInfo(conf *config.PlayerLogRecord) (id, name string) {
	playerData := strings.Split(r.Data, conf.IdDelimiter)
	if len(playerData) < 2 {
		return
	}

	idData := playerData[0]
	nameData := strings.Split(playerData[1], conf.NameDelimiter)
	if idData == "" || len(nameData) == 0 || nameData[0] == "" {
		return
	}

	return strings.Trim(idData, " "), strings.Trim(nameData[0], " ")
}

func (r *Record) parseKill(conf *config.KillLogRecord) (killerId, deadId, causeId string) {
	killData := strings.Split(strings.Trim(r.Data, " "), conf.Separator)
	if len(killData) < 3 {
		return
	}

	return killData[conf.KillerIdx], killData[conf.DeadIdx], killData[conf.CauseIdx]
}
