// Copyright 2025 samber.
//
// Licensed as an Enterprise License (the "License"); you may not use
// this file except in compliance with the License. You may obtain
// a copy of the License at:
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.ee.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package introspection

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"strings"

	"github.com/samber/lo"
)

const unsupportedExpr = "?"

type FunDesc struct {
	Name      string
	Pos       string
	Arguments []FunDescArgument
}

type FunDescArgument struct {
	Name string
	Pos  string
}

func (f *FunDescArgument) String() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Pos)
}

func GetFunctionDescription(skipCaller int, skipArgs int) (*FunDesc, error) {
	// Get caller information
	_, file, line, ok := runtime.Caller(2 + skipCaller)
	if !ok {
		return nil, fmt.Errorf("failed to collect caller description")
	}

	// Parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Prepare output
	output := FunDesc{
		Name:      "",
		Pos:       fmt.Sprintf("%s:%d", file, line),
		Arguments: make([]FunDescArgument, 0),
	}

	// Find the specific call in the AST
	ast.Inspect(node, func(n ast.Node) bool {
		if ce, ok := n.(*ast.CallExpr); ok {
			if fset.Position(ce.Pos()).Line == line {
				// Extract function name
				output.Name = getArgName(ce.Fun)
				output.Pos = fset.Position(ce.Pos()).String()

				// Extract arguments
				args := lo.Slice(ce.Args, skipArgs, len(ce.Args))
				output.Arguments = make([]FunDescArgument, len(args))
				for i, arg := range args {
					name := getArgName(arg)

					output.Arguments[i] = FunDescArgument{
						Name: name,
						Pos:  fset.Position(arg.Pos()).String(),
					}

				}

				return false
			}
		}
		return true
	})

	return &output, nil
}

func getArgName(arg ast.Expr) string {
	switch argType := arg.(type) {

	///////////////// Terminal args
	case *ast.Ident:
		// eg:
		// PipeX(
		//   source,
		//   myOperator,
		//   myPkg.myOperators.map,
		// }
		return argType.Name
	case *ast.BasicLit:
		// eg:
		// PipeX(
		//   source,
		//   myOperators["map"],
		//   Map[int],
		// }
		return argType.Value
	case *ast.FuncLit:
		// eg:
		// PipeX(
		//   source,
		//   func() { ... },
		// }
		return "func(...) { ... }"

	///////////////// Recursive args
	case *ast.ParenExpr:
		// eg:
		// PipeX(
		//   source,
		//   (myOperator),
		// }
		return getArgName(argType.X)
	case *ast.SelectorExpr:
		// eg:
		// PipeX(
		//   source,
		//   myPkg.Mappers,
		//   myPkg.operators.map,
		// }
		return fmt.Sprintf("%s.%s", getArgName(argType.X), getArgName(argType.Sel))
	case *ast.IndexListExpr:
		// eg: PipeX[A any, B any]()
		indices := lo.Map(argType.Indices, func(index ast.Expr, i int) string {
			return getArgName(index)
		})
		return fmt.Sprintf("%s[%s]", getArgName(argType.X), strings.Join(indices, ", "))
	case *ast.IndexExpr:
		// eg:
		// PipeX(
		//   source,
		//   myPkg.Mappers["foo"],
		//   myPkg.operators()["map"].func,
		// }
		return fmt.Sprintf("%s[%s]", getArgName(argType.X), getArgName(argType.Index))
	case *ast.CallExpr:
		// eg:
		// PipeX(
		//   source,
		//   Map(...),
		//   pkg.Map(...),
		//   pkg.operators.Map(...),
		//   pkg.operators["map"].fn(...),
		// }
		return fmt.Sprintf("%s(...)", getArgName(argType.Fun))

		///////////////// Fallback: please open an issue if a case is not properly supported
	default:
		return unsupportedExpr
	}
}
