package filestorer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantchen2003/chunky/internal/util"
)

type LocalFileStore struct {
	dirPath string
}

// since we only check the dir exists in constructor, we assume
// dir is not deleted throughout all of server's lifetime
func NewLocalFileStore(dirPath string) (*LocalFileStore, error) {
	dirExists, err := util.DirExists(dirPath)
	if err != nil {
		return nil, err
	}

	if !dirExists {
		if err := os.Mkdir(dirPath, os.ModeDir); err != nil {
			return nil, err
		}
	}

	return &LocalFileStore{
		dirPath: dirPath,
	}, nil
}

func (lfs *LocalFileStore) Store(data []byte) (chunkId string, err error) {
	chunkId, err = util.GenerateRandomHexString(16)
	if err != nil {
		return "", err
	}

	chunkPath := filepath.Join(lfs.dirPath, chunkId)

	file, err := os.Create(chunkPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing bytes:", err)
		return
	}

	return chunkId, err
}
