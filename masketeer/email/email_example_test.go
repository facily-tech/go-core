package email

import "fmt"

func ExampleMask() {
	emailMasked := Mask("test@example.com", nil)
	fmt.Println(emailMasked)
	// Output:
	// tes...@example.com
}
