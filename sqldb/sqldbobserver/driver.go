package sqldbobserver

import (
	"database/sql/driver"
)

// conn defines a tracing wrapper for driver.Driver.
type observerDriver struct {
	driver driver.Driver
}

// Open implements driver.Driver Open.
func (d *observerDriver) Open(name string) (driver.Conn, error) {
	c, err := d.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &conn{conn: c}, nil
}

// TracingDriver creates and returns a new SQL driver with tracing capabilities.
func NewObserverDriver(d driver.Driver, options ...func(*observerDriver)) driver.Driver {
	td := &observerDriver{driver: d}
	for _, option := range options {
		option(td)
	}
	return td
}
