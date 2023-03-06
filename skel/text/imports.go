// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import "github.com/mcrawfo2/go-jsonschema/pkg/codegen"

const PkgApp = "cto-github.cisco.com/NFV-BU/go-msx/app"
const PkgRestops = "cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
const PkgTypes = "cto-github.cisco.com/NFV-BU/go-msx/types"
const PkgContext = "context"
const PkgOpenapi = "cto-github.cisco.com/NFV-BU/go-msx/schema/openapi"
const PkgLog = "cto-github.cisco.com/NFV-BU/go-msx/log"
const PkgRestopsV2 = PkgRestops + "/v2"
const PkgRestopsV8 = PkgRestops + "/v8"
const PkgSqldb = "cto-github.cisco.com/NFV-BU/go-msx/sqldb"
const PkgOpenApi3 = "github.com/swaggest/openapi-go/openapi3"
const PkgUuid = "github.com/google/uuid"
const PkgPrepared = "cto-github.cisco.com/NFV-BU/go-msx/sqldb/prepared"
const PkgPaging = "cto-github.cisco.com/NFV-BU/go-msx/paging"
const PkgRepository = "cto-github.cisco.com/NFV-BU/go-msx/repository"
const PkgHttp = "net/http"
const PkgGoqu = "github.com/doug-martin/goqu/v9"

const PkgTestHelpers = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
const PkgControllerTest = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers/controllertest"
const PkgWebServiceTest = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
const PkgServiceTest = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers/servicetest"
const PkgRepositoryTest = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers/repositorytest"
const PkgSqlDbTest = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers/sqldbtest"
const PkgSecurityTest = "cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
const PkgTesting = "testing"
const PkgGithubErrors = "github.com/pkg/errors"
const PkgTestifyMock = "github.com/stretchr/testify/mock"
const PkgTestifyAssert = "github.com/stretchr/testify/assert"

var ImportApp = codegen.Import{QualifiedName: PkgApp}
var ImportRestOps = codegen.Import{QualifiedName: PkgRestops}
var ImportTypes = codegen.Import{QualifiedName: PkgTypes}
var ImportContext = codegen.Import{QualifiedName: PkgContext}
var ImportOpenapi = codegen.Import{QualifiedName: PkgOpenapi}
var ImportLog = codegen.Import{QualifiedName: PkgLog}
var ImportRestOpsV2 = codegen.Import{QualifiedName: PkgRestopsV2}
var ImportRestOpsV8 = codegen.Import{QualifiedName: PkgRestopsV8}
var ImportOpenApi3 = codegen.Import{QualifiedName: PkgOpenApi3}
var ImportSqldb = codegen.Import{QualifiedName: PkgSqldb}
var ImportUuid = codegen.Import{QualifiedName: PkgUuid}
var ImportPrepared = codegen.Import{QualifiedName: PkgPrepared, Name: "db"}
var ImportPaging = codegen.Import{QualifiedName: PkgPaging}
var ImportRepository = codegen.Import{QualifiedName: PkgRepository}
var ImportTestHelpers = codegen.Import{QualifiedName: PkgTestHelpers}
var ImportControllerTest = codegen.Import{QualifiedName: PkgControllerTest}
var ImportWebServiceTest = codegen.Import{QualifiedName: PkgWebServiceTest}
var ImportSecurityTest = codegen.Import{QualifiedName: PkgSecurityTest}
var ImportServiceTest = codegen.Import{QualifiedName: PkgServiceTest}
var ImportRepositoryTest = codegen.Import{QualifiedName: PkgRepositoryTest}
var ImportSqlDbTest = codegen.Import{QualifiedName: PkgSqlDbTest}
var ImportTesting = codegen.Import{QualifiedName: PkgTesting}
var ImportHttp = codegen.Import{QualifiedName: PkgHttp}
var ImportTestifyMock = codegen.Import{QualifiedName: PkgTestifyMock}
var ImportTestifyAssert = codegen.Import{QualifiedName: PkgTestifyAssert}
var ImportErrors = codegen.Import{QualifiedName: PkgGithubErrors}
var ImportGoqu = codegen.Import{QualifiedName: PkgGoqu}
