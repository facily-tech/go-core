package env

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Example of usage
func ExampleOfUsage() {
	// Manuall setting the env environment variable to make test works
	os.Setenv("HOST", "0.0.0.0:8080")

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

func ExampleOfUsageWithPrefix() {
	// Manuall setting the env environment variable to make test works
	os.Setenv("HTTP_HOST", "0.0.0.0:8080")
	os.Setenv("HTTP_TIMEOUT", "5s") // this overrides it's default value

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

func ExampleOfMissingParameters() {
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
