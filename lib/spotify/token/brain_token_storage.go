package token

import (
	"github.com/m-nny/universe/lib/brain"
)

type BrainTokenStorage = brain.Brain

var _ TokenStorage = (*BrainTokenStorage)(nil)
