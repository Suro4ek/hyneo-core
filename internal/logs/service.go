package logs

type Service struct {
	storage LogsStorage
}

func NewLogsService(storage LogsStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) Join(
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

func (s Service) Quit(
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

func (s Service) Message(
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

func (s Service) Command(
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
