package postgres

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

type ListenerFactory struct {
	cfg config.Database
	db  *postgres.Storage
}

func NewListenerFactory(cfg config.Database, db *postgres.Storage) *ListenerFactory {
	return &ListenerFactory{
		cfg: cfg,
		db:  db,
	}
}

func (factory *ListenerFactory) CreateListener() storage.Listener {
	return NewNotificator(factory.cfg, factory.db)
}
