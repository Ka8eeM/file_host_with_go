package memory

import (
	"context"
	"errors"
	"host_final_release/db"
	"time"

	kitlog "github.com/go-kit/log"
)

type Repo struct {
	logger kitlog.Logger
	db     map[string]db.Model
}

func New(l kitlog.Logger) (*Repo, error) {

	logger := kitlog.With(l, "Service", "In memory DB")
	return &Repo{logger: logger, db: make(map[string]db.Model)}, nil
}

func (r *Repo) Insert(ctx context.Context, file *db.Model) error {

	r.db[file.Id] = *file
	return nil
}

func (r *Repo) Get(ctx context.Context, id string) (file db.Model, err error) {

	file, ok := r.db[id]
	if !ok {
		// error handling
		err = errors.New("file not found")
	}
	return
}

// Delete remove entry from DB
func (r *Repo) Delete(ctx context.Context, id string) error {
	file, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	if time.Now().Before(file.Expiration) {
		return errors.New("file did not expire yet")
	}
	delete(r.db, id)
	return nil
}
