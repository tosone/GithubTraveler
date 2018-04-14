# GithubTraveler [![Build Status](https://travis-ci.org/EffDataAly/GithubTraveler.svg?branch=master)](https://travis-ci.org/EffDataAly/GithubTraveler)

Travel all the github users, orgs, repos.

### How to compile

- Need golang 1.9.0 or later version.
- Need gcc.
- Compile on Linux `make Linux`.
- Compile on MacOS `make Darwin`.
- Compile on Windows `make windows`.

### How to run

Get the help info:

``` bash
$ GithubTraveler-drawin help
Travel all of the github organizations, users and repository.

Usage:
  GithubTraveler [command]

Available Commands:
  crawler     Travel all of the github organizations, users and repository.
  help        Help about any command
  version     Get version

Flags:
  -h, --help   help for GithubTraveler

Use "GithubTraveler [command] --help" for more information about a command.
```

Get the compile version info:

``` bash
$ GithubTraveler-drawin version
GithubTraveler darwin/amd64 go1.10.1
Version: e03285c96fb2136ff18a2b39b532094c09a03cf2
BuildDate: 2018-03-30_02:58:58PM
BuildHash: e03285c96fb2136ff18a2b39b532094c09a03cf2
```

Run the crawler:

``` bash
$ GithubTraveler-drawin crawler
```

Run in docker:

``` bash
$ docker-compose up -d
```