// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"bytes"
	"github.com/pkg/errors"
	"strings"
	"unicode/utf8"
)

type ExpressionResolver interface {
	ResolveByName(name string) (ResolvedEntry, error)
}

type expression interface {
	Resolve(r ExpressionResolver) (string, error)
}

type literalExpression string

func (l literalExpression) Resolve(_ ExpressionResolver) (string, error) {
	return string(l), nil
}

type concatenateExpression struct {
	Parts []expression
}

func (c concatenateExpression) Resolve(r ExpressionResolver) (string, error) {
	results := bytes.Buffer{}
	for _, part := range c.Parts {
		result, err := part.Resolve(r)
		if err != nil {
			return "", err
		}
		results.WriteString(result)
	}
	return results.String(), nil
}

type variableExpression struct {
	Name          string
	Default       expression
	RequiredKey   bool
	RequiredValue bool
}

func (v variableExpression) Resolve(r ExpressionResolver) (string, error) {
	entry, err := r.ResolveByName(v.Name)
	if errors.Is(err, ErrNotFound) {
		if v.Default != nil {
			return v.Default.Resolve(r)
		} else if v.RequiredKey {
			return "", err
		} else {
			return "", nil
		}
	} else if err != nil {
		return "", err
	}

	resolvedValue := entry.ResolvedValue.String()
	if resolvedValue == "" && v.RequiredValue {
		return "", ErrEmptyValue
	}

	return resolvedValue, nil
}

type expressionScanner struct {
	Pos   int
	Value string
}

func (s expressionScanner) Remaining() string {
	return s.Value[s.Pos:]
}

func (s expressionScanner) Current() rune {
	return s.Lookahead(0)
}

// Matches compares the bytes of the supplied string to the bytes at the current cursor
func (s expressionScanner) Matches(other string) bool {
	if s.Pos+len(other) > len(s.Value) {
		return false
	}

	for i := s.Pos; i < s.Pos+len(other); i++ {
		if s.Value[i] != other[i-s.Pos] {
			return false
		}
	}

	return true
}

func (s expressionScanner) Lookahead(n int) rune {
	var r rune = 0
	for i, w := s.Pos, 0; i < len(s.Value); i += w {
		r, w = utf8.DecodeRuneInString(s.Value[i:])
		if n == 0 {
			return r
		}
		n--
	}

	return 0
}

func (s expressionScanner) Skip() expressionScanner {
	return expressionScanner{
		Pos:   s.Pos + 1,
		Value: s.Value,
	}
}

func (s expressionScanner) SkipOver(other string) (expressionScanner, error) {
	if !s.Matches(other) {
		return s, errors.Wrapf(ErrParseUnexpectedInput, "Found: %q, Expected: %q", s.Value[s.Pos:], other)
	}

	return expressionScanner{
		Pos:   s.Pos + len(other),
		Value: s.Value,
	}, nil
}

func parseExpression(value string) (expr expression, err error) {
	if !strings.Contains(value, "${") {
		return literalExpression(value), nil
	}

	s := expressionScanner{Value: value}
	expr, _, err = parseExpressionAtDepth(s, 0)
	return
}

func parseExpressionAtDepth(si expressionScanner, depth int) (expr expression, s expressionScanner, err error) {
	s = si
	result := concatenateExpression{}
	for s.Current() != 0 {
		if s.Matches("${") {
			// Parse Variable
			var varExpr variableExpression
			varExpr, s, err = parseVariable(s, depth+1)
			if err != nil {
				return
			}
			result.Parts = append(result.Parts, varExpr)
			continue
		} else if s.Current() == '}' && depth > 0 {
			break
		} else {
			var litExpr literalExpression
			litExpr, s, err = parseLiteral(s, depth)
			if err != nil {
				return
			}
			result.Parts = append(result.Parts, litExpr)
		}
	}

	if len(result.Parts) == 0 {
		expr = literalExpression("")
	} else if len(result.Parts) == 1 {
		expr = result.Parts[0]
	} else {
		expr = result
	}

	return
}

func parseVariable(si expressionScanner, depth int) (expr variableExpression, s expressionScanner, err error) {
	s = si

	s, err = s.SkipOver("${")
	if err != nil {
		return
	}

	var name string
	name, s, err = parseKeyName(s)
	if err != nil {
		return
	}

	var defaultExpr expression
	var requiredKey, requiredValue bool
	switch {
	case s.Matches(":?}"):
		s = s.Skip()
		s = s.Skip()
		requiredKey = true

	case s.Matches(":!}"):
		s = s.Skip()
		s = s.Skip()
		requiredKey = true
		requiredValue = true

	case s.Matches(":"):
		s = s.Skip()
		defaultExpr, s, err = parseExpressionAtDepth(s, depth)
		if err != nil {
			return
		}
	}

	s, err = s.SkipOver("}")
	if err != nil {
		return
	}

	expr = variableExpression{
		Name:          name,
		Default:       defaultExpr,
		RequiredKey:   requiredKey,
		RequiredValue: requiredValue,
	}
	return
}

func parseLiteral(si expressionScanner, depth int) (expr literalExpression, s expressionScanner, err error) {
	s = si
	for s.Current() != 0 {
		if s.Matches("${") {
			expr = literalExpression(s.Value[si.Pos:s.Pos])
			return
		} else if s.Current() == '}' {
			if depth > 0 {
				expr = literalExpression(s.Value[si.Pos:s.Pos])
				return
			}
		}

		s.Pos++
	}

	expr = literalExpression(s.Value[si.Pos:])
	return
}

func parseKeyName(si expressionScanner) (name string, s expressionScanner, err error) {
	s = si
	for s.Current() != 0 {
		switch s.Current() {
		case ':', '}':
			name = s.Value[si.Pos:s.Pos]
			return
		}
		if s.Matches("${") {
			err = ErrParseInvalidVariableReference
			return
		}
		s.Pos++
	}

	name = s.Value[si.Pos:]
	return
}
