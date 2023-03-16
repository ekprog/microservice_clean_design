package app

import (
	"git.mills.io/prologic/bitcask"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
	"strconv"
)

func InitStorage() error {

	// Create dir of not exists
	p := viper.GetString("storage.path")
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(p, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "cannot mkdir for storage")
		}
	}
	return nil
}

type Storage struct {
	name string
	db   *bitcask.Bitcask
}

func NewStorage(name string) (*Storage, error) {
	cachePath := viper.GetString("storage.path")
	file := path.Join(cachePath, name)
	db, err := bitcask.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create storage")
	}
	return &Storage{
		db:   db,
		name: name,
	}, nil
}

func (s *Storage) PutString(key, value string) error {
	err := s.db.Put([]byte(key), []byte(value))
	if err != nil {
		return errors.Wrapf(err, "cannot put to storage %s", s.name)
	}
	return nil
}

func (s *Storage) GetString(key string) (*string, error) {
	has := s.db.Has([]byte(key))
	if !has {
		return nil, nil
	}
	value, err := s.db.Get([]byte(key))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get value (key=%s) from storage %s", key, s.name)
	}
	if value == nil {
		return nil, nil
	}
	r := string(value)
	return &r, nil
}

func (s *Storage) PutInt64(key string, value int64) error {
	return s.PutString(key, strconv.FormatInt(value, 10))
}

func (s *Storage) GetInt64(key string) (*int64, error) {
	str, err := s.GetString(key)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get string before getting int64 from storage %s", s.name)
	}
	if str == nil {
		return nil, nil
	}
	r, err := strconv.ParseInt(*str, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse int64 from value %s from storage %s", *str, s.name)

	}
	return &r, nil
}
