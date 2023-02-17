// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/iancoleman/strcase"
	"path"
	"strings"
)

type DomainMigrationGenerator struct {
	Domain string
	Tenant string
	Folder string
	Spec   Spec
	*text.TextFile
}

func (g DomainMigrationGenerator) createTableSnippet() error {
	return g.AddNewText(
		"Table",
		"createTable",
		`
			CREATE TABLE lower_snake_singular (
				lower_snake_singular_id uuid PRIMARY KEY,
			--#if TENANT_DOMAIN
				tenant_id uuid,
			--#endif TENANT_DOMAIN
				data text
			);
			`)
}

func (g DomainMigrationGenerator) createIndexSnippet() error {
	tenantDomain := types.ComparableSlice[string]{TenantSingle, TenantHierarchy}.Contains(g.Tenant)
	if !tenantDomain {
		return nil
	}

	return g.AddNewText(
		"Index",
		"createIndex",
		`
			CREATE INDEX ON lower_snake_singular(tenant_id);
			`)
}

func (g DomainMigrationGenerator) Generate() error {
	errs := types.ErrorList{
		g.createTableSnippet(),
		g.createIndexSnippet(),
	}

	return errs.Filter()
}

func (g DomainMigrationGenerator) Filename() string {
	description := strcase.ToScreamingSnake(g.Inflector.Inflect("Create Table Title Singular"))
	prefix, _ := skel.NextMigrationPrefix(g.Folder)
	filename := fmt.Sprintf("%s__%s.sql", prefix, description)
	target := path.Join(g.Folder, filename)
	return g.Inflector.Inflect(target)
}

func NewDomainMigrationGenerator(spec Spec) ComponentGenerator {
	appVersion := skel.Config().AppVersion
	appVersionFolder := "V" + strings.ReplaceAll(appVersion, ".", "_")
	folder := path.Join("internal", "migrate", appVersionFolder)

	return DomainMigrationGenerator{
		Domain: generatorConfig.Domain,
		Tenant: generatorConfig.Tenant,
		Folder: folder,
		Spec:   spec,
		TextFile: text.NewTextFile(
			skel.FileFormatSql,
			skel.NewInflector(generatorConfig.Domain),
			"Migration for "+generatorConfig.Domain,
			text.NewSections[text.Snippet](
				"Table",
				"Index",
			)),
	}
}
