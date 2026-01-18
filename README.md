![Welcome to oto project](oto-logo.png)

Go driven rpc code generation tool for right now.

- 100% Go
- Describe services with Go interfaces (with structs, methods, comments, etc.)
- Generate server and client code
- Production ready templates (or copy and modify)

## Who's using Oto?

- [Grafana Labs, IRM tool](https://grafana.com/docs/grafana-cloud/alerting-and-irm/incident/api/)
- [Pace.dev](https://pace.dev/docs/)
- [Firesearch.dev](https://firesearch.dev/docs/api)

## Templates

These templates are already being used in production.

- There are some [official Oto templates](https://github.com/pacedotdev/oto/tree/master/otohttp/templates)
- The [Pace CLI tool](https://github.com/pacedotdev/pace/blob/master/oto/cli.go.plush) is generated from an open-source CLI template

## Learn

![](oto-video-preview.jpg)

- VIDEO: [Mat Ryer gives an overview of Oto at the Belfast Gophers meetup](https://www.youtube.com/watch?feature=youtu.be&v=DUg4ZITwMys)

- BLOG: [How code generation wrote our API and CLI](https://pace.dev/blog/2020/07/27/how-code-generation-wrote-our-api-and-cli.html)

## Tutorial

Install the project:

```
go install github.com/pacedotdev/oto@latest
```

Create a project folder, and write your service definition as a Go interface:

```go
// definitions/definitons.go
package definitions

// GreeterService makes nice greetings.
type GreeterService interface {
    // Greet makes a greeting.
    Greet(GreetRequest) GreetResponse
}

// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
    // Name is the person to greet.
    // example: "Mat Ryer"
    Name string
}

// GreetResponse is the response object containing a
// person's greeting.
type GreetResponse struct {
    // Greeting is the greeting that was generated.
    // example: "Hello Mat Ryer"
    Greeting string
}
```

Download templates from otohttp

```bash
mkdir templates \
    && wget https://raw.githubusercontent.com/pacedotdev/oto/master/otohttp/templates/server.go.plush -q -O ./templates/server.go.plush \
    && wget https://raw.githubusercontent.com/pacedotdev/oto/master/otohttp/templates/client.js.plush -q -O ./templates/client.js.plush
```

Use the `oto` tool to generate a client and server:

```bash
mkdir generated
oto -template ./templates/server.go.plush \
    -out ./generated/oto.gen.go \
    -ignore Ignorer \
    -pkg generated \
    ./definitions
gofmt -w ./generated/oto.gen.go ./generated/oto.gen.go
oto -template ./templates/client.js.plush \
    -out ./generated/oto.gen.js \
    -ignore Ignorer \
    ./definitions
```

- Run `oto -help` for more information about these flags

### Additional flags

- `-ignore`: Comma-separated list of interfaces to exclude from generation
- `-fields`: Comma-separated list of fields to exclude from all objects
- `-suppressErrorField`: When set, suppresses the error field in response objects
- `-types`: Path to a YAML config file for type overrides (see "Type Overrides" below)

Implement the service in Go:

```go
// greeter_service.go
package main

// GreeterService makes nice greetings.
type GreeterService struct{}

// Greet makes a greeting.
func (GreeterService) Greet(ctx context.Context, r GreetRequest) (*GreetResponse, error) {
    resp := &GreetResponse{
        Greeting: "Hello " + r.Name,
    }
    return resp, nil
}
```

Use the generated Go code to write a `main.go` that exposes the server:

```go
// main.go
package main

func main() {
    g := GreeterService{}
    server := otohttp.NewServer()
    generated.RegisterGreeterService(server, g)
    http.Handle("/oto/", server)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

- The `otohttp.Server` has a `Basepath` field (default: `/oto/`) that determines the URL prefix for service routes
- When mounting the server, use the same path as `Basepath` (e.g., `http.Handle("/oto/", server)`)
- To customize the basepath: `server.Basepath = "/api/"` then `http.Handle("/api/", server)`

Use the generated client to access the service in JavaScript:

```javascript
import { GreeterService } from "oto.gen.js";

const greeterService = new GreeterService();

greeterService
  .greet({
    name: "Mat",
  })
  .then((response) => alert(response.greeting))
  .catch((e) => alert(e));
```

## Use `json` tags to control the front-end facing name

You can control the name of the field in JSON and in front-end code using `json` tags:

```go
// Thing does something.
type Thing struct {
    SomeField string `json:"some_field"`
}
```

- The `SomeField` field will appear as `some_field` in json and front-end code
- The name must be a valid JavaScript field name

## Type Overrides

You can configure how Go types map to target language types using a YAML configuration file with the `-types` flag:

```bash
oto -template ./templates/client.ts.plush \
    -out ./generated/client.gen.ts \
    -types ./types.yaml \
    ./definitions
```

Create a `types.yaml` file to define your type mappings:

```yaml
version: "1"
overrides:
  - go_type: "time.Time"
    js_type: "Date"
    ts_type: "Date"
    swift_type: "Date"
    dart_type: "DateTime"
  - go_type: "github.com/google/uuid.UUID"
    js_type: "string"
    ts_type: "string"
    swift_type: "UUID"
    dart_type: "String"
  - go_type: "int64"
    swift_type: "Int64"  # Partial override - only specify languages you need
```

Each override can specify any combination of target language types:
- `go_type`: The fully qualified Go type to match (required)
- `js_type`: JavaScript type override
- `ts_type`: TypeScript type override
- `swift_type`: Swift type override
- `dart_type`: Dart type override

## Specifying additional template data

You can provide strings to your templates via the `-params` flag:

```bash
oto \
    -template ./templates/server.go.plush \
    -out ./oto.gen.go \
    -params "key1:value1,key2:value2" \
    ./path/to/definition
```

Within your templates, you may access these strings with `<%= params["key1"] %>`.

## Comment metadata

It's possible to include additional metadata for services, methods, objects, and fields
in the comments.

```go
// Thing does something.
// field: "value"
type Thing struct {
    //...
}
```

The `Metadata["field"]` value will be the string `value`.

- The value must be valid JSON (for strings, use quotes)

Examples are officially supported, but all data is available via the `Metadata` map fields.

### Examples

To provide an example value for a field, you may use the `example:` prefix line
in a comment.

```go
// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
    // Name is the person to greet.
    // example: "Mat Ryer"
    Name string
}
```

- The example must be valid JSON

The example is extracted and made available via the `Field.Example` field.

### Open API

To work on the Open API spec, you might find this command helpful:

```
oto -template ./otohttp/templates/openapi.yaml.plush -out openapi.yaml -v -ignore Ignorer ./parser/testdata/services/pleasantries
```

## Contributions

Special thank you to:

- @mgutz - for struct tag support
- @sethcenterbar - for comment metadata support

![A PACE. project](pace-footer.png)
