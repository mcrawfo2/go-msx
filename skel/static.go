// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// This file loads static files from an embed.FS and returns them in the structure
// provided by bou.ke/staticfiles (staticFiles) and expected by the rest of skel

package skel

import (
	"embed"
	"path"
)

//go:embed all:_templates/*
var statics embed.FS

const (
	staticRoot = "_templates"
)

func ReadStaticFile(filename string) ([]byte, error) {
	staticFileName := path.Join(staticRoot, filename)
	return statics.ReadFile(staticFileName)
}
