# Go 1.22 tracing example

It's a simple example that demonstrates how to use the new tracing enhancements in Go 1.22.
using FlightRecorder from `golang.org/x/exp/trace` package we record the execution trace of the application within
a specified time/size window and then analyze it to identify performance bottlenecks and optimize the code.

### Installation

This application requires Go 1.22 and `golang.org/x/exp/trace` pkg.

Please make sure that you executed the following command before running the application:

```bash
go mod tidy
```

### Documentation

For more details about the tracing enhancements in Go 1.22, please refer to:

- https://go.dev/blog/execution-traces-2024
- https://go.dev/doc/go1.22#runtime/trace
- https://pkg.go.dev/golang.org/x/exp/trace#FlightRecorder

### Usage

In the example below, we create a new `FlightRecorder` instance and start tracing the execution of the program.
Due to the tracing adjustments in Go 1.22, we shall have approx 1-2% overhead instead of 10-20% as it was before.

A specific of the `FlightRecorder` is that it allows us to keep the data of the specific data/size window. 
For example, we can keep the data for the last 24 hours, removing the older data. This allows us to keep 
a relevant amount of data for the analysis in case of an issue and having a small overhead at the same time.

Starting a web server with a route handler for the `/trace` endpoint, we can write the trace data to a file.

```go
// Create a new FlightRecorder instance.
fr := trace.NewFlightRecorder()

// We want to trace the execution of the whole program, so we start the FlightRecorder at this early point.
// However, we can start it later, for example, when we want to trace a specific part of the program execution.
if err := fr.Start(); err != nil {
    // handle error
}

defer func () {
    if err := fr.Stop(); err != nil {
	// handle error
    }
}()

// Create a new Fiber instance to run a web server.
app := fiber.New()

// Define a route handler for the /trace endpoint. When the endpoint is hit, the trace data is written to a file.
app.Get("/trace", func (ctx fiber.Ctx) error {
    f, err := os.OpenFile("file.out", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
    if err != nil {
        // handle error
    }

    if _, err := fr.WriteTo(f); err != nil {
        // handle error
    }
    return ctx.SendStatus(fiber.StatusOK)
})

// We run the game in a separate goroutine to make a scenario where a few goroutines are competing
// for the same resource - channel in our case. This will help us to see how the trace tool works.
// The implementation of the runGame function can be found in the game.go file.
go runGame()

if err := app.Listen(":8080"); err != nil {
    // handle error
    app.Shutdown()
}
```

Once the data is written, it can be analyzed using the trace tool:

```bash
go tool trace file.out
```

