package historyrepositories

import "github.com/uptrace/bun"

type HistoryRepository interface{}

type historyRepository struct {
	db *bun.DB
}

func NewHistoryRepository(db *bun.DB) HistoryRepository {
	return &historyRepository{db: db}
}
