gomon
=====

[![Build Status](https://travis-ci.org/c9s/gomon.png)](https://travis-ci.org/c9s/gomon)

go source file monitor, which restarts/rebuilds your go package automatically
while you are changing it.

Install
-------

    go get -u github.com/c9s/gomon

Usage
-----

    gomon [dir] -- [cmd]

    gomon     # watch current directory and build the package (the default behavior)

Monitoring Current Directory And Build Automatically:

    gomon -b

Monitoring Current Directory And Test Automatically:

    gomon -t

Monitoring Directory And Build Automatically:

    gomon -b path/to/package -b

Monitoring Directory And Build Automatically With Verbose Messages:

    gomon -b -x path/to/package

Monitoring With Custom Command:

    gomon src -- go run -x server.go # execute go run -x server.go
    gomon src -- go build -x package # execute go build -x package



Todo
-----

- Add configration file support.
- Command queue support.


Related Product
---------------

GoTray <http://gotray.extremedev.org/>


Contributors
------------

- Ask Bjørn Hansen
- Mattn mattn.jp

License
--------

MIT License

