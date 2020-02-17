package dir

import "os"

type FileEntry struct {
	Name     string
	DirName  string
	PathName string
	Mode     os.FileMode
	Size     int64
}

type DirEntry struct {
	Name     string
	DirName  string
	PathName string
	files    []*FileEntry
	dirs     []*DirEntry
}

func newFileEntry(dirPath, pathSep string, f os.FileInfo) (e *FileEntry) {

	return &FileEntry{
		Name:     f.Name(),
		DirName:  dirPath,
		PathName: dirPath + pathSep + f.Name(),
		Size:     f.Size(),
		Mode:     f.Mode(),
	}
}

func newDirEntry(dirPath, pathSep string, f os.FileInfo) (e *DirEntry) {
	return &DirEntry{
		Name:     f.Name(),
		DirName:  dirPath,
		PathName: dirPath + pathSep + f.Name(),
	}
}

func (d *DirEntry) GetFiles() []*FileEntry {

	return d.files
}

func (d *DirEntry) GetDirs() []*DirEntry {

	return d.dirs
}
