package object

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store, nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	store := make(map[string]Object)
	return &Environment{store, outer}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(key string) (Object, bool) {
	val, ok := e.store[key]
	if !ok && e.outer != nil {
		return e.outer.Get(key)
	}
	return val, ok
}

func (e *Environment) Set(key string, val Object) Object {
	e.store[key] = val
	return val
}
