package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

/*
* main is the entry point for the types generator
 */
func main() {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes,
	}

	pkgs, err := packages.Load(cfg, "github.com/souravsspace/texly.chat/internal/models")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	fmt.Println("/*")
	fmt.Println(" * This file is auto-generated. Do not edit directly.")
	fmt.Println(" */")
	fmt.Println("")

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					/*
					 * Only process exported structs
					 */
					if !typeSpec.Name.IsExported() {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					/*
					 * Extract comments from GenDecl (usually where they are) or TypeSpec
					 */
					doc := genDecl.Doc
					if doc == nil {
						doc = typeSpec.Doc
					}

					generateInterface(typeSpec.Name.Name, structType, doc)
					fmt.Println()
				}
			}
		}
	}
}

/*
* generateInterface outputs the TypeScript interface definition with comments
 */
func generateInterface(name string, structType *ast.StructType, doc *ast.CommentGroup) {
	if doc != nil {
		fmt.Println("/*")
		for _, comment := range doc.List {
			/*
			 * Clean up comment markers to generate clean JSdoc-style comments
			 */
			text := strings.TrimSpace(comment.Text)
			text = strings.TrimPrefix(text, "//")
			text = strings.TrimPrefix(text, "/*")
			text = strings.TrimSuffix(text, "*/")
			text = strings.TrimSpace(text)
			
			// Remove leading asterisks if present in multi-line comments
			text = strings.TrimPrefix(text, "*")
			text = strings.TrimSpace(text)

			if text != "" {
				fmt.Printf(" * %s\n", text)
			}
		}
		fmt.Println(" */")
	}

	fmt.Printf("export interface %s {\n", name)

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name
		jsonTag := ""
		if field.Tag != nil {
			tagValue := strings.Trim(field.Tag.Value, "`")
			structTag := reflect.StructTag(tagValue)
			jsonTag = structTag.Get("json")
		}

		if jsonTag == "-" {
			continue
		}

		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		tsType := toTypeScriptType(field.Type)
		fmt.Printf("  %s: %s;\n", fieldName, tsType)
	}

	fmt.Println("}")
}

/*
* toTypeScriptType converts Go AST types to TypeScript types
 */
func toTypeScriptType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string":
			return "string"
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64":
			return "number"
		case "bool":
			return "boolean"
		case "Time":
			return "string"
		default:
			return t.Name
		}
	case *ast.SelectorExpr:
		if t.Sel.Name == "Time" {
			return "string"
		}
		return "any"
	case *ast.ArrayType:
		return toTypeScriptType(t.Elt) + "[]"
	case *ast.StarExpr:
		return toTypeScriptType(t.X)
	default:
		return "any"
	}
}
