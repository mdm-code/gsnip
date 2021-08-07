# Gsnip

Manage snippets through a bespoke interface.


## Usage

Spin up snippet server and interact with it using the client.
Maybe this would work:

```sh
nohup ./gserve &
./gsnip find class
```


## Plan

Here is an outline of how I conceive this piece of software.

1. I want the program to take input through standard input, FD1.
	- This implies server-client architecture.
	- The server communication would rely on Unix Domain Sockets.

2. The server has to be able to parse a file with snippets.
	- The server is able to monitor changes done to the file.
	- Snippets are loaded into a dictionary of snippets.

3. A snippet has to meet some formal qualities to be parsed.
	- It would most likely be enclosed in this sort of flat syntax.
		```
		startsnippet [hello]
		func helloWorld () {
			fmt.Println("Hello, world!")
		}
		endsnippet
		```

4. Core server functionality:
	- Register source file
	- Parse files
	- Expose contents = List all loaded snippets
	- Server takes multiple directives: list, find, add?, delete?

