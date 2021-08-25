# Gsnip

This is my personal snippet manager. It offers some basic functionality. It
lets you find and recover snippets found in the snippet source file. Insertion
and deletion---at this stage---is done by editing the source file in an of
choice.  You can list out all of your snippets with the `list` command. You
would normally add your snippets to a text file where each snippet has to
adhere to a predefined syntax that goes like this:

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
2. `NAME` must not be a reserved `gsnip` command (e.g., `list` would list out
   all the snippets found in the file).
3. `COMMENT` should always be enclosed in double quotes.
4. Finally, `BODY` can be pretty much anything.

Here is a list of restricted names (`gsnip` commands):
1. list


## Usage

Type in `gsnip -help` to show the help message. You can type in `gsnip list` to
get a list of all snippets. To get the body of your snippet, type in `gsnip
NAME` where `NAME` is of course the name identifier of the snippet. And that's
it at this point.


## Installation

Consult `Makefile`: there is the `install` directive that you would call with
`make install`, and that's pretty much it when it comes to installation of a Go
program.


## Testing

I am using the `testing` package from the Go standard library, so you can call
`test ./... -v` or you can resort to the `Makefile` and use `make test`
command. I am really not worried about the coverage at this stage; the program
is way to simple and the test as they currenly are, they cover all of the
functional bottlenecks of the program.


## Future plans

1. Common interface for flat files and databases. Right now it just parses a
   text file. Databases would require some more abstarction in the code.
2. Make it a server plus a client where it spin up a server and the the client
   and server talk to each other over a Unix Domain Socket. I'd like to keep
   the commands as they are, therfore, I would rather refrain from TCP.
3. I want the client to consequently rely on STDIN.
4. I have to ponder over ADD and DELETE commands and how to implement them, if
   at all.
