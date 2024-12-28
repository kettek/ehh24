package game

import (
	"fmt"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// InterpExports are our yaegi exports to expose to interpreted code.
var InterpExports = make(interp.Exports)

func init() {
	InterpExports["game/game"] = map[string]reflect.Value{
		"Place": reflect.ValueOf((*Place)(nil)),
	}
}

func setupInterp(i *interp.Interpreter, src string) error {
	i.Use(stdlib.Symbols)

	if err := i.Use(InterpExports); err != nil {
		return err
	}

	_, err := i.Eval(fmt.Sprintf(`
		import (
			"game"
		)
		
		%s
	`, src))
	if err != nil {
		return err
	}
	return nil
}
