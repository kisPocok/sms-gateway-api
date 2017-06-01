# SMS Gateway API

Example application to send SMS via messagebird, written in Go.

## Install & Run

On Mac OS X you can install & run with the following command:
```
make start_macos
```

Similar command for Linux:
```
make start_linux
```

Docker containeraized version is also available. Run `make docker` to build your container locally then run it with `make docker_run`.

Cleanup also available with `make docker_clean` command.

For more information check `Makefile` or run `make help`.

## Manual

Send a simple request to the application:

```
curl -d "recipient=3670..." \
    -d "originator=MessageBird" \
    -d "message=Hello world!" \
    -X POST http://127.0.0.1:8080/api/v1/message
```

Available params for manual use:
- `-port` server port, default is `8080`

## Docs

Swagger API documentation located in `api-doc.yml` file.

