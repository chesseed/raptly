package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	aptly "raptly/pkg/rest-aptly"

	"pault.ag/go/debian/control"
)

type RepoCLI struct {
	List    RepoListCmd    `kong:"cmd,help='get list of all local package repositories on server'"`
	Show    RepoShowCmd    `kong:"cmd,help='display information about local repository, possibly listing all packages in the repository'"`
	Create  RepoCreateCmd  `kong:"cmd,help='create local package repository'"`
	Edit    RepoEditCmd    `kong:"cmd,help='change metadata of local repository'"`
	Drop    RepoDropCmd    `kong:"cmd,help='deletes information about local package repository'"`
	Rename  RepoRenameCmd  `kong:"cmd,help='change name of the local repository'"`
	Add     RepoAddCmd     `kong:"cmd,help='add file(s) to repository'"`
	Include RepoIncludeCmd `kong:"cmd,help='process .changes file or directory for upload'"`
}

type RepoListCmd struct{}

func (c *RepoListCmd) Run(ctx *Context) error {
	repos, err := ctx.client.ReposList()

	if err != nil {
		return err
	}

	fmt.Println("List of local repos:")
	for _, repo := range repos {
		fmt.Printf(" * [%s]\n", repo.Name)
	}
	return nil
}

type RepoShowCmd struct {
	WithPackages bool   `kong:"name='with-packages'"`
	Newest       bool   `kong:"name='newest',help='only show the newest version of each package, implies with-packages'"`
	Name         string `kong:"arg"`
}

func (c *RepoShowCmd) Run(ctx *Context) error {
	repo, err := ctx.client.ReposShow(c.Name)
	if err != nil {
		return err
	}

	conf := aptly.ListPackagesOptions{
		MaximumVersion: c.Newest,
	}
	packages, err := ctx.client.ReposListPackages(c.Name, conf)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", repo.Name)
	fmt.Printf("Comment: %s\n", repo.Comment)
	fmt.Printf("Default Distribution: %s\n", repo.DefaultDistribution)
	fmt.Printf("Default Component: %s\n", repo.DefaultComponent)
	fmt.Printf("Number of packages: %v\n", len(packages))
	if c.WithPackages || c.Newest {
		for _, pkg := range packages {
			fmt.Printf("  %s\n", pkg.Key)
		}
	}

	// TODO packages flag
	return nil
}

type RepoCreateCmd struct {
	Name         string  `kong:"arg"`
	Comment      *string `kong:"name='comment'"`
	Component    *string `kong:"name='component'"`
	Distribution *string `kong:"name='distribution'"`
}

func (c *RepoCreateCmd) Run(ctx *Context) error {
	opts := aptly.RepoCreateOptions{
		Comment:             *c.Comment,
		DefaultComponent:    *c.Component,
		DefaultDistribution: *c.Distribution,
	}

	repo, err := ctx.client.ReposCreate(c.Name, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Repo [%s] successfully added.\n", repo.Name)
	return nil
}

type RepoEditCmd struct {
	Name         string  `kong:"arg"`
	Comment      *string `kong:"name='comment'"`
	Component    *string `kong:"name='component'"`
	Distribution *string `kong:"name='distribution'"`
}

func (c *RepoEditCmd) Run(ctx *Context) error {
	opts := aptly.RepoUpdateOptions{
		Comment:             *c.Comment,
		DefaultComponent:    *c.Component,
		DefaultDistribution: *c.Distribution,
	}

	repo, err := ctx.client.ReposEdit(c.Name, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Repo [%s] successfully updated.\n", repo.Name)
	return nil
}

type RepoRenameCmd struct {
	Name    string `kong:"arg"`
	NewName string `kong:"arg,name='new-name'"`
}

func (c *RepoRenameCmd) Run(ctx *Context) error {
	opts := aptly.RepoUpdateOptions{
		Name: c.NewName,
	}

	repo, err := ctx.client.ReposEdit(c.Name, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Local repository [%s] -> [%s] has been successfully renamed.\n", c.Name, repo.Name)
	return nil
}

type RepoDropCmd struct {
	Force bool   `kong:"optional"`
	Name  string `kong:"arg"`
}

func (c *RepoDropCmd) Run(ctx *Context) error {
	err := ctx.client.ReposDrop(c.Name, c.Force)
	if err != nil {
		return err
	}
	fmt.Printf("Local repo [%s] has been removed.\n", c.Name)
	return nil
}

type RepoAddCmd struct {
	ForceReplace bool `kong:"name='force-replace'"`
	// RemoveFiles  bool   `kong:"name='remove-files'"`
	Name string `kong:"arg"`
	Path string `kong:"arg"`
}

func (c *RepoAddCmd) Run(ctx *Context) error {
	dir := fmt.Sprintf("upload_%s", randSeq(8))
	filesToUpload := []string{}

	fi, err := os.Stat(c.Path)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// upload all debian files in directory
		err := filepath.Walk(c.Path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				extension := filepath.Ext(path)
				if isDebianFile(extension) {
					filesToUpload = append(filesToUpload, path)
					if extension == ".dsc" {
						referenced, err := getFilesFromDsc(path)
						if err != nil {
							return err
						}
						filesToUpload = append(filesToUpload, referenced...)
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		if len(filesToUpload) == 0 {
			return fmt.Errorf("no debian package files (*.[u]deb | *.dsc) found in dir '%s'", c.Path)
		}
	case mode.IsRegular():
		extension := filepath.Ext(c.Path)
		if isDebianFile(extension) {
			filesToUpload = append(filesToUpload, c.Path)
			if extension == ".dsc" {
				referenced, err := getFilesFromDsc(c.Path)
				if err != nil {
					return err
				}
				filesToUpload = append(filesToUpload, referenced...)
			}
		} else {
			return fmt.Errorf("file '%s' not a debian package (*.[u]deb | *.dsc)", c.Path)
		}
	}

	_, err = ctx.client.FilesUpload(dir, filesToUpload)
	if err != nil {
		return err
	}
	res, err := ctx.client.ReposAddDirectory(c.Name, dir, aptly.RepoAddOptions{ForceReplace: c.ForceReplace})
	if err != nil {
		return err
	}
	if len(res.FailedFiles) > 0 {
		return fmt.Errorf("failed files:\n%v", res.FailedFiles)
	}

	fmt.Printf("Added %s\n", c.Path)
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func isDebianFile(extension string) bool {
	return extension == ".deb" || extension == ".udeb" || extension == ".dsc"
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func getFilesFromDsc(path string) ([]string, error) {
	files := []string{}
	fileDir, filename := filepath.Split(path)
	dsc, err := control.ParseDscFile(path)
	if err != nil {
		return nil, err
	}
	for _, file := range dsc.Files {
		path := filepath.Join(fileDir, file.Filename)
		if checkFileExists(path) {
			files = append(files, path)
		} else {
			return nil, fmt.Errorf("file '%s' referenced by '%s' does not exist", file.Filename, filename)
		}
	}
	return files, nil
}

type RepoIncludeCmd struct {
	ForceReplace     bool   `kong:"name='force-replace'"`
	AcceptUnsigned   bool   `kong:"name='accept-unsigned'"`
	IgnoreSignatures bool   `kong:"name='ignore-signatures'"`
	Name             string `kong:"arg"`
	Path             string `kong:"arg"`
}

func (c *RepoIncludeCmd) Run(ctx *Context) error {
	dir := fmt.Sprintf("upload_%s", randSeq(8))
	filesToUpload := []string{}

	fi, err := os.Stat(c.Path)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// upload all debian changes files in directory
		err := filepath.Walk(c.Path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				extension := filepath.Ext(path)
				if extension == ".changes" {
					filesToUpload = append(filesToUpload, path)
					referenced, err := getFilesFromChanges(path)
					if err != nil {
						return err
					}
					filesToUpload = append(filesToUpload, referenced...)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		if len(filesToUpload) == 0 {
			return fmt.Errorf("no debian *.changes files found in dir '%s'", c.Path)
		}
	case mode.IsRegular():
		extension := filepath.Ext(c.Path)
		if extension == ".changes" {
			filesToUpload = append(filesToUpload, c.Path)
			referenced, err := getFilesFromChanges(c.Path)
			if err != nil {
				return err
			}
			filesToUpload = append(filesToUpload, referenced...)
		} else {
			return fmt.Errorf("file '%s' not a debian *.changes files", c.Path)
		}
	}

	_, err = ctx.client.FilesUpload(dir, filesToUpload)
	if err != nil {
		return err
	}
	res, err := ctx.client.ReposAddDirectory(c.Name, dir, aptly.RepoAddOptions{ForceReplace: c.ForceReplace})
	if err != nil {
		return err
	}
	if len(res.FailedFiles) > 0 {
		return fmt.Errorf("failed files:\n%v", res.FailedFiles)
	}

	fmt.Printf("Added %s\n", c.Path)
	return nil
}

func getFilesFromChanges(path string) ([]string, error) {
	files := []string{}
	fileDir, filename := filepath.Split(path)
	changes, err := control.ParseChangesFile(path)
	if err != nil {
		return nil, err
	}
	for _, file := range changes.Files {
		path := filepath.Join(fileDir, file.Filename)
		if checkFileExists(path) {
			files = append(files, path)
		} else {
			return nil, fmt.Errorf("file '%s' referenced by '%s' does not exist", file.Filename, filename)
		}
	}
	return files, nil
}
