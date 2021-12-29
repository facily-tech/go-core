package phone

import (
	"strings"
)

// Option is the options of masked used to mask phone.
type Option struct {
	NumberOfVisibleCharsOnSufix int

	// UseAsVisibleChars  holds the chars that will be used to composite the final string
	// for exemplo stuffs like + () - will not be on the final string and will not be
	// taken into visible chars account
	UseAsVisibleChars string
}

// Mask a phone number with the option given or default parameters.
func Mask(phone string, opt *Option) string {
	opt = getDefault(opt)
	phone = strings.TrimSpace(phone)
	count := 0
	maskedPhone := ""
	size := len(phone)

	for i := 0; i < size; i++ {
		reverseIndex := size - i
		ch := phone[reverseIndex-1 : reverseIndex]
		if !strings.Contains(opt.UseAsVisibleChars, ch) {
			continue
		}

		count++
		if count > opt.NumberOfVisibleCharsOnSufix {
			break
		}

		maskedPhone = ch + maskedPhone
	}

	if len(maskedPhone) > 0 {
		return "..." + maskedPhone
	}

	return ""
}
