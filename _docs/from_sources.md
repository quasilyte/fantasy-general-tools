## Building from Source

To build `fgtool` from sources you would need a Go compiler.

The simplest way to install `fgtool` would then be:

```bash
go install -v github.com/quasilyte/fantasy-general-tools@latest
```

Alternatively:

```bash
git clone https://github.com/quasilyte/fantasy-general-tools.git
cd fantasy-general-tools
mkdir bin
go build -o bin/fgtool ./cmd/fgtool
```
