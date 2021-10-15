package email

import "fmt"

func ExampleMask() {
	emailMasked := Mask("test@example.com", nil)
	fmt.Println(emailMasked)
	// Output:
	// tes...@example.com
}

func ExampleMask_with_options() {
	emailMasked := Mask("test@example.com", &Option{
		NumberOfVisibleCharsOnPrefix: 2,
	})
	fmt.Println(emailMasked)
	// Output:
	// te...@example.com
}
