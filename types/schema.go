package types

import (
	"fmt"
)

/**
Schema Definition
A Schema is created by supplying the root types of each type of operation,
query and mutation (optional). A schema definition is then supplied to the
validator and executor.
Example:
    myAppSchema, err := NewGraphQLSchema(GraphQLSchemaConfig({
      Query: MyAppQueryRootType
      Mutation: MyAppMutationRootType
    });
*/
type GraphQLSchemaConfig struct {
	Query    *GraphQLObjectType
	Mutation *GraphQLObjectType
}

// chose to name as GraphQLTypeMap instead of TypeMap
type GraphQLTypeMap map[string]GraphQLType

type GraphQLSchema struct {
	schemaConfig GraphQLSchemaConfig
	typeMap      GraphQLTypeMap
	directives   []GraphQLDirective
}

func NewGraphQLSchema(config GraphQLSchemaConfig) (GraphQLSchema, error) {
	var err error

	schema := GraphQLSchema{}

	// if schema config contains error at creation time, return those errors
	if config.Query != nil && config.Query.err != nil {
		return schema, config.Query.err
	}
	if config.Mutation != nil && config.Mutation.err != nil {
		return schema, config.Mutation.err
	}

	schema.schemaConfig = config

	// Build type map now to detect any errors within this schema.
	typeMap := GraphQLTypeMap{}
	objectTypes := []*GraphQLObjectType{
		schema.GetQueryType(),
		schema.GetMutationType(),
		__Type,
		__Schema,
	}
	for _, objectType := range objectTypes {
		if objectType == nil {
			continue
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
		case *GraphQLObjectType:
			for _, iface := range ttype.GetInterfaces() {
				assertObjectImplementsInterface(ttype, iface)
			}
		}
	}

	return schema, nil
}

func (gq *GraphQLSchema) GetQueryType() *GraphQLObjectType {
	return gq.schemaConfig.Query
}

func (gq *GraphQLSchema) GetMutationType() *GraphQLObjectType {
	return gq.schemaConfig.Mutation
}

func (gq *GraphQLSchema) GetDirectives() []GraphQLDirective {
	return gq.directives
}

func (gq *GraphQLSchema) GetTypeMap() GraphQLTypeMap {
	return gq.typeMap
}

func (gq *GraphQLSchema) GetType(name string) GraphQLType {
	return gq.GetTypeMap()[name]
}

func typeMapReducer(typeMap GraphQLTypeMap, objectType GraphQLType) (GraphQLTypeMap, error) {
	var err error
	if objectType == nil || objectType.GetName() == "" {
		return typeMap, nil
	}

	switch objectType := objectType.(type) {
	case *GraphQLList:
		if objectType.OfType != nil {
			return typeMapReducer(typeMap, objectType.OfType)
		}
	case *GraphQLNonNull:
		if objectType.OfType != nil {
			return typeMapReducer(typeMap, objectType.OfType)
		}
	case *GraphQLObjectType:
		if objectType.GetName() == "__Type" && objectType.err != nil {
			return typeMap, nil
		}
	}

	if mappedObjectType, ok := typeMap[objectType.GetName()]; ok {
		err := invariant(
			mappedObjectType == objectType,
			fmt.Sprintf(`Schema must contain unique named types but contains multiple types named "%v".`, objectType.GetName()),
		)
		if err != nil {
			return typeMap, err
		}
		return typeMap, err
	}
	if objectType.GetName() == "" {
		return typeMap, nil
	}

	typeMap[objectType.GetName()] = objectType

	switch objectType := objectType.(type) {
	case *GraphQLUnionType:
		for _, innerObjectType := range objectType.GetPossibleTypes() {
			typeMap, err = typeMapReducer(typeMap, innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	case *GraphQLInterfaceType:
		for _, innerObjectType := range objectType.GetPossibleTypes() {
			typeMap, err = typeMapReducer(typeMap, innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	case *GraphQLObjectType:
		for _, innerObjectType := range objectType.GetInterfaces() {
			typeMap, err = typeMapReducer(typeMap, innerObjectType)
			if err != nil {
				return typeMap, err
			}
		}
	}

	switch objectType := objectType.(type) {
	case *GraphQLObjectType:
		fieldMap := objectType.GetFields()
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
	case *GraphQLInterfaceType:
		fieldMap := objectType.GetFields()
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
		//	case *GraphQLInputObjectType:
	}
	return typeMap, nil
}

func assertObjectImplementsInterface(object *GraphQLObjectType, iface *GraphQLInterfaceType) error {
	objectFieldMap := object.GetFields()
	ifaceFieldMap := iface.GetFields()

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
			argName := ifaceArg.Name
			var objectArg *GraphQLArgument
			for _, arg := range objectField.Args {
				if arg.Name == argName {
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
						`but %v.%v($%v:) provides `+
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
			argName := objectArg.Name
			var ifaceArg *GraphQLArgument
			for _, arg := range ifaceField.Args {
				if arg.Name == argName {
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

func isEqualType(typeA GraphQLType, typeB GraphQLType) bool {
	switch typeA := typeA.(type) {
	case *GraphQLNonNull:
		switch typeB := typeB.(type) {
		case *GraphQLNonNull:
			return isEqualType(typeA.OfType, typeB.OfType)
		default:
			return typeA.GetName() == typeB.GetName()
		}
	case *GraphQLList:
		switch typeB := typeB.(type) {
		case *GraphQLList:
			return isEqualType(typeA.OfType, typeB.OfType)
		default:
			return typeA.GetName() == typeB.GetName()
		}
	default:
		return typeA.GetName() == typeB.GetName()
	}
}
