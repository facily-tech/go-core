package masketeer

import (
	"fmt"

	"github.com/facily-tech/go-core/masketeer/email"
	"github.com/facily-tech/go-core/masketeer/phone"
)

func ExampleNew() {
	mask := New(&Option{})
	emailMasked := mask.Email("test@example.com")
	phoneMasked := mask.Phone("+12 123 123-456-789")
	fmt.Println("masked email:", emailMasked)
	fmt.Println("masked phone:", phoneMasked)
	// Output:
	// masked email: tes...@example.com
	// masked phone: ...6789
}

func ExampleNew_email_with_options() {
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

func ExampleNew_phone_with_options() {
	mask := New(&Option{
		Phone: &phone.Option{
			NumberOfVisibleCharsOnSufix: 3,
		},
	})
	phoneMasked := mask.Phone("+12 123 123-456-789")
	fmt.Println(phoneMasked)
	// Output:
	// ...789
}
