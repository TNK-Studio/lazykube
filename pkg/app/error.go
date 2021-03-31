package app

import "errors"

var (
	resourceNotFoundErr   = errors.New("Resource not found. ")
	noResourceSelectedErr = errors.New("No resource selected. ")
	noNamespaceSelectedErr = errors.New("No namespace selected. ")
)
