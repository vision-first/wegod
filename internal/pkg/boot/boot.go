package boot

type ServiceProvider interface {
	Boot() error
}

type Bootstrapper struct {
	providers []ServiceProvider
	isBooted bool
}

func NewBootstrapper(providers []ServiceProvider) *Bootstrapper {
	return &Bootstrapper{
		providers: providers,
	}
}

func (b *Bootstrapper) Boot() error {
	if b.isBooted {
		return nil
	}
	for _, provider := range b.providers {
		if err := provider.Boot(); err != nil {
			return err
		}
	}
	b.isBooted = true
	return nil
}
