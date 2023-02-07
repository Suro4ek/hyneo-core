package storage

import (
	"hyneo/internal/logs"
	"hyneo/pkg/mysql"
)

type logsStorage struct {
	client *mysql.Client
}

func NewLogsStorage(client *mysql.Client) logs.LogsStorage {
	return &logsStorage{
		client: client,
	}
}

func (l logsStorage) Logs(logs *logs.Logs) error {
	return l.client.DB.Create(logs).Error
}
