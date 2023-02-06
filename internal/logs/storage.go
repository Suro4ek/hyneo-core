package logs

type LogsStorage interface {
	/*
		Logs - логирует действие игрока
	*/
	Logs(logs *Logs) error
}
