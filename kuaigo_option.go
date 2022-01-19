package kuaigo

import "github.com/LuoHongLiang0921/kuaigo/kpkg/conf"

type Option func(a *App)

type Disable int

func (app *App) WithOptions(options ...Option) {
	for _, option := range options {
		option(app)
	}
}

func WithConfigParser(unmarshaller conf.Unmarshaller) Option {
	return func(a *App) {
		a.configParser = unmarshaller
	}
}

func WithDisable(d Disable) Option {
	return func(a *App) {
		a.disableMap[d] = true
	}
}

//New new a Application
func New(fns ...func() error) (*App, error) {
	app := &App{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}
