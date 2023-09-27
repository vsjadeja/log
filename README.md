# Log package

This package is a wrapper on zap. Please see documentation in https://github.com/uber-go/zap

## 1. Add to project

Launch in the your project directory

```sh
go get github.com/vsjadeja/log
```

## 2. Setup package

```go
package main

import (
	lg "github.com/vsjadeja/log"
	"context"
)

func main() {
	logger := lg.L() //as singleton
	//or
	logger = lg.NewLogger() //normal method
}
```

## 3. Usage logger
```go

package main

import (
	lg "github.com/vsjadeja/log"
	"go.uber.org/zap/zapcore"
	
	"context"
)

func main() {
	ctx := context.Background()
	logger := lg.NewLogger()
	
	//prints log and traceId if it exists in context
	logger.Info(ctx, "info log", "key1", 1, "key2", 2)
	//example: {"level":"info","time":"2022-06-09T16:24:23.159+0300","caller":"example/main2.go:11","message":"info log","key1":1,"key2":2,"traceId":"unknown"}

	//DO NOT USE THIS DEPRECATED METHOD
	logger.Infof("infof log %s:%d", "key", 1)
	//example: {"level":"info","time":"2022-06-09T16:24:23.159+0300","caller":"example/main2.go:13","message":"infof log key:1"}

	//prints log and append key value pairs
	logger.Infow("infow log", "key1", 1, "key2", 2)
	//example: {"level":"info","time":"2022-06-09T16:38:44.348+0300","caller":"example/main2.go:15","message":"infow log","key1":1,"key2":2}

	//set name for logger
	logger.Named("test")
	logger.Info(ctx, "info log", "key1", 1, "key2", 2)
	//example: {"level":"info","time":"2022-06-09T16:45:01.886+0300","logger":"test","caller":"example/main2.go:11","message":"info log","key1":1,"key2":2,"traceId":"unknown"}

	//append some field to logger and creates sublogger with new parameter
	subLogger := logger.With(zmlog.Field{"some_field", zapcore.StringType, 0, "test", nil})
	subLogger.Info(ctx, "info log", "key1", 1, "key2", 2)
	//example: {"level":"info","time":"2022-06-09T16:56:18.712+0300","logger":"test","caller":"example/main2.go:19","message":"info log","some_field":"test","key1":1,"key2":2,"traceId":"unknown"}
}
Please, don`t use printf method in logger. 
```

### 4. Syntax sugar

#### 4.1 Format error in log
```go
    err := errors.New("test error")
    logger.Error(ctx, "error log", "key1", 1, zmlog.Error(err)) 
	//example: {"level":"error","time":"2022-06-09T18:00:56.196+0300","caller":"example/main.go:22","message":"error log","key1":1,"error":"test error","traceId":"unknown","stacktrace":"log/example/main.go:22\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:255"}
```
