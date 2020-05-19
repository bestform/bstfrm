# bstfrm

This is a hobby project to explore the world of lexers, parsers and interpreters.
It implements a useless programming language in go and provides a vm and a repl to run code.

It is in a very early state, but a few things already work, like printing strings to the repl.

Try it out:

`go run cmd/bstfrm.go`

```
Welcome to bstfrm.
> set #foo="World!";

ok
> print "Hello " #foo;
Hello World!
ok
> calc (1+2)*3;
9

ok
> 

```