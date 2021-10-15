package masketeer

import "github.com/facily-tech/go-core/masketeer/email"

// IMasketeer is Masketeer interface to use when you wanna make a default setup on its methods.
type IMasketeer interface {
	Email(eml string) string
}

type Masketeer struct {
	opt *Option
}

type Option struct {
	Email *email.Option
}

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
