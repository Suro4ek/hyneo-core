package logs

type LogsService struct {
	storage LogsStorage
}

func NewLogsService() *LogsService {
	return &LogsService{}
}

func (s LogsService) Join(
	playerId int64,
	serverName string,
	message string) error {
	return s.storage.Logs(&Logs{
		UserID:     playerId,
		ServerName: serverName,
		ActionType: "join",
		Message:    message,
	})
}

func (s LogsService) Quit(
	playerId int64,
	serverName string,
	message string) error {
	return s.storage.Logs(&Logs{
		UserID:     playerId,
		ServerName: serverName,
		ActionType: "quit",
		Message:    message,
	})
}

func (s LogsService) Message(
	playerId int64,
	serverName string,
	message string) error {
	return s.storage.Logs(&Logs{
		UserID:     playerId,
		ServerName: serverName,
		ActionType: "message",
		Message:    message,
	})
}

func (s LogsService) Command(
	playerId int64,
	serverName string,
	message string) error {
	return s.storage.Logs(&Logs{
		UserID:     playerId,
		ServerName: serverName,
		ActionType: "command",
		Message:    message,
	})
}
