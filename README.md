# Terraform Provider JqChain

JqChain is a custom Terraform provider that mimics the behavior of functional reduce or fold operations, enabling you to process and transform JSON data directly within your Terraform configuration using a sequence of chained JQ expressions.

This provider includes:
- A data source `jqchain_reduce` that applies a sequence of JQ reducers to an initial JSON value
- A function `jqchain.reduce()` for inline use in locals, modules, and outputs

Use Cases:
- Transform JSON data by applying step-by-step filters, mappings, or reshaping operations.
- Build derived values like computed tags, names, or labels by folding input data through a transformation chain.
- Aggregate nested JSON fields by sequentially reducing a complex object into a simplified result using chained JQ logic.
- Clean or normalize configuration inputs by removing unused fields, adjusting formats, or injecting defaults in a declarative way.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

The following documentation provides examples of how to use the provider:
- [Reduce data source](docs/data-sources/reduce.md) – use JQ expressions to transform JSON via a data block
- [Reduce function](docs/functions/reduce.md) – apply chained JQ logic inline using a Terraform function

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
