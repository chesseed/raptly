package main

import (
	"fmt"
)

type PkgsCLI struct {
	Search PkgSearchCmd `kong:"cmd,help='Search whole package database for packages matching query. Requires Aptly Server 1.6.0'"`
	// Show PkgShowcmd `kong:"cmd,help='Display details about packages from whole package database. Like 'search' with more information'"`
}

type PkgSearchCmd struct {
	Query string `kong:"arg"`
}

func (c *PkgSearchCmd) Run(ctx *Context) error {
	pkgs, err := ctx.client.PackagesSearch(c.Query)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		fmt.Printf("%s\n", pkg)
	}
	return nil
}

// TODO
// type PkgShowcmd struct {
// 	Query          string `kong:"arg"`
// 	WithFiles      bool   `kong:"name='with-files'"`
// 	WithReferences bool   `name:"name='with-references'"`
// }

// func (c *PkgShowcmd) Run(ctx *Context) error {
// 	pkgs, err := ctx.client.PackagesSearchDetailed(c.Query)
// 	if err != nil {
// 		return err
// 	}
// 	// TODO handle flags
// 	for _, pkg := range pkgs {
// 		fmt.Printf("%s\n", pkg.ShortKey)
// 	}
// 	return nil
// }
