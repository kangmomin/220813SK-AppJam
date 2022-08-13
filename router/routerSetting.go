package router

import (
	"appjam/domain"
	"context"
)

var (
	argonConfig = domain.ArgonConfig{
		Time:   10,
		Memory: 64 * 1024,
		Thread: 4,
		KeyLen: 32,
	}

	ctx = context.Background()
)
