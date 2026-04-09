package storage

type Repository interface {
	InsertRate(timestamp int64, bestAsk, bestBid string) (int64, error)
	Close() error
}
