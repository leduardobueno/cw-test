package config

type LogFileConfig struct {
	FilePath     string
	Record       LogRecord
	RecordTypes  LogRecordTypes
	PlayerRecord PlayerLogRecord
	KillRecord   KillLogRecord
}

type LogRecord struct {
	Delimiter byte
	Offset    int
	Separator string
}

type LogRecordTypes struct {
	NewMatch   string
	PlayerInfo string
	Kill       string
}

type PlayerLogRecord struct {
	IdDelimiter   string
	NameDelimiter string
}

type KillLogRecord struct {
	Separator string
	KillerIdx int
	DeadIdx   int
	CauseIdx  int
	WorldId   string
}
