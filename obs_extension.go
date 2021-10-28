package uoa

type extensionOptions interface{}

type extensionHeaders func(headers map[string][]string, isObs bool) error
