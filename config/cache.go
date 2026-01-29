package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache = cache.New(10*time.Minute, 15*time.Minute)
