package service

import (
	"context"
	"errors"
	"host_final_release/db"
	"io"
	"strconv"
	"strings"
	"time"

	kitlog "github.com/go-kit/log"
)

// Item struct represent the strored items
type Item struct {
	Id         string
	Body       io.Reader
	FileName   string
	Expiration time.Time
	Length     int64
}

// Srvs implement service interface
type Srvs struct {
	logger  kitlog.Logger
	storage storage.Storage
	db      db.Database
}

// service is interface define the service function
type Service interface {
	Get(ctx context.Context, id string) (Item, error)
	Upload(ctx context.Context, opts Item) (string, error)
}

func (s Srvs) Upload(ctx context.Context, opts Item) (string, error) {

	// TODO validate file size

	f := &storage.File{
		Content:    opts.Body,
		Expiration: opts.Expiration,
	}

	err := s.storage.Save(ctx, f)
	if err != nil {
		return "", err
	}

	err = s.db.Insert(ctx, &db.Model{
		FileName:   opts.FileName,
		Id:         f.ID,
		Expiration: f.Expiration,
		Length:     opts.Length,
	})

	if err != nil {
		return "", err
	}
	return f.ID, nil
}

func (s Srvs) Get(ctx context.Context, id string) (Item, error) {
	// get expiration from id
	splitedID := strings.Split(id, "_")
	if len(splitedID) != 2 {
		s.logger.Log(
			"message", "invalid id",
			"id", id,
		)
		return Item{}, errors.New("invalid id")
	}
	unixNano, err := strconv.ParseInt(splitedID[0], 10, 64)
	if err != nil {
		s.logger.Log(
			"message", "failed ParseInt id",
			"error", err,
			"id", id,
		)
		return Item{}, err
	}

	// validate expiration
	// if expired delete and ignore errors
	if time.Now().After(time.Unix(0, unixNano)) {
		return Item{}, errors.New("File expired")
	}

	// get from DB
	fileModel, err := s.db.Get(ctx, id)
	if err != nil {
		return Item{}, err
	}
	// get from storage
	file, err := s.storage.Get(ctx, id)
	if err != nil {
		return Item{}, err
	}

	return Item{
		Id:         id,
		Body:       file,
		Expiration: fileModel.Expiration,
		FileName:   fileModel.FileName,
		Length:     fileModel.Length,
	}, nil

}

func New(logger kitlog.Logger, storage storage.Storage, db db.Database) *Srvs {
	return &Srvs{
		logger:  logger,
		storage: storage,
		db:      db,
	}
}
