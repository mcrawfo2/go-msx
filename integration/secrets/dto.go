// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package secrets

import "cto-github.cisco.com/NFV-BU/go-msx/integration"

type Pojo integration.Pojo
type PojoArray integration.PojoArray

type EncryptSecretsDTO struct {
	Scope       map[string]string `json:"scope"`
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	SecretNames []string          `json:"secretNames"`
}

type GetSecretRequestDTO struct {
	Names   []string          `json:"names"`
	Encrypt EncryptSecretsDTO `json:"encrypt"`
}

type GenerateSecretRequestDTO struct {
	Names   []string `json:"names"`
	Save    bool     `json:"save"`
	Encrypt *Pojo    `json:"encrypt"`
}

type SecretPolicySetRequest struct {
	AgingRule      AgingRule      `json:"agingRule"`
	CharacterRule  CharacterRule  `json:"characterRule"`
	DictionaryRule DictionaryRule `json:"dictionaryRule"`
	HistoryRule    HistoryRule    `json:"historyRule"`
	KeyRule        KeyRule        `json:"keyRule"`
	LengthRule     LengthRule     `json:"lengthRule"`
}

type SecretPolicyResponse struct {
	SecretPolicySetRequest
	Name string `json:"name"`
}

type AgingRule struct {
	Enabled          bool `json:"enabled"`
	ExpireWarningSec int  `json:"expireWarningSec"`
	GraceAuthNLimit  int  `json:"graceAuthNLimit"`
	MaxAgeSec        int  `json:"maxAgeSec"`
	MinAgeSec        int  `json:"minAgeSec"`
}

type CharacterRule struct {
	ASCIICharactersonly bool   `json:"asciiCharactersonly"`
	Enabled             bool   `json:"enabled"`
	MinDigit            int    `json:"minDigit"`
	MinLowercasechars   int    `json:"minLowercasechars"`
	MinSpecialchars     int    `json:"minSpecialchars"`
	MinUppercasechars   int    `json:"minUppercasechars"`
	SpecialCharacterSet string `json:"specialCharacterSet"`
}

type DictionaryRule struct {
	Enabled              bool `json:"enabled"`
	TestReversedPassword bool `json:"testReversedPassword"`
}

type HistoryRule struct {
	Enabled                    bool `json:"enabled"`
	Passwdhistorycount         int  `json:"passwdhistorycount"`
	PasswdhistorydurationMonth int  `json:"passwdhistorydurationMonth"`
}

type KeyRule struct {
	Algorithm string `json:"algorithm"`
	Enabled   bool   `json:"enabled"`
	Format    string `json:"format"`
}

type LengthRule struct {
	Enabled   bool `json:"enabled"`
	MaxLength int  `json:"maxLength"`
	MinLength int  `json:"minLength"`
}

type SecretsResponse map[string]string

func (s SecretsResponse) Value(key string) string {
	return s[key]
}
