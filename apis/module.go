package apis

import "go.uber.org/fx"

var Module = fx.Provide(NewCompanyHandler)
