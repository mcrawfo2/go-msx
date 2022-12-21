// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/pkg/errors"
	"reflect"
)

const messageNoValueFound = "No value or default found for field peer %q"

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
	if p.portStruct == nil {
		return nil
	}

	for _, portField := range p.port.Fields {
		switch portField.Type.Shape {
		case FieldShapePrimitive:
			err = p.populatePrimitive(portField)
		case FieldShapeArray:
			err = p.populateArray(portField)
		case FieldShapeObject:
			err = p.populateObject(portField)
		case FieldShapeFile:
			err = p.populateFile(portField)
		case FieldShapeFileArray:
			err = p.populateFileArray(portField)
		case FieldShapeContent:
			err = p.populateContent(portField)
		case FieldShapeAny:
			err = p.populateAny(portField)
		}

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
				messageNoValueFound,
				pf.Peer)
		} else {
			return nil
		}
	}

	return p.injector(pf).InjectPrimitive(optionalValue.Value())
}

func (p InputsPopulator) populateArray(pf *PortField) error {
	values, err := p.decoder.DecodeArray(pf)
	if err != nil {
		return err
	}

	if len(values) == 0 && pf.Optional {
		return nil
	}

	return p.injector(pf).InjectArray(values)
}

func (p InputsPopulator) populateObject(pf *PortField) error {
	value, err := p.decoder.DecodeObject(pf)
	if err != nil {
		return err
	}

	if value == nil && pf.Optional {
		return nil
	}

	return p.injector(pf).InjectObject(value)
}

func (p InputsPopulator) populateContent(pf *PortField) error {
	contentSource, err := p.decoder.DecodeContent(pf)
	if err != nil {
		return err
	}

	if !contentSource.IsPresent() {
		if !pf.Optional {
			return errors.Wrapf(ErrMissingRequiredValue,
				messageNoValueFound,
				pf.Peer)
		} else {
			return nil
		}
	}

	return p.injector(pf).InjectContent(contentSource)
}

func (p InputsPopulator) populateFile(pf *PortField) error {
	file, err := p.decoder.DecodeFile(pf)
	if err != nil {
		return err
	}

	if file == nil {
		if !pf.Optional {
			return errors.Wrapf(ErrMissingRequiredValue,
				messageNoValueFound,
				pf.Peer)
		} else {
			return nil
		}
	}

	return p.injector(pf).InjectFile(file)
}

func (p InputsPopulator) populateFileArray(pf *PortField) error {
	files, err := p.decoder.DecodeFileArray(pf)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	return p.injector(pf).InjectFileArray(files)
}

func (p InputsPopulator) populateAny(pf *PortField) error {
	optionalValue, err := p.decoder.DecodeAny(pf)
	if err != nil {
		return err
	}

	if !optionalValue.IsPresent() {
		if !pf.Optional {
			return errors.Wrapf(ErrMissingRequiredValue,
				messageNoValueFound,
				pf.Peer)
		} else {
			return nil
		}
	}

	return p.injector(pf).InjectAny(optionalValue.Value())
}

func NewInputsPopulator(port *Port, decoder InputDecoder) InputsPopulator {
	var portStruct *interface{}
	if port != nil {
		portStruct = port.NewStruct()
	}

	return InputsPopulator{
		portStruct: portStruct,
		port:       port,
		decoder:    decoder,
	}
}
