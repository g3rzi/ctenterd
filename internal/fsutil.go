package internal

import (
	"os"
	"time"
)

type DirEntry struct {
	Name    string
	IsDir   bool
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

func ReadDir(path string) ([]*DirEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}
	
	var entries []*DirEntry
	for _, info := range fileInfos {
		entry := &DirEntry{
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}
