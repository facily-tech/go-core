package env

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// Example of usage.
func Example_ofUsage() {
	// Manuall setting the env environment variable to make test works
	if err := os.Setenv("HOST", "0.0.0.0:8080"); err != nil {
		panic(err)
	}

	// Example of usage start below.
	ctx := context.Background()

	type config struct {
		Host    string        `env:"HOST,required"`
		Timeout time.Duration `env:"TIMEOUT,default=10s"`
	}

	cnf := &config{}
	prefix := ""

	if err := LoadEnv(ctx, cnf, prefix); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Host: %s\n", cnf.Host)
	fmt.Printf("Timeout: %s\n", cnf.Timeout)

	// Output:
	// Host: 0.0.0.0:8080
	// Timeout: 10s
	//
}

func Example_ofUsageWithPrefix() {
	// Manuall setting the env environment variable to make test works
	if err := os.Setenv("HTTP_HOST", "0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}

	// this overrides it's default value
	if err := os.Setenv("HTTP_TIMEOUT", "5s"); err != nil {
		log.Fatal(err)
	}

	// Example of usage start below.
	ctx := context.Background()

	type config struct {
		Host    string        `env:"HOST,required"`
		Timeout time.Duration `env:"TIMEOUT,default=10s"`
	}

	cnf := &config{}
	prefix := "HTTP_"

	if err := LoadEnv(ctx, cnf, prefix); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Host: %s\n", cnf.Host)
	fmt.Printf("Timeout: %s\n", cnf.Timeout)

	// Output:
	// Host: 0.0.0.0:8080
	// Timeout: 5s
	//
}

func Example_ofMissingParameters() {
	// Example of usage start below.
	ctx := context.Background()

	type config struct {
		UndeclaredParameter string        `env:"UNDECLARED_PARAMETER,required"`
		Timeout             time.Duration `env:"TIMEOUT,default=10s"`
	}

	cnf := &config{}
	prefix := ""

	if err := LoadEnv(ctx, cnf, prefix); err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("UndeclaredParameter: %s\n", cnf.UndeclaredParameter)
	fmt.Printf("Timeout: %s\n", cnf.Timeout)

	// Output:
	// error while creating config from environment: UndeclaredParameter: missing required value: UNDECLARED_PARAMETER
}
