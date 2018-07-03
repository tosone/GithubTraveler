# GithubTraveler

[![Build Status](https://travis-ci.org/EffDataAly/GithubTraveler.svg?branch=master)](https://travis-ci.org/EffDataAly/GithubTraveler) [![Build status](https://ci.appveyor.com/api/projects/status/o9o5x1i2j9smx11o/branch/master?svg=true)](https://ci.appveyor.com/project/tosone21763/githubtraveler/branch/master) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/c2081378a1b84f0da0855a10787981cf)](https://www.codacy.com/app/EffDataAly/GithubTraveler?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=EffDataAly/GithubTraveler&amp;utm_campaign=Badge_Grade)

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
% cd docker/sqlite3 && docker-compose up -d --build
```
