# Gsnip

This my personal snippet manager. It lets you find and print out snippets
stored in a text file and written with straightforward syntax rules. My goal
was to keep the program as simple as possible, so there aren't any features
that would allow you to insert or delete snippet definitions other than by
manually editing the source file. It simply scans the source file with snippets
and offers an interface to retrieve them.

There are two parts of this workflow: one, `gsnipd`, a UDP-based server
handling connections, and `gsnip`, which is the client that relies on `FD0` or
`SDTIN` to message the server.

Both `gsnip` and `gsnipd` print out moderately useful information with `--help`
attribute.

You want to first spin up the server with either of these commands:

```sh
gsnipd
gsnipd &>/dev/null &
```

The first one will write `STDERR` to the terminal so that you can see server
messages. The other sends all messages to `/dev/null` and gets detached from
the current session.

Then you want to message the server like this:

```sh
echo [snip-name] | gsnip
echo list | gsnip
```

You can query the server for any snippet stored in the source file.
Alternatively, you can ask the server to list out all available snippets with
the `@LIST` command.

You can reload the source snippet file at the server runtime by piping the
keyword `@RELOAD` to `gsnip` client, which is the equivalent of sending `SIGHUP`
to the process using `kill -1 [pid]`. The latter is annoying because you have
to find the process id with `ps` before sending the signal. Another way would
be to write PID to a known file that the client could access.

Here is a list of restricted names (`gsnip` commands):
1. @LIST
2. @RELOAD

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
2. `NAME` must not be a reserved `gsnip` command (e.g., `@LIST` would list out
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

Consult `Makefile`; there is the `install` directive that you would call with
`make install`, and that's pretty much it when it comes to installation of a Go
program.


## Testing

I am using the `testing` package from the Go standard library, so you can call
`test ./... -v`, or you can resort to `Makefile` and use `make test` command. I
am really not worried about the coverage at this stage; the program is way to
simple and the test as they currently are, they cover all of the functional
bottlenecks of the program. However, if you want to peek at the coverage, type
`make cover` and the `make clean` once you're done.
