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

func (e *Environment) GetCopyOfEnvWithEmptyOuter() *Environment {
	newEnv := &Environment{}

	newStore := map[string]Object{}
	newOuter := Environment{}

	for key, value := range e.store {
		newStore[key] = value
	}

	newEnv.store = newStore
	newEnv.outer = &newOuter

	return newEnv
}
