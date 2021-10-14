package masketeer

import "fmt"

func ExampleNew() {
	mask := New(&Option{})
	emailMasked := mask.Email("test@example.com")
	fmt.Println(emailMasked)
	// Output:
	// tes...@example.com
}
