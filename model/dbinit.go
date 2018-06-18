package model

import (
	"github.com/go-pg/pg"
)

var Db *pg.DB

func CloseDB() {
	if Db != nil {
		Db.Close()
	}
}

func InitDB(dbConfig pg.Options) error {
	Db = pg.Connect(&dbConfig)

	// Db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
	// 	query, err := event.FormattedQuery()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	log.Debugf("%s %s", time.Since(event.StartTime), query)
	// })

	return nil
}
