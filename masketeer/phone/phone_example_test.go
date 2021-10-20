package phone

import "fmt"

func ExampleMask() {
	phoneMasked := Mask("+55 123 1234567", nil)
	fmt.Println(phoneMasked)
	// Output:
	// ...4567
}

func ExampleMask_with_options() {
	phoneMasked := Mask("+55 123 1234567", &Option{
		NumberOfVisibleCharsOnSufix: 3,
	})
	fmt.Println(phoneMasked)
	// Output:
	// ...567
}

func ExampleMask_with_more_used_chars() {
	phoneMasked := Mask("+55 (123) 123-456-789", &Option{
		NumberOfVisibleCharsOnSufix: 6,
		UseAsVisibleChars:           DefaultUseAsVisibleChars + "-",
	})
	fmt.Println(phoneMasked)
	// Output:
	// ...56-789
}
