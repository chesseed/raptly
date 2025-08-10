# raptly

CLI tool to manage aptly repositories via their REST API. The commands are similar to the [local aptly commands](https://www.aptly.info/doc/commands/)

## Important differnces

All flags use the long form with double hyphen `--` instead of single hypen `-`

## Currently not implemented

* signing options
* all mirror commands
* repo include/import/copy/move/remove/search
* snapshot verify/pull/filter/merge/search
* db commands
* task commands
* graph command

## Usage

Every raptly invocation requires the server url+port. The url can be passed via the `--url` flag or the `RAPTLY_URL` environment variable.

### HTTP basic authentification

HTTP basic auth is supported. Pass the user via `--user` flag or the `RAPTLY_USER` environment variable.  
The password is passed with `--basic-pass` flag or the `RAPTLY_BASIC_PASS` environment variable.  

### Self-signed HTTPS certificates

Currently only the option to ignore SSL errors is implemented `--insecure`
