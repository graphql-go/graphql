package graphql

import (
	"fmt"
)

/**
Schema Definition
A Schema is created by supplying the root types of each type of operation,
query and mutation (optional). A schema definition is then supplied to the
validator and executor.
Example:
    myAppSchema, err := NewSchema(SchemaConfig({
      Query: MyAppQueryRootType
      Mutation: MyAppMutationRootType
    });
*/
type SchemaConfig struct {
	Query    *Object
	Mutation *Object
}

// chose to name as TypeMap instead of TypeMap
type TypeMap map[string]Type

type Schema struct {
	schemaConfig SchemaConfig
	typeMap      TypeMap
	directives   []*Directive
}

func NewSchema(config SchemaConfig) (Schema, error) {
	var err error

	schema := Schema{}

	err = invariant(config.Query != nil, "Schema query must be Object Type but got: nil.")
	if err != nil {
		return schema, err
	}

	// if schema config contains error at creation time, return those errors
	if config.Query != nil && config.Query.err != nil {
		return schema, config.Query.err
	}
	if config.Mutation != nil && config.Mutation.err != nil {
		return schema, config.Mutation.err
	}

	schema.schemaConfig = config

	// Build type map now to detect any errors within this schema.
	typeMap := TypeMap{}
	objectTypes := []*Object{
		schema.QueryType(),
		schema.MutationType(),
		__Type,
		__Schema,
	}
	for _, objectType := range objectTypes {
		if objectType == nil {
			continue
		}
		if objectType.err != nil {
			return schema, objectType.err
		}
		typeMap, err = typeMapReducer(typeMap, objectType)
		if err != nil {
			return schema, err
		}
	}
	schema.typeMap = typeMap
	// Enforce correct interface implementations
	for _, ttype := range typeMap {
		switch ttype := ttype.(type) {
		case *Object:
			for _, iface := range ttype.Interfaces() {
				err := assertObjectImplementsInterface(ttype, iface)
				if err != nil {
					return schema, err
				}
			}
		}
	}

	return schema, nil
}

func (gq *Schema) QueryType() *Object {
	return gq.schemaConfig.Query
}

func (gq *Schema) MutationType() *Object {
	return gq.schemaConfig.Mutation
}

func (gq *Schema) Directives() []*Directive {
	if len(gq.directives) == 0 {
		gq.directives = []*Directive{
			IncludeDirective,
			SkipDirective,
		}
	}
	return gq.directives
}

func (gq *Schema) Directive(name string) *Directive {
	for _, directive := range gq.Directives() {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (gq *Schema) TypeMap() TypeMap {
	return gq.typeMap
}

func (gq *Schema) Type(name string) Type {
	return gq.TypeMap()[name]
}

func typeMapReducer(typeMap TypeMap, objectType Type) (TypeMap, error) {
	var err error
	if objectType == nil || objectType.Name() == "" {
		return typeMap, nil
	}

	switch objectType := objectType.(type) {
	case *List:
		if objectType.OfType != nil {
			return typeMapReducer(typeMap, objectType.OfType)
		}
	case *NonNull:
		if objectType.OfType != nil {
			return typeMapReducer(typeMap, objectType.OfType)
		}
	case *Object:
		if objectType.err != nil {
			return typeMap, objectType.err
		}
	}

	if mappedObjectType, ok := typeMap[objectType.Name()]; ok {
		err := invariant(
			mappedObjectType == objectType,
			fmt.Sprintf(`Schema must contain unique named types but contains multiple types named "%v".`, objectType.Name()),
		)
		if err != nil {
			return typeMap, err
		}
		return typeMap, err
	}
	if objectType.Name() == "" {
		return typeMap, nil
	}

	typeMap[objectType.Name()] = objectType

	switch objectType := objectType.(type) {
	case *Union:
		types := objectType.PossibleTypes()
		if objectType.err != nil {
			return typeMap, objectType.err
		}
		for _, innerObjectType := range types {
			if innerObjectType.err != nil {
				return typeMap, innerObjectType.err
			}
			typeMap, err = typeMapReducer(typeMap, innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	case *Interface:
		types := objectType.PossibleTypes()
		if objectType.err != nil {
			return typeMap, objectType.err
		}
		for _, innerObjectType := range types {
			if innerObjectType.err != nil {
				return typeMap, innerObjectType.err
			}
			typeMap, err = typeMapReducer(typeMap, innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	case *Object:
		interfaces := objectType.Interfaces()
		if objectType.err != nil {
			return typeMap, objectType.err
		}
		for _, innerObjectType := range interfaces {
			if innerObjectType.err != nil {
				return typeMap, innerObjectType.err
			}
			typeMap, err = typeMapReducer(typeMap, innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	}

	switch objectType := objectType.(type) {
	case *Object:
		fieldMap := objectType.Fields()
		if objectType.err != nil {
			return typeMap, objectType.err
		}
		for _, field := range fieldMap {
			for _, arg := range field.Args {
				typeMap, err = typeMapReducer(typeMap, arg.Type)
				if err != nil {
					return typeMap, err
				}
			}
			typeMap, err = typeMapReducer(typeMap, field.Type)
			if err != nil {
				return typeMap, err
			}
		}
	case *Interface:
		fieldMap := objectType.Fields()
		if objectType.err != nil {
			return typeMap, objectType.err
		}
		for _, field := range fieldMap {
			for _, arg := range field.Args {
				typeMap, err = typeMapReducer(typeMap, arg.Type)
				if err != nil {
					return typeMap, err
				}
			}
			typeMap, err = typeMapReducer(typeMap, field.Type)
			if err != nil {
				return typeMap, err
			}
		}
	case *InputObject:
		fieldMap := objectType.Fields()
		if objectType.err != nil {
			return typeMap, objectType.err
		}
		for _, field := range fieldMap {
			typeMap, err = typeMapReducer(typeMap, field.Type)
			if err != nil {
				return typeMap, err
			}
		}
	}
	return typeMap, nil
}

func assertObjectImplementsInterface(object *Object, iface *Interface) error {
	objectFieldMap := object.Fields()
	ifaceFieldMap := iface.Fields()

	// Assert each interface field is implemented.
	for fieldName, _ := range ifaceFieldMap {
		objectField := objectFieldMap[fieldName]
		ifaceField := ifaceFieldMap[fieldName]

		// Assert interface field exists on object.
		err := invariant(
			objectField != nil,
			fmt.Sprintf(`"%v" expects field "%v" but "%v" does not `+
				`provide it.`, iface, fieldName, object),
		)
		if err != nil {
			return err
		}

		// Assert interface field type matches object field type. (invariant)
		err = invariant(
			isEqualType(ifaceField.Type, objectField.Type),
			fmt.Sprintf(`%v.%v expects type "%v" but `+
				`%v.%v provides type "%v".`,
				iface, fieldName, ifaceField.Type,
				object, fieldName, objectField.Type),
		)
		if err != nil {
			return err
		}

		// Assert each interface field arg is implemented.
		for _, ifaceArg := range ifaceField.Args {
			argName := ifaceArg.PrivateName
			var objectArg *Argument
			for _, arg := range objectField.Args {
				if arg.PrivateName == argName {
					objectArg = arg
					break
				}
			}
			// Assert interface field arg exists on object field.
			err = invariant(
				objectArg != nil,
				fmt.Sprintf(`%v.%v expects argument "%v" but `+
					`%v.%v does not provide it.`,
					iface, fieldName, argName,
					object, fieldName),
			)
			if err != nil {
				return err
			}

			// Assert interface field arg type matches object field arg type.
			// (invariant)
			err = invariant(
				isEqualType(ifaceArg.Type, objectArg.Type),
				fmt.Sprintf(
					`%v.%v(%v:) expects type "%v" `+
						`but %v.%v(%v:) provides `+
						`type "%v".`,
					iface, fieldName, argName, ifaceArg.Type,
					object, fieldName, argName, objectArg.Type),
			)
			if err != nil {
				return err
			}
		}
		// Assert argument set invariance.
		for _, objectArg := range objectField.Args {
			argName := objectArg.PrivateName
			var ifaceArg *Argument
			for _, arg := range ifaceField.Args {
				if arg.PrivateName == argName {
					ifaceArg = arg
					break
				}
			}
			err = invariant(
				ifaceArg != nil,
				fmt.Sprintf(`%v.%v does not define argument "%v" but `+
					`%v.%v provides it.`,
					iface, fieldName, argName,
					object, fieldName),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isEqualType(typeA Type, typeB Type) bool {
	if typeA, ok := typeA.(*NonNull); ok {
		if typeB, ok := typeB.(*NonNull); ok {
			return isEqualType(typeA.OfType, typeB.OfType)
		}
	}
	if typeA, ok := typeA.(*List); ok {
		if typeB, ok := typeB.(*List); ok {
			return isEqualType(typeA.OfType, typeB.OfType)
		}
	}
	return typeA == typeB
}
