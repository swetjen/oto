# otohttp Example

This example demonstrates how to use Oto to generate server and client code.

## Structure

* `def/` - Service definitions as Go interfaces
* `server.go.plush` - Template for generating Go server code
* `client.js.plush` - Template for generating JavaScript client
* `client.swift.plush` - Template for generating Swift client
* `generate.sh` - Script to regenerate all code (run with `go generate`)
* `main.go` - Server implementation using the generated code
* `index.html` - Web UI that uses the generated JavaScript client
* `swift/SwiftCLIExample/` - Xcode project using the generated Swift client

## Generated files

* `server.gen.go` - Generated Go server handlers
* `client.gen.js` - Generated JavaScript client
* `swift/SwiftCLIExample/SwiftCLIExample/client.gen.swift` - Generated Swift client

## Run the example

```bash
go run *.go
```

Open http://localhost:8080

## Regenerate code

```bash
go generate
```

Or directly:

```bash
./generate.sh
```

## Exercise

1. Add a new service interface to the `def` package
2. Run `go generate` to regenerate the server and client code
3. Implement the new service in `main.go`
4. Call your new service from `index.html`
