package masketeer

import (
	"github.com/facily-tech/go-core/masketeer/email"
	"github.com/facily-tech/go-core/masketeer/phone"
)

// IMasketeer is Masketeer interface to use when you wanna make a default setup on its methods.
type IMasketeer interface {
	// Email returns a string with the email masked with options given on New
	// if "@" was not found then it will return a empty string
	Email(eml string) string

	// Phone returns a string with the phone masked with options given on New
	Phone(pho string) string
}

type Masketeer struct {
	opt *Option
}

type Option struct {
	Email *email.Option
	Phone *phone.Option
}

// New returns a new Masketeer struct
func New(opt *Option) *Masketeer {
	if opt == nil {
		opt = &Option{}
	}

	return &Masketeer{
		opt: opt,
	}
}

func (m *Masketeer) Email(eml string) string {
	return email.Mask(eml, m.opt.Email)
}

func (m *Masketeer) Phone(pho string) string {
	return phone.Mask(pho, m.opt.Phone)
}
