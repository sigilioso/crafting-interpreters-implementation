
# Crafting interpreters implementation

## Running tests from the book repository

Book repository: <https://github.com/munificent/craftinginterpreters>

Build docker image:

```bash
$ cd /home/christian/code/repos/github.com/munificient/craftinginterpreters
$ docker build -t craftinginterpreters . -f /home/christian/code/repos/github.com/sigilioso/crafting-interpreters-implementation/Dockerfile
$ cd -
```

Run tests:

```bash
❯ docker run -t --volume $(pwd):/code craftinginterpreters dart tool/bin/test.dart chap04_scanning --interpreter /code/glox/glox
All 6 tests passed (59 expectations).
```
