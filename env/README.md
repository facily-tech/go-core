# Env

[![Go Reference](https://pkg.go.dev/badge/github.com/facily-tech/go-core/env.svg)](https://pkg.go.dev/github.com/facily-tech/go-core/env)

Package env provides helper functions to load environment variables with some degree of control of behavior like: empty,
prefixes, default and mutators. All these beneficits come from github.com/sethvargo/go-envconfig.

## Usage

Take a look into this [usage examples](./env_example_test.go).

But, in a nutshell:
1. Create a structure with "env" anotation someting like this:
    ```go
        type config struct {
            Foo string `env:"FOO,required"`
            Boo string `env:"BOO,default=xpto"`
        }
    ```
2. Load Env Variables on it structure with LoadEnv function:
    ```go
        cnf := config{}
        LoadEnv(ctx, cnf, "")
    ```
