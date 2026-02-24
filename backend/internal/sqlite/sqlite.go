// Package sqlite
// @author: liuyanfeng
// @date: 2026/1/9 15:16
// @description:
package sqlite

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init 初始化全局 DB
func Init(path string) {
	var err error
	DB, err = sql.Open("sqlite3", "file:"+path+"?_journal_mode=WAL&_foreign_keys=ON&_busy_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

// WithTx 强制事务执行
func WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
