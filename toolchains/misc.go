package toolchains

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

func CreateBox(workdir string) (string, error) {
	tempDir, err := os.MkdirTemp(workdir, "box_*")
	if err != nil {
		return tempDir, errors.New("failed to create temp dir")
	}

	if err := os.Chown(tempDir, 1000, 1000); err != nil {
		return tempDir, errors.New("failed to set directory permission")
	}

	return tempDir, nil
}

func WriteSourceCodes(dir string, sourceFiles []model.SourceFile) error {
	for _, src := range sourceFiles {
		filePath := filepath.Join(dir, src.Filename)
		err := util.CreateROFile(filePath, src.SourceCode)
		if err != nil {
			return fmt.Errorf("failed to write %s", src.Filename)
		}
	}
	return nil
}
