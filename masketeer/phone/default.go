/*
Package phone have functions to mask phone number.
*/
package phone

const (
	defaultNumberOfVisibleCharsOnSufix = 4
	// DefaultUseAsVisibleChars is the default number of visible characters.
	DefaultUseAsVisibleChars = "0123456789"
)

func getDefault(opt *Option) *Option {
	if opt == nil {
		opt = &Option{}
	}

	if opt.NumberOfVisibleCharsOnSufix == 0 {
		opt.NumberOfVisibleCharsOnSufix = defaultNumberOfVisibleCharsOnSufix
	}

	if opt.UseAsVisibleChars == "" {
		opt.UseAsVisibleChars = DefaultUseAsVisibleChars
	}

	return opt
}
