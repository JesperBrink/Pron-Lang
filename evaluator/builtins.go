package evaluator

import (
	"Pron-Lang/object"
	"fmt"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"add": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if args[0].Type() == object.ARRAY_OBJ {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)

				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}

			} else if args[0].Type() == object.HASH_OBJ {
				if len(args) != 3 {
					return newError("wrong number of arguments. got=%d, want=3", len(args))
				}

				hash := args[0].(*object.Hash)
				newElements := make(map[object.HashKey]object.HashPair)

				for key, value := range hash.Pairs {
					newElements[key] = value
				}
				key := args[1].(object.Hashable)
				newElements[key.HashKey()] = object.HashPair{Key: args[1], Value: args[2]}

				return &object.Hash{Pairs: newElements}

			} else {
				return newError("argument to `add` must be ARRAY or MAP, got %s", args[0].Type())
			}
		},
	},
	"remove": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if args[0].Type() == object.ARRAY_OBJ {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				arr := args[0].(*object.Array)

				length := len(arr.Elements)
				if length == 0 {
					return newError("length of array must be greater than 0")
				}

				removeIndex := args[1].(*object.Integer)
				if removeIndex.Value < int64(0) || int64(length-1) < removeIndex.Value {
					return newError("index parameter must be between 0 and length of arr - 1")
				}

				newElements := []object.Object{}

				for i, elem := range arr.Elements {
					if int64(i) != removeIndex.Value {
						newElements = append(newElements, elem)

					}
				}

				return &object.Array{Elements: newElements}

			} else if args[0].Type() == object.HASH_OBJ {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				hash := args[0].(*object.Hash)

				length := len(hash.Pairs)
				if length == 0 {
					return newError("cannot remove from empty map")
				}

				removeKey := args[1].(object.Hashable)
				_, ok := hash.Pairs[removeKey.HashKey()]
				if ok {
					delete(hash.Pairs, removeKey.HashKey())
				} else {
					return newError("key not found in map")
				}

				return hash
			} else {
				return newError("argument to `add` must be ARRAY or MAP, got %s", args[0].Type())
			}

		},
	},
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
