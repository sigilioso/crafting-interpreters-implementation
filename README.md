
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

The rust version needs the _musl_ version of the binary as the `dart:2` image has an old version of glic.
Ensure that the corresponding alias is installed:

```bash
$ rustup target add x86_64-unknown-linux-musl
```

And compile using the alias configured in [rlox/.cargo/config.toml](rlox/.cargo/config.toml)

```bash
$ cargo rel
```
