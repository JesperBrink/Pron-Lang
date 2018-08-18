package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) GetOuterMost(name string) (Object, bool) {
	if e.outer != nil {
		return e.outer.GetOuterMost(name)
	}

	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Update(name string, val Object) bool {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		ok = e.outer.Update(name, val)
	} else {
		e.Set(name, val)
	}
	return ok
}

func (e *Environment) UpdateOuterMost(name string, val Object) bool {
	if e.outer != nil {
		return e.outer.UpdateOuterMost(name, val)
	}

	_, ok := e.store[name]
	if ok {
		e.Set(name, val)
	}

	return ok
}

func (e *Environment) GetCopyOfEnvWithOuterEnvNil() *Environment {
	newEnv := &Environment{}

	newStore := map[string]Object{}

	for key, value := range e.store {
		newStore[key] = value
	}

	newEnv.store = newStore

	return newEnv
}
