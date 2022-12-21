// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	"mime/multipart"
	"sort"
	"strings"
)

const (
	FieldStyleSimple         = "simple"
	FieldStyleForm           = "form"
	FieldStyleSpaceDelimited = "spaceDelimited"
	FieldStylePipeDelimited  = "pipeDelimited"
	FieldStyleDeepObject     = "deepObject"
)

type EndpointRequestDecoder struct {
	DataSource         RequestDataSource
	Body               []byte
	defaultContentType string
	defaultEncoding    string
}

func (d *EndpointRequestDecoder) DecodeFile(pf *ops.PortField) (result *multipart.FileHeader, err error) {
	name := pf.Peer
	optional := pf.Optional

	values, err := d.decodeFormMultiFile(name)
	if err != nil {
		return nil, err
	}

	if len(values) > 0 {
		return values[0], nil
	}

	if !optional {
		return nil, errors.Wrap(ErrMissingRequiredFile, name)
	}

	return nil, nil
}

func (d *EndpointRequestDecoder) DecodeFileArray(pf *ops.PortField) (result []*multipart.FileHeader, err error) {
	return d.decodeFormMultiFile(pf.Peer)
}

func (d *EndpointRequestDecoder) DecodeContent(pf *ops.PortField) (result ops.Content, err error) {
	switch pf.Group {
	case FieldGroupHttpBody:
		var payload []byte

		bodyContentOptions := d.DataSource.BodyContentOptions(d.defaultContentType, d.defaultEncoding)
		payload, err = d.DataSource.Body()
		if err != nil {
			return
		}

		result = ops.NewContentFromBytes(bodyContentOptions, payload)

	default:
		err = errors.Errorf("Cannot retrieve %q value from field group %q", pf.Type.Shape, pf.Group)
	}

	return
}

func (d *EndpointRequestDecoder) DecodePrimitive(pf *ops.PortField) (result types.Optional[string], err error) {
	in := pf.Group
	name := pf.Peer
	specifiedStyle := types.
		NewOptionalStringFromString(pf.Options["style"]).
		NilIfEmpty()

	switch in {
	case FieldGroupHttpHeader:
		values := d.DataSource.Headers().Values(name)

		switch specifiedStyle.OrElse(FieldStyleSimple) {
		case FieldStyleSimple:
			return d.decodeSimplePrimitive(values), nil
		}

	case FieldGroupHttpPath:
		pathParameters := d.DataSource.PathParameters()
		pathParameter, ok := pathParameters[name]
		if !ok {
			return types.Optional[string]{}, nil
		}

		values := []string{pathParameter}
		switch specifiedStyle.OrElse(FieldStyleSimple) {
		case FieldStyleSimple:
			return d.decodeSimplePrimitive(values), nil
		}

	case FieldGroupHttpQuery:
		queryParameters := d.DataSource.Query()
		values, ok := queryParameters[name]
		if !ok {
			return types.Optional[string]{}, nil
		}

		switch specifiedStyle.OrElse(FieldStyleForm) {
		case FieldStyleForm, FieldStyleSpaceDelimited, FieldStylePipeDelimited:
			return d.decodeSimplePrimitive(values), nil
		}

	case FieldGroupHttpCookie:
		cookies := d.DataSource.Cookies()
		var values []string
		for _, cookie := range cookies {
			if cookie.Name == name {
				values = append(values, cookie.Value)
				break
			}
		}

		switch specifiedStyle.OrElse(FieldStyleForm) {
		case FieldStyleForm:
			return d.decodeSimplePrimitive(values), nil
		}

	case FieldGroupHttpForm: // Form body field or struct
		form, _, err := d.DataSource.Form()
		if err != nil {
			return types.Optional[string]{}, err
		}
		values := form[name]

		// Style is not used for form fields
		return d.decodeSimplePrimitive(values), nil

	default:
		return types.Optional[string]{}, errors.Wrapf(ErrUnknownFieldSource, in)
	}

	return types.Optional[string]{}, errors.Wrapf(ErrUnsupportedStyle, specifiedStyle.String())
}

func (d *EndpointRequestDecoder) DecodeArray(pf *ops.PortField) (result []string, err error) {
	in := pf.Group
	name := pf.Peer
	specifiedStyle := types.NewOptionalStringFromString(pf.Options["style"]).NilIfEmpty()
	explode, _ := pf.BoolOption("explode")
	style := ""

	switch in {
	case FieldGroupHttpHeader:
		values := d.DataSource.Headers().Values(name)
		style = specifiedStyle.OrElse(FieldStyleSimple)
		switch style {
		case FieldStyleSimple: // ignore explode, always csv
			return d.decodeFormArray(values, false), nil
		}

	case FieldGroupHttpPath:
		pathParameters := d.DataSource.PathParameters()
		pathParameter, ok := pathParameters[name]
		if !ok {
			return nil, nil
		}

		values := []string{pathParameter}
		style = specifiedStyle.OrElse(FieldStyleSimple)
		switch style {
		case FieldStyleSimple:
			return d.decodeFormArray(values, false), nil
		}

	case FieldGroupHttpQuery:
		queryParameters := d.DataSource.Query()
		values, ok := queryParameters[name]
		if !ok {
			return nil, nil
		}

		style = specifiedStyle.OrElse(FieldStyleForm)
		switch style {
		case FieldStyleForm:
			return d.decodeFormArray(values, explode), nil
		case FieldStyleSpaceDelimited:
			return d.decodeSeparatedArray(values, " "), nil
		case FieldStylePipeDelimited:
			return d.decodeSeparatedArray(values, "|"), nil
		}

	case FieldGroupHttpCookie:
		cookies := d.DataSource.Cookies()
		var values []string
		for _, cookie := range cookies {
			if cookie.Name == name {
				values = append(values, cookie.Value)
			}
		}

		style = specifiedStyle.OrElse(FieldStyleForm)
		switch style {
		case FieldStyleForm:
			return d.decodeFormArray(values, explode), nil
		}

	case FieldGroupHttpForm: // Form body field or struct
		form, _, err := d.DataSource.Form()
		if err != nil {
			return nil, err
		}
		values := form[name]

		// Style is not used for form fields, but explode is
		return d.decodeFormArray(values, explode), nil

	default:
		return nil, errors.Wrapf(ErrUnknownFieldSource, in)
	}

	return nil, errors.Wrapf(ErrUnsupportedStyle, style)
}

func (d *EndpointRequestDecoder) DecodeObject(pf *ops.PortField) (result types.Pojo, err error) {
	in := pf.Group
	name := pf.Peer
	specifiedStyle := types.NewOptionalStringFromString(pf.Options["style"]).NilIfEmpty()
	explode, _ := pf.BoolOption("explode")
	style := ""

	switch in {
	case FieldGroupHttpHeader:
		values := d.DataSource.Headers().Values(name)

		style = specifiedStyle.OrElse(FieldStyleSimple)
		switch style {
		case FieldStyleSimple:
			return d.decodeSeparatedObject(values, ",", explode), nil
		}

	case FieldGroupHttpPath:
		pathParameters := d.DataSource.PathParameters()
		pathParameter, ok := pathParameters[name]
		if !ok {
			return nil, nil
		}

		values := []string{pathParameter}
		style = specifiedStyle.OrElse(FieldStyleSimple)
		switch style {
		case FieldStyleSimple:
			return d.decodeSeparatedObject(values, ",", explode), nil
		}

	case FieldGroupHttpQuery:
		queryParameters := d.DataSource.Query()
		var values []string
		var pairs types.StringPairSlice
		style = specifiedStyle.OrElse(FieldStyleForm)
		if style == FieldStyleDeepObject {
			prefix := name + "["
			for k, v := range queryParameters {
				if !strings.HasPrefix(k, prefix) {
					continue
				}
				if len(v) == 0 {
					continue
				}

				pairs = append(pairs, types.StringPair{
					Left:  strings.TrimPrefix(k, name),
					Right: v[0],
				})
			}
			if len(pairs) == 0 {
				return nil, nil
			}
			return d.decodeDeepObjectExplodeObject(pairs), nil
		} else if style == FieldStyleForm && explode {
			for k, v := range queryParameters {
				if len(v) > 0 {
					pairs = append(pairs, types.StringPair{
						Left:  k,
						Right: v[0],
					})
				}
			}
			if len(pairs) == 0 {
				return nil, nil
			}
			return d.decodeFormExplodeObject(pairs), nil
		} else {
			var ok bool
			values, ok = queryParameters[name]
			if !ok {
				return nil, nil
			}
		}

		switch style {
		case FieldStyleForm: // !explode
			return d.decodeSeparatedObject(values, ",", explode), nil
		case FieldStyleSpaceDelimited:
			return d.decodeSeparatedObject(values, " ", explode), nil
		case FieldStylePipeDelimited:
			return d.decodeSeparatedObject(values, "|", explode), nil
		}

	case FieldGroupHttpCookie:
		style = specifiedStyle.OrElse(FieldStyleForm)
		cookies := d.DataSource.Cookies()
		var values []string
		if style == FieldStyleForm && explode {
			var pairs types.StringPairSlice
			for _, v := range cookies {
				pairs = append(pairs, types.StringPair{
					Left:  v.Name,
					Right: v.Value,
				})
			}
			return d.decodeFormExplodeObject(pairs), nil

		} else {
			for _, cookie := range cookies {
				if cookie.Name == name {
					values = append(values, cookie.Value)
				}
			}
		}

		switch style {
		case FieldStyleForm: // !explode
			return d.decodeSeparatedObject(values, ",", explode), nil
		}

	case FieldGroupHttpForm: // Form body field or struct
		form, _, err := d.DataSource.Form()
		if err != nil {
			return nil, err
		}
		values := form[name]

		// swagger-ui encodes form objects as json
		return d.decodeFormObject(values)

	default:
		return nil, errors.Wrapf(ErrUnknownFieldSource, in)
	}

	return nil, errors.Wrapf(ErrUnsupportedStyle, style)
}

func (d *EndpointRequestDecoder) DecodeAny(pf *ops.PortField) (result types.Optional[any], err error) {
	return types.OptionalEmpty[any](), errors.Wrap(ErrNotImplemented, "Any types not supported by rest ops")
}

// decodeSimplePrimitive decodes a single value from the list of values.
// For example: `blue`
func (d *EndpointRequestDecoder) decodeSimplePrimitive(values []string) types.Optional[string] {
	if len(values) == 0 {
		return types.OptionalEmpty[string]()
	}
	return types.OptionalOf(values[0])
}

// decodeSeparatedArray decodes a set of separated values
// For example: `blue,black,brown`
func (d *EndpointRequestDecoder) decodeSeparatedArray(values []string, separator string) []string {
	if len(values) == 0 {
		return nil
	}
	value := values[0]
	return strings.Split(value, separator)
}

// decodeFormArray decodes a set of CSV or multi values
// For example: `blue,black,brown`
func (d *EndpointRequestDecoder) decodeFormArray(values []string, explode bool) []string {
	if explode {
		return values
	} else {
		return d.decodeSeparatedArray(values, ",")
	}
}

// decodeSeparatedObject decodes a set of K/V pairs
// Explode=true example: `R=100,G=200,B=150`
// Explode=false example: `R,100,G,200,B,150`
func (d *EndpointRequestDecoder) decodeSeparatedObject(values []string, separator string, explode bool) types.Pojo {
	if len(values) == 0 {
		return nil
	}
	value := values[0]

	result := make(types.Pojo)
	if explode {
		for _, pair := range strings.Split(value, ",") {
			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				continue
			}
			if len(parts[0]) == 0 {
				continue
			}
			result[parts[0]] = parts[1]
		}
	} else {
		var key string
		for i, part := range strings.Split(value, separator) {
			if i%2 == 0 {
				key = part
				continue
			}

			if len(key) == 0 {
				continue
			}
			result[key] = part
		}
	}

	return result
}

func (d *EndpointRequestDecoder) decodeFormExplodeObject(pairs types.StringPairSlice) types.Pojo {
	if len(pairs) == 0 {
		return nil
	}

	var result = make(types.Pojo)
	for _, pair := range pairs {
		result[pair.Left] = pair.Right
	}
	return result
}

func (d *EndpointRequestDecoder) decodeDeepObjectExplodeObject(pairs types.StringPairSlice) types.Pojo {
	if len(pairs) == 0 {
		return nil
	}

	// Convert path specs to path part slices
	keyCache := map[string][]string{}
	keys := func(k string) (result []string) {
		if cached, ok := keyCache[k]; ok {
			return cached
		}

		parts := strings.Split(k, "][")
		for _, part := range parts {
			result = append(result, strings.TrimSuffix(strings.TrimPrefix(part, "["), "]"))
		}

		keyCache[k] = result
		return
	}

	// Ensure we walk the tree in a rational order
	sort.Slice(pairs, func(i, j int) bool {
		ik := keys(pairs[i].Left)
		jk := keys(pairs[j].Left)

		for p := 0; p < len(ik) && p < len(jk); p++ {
			ip := ik[p]
			jp := jk[p]

			if ip > jp {
				return false
			} else if jp < ip {
				return true
			}
		}

		if len(ik) > len(jk) {
			return false
		} else {
			return true
		}
	})

	var result = types.Pojo{}

	for _, pair := range pairs {
		k := keys(pair.Left)
		here := result
		for i, kp := range k {
			if i == len(k)-1 {
				here[kp] = pair.Right
			} else {
				next, err := here.ObjectValue(kp)
				if err != nil {
					next = types.Pojo{}
					here[kp] = next
				}
				here = next
			}
		}
	}

	return result
}

func (d *EndpointRequestDecoder) decodeFormObject(values []string) (types.Pojo, error) {
	if len(values) == 0 {
		return nil, nil
	}

	value := values[0]
	var result types.Pojo
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *EndpointRequestDecoder) decodeFormMultiFile(name string) ([]*multipart.FileHeader, error) {
	_, multipartForm, err := d.DataSource.Form()
	if err != nil {
		return nil, err
	}

	return multipartForm.File[name], nil
}

func NewRequestDecoder(source RequestDataSource) ops.InputDecoder {
	return &EndpointRequestDecoder{
		DataSource:         source,
		defaultContentType: MediaTypeJson,
	}
}
