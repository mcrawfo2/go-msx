// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"github.com/pkg/errors"
	"path/filepath"
)

func GoGenerate(dirs []string) (err error) {
	skeletonConfig := skel.Config()
	for _, gendir := range dirs {
		fp, _ := filepath.Abs(filepath.Join(skeletonConfig.TargetDirectory(), gendir))
		err = skel.GoGenerate(fp)
		if err != nil {
			return errors.Wrap(err, "Error running go generate")
		}
	}
	return nil
}
