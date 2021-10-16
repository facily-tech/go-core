package phone

const defaultNumberOfVisibleCharsOnSufix = 4
const DefaultUseAsVisibleChars = "0123456789"

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
