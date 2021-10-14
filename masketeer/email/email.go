package email

import "strings"

type Option struct {

	/**
	NumberOfVisibleCharsOnPrefix is the number of visible chars that will be showed
	default value 3, if the email doesn't have more than it number before @, then it
	will only show the first char of user email.
	The domain will not be hidden on this version.
	*/
	NumberOfVisibleCharsOnPrefix int
}

/*
	Mask will mask the email with given options or it default values
	if "@" was not found then it will return a empty string
*/
func Mask(email string, opt *Option) string {
	opt = getDefault(opt)

	if idx := strings.Index(email, "@"); idx > 0 {
		prefix := email[0:idx]

		if len(prefix) > opt.NumberOfVisibleCharsOnPrefix {
			prefix = prefix[0:opt.NumberOfVisibleCharsOnPrefix]
		} else {
			prefix = prefix[0:1]
		}

		e := prefix + "..." + email[idx:]

		return e
	}

	return ""
}
