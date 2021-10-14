package masketeer

import "github.com/facily-tech/go-core/masketeer/email"

type IMasketeer interface {
	Email(eml string) string
}

type Masketer struct {
	opt *Option
}

type Option struct {
	Email *email.Option
}

func New(opt *Option) *Masketer {
	if opt == nil {
		opt = &Option{}
	}

	return &Masketer{
		opt: opt,
	}
}

func (m *Masketer) Email(eml string) string {
	return email.Mask(eml, m.opt.Email)
}
