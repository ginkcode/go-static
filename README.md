# go-static

Serving http server with static files, leveraged by `cli` and `gofiber`

```
go get -u github.com/ginkcode/go-static
```

- Default port: 8000
- Api default path: /api
- Api default proxy: localhost:3000

```shell
static-server - Start http server with static files!

USAGE:
   go-static [global options] command [command options] Path of static files (default: ".")

VERSION:
   v1.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value   listening port of http server (default: "8000")
   --index, -i              allow index all directory (default: false)
   --spa, -s                support SPA mode, this option will be ignored when allowing index (default: false)
   --api, -a                enable proxy for /api, combine with --proxy to config (default: false)
   --proxy value, -x value  list of proxy for /api, separated by comma (default: "http://localhost:3000")
   --help, -h               show help (default: false)
   --version, -v            print the version (default: false)

```
