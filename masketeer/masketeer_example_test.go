package masketeer

import (
	"fmt"

	"github.com/facily-tech/go-core/masketeer/email"
)

func ExampleNew() {
	mask := New(&Option{})
	emailMasked := mask.Email("test@example.com")
	fmt.Println(emailMasked)
	// Output:
	// tes...@example.com
}

func ExampleNew_with_options() {
	mask := New(&Option{
		Email: &email.Option{
			NumberOfVisibleCharsOnPrefix: 2,
		},
	})
	emailMasked := mask.Email("test@example.com")
	fmt.Println(emailMasked)
	// Output:
	// te...@example.com
}
