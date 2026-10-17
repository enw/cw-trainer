package library

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

const maxFileSize = 2 << 20

type Store struct {
	path string
	lock *os.File
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "cw-looper", "library.json"), nil
}

// Open takes a lock on a separate file so atomic library replacement retains it.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		lock.Close()
		return nil, fmt.Errorf("library is already open, or cannot be locked: %w", err)
	}
	return &Store{path, lock}, nil
}
func (s *Store) Close() error {
	if s.lock == nil {
		return nil
	}
	err := s.lock.Close()
	s.lock = nil
	return err
}
func (s *Store) Load() (Library, error) {
	var data Library
	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		data = Fixture()
		return data, s.Save(data)
	}
	if err != nil {
		return data, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return data, err
	}
	if info.Size() > maxFileSize {
		return data, fmt.Errorf("library exceeds 2 MiB")
	}
	decoder := json.NewDecoder(io.LimitReader(file, maxFileSize+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("read library; existing file preserved: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return data, fmt.Errorf("library must contain one JSON document")
	}
	return data, data.Validate()
}
func (s *Store) Save(data Library) error {
	if s.lock == nil {
		return fmt.Errorf("library store is closed")
	}
	if err := data.Validate(); err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if len(encoded) > maxFileSize {
		return fmt.Errorf("library exceeds 2 MiB")
	}
	dir := filepath.Dir(s.path)
	file, err := os.CreateTemp(dir, ".library-*")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(name, s.path); err != nil {
		return err
	}
	folder, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer folder.Close()
	return folder.Sync()
}
