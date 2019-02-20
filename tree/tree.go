package tree

import (
	"fmt"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	startPath  = "/"
	RootFolder = newFolder(startPath, startPath)
)

type File struct {
	Name string
	Path string
}

type Folder struct {
	Name    string
	Path    string
	Files   []File
	Folders map[string]*Folder
}

func newFolder(name, path string) *Folder {
	return &Folder{name, path, []File{}, make(map[string]*Folder)}
}

func (f *Folder) getFolder(name string) *Folder {
	if nextF, ok := f.Folders[name]; ok {
		return nextF
	} else if f.Name == name {
		return f
	} else {
		return &Folder{}
	}
}

func (f *Folder) existFolder(name string) bool {
	for _, v := range f.Folders {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (f *Folder) addFolder(folderName, folderPath string) {
	if !f.existFolder(folderName) {
		f.Folders[folderName] = newFolder(folderName, folderPath)
	}
}

func (f *Folder) addFile(fileName, filePath string) {
	f.Files = append(f.Files, File{fileName, filePath})
}

func isFile(str string) bool {
	if path.Ext(str) != "" {
		return true
	}
	return false
}

func DeleteEmptyElements(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func CreateTree(s3Output *s3.ListObjectsV2Output) {

	for _, obj := range s3Output.Contents {
		filePath := obj.Key
		splitPath := DeleteEmptyElements(strings.Split(*filePath, "/"))
		tmpFolder := RootFolder
		for _, item := range splitPath {
			if isFile(item) {
				tmpFolder.addFile(item, *filePath)
			} else {
				if item != startPath {
					tmpFolder.addFolder(item, *filePath)
				}
				tmpFolder = tmpFolder.getFolder(item)
			}
		}
	}
}

func FindNode(rootFolder *Folder, findItem string) *Folder {
	found := newFolder("", "")

	if rootFolder.Name == findItem {
		return rootFolder
	}

	for _, folder := range rootFolder.Folders {
		if folder.Name == findItem {
			return folder
		}
		found = FindNode(folder, findItem)
	}

	return found
}

func OldTree() {
	arrayPaths := []map[string]string{
		{
			"id":       "1",
			"filePath": "",
		},
		{
			"id":       "1",
			"filePath": "peelz.here",
		},
		{
			"id":       "2",
			"filePath": "test1/",
		},
		{
			"id":       "3",
			"filePath": "test1/Nene_noises_for_1_32_minutes.mp4",
		},
		{
			"id":       "3",
			"filePath": "test1/neptune_all_the_meme.jpg",
		},
		{
			"id":       "3",
			"filePath": "test2/",
		},
		{
			"id":       "3",
			"filePath": "test2/america_chan_seijouki.png",
		},
		{
			"id":       "3",
			"filePath": "test2/bongo_cat_levan_polka_miku.mp4",
		},
		{
			"id":       "3",
			"filePath": "test3/",
		},
		{
			"id":       "3",
			"filePath": "test3/inside_test3.jpg",
		},
		{
			"id":       "3",
			"filePath": "test3/test4/",
		},
		{
			"id":       "3",
			"filePath": "test3/test4/second_level.jpg",
		},
		{
			"id":       "3",
			"filePath": "test3/test4/another_s2.mp3",
		},
	}

	startPath := "/"
	rootFolder := newFolder(startPath, startPath)

	for _, path := range arrayPaths {
		filePath := path["filePath"]
		splitPath := DeleteEmptyElements(strings.Split(filePath, "/"))
		tmpFolder := rootFolder
		for _, item := range splitPath {
			if isFile(item) {
				tmpFolder.addFile(item, filePath)
			} else {
				if item != startPath {
					tmpFolder.addFolder(item, filePath)
				}
				tmpFolder = tmpFolder.getFolder(item)
			}
		}
	}
	fmt.Println(rootFolder)
}
