package thumbor

import "github.com/google/wire"

var ProviderSet = wire.NewSet(ProvideThumbor)
