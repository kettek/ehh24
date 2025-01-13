package res

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteFile writes a file to the disk.
func WriteFile(name string, data []byte) error {
	p := filepath.Join("res", name)
	if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// ReadFile reads a file from the disk, and if that fails, from an embedded resource.
func ReadFile(name string) ([]byte, error) {
	p := filepath.Join("res", name)
	file, err := os.Open(p)
	if err != nil {
		embed, err := f.ReadFile(name)
		if err != nil {
			return nil, err
		}
		return embed, nil
	}
	defer file.Close()

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadDir reads a directory from the disk, and if that fails, from an embedded directory.
func ReadDir(name string) ([]fs.DirEntry, error) {
	p := filepath.Join("res", name)
	entries, err := os.ReadDir(p)
	if err != nil {
		fmt.Println(err)
	}

	embedEntries, err := f.ReadDir(name)
	if err != nil {
		return nil, err
	}

	for _, e := range embedEntries {
		exists := false
		for _, e2 := range entries {
			if e.Name() == e2.Name() {
				exists = true
				continue
			}
		}
		if !exists {
			entries = append(entries, e)
		}
	}

	return entries, nil
}
