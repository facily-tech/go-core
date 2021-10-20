package email

const defaultNumberOfVisibleCharsOnPrefix = 3

func getDefault(opt *Option) *Option {
	if opt != nil {
		return opt
	}

	return &Option{
		NumberOfVisibleCharsOnPrefix: defaultNumberOfVisibleCharsOnPrefix,
	}
}
