package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type FilesCLI struct {
	List   FileListCmd   `kong:"cmd,help='get list of all created snapshots'"`
	Show   FileShowCmd   `kong:"cmd,help='display detailed information about snapshot'"`
	Upload FileUploadCmd `kong:"cmd,help='upload file or directory, does not check file type'"`
	Delete FileDeleteCmd `kong:"cmd,help='delete uploaded file or directory'"`
}

type FileListCmd struct{}

func (c *FileListCmd) Run(ctx *Context) error {
	dirs, err := ctx.client.FilesListDirs()
	if err != nil {
		return err
	}

	fmt.Println("List of directories:")
	for _, dir := range dirs {
		fmt.Printf(" * %s\n", dir)
	}
	return nil
}

type FileShowCmd struct {
	Name string `kong:"arg"`
}

func (c *FileShowCmd) Run(ctx *Context) error {
	files, err := ctx.client.FilesListFiles(c.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Files in '%s':\n", c.Name)
	for _, file := range files {
		fmt.Printf(" * %s\n", file)
	}

	return nil
}

type FileUploadCmd struct {
	Dir  string `kong:"arg"`
	Path string `kong:"arg"`
}

func (c *FileUploadCmd) Run(ctx *Context) error {
	var files []string

	fi, err := os.Stat(c.Path)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// upload all files in directory
		err := filepath.Walk(c.Path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	case mode.IsRegular():
		files = append(files, c.Path)
	}

	uploaded, err := ctx.client.FilesUpload(c.Dir, files)
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded files:\n")
	for _, file := range uploaded {
		fmt.Printf(" * %s\n", file)
	}
	return nil
}

type FileDeleteCmd struct {
	Dir  string  `kong:"arg"`
	File *string `kong:"arg,optional"`
}

func (c *FileDeleteCmd) Run(ctx *Context) error {
	var err error

	if c.File != nil {
		err = ctx.client.FilesDeleteFile(c.Dir, *c.File)
		if err != nil {
			return err
		}
		fmt.Printf("Delete file '%s/%s'\n", c.Dir, *c.File)
	} else {
		err = ctx.client.FilesDeleteDir(c.Dir)
		if err != nil {
			return err
		}
		fmt.Printf("Delete directory '%s'\n", c.Dir)
	}
	return nil
}
