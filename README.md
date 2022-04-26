<h1 align="center">Gin Example</h1>

<div align="center">
this a message.
</div>

## Quick Start

#### 1. Setter goproxy (if in china)

```bash
$ go generate
```


#### 2. Install pkgs (you can need `make`)

```bash
$ go mod download
$ go install github.com/swaggo/swag/cmd/swag@latest
$ go install github.com/mitchellh/gox@latest
```

#### 3. Start App

```bash
$ make start
```

---

## Other

### Project Config

```bash
# self copy
$ cp ./config/config.example.yml ./config/config.yml
# or .. auto generate
$ make genconfig
```