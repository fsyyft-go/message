// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"github.com/google/wire"
)

var (
	ProviderSet = wire.NewSet(NewWebServer)
)
