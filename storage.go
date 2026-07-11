package ormc

import (
	"go/ast"
	"strings"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/model"
)

var builtinKinds = map[string]model.FieldType{
	"Text":        model.FieldText,
	"Int":         model.FieldInt,
	"Float":       model.FieldFloat,
	"Bool":        model.FieldBool,
	"Blob":        model.FieldBlob,
	"Raw":         model.FieldRaw,
	"Struct":      model.FieldStruct,
	"IntSlice":    model.FieldIntSlice,
	"StructSlice": model.FieldStructSlice,
}

func (g *Generator) resolveStorage(infos []StructInfo, file *ast.File) error {
	importMap := make(map[string]string) // alias -> path
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			alias = parts[len(parts)-1]
		}
		importMap[alias] = path
	}

	var toProbe []fieldProbe
	for i := range infos {
		for j := range infos[i].Fields {
			fi := &infos[i].Fields[j]
			if fi.KindConstructor == "" {
				continue
			}

			// Parse the constructor expression
			selector, constructor := parseConstructor(fi.KindConstructor)
			path := importMap[selector]
			if path == "" {
				// Try common aliases if not found in imports
				if selector == "model" { path = "github.com/tinywasm/model" }
			}

			if path == "github.com/tinywasm/model" || (path == "" && selector == "model") {
				if t, ok := builtinKinds[constructor]; ok {
					fi.Type = t
					continue
				}
			}

			if selector == "" && constructor != "" {
				// Local kind? Plan says custom kinds live in their own package.
				return fmt.Errf("field %s (struct %s): custom kinds must live in their own package (found local kind %s)", fi.ColumnName, infos[i].Name, constructor)
			}

			if path == "" {
				return fmt.Errf("field %s (struct %s): unknown package alias %s in kind %s", fi.ColumnName, infos[i].Name, selector, fi.KindConstructor)
			}

			toProbe = append(toProbe, fieldProbe{
				infoIdx:  i,
				fieldIdx: j,
				expr:     fi.KindConstructor,
				pkgPath:  path,
			})
		}
	}

	if len(toProbe) > 0 {
		results, err := g.runProbe(toProbe, infos)
		if err != nil {
			return err
		}
		for _, res := range results {
			infos[res.infoIdx].Fields[res.fieldIdx].Type = res.storage
		}
	}

	// Final validation and GoType update
	for i := range infos {
		for j := range infos[i].Fields {
			fi := &infos[i].Fields[j]
			if fi.Type == model.FieldStruct || fi.Type == model.FieldStructSlice {
				if fi.Ref == "" {
					return fmt.Errf("field %s (struct %s): composition kind %s requires a non-nil Definition argument", fi.ColumnName, infos[i].Name, fi.KindConstructor)
				}
			} else if fi.Type == 0 {
				// This should have been caught by probe failure, but just in case
				return fmt.Errf("field %s (struct %s): failed to resolve storage for kind %s", fi.ColumnName, infos[i].Name, fi.KindConstructor)
			}
			fi.GoType = FieldTypeToGoType(fi.Type, fi.Ref)
		}
	}

	return nil
}

type fieldProbe struct {
	infoIdx  int
	fieldIdx int
	expr     string
	pkgPath  string
}

type probeResult struct {
	infoIdx int
	fieldIdx int
	storage model.FieldType
}

func parseConstructor(expr string) (selector, constructor string) {
	// Simple parser for "pkg.Kind()" or "Kind()"
	if idx := strings.Index(expr, "("); idx != -1 {
		expr = expr[:idx]
	}
	if idx := strings.LastIndex(expr, "."); idx != -1 {
		return expr[:idx], expr[idx+1:]
	}
	return "", expr
}
