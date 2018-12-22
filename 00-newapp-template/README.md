<img src="https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/kphgopherit.png" width="350">

# Welcome!
Do you want to build a modern Go tool with a **C**ommand **L**ine **I**nterface (**CLI**) similar too Docker, Kubernetes, aws-cli, etc.?  ? Start here. :- ) 

This project is a starting set of Go files and directories. Simply 'Copy & Paste', 'Find & Replace', tweak a few default values and you can be up and running. 

This package has four major parts to it:
  1) CLI invocation and configuration ([`cobra`](https://github.com/spf13/cobra) / [`viper`](https://github.com/spf13/viper) ) 
  2) A service library for making HTTP ACME API service calls - `pkg.acme.Service` that return JSON arrays of  and `pkg.acme.Thing` 
  3) A client library (`internal.pkg.adapter.Adapter`) using the service library to convert ACME API results (`pkg.acme.Gopher`) to our Go structures (`internal.pkg.adapter.Gopher`) 
  4) An HTTP server (using [`go-chi`](https://github.com/go-chi/chi)) to provide an ACME HTTP API server- `interal.pkg.server.Server`

I've [curated a YouTube playlist](https://www.youtube.com/playlist?list=PLa1qVAzg1FHthbIaRRbLyA4sNE4PmLmn6) of videos which help explain how I ended up with this structure and 'why things are the way they are.' I've leveraged 'best practices' I've seen and that have been explicted called out by others. Of course **THERE ARE SOME WRINKLES** and few **PURELY DEMONSTRATION** portions of code. I hope to be able to keep improving on this.

[![Go Report Card](https://goreportcard.com/badge/github.com/whereiskurt/gopherit)](https://goreportcard.com/report/github.com/whereiskurt/gopherit)

# Go version 1.11 or greater required!

A lot has happened in the Go ecosystem in the last year two-years. This project is:
- Using go modules proper (ie. `go.mod`, `go.sum`, `vendor` folder) 
  - Works outside of `$GOPATH`
  - `go test -v ./...` to server start server and test client
  - 'Hermetic build/run/test' with `vendor` folder checked-in 
  - **NOTE:** still need `GOFLAGS="-mod=vendor"` until Go 1.12

The of code includes:
- [x] Fundamental Go features like tests, templates, go routines, contexts, channels, HTTP routing
  - The `config\template\*\*.tmpl` contain all templates
- [x] Built using [`cobra`](https://github.com/spf13/cobra) and [`viper`](https://github.com/spf13/viper) (without func inits!!!)
  - A lot of sample Cobra/Viper code rely on `func init()` making it more difficult to reuse. 
  - This code cleanly separates CLI/configruation invocation from client library calls - using `viper.Unmarshal` to transfer from Viper to `pkg.Config` structure.
- [X] Logging from the [`logrus`](https://github.com/sirupsen/logrus) library
- [X] HTTP serving/routing with middleware from [`go-chi`](https://github.com/go-chi/chi)
    - Demonstrate the `NewStructuredLogger` middleware, decorating each route
    - Custom middleware handlers to create various Contexts
- [x] An example Dockerfile for a docker workflow

# The Story of 00-newapp-template
There is a vendor named ACME who provides API to access to `Gophers` and `Things`. Because I use the ACME API **all the time** to track `Gophers` and their `Things` and I have decided to create a CLI tool to perform the HTTP API calls needed and output a text table or JSON structure. Ideally using a simple command like:
```
  ./gopherit list 
```
The client would make all the necessary calls to the HTTP ACME services, convert the JSON responses to Go structures andoutput something like:
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
At ACME each `Thing` has a `Gopher` but each `Gopher` does not have `Things`. That means given a `Gopher` we don't know their collection of `Things`. Also, given a `Thing` we only have the associated `Gopher` `ID` and not the `Gopher` `Name` or `Gopher` `Description`. 
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
These JSON structures are defined in `pkg\acme\json.go` and are what the ACME HTTP API return.

While ACME's data structure and JSON schema is useful for their purposes - it's imcompleted for our application's needs. A much more useful structure for our application (listing `Gophers` and their `Things`) looks like this:
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
Here each `Gopher` has a collection of their `Things` and each `Thing` has a reference to a complete `Gopher` (not just the `Gopher` `ID`.)

## The Solution - Build an adapter!
We will write a client library to make ACME HTTP service calls, unmarshal the JSON, filter and convert it into our Go structures. All of this adapter code is in the `pkg\adapter\` folder. 

Our adapter uses services provide by ACME as defined in the `pkg\acme\service.go` file.

Because ACME API charges us per API call, we decided to mock implement our own ACME HTTP API server. Our server will listen for connections on local host and respond just like an ACME server would. All the code to provide a full HTTP API server is in `pkg\acme\server\` and to invoke our ACME server:
```
  ./gopherit server start
```
![Server Startup](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/serverclient.gif)

# Complete Process - Looping GIF Video
This video shows download Go version 1.11.4 and checking out from scratch:

![Quick overview](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/demo.gif)

# MP4 videos demos
- [Download and install Go](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/vokoscreen-2018-12-20_19-00-52.mp4)
- [Clone Gihub Repository](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/vokoscreen-2018-12-20_19-01-56.mp4)
- [Review go.mod Dependencies](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/vokoscreen-2018-12-20_19-03-05.mp4)
- [go test](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/vokoscreen-2018-12-20_19-04-22.mp4)
- [Dockerfile and Docker workflow](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/vokoscreen-2018-12-20_19-05-58.mp4)
- [Hermetic Build](https://github.com/whereiskurt/gopherit/blob/master/00-newapp-template/docs/videos/vokoscreen-2018-12-20_19-07-32.mp4)

## How-to Run:
```
 $ go run cmd/gopherit.go 

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
    gopherit [COMMAND] [SUBCOMMAND] [OPTIONS]

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
      $ gopherit client list

      $ gopherit list
      $ gopherit list -g1 -t2,3
 
      $ gopherit server start
      $ gopherit server start --port=102102 --docroot=./config/docroot
 
      $ gopherit server stop
 
      $ gopherit version

```

## Package Relationships
TODO
