<img src="https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/kphgopherit.png" width="350">

# Welcome!
Do you want to build a modern Go tool with a **C**ommand **L**ine **I**nterface (**CLI**) similar to Docker, Kubernetes, aws-cli, etc.?  ? Start here. :- ) 

This project is a starting set of Go files and directories usually needed to build a CLI, that calls a remote service, and needs to interpret/convert the results.

Simply 'Copy & Paste', 'Find & Replace', tweak a few default values and you can be up and running.

## Overview
This package has four major parts to it:
  1) A Server implementation of ACME Services such as `/gophers`, `/gopher/1/things` (built using [`go-chi`](https://github.com/go-chi/chi) ) 
  2) ACME Services Library (`pkg.acme.Service`) for making the HTTP ACME API service calls against a server   
  3) Client which uses the service library to call the server and converts the returned ACME Gophers (`pkg.acme.Gopher`) to Client.Gophers (`pkg.client.Gopher`) 
  4) A CLI invocation and configuration framework with [`cobra`](https://github.com/spf13/cobra) and [`viper`](https://github.com/spf13/viper)

## This code includes:
- [x] Fundamental Go features like tests, generate, templates, go routines, contexts, channels, OS signals, HTTP routing, build/run tags, ldflags, 
- [x] Uses [`cobra`](https://github.com/spf13/cobra) and [`viper`](https://github.com/spf13/viper) (without func inits!!!)
  - Cleanly separated CLI/configuration invocation from client library calls - by calling `viper.Unmarshal` to transfer our `pkg.Config`
  - **NOTE**: A lot of sample Cobra/Viper code rely on `func init()` making it more difficult to reuse. 
- [x] Using [`vfsgen`](https://github.com/shurcooL/vfsgen) in to embed templates into binary
    - The `config\template\*` contain all text output and is compiled into a `templates_generate.go` via [`vfsgen`](https://github.com/shurcooL/vfsgen) for the binary build
- [X] Logging from the [`logrus`](https://github.com/sirupsen/logrus) library
- [x] Cached response folder `.cache` with entries from the Server, Client and Services
  - The server uses entries in `.cache` instead of making DB calls (when present.)
- [x] [Retry](https://github.com/matryer/try) using @matryer's idiomatic `try.Do(..)`
- [X] Instrumentation with [`prometheus`](https://prometheus.io/) in the server and client library
  - [Tutorials](https://pierrevincent.github.io/2017/12/prometheus-blog-series-part-4-instrumenting-code-in-go-and-java/)
- [X] HTTP serving/routing with middleware from [`go-chi`](https://github.com/go-chi/chi)
    - Using `NewStructuredLogger` middleware to decorate each route with log output
    - `ResponseHandler` to pretty print JSON with [`jq`](https://stedolan.github.io/jq/)
    - Custom middlewares (`GopherCtx`,`ThingCtx`) to handle creating Context from HTTP requests
- [x] An example Dockerfile and build recipe `(docs/recipe/)` for a docker workflow
  - Use `docker build --tag gophercli:v1 .` to create a full golang image
  - Use `docker run -it --rm gophercli:v1` to work from with the container

I've [curated a YouTube playlist](https://www.youtube.com/playlist?list=PLa1qVAzg1FHthbIaRRbLyA4sNE4PmLmn6) of videos which help explain how I ended up with this structure and 'why things are the way they are.' I've leveraged 'best practices' I've seen and that have been explicted called out by others. Of course **THERE ARE SOME WRINKLES** and few **PURELY DEMONSTRATION** portions of code. I hope to be able to keep improving on this.

[![Go Report Card](https://goreportcard.com/badge/github.com/whereiskurt/gopherit)](https://goreportcard.com/report/github.com/whereiskurt/gopherit)

## Go version 1.11 or greater required!
A lot has happened in the Go ecosystem in the last year two-years:
- Using go modules proper (ie. `go.mod`, `go.sum`, `vendor` folder) 
  - **Works outside of `$GOPATH`**
- `go test -v ./...` showing server start / stop, add/update/delete gopher and things (95% file coverage, 78% statements coverage) 
- `go build -tags release cmd/gophercli.go` builds a self-contained executable with text templates and default configuration files embedded in the binary ([`vfsgen`](https://github.com/shurcooL/vfsgen)) 
    - Hermetic build/run/test with `vendor` folder checked-in:
      - **NOTE: still need `GOFLAGS="-mod=vendor"` until Go 1.12**
    - Self-contained build with `Dockerfile` 
      - Uses `--ldflags "-X 00-newapp-template/internal/app/cmd.ReleaseVersion=$VERSION`
- `go generate ./...` to embed text templates into a Go source file of `[]byte` 

# The Story of 00-newapp-template
There is a vendor named ACME who provides API to access to `Gophers` and `Things`. Because I use the ACME API **all the time** to track `Gophers` and their `Things` and I have decided to create a CLI tool to perform the HTTP API calls needed and output a text table or JSON structure. Ideally using a simple command like:
```
  ./gophercli list 
```
The client would make all the necessary calls to the HTTP ACME services, convert the JSON responses to Go structures and output something like:
```
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++++                    _               _   _     _              +++++++++++++++
++++   __ _  ___  _ __ | |__   ___ _ __| |_| |__ (_)_ __   __ _  +++++++++++++++
++++  / _` |/ _ \| '_ \| '_ \ / _ \ '__| __| '_ \| | '_ \ / _` | +++++++++++++++
++++ | (_| | (_) | |_) | | | |  __/ |  | |_| | | | | | | | (_| | +++++++++++++++
++++  \__, |\___/| .__/|_| |_|\___|_|   \__|_| |_|_|_| |_|\__, | +++++++++++++++
++++  |___/      |_|                                      |___/  +++++++++++++++
++++                                                             +++++++++++++++
+------------------------------------------------------------------------------+
| ID  |Name      |Description                   | ID  | Thing.Name |Description|
+------------------------------------------------------------------------------+
| 1   |Gopher1   |The first Gopher (#1st)       | 1   | Head       |Hat        |
| 1   |Gopher1   |The first Gopher (#1st)       | 5   | Feet       |Shoes      |
| 1   |Gopher1   |The first Gopher (#1st)       | 9   | Waist      |Belt       |
| 2   |Gopher2   |The second Gopher (#2nd)      | 10  | Waist      |Belt       |
| 2   |Gopher2   |The second Gopher (#2nd)      | 2   | Head       |Hat        |
| 2   |Gopher2   |The second Gopher (#2nd)      | 6   | Feet       |Shoes      |
+------------------------------------------------------------------------------+
```
## The Challenge: ACME Data Types
At ACME each `Thing` has a `Gopher` but each `Gopher` does not have `Things`. If we have a `Gopher` we don't know their collection of `Things`. And, if we have a `Thing` we only have the associated `GopherID` (not the `Gopher` `Name` or `Gopher` `Description`.) 

Here is the data structure/JSON from the ACME data type:
```
  package acme
  type Gopher struct {
    ID          json.Number `json:"gopher_id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
  }

  type Thing struct {
    ID          json.Number `json:"thing_id"`
    GopherID    json.Number `json:"gopher_id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
  }
```
(These JSON structures are defined in `pkg\acme\json.go` and are what the ACME HTTP API return.)

While ACME's data structure and JSON schema is useful for their purposes - it's incomplete for our application's needs. 

A much more useful structure for our application of listing `Gophers` and their `Things` looks like this:
```
  package adapter
  type Gopher struct {
    ID          string           `json:"id"`
    Name        string           `json:"name"`
    Description string           `json:"description"`
    Things      map[string]Thing `json:"things"`
  }

  type Thing struct {
    Gopher      Gopher `json:"gopher"`
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
  }
```
Above, each `Gopher` has a collection of their `Things` and each `Thing` has a reference to a complete `Gopher` (not just the `GopherID`.)

## The Solution - Build an adapter!
We will write a client library to make ACME HTTP service calls, unmarshal the JSON, filter and convert it into our Go structures. All of this adapter code is in the `pkg\adapter\` folder. 

Our adapter uses services provide by ACME as defined in the `pkg\acme\service.go` file.

Because ACME API charges us per API call, we decided to mock implement our own ACME HTTP API server. Our server will listen for connections on local host and respond just like an ACME server would. All the code to provide a full HTTP API server is in `pkg\acme\server\` and to invoke our ACME server:
```
  ./gophercli server start
```
![Server Startup](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/serverstart.gif)
[(Download MP4 of Server Startup)](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/serverstart.mp4)

# Complete Process - Looping GIF Video
This video shows download Go version 1.11.4 and checking out from scratch:

![Download Go v1.11.4](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/getgo.gif)
[(Download MP4 of Go v1.11.4)](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/getgo.gif)

![git clone gopherit repo](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/getgopherit.gif)
[(Download MP4 of git clone gopherit repo)](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/getgopherit.gif)

![docker build/run](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/dockerworkflow.gif)
[(Download MP4 of docker build/run)](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/dockerworkflow.mp4)

## How-to Run:
```
 $ go run cmd/gophercli.go 

              ,_---~~~~~----._
        _,,_,*^____      _____''*g*\"*,
        / __/ /'     ^.  /      \ ^@q   f
        [  @f | @))    |  | @))   l  0 _/
        \'/   \~____ / __ \_____/    \
        |           _l__l_           I
        }          [______]           I
        ]            | | |            |
        ]             ~ ~             |
        |                            |
      [@https://gist.github.com/belbomemo]

GopherIT! uses ACME(TM) API to review and modify Gophers and their Things.

Find more information at:
    https://github.com/whereiskurt/gopherit/00-newapp-template/

Usage:
    gophercli [COMMAND] [SUBCOMMAND] [OPTIONS]

Global Options:
    Verbosity:
      --silent,  -s     Set logging/output level [level1]
      --quiet,   -q     Set logging/output level [level2]
      --info,    --v    Set logging/output level [level3-default]
      --debug,   --vv   Set logging/output level [level4]
      --trace,   --vvv  Set logging/output level [level5]
      --level=3         Sets the output verbosity level numerically [default]

Examples:
    
      ## The default command is 'client' and is optional before subcommands.
      $ gophercli client list

      $ gophercli list
      $ gophercli list -g1 -t2,3
 
      $ gophercli server start
      $ gophercli server start --port=102102 --docroot=./config/docroot
 
      $ gophercli server stop
 
      $ gophercli version

```

## Package Relationships
TODO
