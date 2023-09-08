## Usage

> <span style="color: red;">How to use it?</span> Please go to the [examples](https://github.com/topazur/sqlc-gen-go-separate/tree/main/examples) folder

## Why do I need this plugin?

> Now let me explain why we don't use the official [codegen/golang](https://github.com/sqlc-dev/sqlc/tree/main/internal/codegen/golang) function. Let's take a look

  1. If there is only one output directory, it makes the code difficult to read and maintain.

  2. If there are multiple output directories, it is difficult for us to manage imports during use. (Must follow the specifications of Golang language)

### Solution ideas

So the plugin was born, by modifying the [templates](https://github.com/sqlc-dev/sqlc/tree/main/internal/codegen/golang/templates) folder to separate types into separate files and aggregating query methods. 

So when using it, only focus on types and aggregation functions.

### Patch Code

- [copy codegen/golang folder form sqlc repo](https://github.com/sqlc-dev/sqlc/tree/main/internal/codegen/golang)

- [templates folder](https://github.com/topazur/sqlc-gen-go-separate/tree/main/internal/templates)

- [patch folder](https://github.com/topazur/sqlc-gen-go-separate/tree/main/internal/patch)
