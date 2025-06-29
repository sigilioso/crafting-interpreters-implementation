
# Crafting interpreters implementation

## Running tests from the book repository

Book repository: <https://github.com/munificent/craftinginterpreters>

Build docker image:

```bash
$ git submodule init --update
$ cd book
$ docker build -t craftinginterpreters . -f ../Dockerfile
$ cd -
```

Run tests (chapter04 example, check book/Makefile for naming):

```bash
❯ docker run -t --volume $(pwd):/code craftinginterpreters dart tool/bin/test.dart chap04_scanning --interpreter /code/glox/glox
All 6 tests passed (59 expectations).
```
