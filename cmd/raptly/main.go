package main

import (
	"crypto/tls"
	"fmt"
	"os"
	aptly "raptly/pkg/rest-aptly"

	"github.com/alecthomas/kong"
)

var Version = "unknown"

type Context struct {
	client *aptly.Client
}

var cli struct {
	Version kong.VersionFlag `name:"version" help:"Print version information and quit"`

	Url      string  `kong:"help='required,Aptly server API URL',env='RAPTLY_URL'"`
	Insecure bool    `kong:"optional,help='Allow insecure HTTPS connections'"`
	User     *string `kong:"help='HTTP basic auth username',env='RAPTLY_USER'"`
	BasicPW  *string `kong:"name='basic-pass',help='HTTP basic auth password',env='RAPTLY_BASIC_PASS'"`

	Repo     RepoCLI     `kong:"cmd,help='Repository management commands',group='Repo'"`
	Publish  publishCLI  `kong:"cmd,help='Published lists commands',group='publish'"`
	Snapshot SnapshotCLI `kong:"cmd,help='Snapshot lists commands',group='snapshot'"`
	Package  PkgsCLI     `kong:"cmd,help='Package search commands',group='package'"`
	Files    FilesCLI    `kong:"cmd,help='Uploaded file management commands',group='Files'"`
	Status   StatusCLI   `kong:"cmd,help='Aptly server status command',group='Status'"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Vars{"version": Version})

	client := aptly.NewClient(cli.Url)
	if cli.Insecure {
		client.GetClient().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if cli.User != nil {
		if cli.BasicPW != nil {
			client.GetClient().SetBasicAuth(*cli.User, *cli.BasicPW)
		} else {
			fmt.Fprintf(os.Stderr, "Basic auth username set but no password, define RAPTLY_BASIC_PASS environment variable or use --basic-pass\n")
			os.Exit(1)
		}
	}

	err := ctx.Run(&Context{client: client})
	ctx.FatalIfErrorf(err)

	os.Exit(0)
}
