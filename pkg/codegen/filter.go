package codegen

import "github.com/getkin/kin-openapi/openapi3"

func filterOperationsByTag(swagger *openapi3.T, opts Configuration) {
	if len(opts.OutputOptions.ExcludeTags) > 0 {
		excludeOperationsWithTags(swagger.Paths, opts.OutputOptions.ExcludeTags)
	}
	if len(opts.OutputOptions.IncludeTags) > 0 {
		includeOperationsWithTags(swagger.Paths, opts.OutputOptions.IncludeTags, false)
	}
}

func excludeOperationsWithTags(paths openapi3.Paths, tags []string) {
	includeOperationsWithTags(paths, tags, true)
}

type OperationPredicate func(op *openapi3.Operation) bool

func includeOperations(paths openapi3.Paths, filter OperationPredicate, exclude bool) {
	for _, pathItem := range paths {
		ops := pathItem.Operations()
		names := make([]string, 0, len(ops))
		for name, op := range ops {
			if filter(op) == exclude {
				names = append(names, name)
			}
		}
		for _, name := range names {
			pathItem.SetOperation(name, nil)
		}
	}
}


func includeOperationsWithTags(paths openapi3.Paths, tags []string, exclude bool) {
	includeOperations(paths, func(op *openapi3.Operation) bool { return operationHasTag(op, tags)}, exclude)
}

// operationHasTag returns true if the operation is tagged with any of tags
func operationHasTag(op *openapi3.Operation, tags []string) bool {
	if op == nil {
		return false
	}
	for _, hasTag := range op.Tags {
		for _, wantTag := range tags {
			if hasTag == wantTag {
				return true
			}
		}
	}
	return false
}

func excludeOperationsIfDepreciated(paths openapi3.Paths) {
	includeOperations(paths, func(op *openapi3.Operation) bool {
		if op == nil {
			return false
		}
		return op.Deprecated
	}, true)
}