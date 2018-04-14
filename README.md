# GithubTraveler [![Build Status](https://travis-ci.org/EffDataAly/GithubTraveler.svg?branch=master)](https://travis-ci.org/EffDataAly/GithubTraveler) [![Build status](https://ci.appveyor.com/api/projects/status/o9o5x1i2j9smx11o?svg=true)](https://ci.appveyor.com/project/tosone21763/githubtraveler)

Travel all the github users, orgs, repos.

### How to compile

- Need golang 1.10.0 or later version.
- Need gcc or Mingw.
- Compile just `make`.

### How to run

Get the help info:

<img src="doc/help.jpeg" width="600">

Get the compile version info:

<img src="doc/version.jpeg" width="400">

Run the crawler:

``` bash
% GithubTraveler-drawin crawler
```

Run in docker:

``` bash
% docker-compose up -d
```