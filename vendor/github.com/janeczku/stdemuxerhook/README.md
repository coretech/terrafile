# A stdout/stderr demuxer hook for Logrus

[Logrus](https://github.com/Sirupsen/logrus) loggers always output all log levels to a single common output, e.g. `stderr`. 
This logrus hook makes it so that logs with a severity below `Error` are written to `stdout` while all the important stuff goes to `stderr`.

You can also use the hook to demux logs to custom IO writers based on severity. Just override the default outputs using the hook's `SetOutput(infoLevel, errorLevel io.Writer)` method.

## Example

Given you have an application that uses a the logrus standard logger similar to this:

```go
import (
    log "github.com/Sirupsen/logrus"
)

func main() {

    log.SetLevel(log.InfoLevel)

    log.Info("A group of penguins emerges from the ocean")

    (...)

```

The only change required is adding in the hook. Make sure to configure the parent logger before initializing the hook as the latter will inherit it's configuration.

```go
import (
    log "github.com/Sirupsen/logrus"
    "github.com/janeczku/stdemuxerhook"
)

func main() {

    log.SetLevel(log.InfoLevel)
    log.AddHook(stdemuxerhook.New(log.StandardLogger()))

    log.Info("A group of penguins emerges from the ocean") // -> logged to stdout
    log.Panic("A group of polar bears emerges from the ocean") // -> logged to stderr

    (...)

}
```

## Benchmarks

```BASH
BenchmarkLoggerWithoutHook        300000          4348 ns/op
BenchmarkLoggerWithHook           200000          5436 ns/op
```