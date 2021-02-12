// Package grate opens tabular data files (such as spreadsheets and delimited plaintext files)
// and allows programmatic access to the data contents in a consistent interface.
package grate

import (
	"errors"
)

// Source represents a set of data collections.
type Source interface {
	// List the individual data tables within this source.
	List() ([]string, error)

	// Get a Collection from the source by name.
	Get(name string) (Collection, error)
}

// OpenFunc defines a Source's instantiation function.
// It should return ErrNotInFormat immediately if filename is not of the correct file type.
type OpenFunc func(filename string) (Source, error)

// ErrNotInFormat is used to auto-detect file types using the defined OpenFunc
// It is returned by OpenFunc when the code does not detect correct file formats.
var ErrNotInFormat = errors.New("grate: file is not in this format")

// Open a tabular data file and return a Source for accessing it's contents.
func Open(filename string) (Source, error) {
	for _, o := range srcTable {
		src, err := o(filename)
		if err == nil {
			return src, nil
		}
		if err != ErrNotInFormat {
			return nil, err
		}
	}
	return nil, errors.New("grate: file format is not known/supported")
}

var srcTable = make(map[string]OpenFunc)

// Register the named source as a grate datasource implementation.
func Register(name string, opener OpenFunc) error {
	if _, ok := srcTable[name]; ok {
		return errors.New("grate: source already registered")
	}
	srcTable[name] = opener
	return nil
}

// Collection represents an iterable collection of records.
type Collection interface {
	// Next advances to the next record of content.
	// It MUST be called prior to any Scan().
	Next() bool

	// Strings extracts values from the current record into a list of strings.
	Strings() []string

	// Scan extracts values from the current record into the provided arguments
	// Arguments must be pointers to one of 5 supported types:
	//     bool, int, float64, string, or time.Time
	Scan(args ...interface{}) error

	// IsEmpty returns true if there are no data values.
	IsEmpty() bool

	// Err returns the last error that occured.
	Err() error
}
