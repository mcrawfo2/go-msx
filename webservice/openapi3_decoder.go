package webservice

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"encoding/xml"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
)

const (
	encodingGzip    = "gzip"
	encodingDeflate = "deflate"
)

var ErrUnknownFieldSource = errors.New("Unknown field source")
var ErrUnsupportedStyle = errors.New("Unsupported style")
var ErrNoContentTypeDecoder = errors.New("No content-type decoder")
var ErrMissingRequiredFile = errors.New("Missing required file")
var ErrMissingRequiredBody = errors.New("Missing required body")

type RequestDataSource interface {
	Cookies() []*http.Cookie
	Headers() http.Header
	Form() (url.Values, *multipart.Form, error)
	Query() url.Values
	PathParameters() map[string]string
	Body() ([]byte, error)
	ReadEntity(e interface{}) (err error)
}

type RestfulRequestDataSource struct {
	Request  *restful.Request
	BodyData []byte
}

func (r *RestfulRequestDataSource) PathParameters() map[string]string {
	return r.Request.PathParameters()
}

func (r *RestfulRequestDataSource) Form() (url.Values, *multipart.Form, error) {
	if r.Request.Request.MultipartForm == nil {
		bodyData, err := r.Body()
		if err != nil {
			return nil, nil, err
		}

		// Generate a new body reader
		r.Request.Request.Body = io.NopCloser(bytes.NewReader(bodyData))

		if r.Request.Request.MultipartForm == nil {
			err = r.Request.Request.ParseMultipartForm(10 * 1024 * 1024)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return r.Request.Request.PostForm, r.Request.Request.MultipartForm, nil
}

func (r *RestfulRequestDataSource) Query() url.Values {
	return r.Request.Request.URL.Query()
}

func (r *RestfulRequestDataSource) Headers() http.Header {
	return r.Request.Request.Header
}

func (r *RestfulRequestDataSource) Cookies() []*http.Cookie {
	return r.Request.Request.Cookies()
}

func (r *RestfulRequestDataSource) Body() (result []byte, err error) {
	if r.BodyData != nil {
		return r.BodyData, nil
	}

	// Store the body
	bodyReader := r.Request.Request.Body
	defer bodyReader.Close()

	contentEncoding := r.Request.Request.Header.Get(headerNameContentEncoding)
	// check if the request body needs decompression
	if encodingGzip == contentEncoding {
		gzipReader, err := gzip.NewReader(bodyReader)
		if err != nil {
			return nil, err
		}

		bodyReader = gzipReader
	} else if encodingDeflate == contentEncoding {
		zlibReader, err := zlib.NewReader(bodyReader)
		if err != nil {
			return nil, err
		}
		bodyReader = zlibReader
	}

	r.BodyData, err = ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	return r.BodyData, nil
}

func (r *RestfulRequestDataSource) ReadEntity(body interface{}) (err error) {
	bodyData, err := r.Body()
	if err != nil {
		return err
	}

	// Generate a new body reader
	r.Request.Request.Body = io.NopCloser(bytes.NewReader(bodyData))

	contentType := r.Request.Request.Header.Get(headerNameContentType)
	switch {
	case strings.Contains(contentType, MIME_JSON):
		decoder := json.NewDecoder(r.Request.Request.Body)
		decoder.UseNumber()
		err = decoder.Decode(body)

	case strings.Contains(contentType, MIME_XML):
		err = xml.NewDecoder(r.Request.Request.Body).Decode(body)

	default:
		err = errors.Wrap(ErrNoContentTypeDecoder, contentType)
	}

	if err != nil {
		if err == io.EOF {
			err = errors.Wrap(err, "Missing required body")
		}
		return NewBadRequestError(err)
	}

	return nil
}

type RequestDecoder interface {
	DecodeFormFile(name string, optional bool) (*multipart.FileHeader, error)
	DecodeFormMultiFile(name string) ([]*multipart.FileHeader, error)
	DecodeBody(body interface{}, optional bool) (err error)
	DecodeBodyToJson(optional bool) (body json.RawMessage, err error)
	DecodePrimitive(in, name, style string, explode bool) (result types.OptionalString, err error)
	DecodeArray(in, name, style string, explode bool) (result []string, err error)
	DecodeObject(in, name, style string, explode bool) (result types.Pojo, err error)
}

type OpenApiRequestDecoder struct {
	DataSource RequestDataSource
	Body       []byte
}

func (d *OpenApiRequestDecoder) DecodeFormFile(name string, optional bool) (*multipart.FileHeader, error) {
	values, err := d.DecodeFormMultiFile(name)
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

func (d *OpenApiRequestDecoder) DecodeFormMultiFile(name string) ([]*multipart.FileHeader, error) {
	_, multipartForm, err := d.DataSource.Form()
	if err != nil {
		return nil, err
	}

	return multipartForm.File[name], nil
}

func (d *OpenApiRequestDecoder) DecodeBodyToJson(optional bool) (body json.RawMessage, err error) {
	if body, err = d.DataSource.Body(); err != nil {
		if errors.Is(err, io.EOF) {
			if optional {
				err = nil
			} else {
				err = errors.Wrap(ErrMissingRequiredBody, "body")
			}
		}
	}

	return
}

func (d *OpenApiRequestDecoder) DecodeBody(body interface{}, optional bool) (err error) {
	// Read the body
	if bodyReader, ok := body.(io.ReadCloser); ok {
		var bodyBytes []byte
		bodyBytes, err = d.DataSource.Body()
		if err != nil {
			return
		}
		bodyReaderValue := reflect.ValueOf(bodyReader)
		bodyReader = io.NopCloser(bytes.NewReader(bodyBytes))
		bodyReaderValue.Set(reflect.ValueOf(bodyReader))
		return
	}

	if err = d.DataSource.ReadEntity(body); err != nil {
		if errors.Is(err, io.EOF) {
			if optional {
				err = nil
			} else {
				err = errors.Wrap(ErrMissingRequiredBody, "body")
			}
		}
	}

	return
}

func (d *OpenApiRequestDecoder) DecodePrimitive(in, name, style string, explode bool) (result types.OptionalString, err error) {
	switch in {
	case "header":
		values := d.DataSource.Headers().Values(name)

		switch style {
		case "simple":
			return d.DecodeSimplePrimitive(values), nil
		}

	case "path":
		pathParameters := d.DataSource.PathParameters()
		pathParameter, ok := pathParameters[name]
		if !ok {
			return types.OptionalString{}, nil
		}

		values := []string{pathParameter}
		switch style {
		case "simple":
			return d.DecodeSimplePrimitive(values), nil
		}

	case "query":
		queryParameters := d.DataSource.Query()
		values, ok := queryParameters[name]
		if !ok {
			return types.OptionalString{}, nil
		}

		switch style {
		case "form", "spaceDelimited", "pipeDelimited":
			return d.DecodeSimplePrimitive(values), nil
		}

	case "cookie":
		cookies := d.DataSource.Cookies()
		var values []string
		for _, cookie := range cookies {
			if cookie.Name == name {
				values = append(values, cookie.Value)
				break
			}
		}

		switch style {
		case "form":
			return d.DecodeSimplePrimitive(values), nil
		}

	case "form": // Form body field or struct
		form, _, err := d.DataSource.Form()
		if err != nil {
			return types.OptionalString{}, err
		}
		values := form[name]

		// Style is not used for form fields
		return d.DecodeSimplePrimitive(values), nil

	default:
		return types.OptionalString{}, errors.Wrapf(ErrUnknownFieldSource, in)
	}

	return types.OptionalString{}, errors.Wrapf(ErrUnsupportedStyle, style)
}

func (d *OpenApiRequestDecoder) DecodeArray(in, name, style string, explode bool) (result []string, err error) {
	switch in {
	case "header":
		values := d.DataSource.Headers().Values(name)

		switch style {
		case "simple": // ignore explode, always csv
			return d.DecodeFormArray(values, false), nil
		}

	case "path":
		pathParameters := d.DataSource.PathParameters()
		pathParameter, ok := pathParameters[name]
		if !ok {
			return nil, nil
		}

		values := []string{pathParameter}
		switch style {
		case "simple":
			return d.DecodeFormArray(values, false), nil
		}

	case "query":
		queryParameters := d.DataSource.Query()
		values, ok := queryParameters[name]
		if !ok {
			return nil, nil
		}

		switch style {
		case "form":
			return d.DecodeFormArray(values, explode), nil
		case "spaceDelimited":
			return d.DecodeSeparatedArray(values, " "), nil
		case "pipeDelimited":
			return d.DecodeSeparatedArray(values, "|"), nil
		}

	case "cookie":
		cookies := d.DataSource.Cookies()
		var values []string
		for _, cookie := range cookies {
			if cookie.Name == name {
				values = append(values, cookie.Value)
			}
		}

		switch style {
		case "form":
			return d.DecodeFormArray(values, explode), nil
		}

	case "form": // Form body field or struct
		form, _, err := d.DataSource.Form()
		if err != nil {
			return nil, err
		}
		values := form[name]

		// Style is not used for form fields
		return d.DecodeFormArray(values, false), nil

	default:
		return nil, errors.Wrapf(ErrUnknownFieldSource, in)
	}

	return nil, errors.Wrapf(ErrUnsupportedStyle, style)
}

func (d *OpenApiRequestDecoder) DecodeObject(in, name, style string, explode bool) (result types.Pojo, err error) {
	switch in {
	case "header":
		values := d.DataSource.Headers().Values(name)

		switch style {
		case "simple":
			return d.DecodeSeparatedObject(values, ",", explode), nil
		}

	case "path":
		pathParameters := d.DataSource.PathParameters()
		pathParameter, ok := pathParameters[name]
		if !ok {
			return nil, nil
		}

		values := []string{pathParameter}
		switch style {
		case "simple":
			return d.DecodeSeparatedObject(values, ",", explode), nil
		}

	case "query":
		queryParameters := d.DataSource.Query()
		var values []string
		var pairs types.StringPairSlice
		if style == "deepObject" {
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
			return d.DecodeDeepObjectExplodeObject(pairs), nil
		} else if style == "form" && explode {
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
			return d.DecodeFormExplodeObject(pairs), nil
		} else {
			var ok bool
			values, ok = queryParameters[name]
			if !ok {
				return nil, nil
			}
		}

		switch style {
		case "form": // !explode
			return d.DecodeSeparatedObject(values, ",", explode), nil
		case "spaceDelimited":
			return d.DecodeSeparatedObject(values, " ", explode), nil
		case "pipeDelimited":
			return d.DecodeSeparatedObject(values, "|", explode), nil
		}

	case "cookie":
		cookies := d.DataSource.Cookies()
		var values []string
		if style == "form" && explode {
			var pairs types.StringPairSlice
			for _, v := range cookies {
				pairs = append(pairs, types.StringPair{
					Left:  v.Name,
					Right: v.Value,
				})
			}
			return d.DecodeFormExplodeObject(pairs), nil

		} else {
			for _, cookie := range cookies {
				if cookie.Name == name {
					values = append(values, cookie.Value)
				}
			}
		}

		switch style {
		case "form": // !explode
			return d.DecodeSeparatedObject(values, ",", explode), nil
		}

	case "form": // Form body field or struct
		form, _, err := d.DataSource.Form()
		if err != nil {
			return nil, err
		}
		values := form[name]

		// swagger-ui encodes form objects as json
		return d.DecodeFormObject(values)

	default:
		return nil, errors.Wrapf(ErrUnknownFieldSource, in)
	}

	return nil, errors.Wrapf(ErrUnsupportedStyle, style)
}

// DecodeSimplePrimitive decodes a single value from the list of values.
// For example: `blue`
func (d *OpenApiRequestDecoder) DecodeSimplePrimitive(values []string) types.OptionalString {
	if len(values) == 0 {
		return types.OptionalString{}
	}
	return types.NewOptionalStringFromString(values[0])
}

// DecodeSeparatedArray decodes a set of separated values
// For example: `blue,black,brown`
func (d *OpenApiRequestDecoder) DecodeSeparatedArray(values []string, separator string) []string {
	if len(values) == 0 {
		return nil
	}
	value := values[0]
	return strings.Split(value, separator)
}

// DecodeFormArray decodes a set of CSV or multi values
// For example: `blue,black,brown`
func (d *OpenApiRequestDecoder) DecodeFormArray(values []string, explode bool) []string {
	if explode {
		return values
	} else {
		return d.DecodeSeparatedArray(values, ",")
	}
}

// DecodeSeparatedObject decodes a set of K/V pairs
// Explode=true example: `R=100,G=200,B=150`
// Explode=false example: `R,100,G,200,B,150`
func (d *OpenApiRequestDecoder) DecodeSeparatedObject(values []string, separator string, explode bool) types.Pojo {
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

func (d *OpenApiRequestDecoder) DecodeFormExplodeObject(pairs types.StringPairSlice) types.Pojo {
	if len(pairs) == 0 {
		return nil
	}

	var result = make(types.Pojo)
	for _, pair := range pairs {
		result[pair.Left] = pair.Right
	}
	return result
}

func (d *OpenApiRequestDecoder) DecodeDeepObjectExplodeObject(pairs types.StringPairSlice) types.Pojo {
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

func (d *OpenApiRequestDecoder) DecodeFormObject(values []string) (types.Pojo, error) {
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
