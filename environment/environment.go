package environment

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tanlian/rulego/object"
)

type Environment struct {
	data   map[string]object.Object
	parent *Environment
}

var Root = &Environment{
	data: make(map[string]object.Object),
}

func New(parent *Environment) *Environment {
	return &Environment{
		data:   make(map[string]object.Object),
		parent: parent,
	}
}

func (env *Environment) Get(key string) (object.Object, bool) {
	if v, ok := env.data[key]; ok {
		return v, true
	}
	if env.parent != nil {
		return env.parent.Get(key)
	}
	return nil, false
}

func (env *Environment) GetKeyEnv(key string) *Environment {
	if _, ok := env.data[key]; ok {
		return env
	}
	if env.parent != nil {
		return env.parent.GetKeyEnv(key)
	}
	return nil
}

func (env *Environment) Set(key string, value object.Object) {
	if parent := env.GetKeyEnv(key); parent != nil {
		parent.data[key] = value
		return
	}
	env.data[key] = value
}

func (env *Environment) SetCurrent(key string, value object.Object) {
	env.data[key] = value
}

func (env *Environment) Inject(key string, value any) {
	env.SetCurrent(key, object.ToObject(reflect.ValueOf(value)))
}

func (env *Environment) Remove(key string) {
	delete(env.data, key)
}

func (env *Environment) String() string {
	var s strings.Builder
	t := env
	for t != nil {
		s.WriteString(fmt.Sprintf("%s => ", t.excludeFunc()))
		t = t.parent
	}
	s.WriteString("None")
	return s.String()
}

func (env *Environment) excludeFunc() string {
	d := make(map[string]any)
	for k, v := range env.data {
		d[k] = v.GetValue()
	}
	v, _ := json.Marshal(d)
	return string(v)
}
