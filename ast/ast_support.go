package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const comment_prefix = string("//")

type AnalysisResult struct {
	PkgName     string
	RecvMethods map[string][]MethodInfo // key RecvName
	Funcs       []FuncInfo
}

type MethodInfo struct {
	PkgName    string
	RecvName   string
	MethodName string
	Comment    []string
}

type FuncInfo struct {
	PkgName  string
	FuncName string
	Comment  []string
}

// print go file ast detail
func PrintAstInfo(fileName, src string, mode parser.Mode) {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, fileName, src, mode)
	if err != nil {
		panic(err)
	}
	ast.Print(fSet, f)
}

// find func and method in go file by target comment
// @see github.com\astaxie\beego\parser.go parserPkg
// @see github.com\astaxie\beego\parser.go parserComments
func ScanFuncDeclByComment(fileName, src, targetComment string) *AnalysisResult {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, fileName, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	result := &AnalysisResult{
		RecvMethods: make(map[string][]MethodInfo),
	}
	result.PkgName = f.Name.String()
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			if decl.Doc != nil {
				if isContainComment(decl.Doc.List, targetComment) {
					// bingo
					analysisComment(f.Name.String(), decl, result, targetComment)
				}
			}
		}
	}
	return result
}

func analysisComment(pkgName string, decl *ast.FuncDecl, result *AnalysisResult, comment string) {
	if decl.Recv != nil {
		// TODO extends and override need for range the Recv.List
		field := decl.Recv.List[0]
		switch f := field.Type.(type) {
		case *ast.StarExpr:
			// TODO convert ast.StarExpr.X
			recvName := fmt.Sprintf("%v", f.X)
			methodInfo := MethodInfo{
				PkgName:    pkgName,
				RecvName:   recvName,
				MethodName: decl.Name.String(),
				Comment:    []string{comment},
			}
			if result.RecvMethods[recvName] == nil {
				list := []MethodInfo{methodInfo}
				result.RecvMethods[recvName] = list
			} else {
				result.RecvMethods[recvName] = append(result.RecvMethods[recvName], methodInfo)
			}
		}
	} else {
		result.Funcs = append(result.Funcs, FuncInfo{
			PkgName:  pkgName,
			FuncName: decl.Name.String(),
			Comment:  []string{comment},
		})
	}
}

func isContainComment(lines []*ast.Comment, targetComment string) bool {
	for _, l := range lines {
		c := strings.TrimSpace(strings.TrimLeft(l.Text, comment_prefix))
		if c == targetComment {
			return true
		}
	}
	return false
}
