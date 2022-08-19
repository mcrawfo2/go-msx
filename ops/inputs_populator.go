// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/pkg/errors"
	"reflect"
)

type InputDecoder interface {
	DecodePrimitive(pf *PortField) (result types.Optional[string], err error)
	DecodeContent(pf *PortField) (content Content, err error)
	// DecodeArray(pf *PortField) (result []string, err error)
	// DecodeObject(pf *PortField) (result types.Pojo, err error)
	// DecodeFile(pf *PortField) (result *multipart.FileHeader, err error)
	// DecodeFileArray(pf *PortField) (result []*multipart.FileHeader, err error)
}

type InputsPopulator struct {
	portStruct *interface{}
	port       *Port
	decoder    InputDecoder
}

func (p InputsPopulator) injector(pf *PortField) PortFieldInjector {
	return NewPortFieldInjector(pf, *p.portStruct)
}

func (p InputsPopulator) PopulateInputs() (interface{}, error) {
	err := p.populateInputPortStruct()
	if err != nil {
		return nil, err
	}

	// Auto-validation for validatable Port Struct
	if p.portStruct != nil {
		portStructValue := reflect.ValueOf(*p.portStruct)
		if err = validate.ValidateValue(portStructValue); err != nil {
			return nil, errors.Wrap(err, "Failed to validate inputs")
		}

		return *p.portStruct, nil
	}

	return nil, nil
}

func (p InputsPopulator) populateInputPortStruct() (err error) {
	for _, portField := range p.port.Fields {
		switch portField.Type.Shape {
		case FieldShapePrimitive:
			err = p.populatePrimitive(portField)
		case FieldShapeArray,
			FieldShapeObject,
			FieldShapeFile,
			FieldShapeFileArray:
			err = errors.Errorf("Unimplemented input shape: %q", portField.Type.Shape)
		case FieldShapeContent:
			err = p.populateContent(portField)
		}

		// TODO: Other shapes

		if err != nil {
			return errors.Wrapf(err, "Failed to populate %q field %q", portField.Group, portField.Name)
		}
	}

	return nil
}

func (p InputsPopulator) populatePrimitive(pf *PortField) error {
	optionalValue, err := p.decoder.DecodePrimitive(pf)
	if err != nil {
		return err
	}

	if !optionalValue.IsPresent() {
		if !pf.Optional {
			return errors.Wrapf(ErrMissingRequiredValue,
				"No value or default found for field peer %q",
				pf.Peer)
		} else {
			return nil
		}
	}

	return p.injector(pf).InjectPrimitive(optionalValue.Value())
}

func (p InputsPopulator) populateContent(pf *PortField) error {
	contentSource, err := p.decoder.DecodeContent(pf)
	if err != nil {
		return err
	}

	if !contentSource.IsPresent() {
		if !pf.Optional {
			return errors.Wrapf(ErrMissingRequiredValue,
				"No value or default found for field peer %q",
				pf.Peer)
		} else {
			return nil
		}
	}

	return p.injector(pf).InjectContent(contentSource)
}

func NewInputsPopulator(port *Port, decoder InputDecoder) InputsPopulator {
	return InputsPopulator{
		portStruct: port.NewStruct(),
		port:       port,
		decoder:    decoder,
	}
}
