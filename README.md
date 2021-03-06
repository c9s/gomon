gomon
=====

[![CI Status](https://github.com/c9s/gomon/workflows/CI/badge.svg)](https://github.com/c9s/gomon/actions)

go source file monitor, which restarts/rebuilds your go package automatically
while you are changing it.

Install
-------

    go get -u github.com/c9s/gomon

Usage
-----

    gomon [dir] -- [cmd]

    gomon     # watch current directory and build the package (the default behavior)

Monitoring Current Directory And Format Automatically:

    gomon -f

Monitoring Current Directory And Build Automatically:

    gomon -b

Monitoring Current Directory And Test Automatically:

    gomon -t

Monitoring Current Directory And Install Automatically:

    gomon -i

You can run commands sequentialy by specifying multiple options above.
Monitoring Current Directory And Format, Build, Test and Install Automatically:

    gomon -f -b -t -i

Monitoring Directory And Build Automatically:

    gomon -b path/to/package

Monitoring Directory And Build Automatically With Verbose Messages:

    gomon -b -x path/to/package

Monitoring With Custom Command:

    gomon src -- go run -x server.go # execute go run -x server.go
    gomon src -- go build -x package # execute go build -x package


Screenshot
----------

![](https://raw.github.com/c9s/gomon/gh-pages/images/screenshot.png)

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
- Yasuhiro Matsumoto (a.k.a mattn)

License
--------

MIT License

