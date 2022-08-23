# rotabot

SlackApp that makes team rotations easy

## Development

[Install Go](https://go.dev/doc/install), you can do this with Brew or your favourite way of installing dependencies.

```shell
  brew install go
```

Run Tests:

```shell
  make test
```

Spin up database and other dependencies:

```shell
  docker-compose up -d
```

Once you don't need them anymore you can run the following command to stop the containers:

```shell
  docker-compose down -v
```

Build and run:

```shell
  make clean build
  make run
```

Check lint and dependencies:

```shell
  make check
  make lint
``` 