package storage

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Manager struct {
	Store Storage
}

func (m *Manager) UploadFolder(bucket, prefix, dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := d.Info()
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		key := filepath.ToSlash(filepath.Join(prefix, rel))

		return m.Store.PutObject(bucket, key, file, info.Size())
	})
}

func (m *Manager) DownloadFile(bucket, key, outputPath string) error {
	reader, err := m.Store.GetObject(bucket, key)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}