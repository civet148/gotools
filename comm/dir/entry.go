package dir

import (
	"os"
	"time"
)

type FileEntry struct {
	Name     string
	DirName  string
	PathName string
	Mode     os.FileMode
	Size     int64
	Time     time.Time
}

type DirEntry struct {
	Name     string
	DirName  string
	PathName string
	Time     time.Time
	files    []*FileEntry
	dirs     []*DirEntry
}

func getSep(dirPath, pathSep string) string {
	if dirPath == pathSep {
		return ""
	}
	return pathSep
}

func newFileEntry(dirPath, pathSep string, f os.FileInfo) (e *FileEntry) {

	return &FileEntry{
		Name:     f.Name(),
		DirName:  dirPath,
		PathName: dirPath + getSep(dirPath, pathSep) + f.Name(),
		Size:     f.Size(),
		Mode:     f.Mode(),
		Time:     f.ModTime(),
	}
}

func newDirEntry(dirPath, pathSep string, f os.FileInfo) (e *DirEntry) {
	return &DirEntry{
		Name:     f.Name(),
		DirName:  dirPath,
		PathName: dirPath + getSep(dirPath, pathSep) + f.Name(),
		Time:     f.ModTime(),
	}
}

func (d *DirEntry) GetFiles() []*FileEntry {

	return d.files
}

func (d *DirEntry) GetDirs() []*DirEntry {

	return d.dirs
}
