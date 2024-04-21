# execs

A CLI that allows ECS Exec to be interactive.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Installation

### Homebrew

```sh
brew install t4kamura/tap/execs
```

### Go Install

```sh
go install github.com/t4kamura/execs
```

## Usage

To get started

```sh
execs -p <AWS Profile>
```

or first select an interactive AWS profile.

```sh
execs
```

Must have a valid AWS Profile set up.

`--help` for other options.

## Todo

- table print (selected resources)
- add spinner when loading
- support SSO format aws profile

## License

MIT
