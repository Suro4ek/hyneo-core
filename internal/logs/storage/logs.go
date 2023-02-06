package storage

import (
	"hyneo/internal/logs"
	"hyneo/pkg/mysql"
)

type logsStorage struct {
	client *mysql.Client
}

func NewLogsStorage() logs.LogsStorage {
	return &logsStorage{}
}

func (l logsStorage) Logs(logs *logs.Logs) error {
	return l.client.DB.Create(logs).Error
}
