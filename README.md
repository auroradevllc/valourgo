# Valour for Go

> **Important:** This library is **very early alpha**. APIs may change without notice, features may be incomplete, and parts of the library may be broken.


valourgo provides a Go-friendly way to interact with the [Valour](https://github.com/Valour-Software/Valour) API, as documented here:

https://app.valour.gg/swagger/index.html

The goal of this project is to make it easier to build bots, services, and integrations for Valour using Go, without having to manually manage HTTP requests and low-level API details.

## Installation

```bash
go get github.com/auroradevllc/valourgo
````

*(Module path may change during early development.)*

## Usage

A basic example bot is available at [examples/bot/main.go](examples/bot/main.go)

The example demonstrates:

* Creating a client
* Authenticating with Valour
* Connecting the Realtime API (SignalR)

This is the best place to start to understand the intended usage of the SDK.

## API Coverage

This SDK is based on the official Valour API:

[https://app.valour.gg/swagger/index.html](https://app.valour.gg/swagger/index.html)

Not all endpoints may be implemented yet.

## Contributing

Issues, pull requests, and feedback are welcome, especially while the project is in its early stages.

Expect rough edges, refactors, and breaking changes.

## License

A copy of the license is provided in LICENSE, and available here:

```
Copyright 2026 Aurora Development LLC

Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED “AS IS” AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
```