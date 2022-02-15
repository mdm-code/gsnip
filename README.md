![(logo)](./logo.png)

# Gsnip

<div align="center">
<p>
	<a href="https://github.com/mdm-code/gsnip/actions?query=workflow%3ACI">
		<img alt="Build status" src="https://github.com/mdm-code/gsnip/workflows/gsnip%20CI/CD/badge.svg"
	</a>
    <a href="https://opensource.org/licenses/GPL-3.0" rel="nofollow">
        <img alt="GPL-3 license" src="https://img.shields.io/github/license/mdm-code/gsnip">
    </a>
</p>
</div>

This my personal snippet manager. It lets you find, insert, delete and list out
all snippets stored in a text file and written with straightforward, I believe,
syntax rules. My goal was to keep the program as simple as possible: it scans
the source file with snippets and offers an interface to interact with it.


## Usage

There are two parts of the workflow: one, `gsnipd`, a server running on Unix
Domain Socket handling connections, and `gsnip`, which is the client that
relies on `FD0` or `SDTIN` and simple sub-commands to send messages to the
server.

`gsnipd`, `gsnip` and all its subcommands print out useful information with
`-h|--help`.

First, you want to spin up the server with either of these commands:

```sh
gsnipd
gsnipd &>/dev/null &
```

The first one will write `STDERR` to the terminal so that you can see server
messages. The other one sends all messages to `/dev/null` and gets detached
from the current session.

You can also enable it as a daemon in `systemd` or `launchd` on MacOS so that
it runs on startup, and you don't have to mess around with it each time you
restart your computer.

The name of the source file used to store snippets can be passed as an argument
to the `gsnipd` program. If it isn't, the program will search for a `snippets`
file at the `gsnip` subdirectory in XDG data directories. It will error out if
it could not find one.

Then you can interact with the server using `gsnip` client like this:

```sh
echo [snip-name] | gsnip find
gsnip list
echo [snip-name] | gsnip delete
gsnip insert
gsnip reload
```

You can query the server with `find` for any snippet stored in the source file.
Alternatively, you can ask the server to `list` out all available snippets.
You can delete existing snippets with `delete` subcommand. You can also `insert`
new snippets through an editor or `STDIN`.

In order to add a new snippet right from the command line, the easy way would be
to use HereDoc like this, for instance:

```sh
gsnip insert << EOF
startsnip test "this is just a test snippet"
func test() bool {
	return true
}
endsnip
EOF
```

You can reload the source snippet file at the server runtime by calling the
`gsnip` client with the `reaload` subcommand, which is the equivalent of
sending `SIGHUP` to the process using `kill -1 [pid]`. The latter is annoying
because you have to find the process id with `ps` before sending the signal.
Another way would be to write PID to a known file that the client could access.

For the sake of clarity, there are five four-byte-long headers that a `gsnipd`
server can understand:

- @LST
- @RLD
- @FND
- @INS
- @DEL

They correspond to the client subcommands.

The idea was to use `gsnip` as an application agnostic tool. Since it operates
on standard file descriptors, it can be used in most Unix pipes and most
importantly `vim` through the use of `!` inside the editor. I do not like other
editors.


## Syntax

The snippet syntax looks like this:

```
startsnip NAME "COMMENT"
BODY
endsnip
```

`startsnip` and `endsnip` delimit the scope of a single snippet. They are used
by the parser to identify the start and the end. There few more rules that have
to be respected:

1. `NAME` could be anything so long as it does not contain any white space
   characters.
2. `NAME` must not be a reserved `gsnip` command (e.g., `@LST` would list out
   all the snippets found in the file).
3. `COMMENT` should always be enclosed in double quotes.
4. Finally, `BODY` can be pretty much anything.


The `snippet` package has a snippet container based on `PostgreSQL` database.
You might want to use it instead of the file-based implementation, but you'd
need to do a little rewrite of the server `gsnipd` command to accommodate for
this change. I personally do not like this because it makes the program clunky.
Database would then contain a simple table with the name, comment and body text
fields. There isn't really more to it---it does the same work as a flat file.


## Installation

```sh
go get github.com/mdm-code/gsnip
```

If you want to build it from source, consult `Makefile`; there is an `install`
directive that you can call with `make install`, and that's pretty much it when
it comes to the installation of this program. Make sure that the directory
where it's installed is on your $PATH. You can move, copy or symlink the
binaries anywhere on the $PATH.


## Development

Have a look at `Makefile` to see how to run the standard code linting and testing
works.

Don't forget to install `golint` before you run the test command from the
`Makefile`:

```sh
go get -u golang.org/x/lint/golint
```


## License

Copyright (c) 2022 Michał Adamczyk.

This project is licensed under the [GNU GPL 3 license](https://opensource.org/licenses/GPL-3.0).
See [LICENSE](LICENSE) for more details.
