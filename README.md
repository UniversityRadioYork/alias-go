alias-go
========

A re-implementation of [aliasgen](https://github.com/UniversityRadioYork/aliasgen) in Go, using [myradio-go](https://github.com/UniversityRadioYork/myradio-go).

URY mailiases generator.
Generates the exim aliases file which is responsible for all @ury.org.uk email addresses that don't have a local mailbox being sent to the right place. 
Using [myradio-go](https://github.com/UniversityRadioYork/myradio-go) it gets:
- Mailing list aliases, e.g. computing@ury.org.uk
- Official aliases, e.g. head.of.computing@ury.org.uk
- User aliases, e.g. sam.w@ury.org.uk
- 'Text' (misc) aliases, e.g. sexual.advances@ury.org.uk (??)

## Installation
```bash
$ git clone https://github.com/UniversityRadioYork/alias-go
$ cd alias-go
$ go install
$ alias-go -e config.toml # will generate an example config
```

## Usage
```
USAGE:
   alias-go [global options] command [command options] [arguments...]

GLOBAL OPTIONS:
   --config-file FILE, --config FILE, -c FILE      Load configuration from FILE (required)
   --out-filename FILE, --out FILE, -o FILE        Write aliases to FILE (default: "aliases")
   --example-config FILE, --example FILE, -e FILE  Write an example config to FILE
   --verbose, -v                                   Output additional information to stdout
   --help, -h                                      show help
```

## Testing
```bash
$ go test ./...
```
