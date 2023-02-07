package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/reporter"
)

var (
	protoFiles         *string
	importPaths        *string
	errorUnusedImports = errors.New("found unused imports")
)

func init() {
	protoFiles = flag.String("proto-files", "", "A list of proto files, in a comma-separated string")
	importPaths = flag.String("import-paths", "", "A list of import paths for the protos to be passed to the compiler, in a comma-separated string.")
}

func parseAndFindUnusedProtos(protos []string) error {
	var impPaths []string
	if *importPaths != "" {
		impPaths = strings.Split(*importPaths, ",")
	}

	// unused imports are reported as warnings by protoc
	// so we need a custom reporter that collects the warnings
	// and processed
	type msg struct {
		pos  ast.SourcePos
		text string
	}

	var msgs []msg
	rep := func(warn reporter.ErrorWithPos) {
		msgs = append(msgs, msg{
			pos: warn.GetPosition(), text: warn.Unwrap().Error(),
		})
	}
	c := protocompile.Compiler{
		MaxParallelism: runtime.NumCPU(),
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: impPaths,
		}),
		Reporter: reporter.NewReporter(nil, rep),
	}

	ctx := context.Background()

	// Parser will use protoc to compile and discover any validation errors
	_, err := c.Compile(ctx, protos...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile protos: %v", err)
		return err
	}

	foundUnusedImports := false
	if len(msgs) > 0 {
		for _, msg := range msgs {
			if strings.HasPrefix(msg.text, "import") && strings.HasSuffix(msg.text, "not used") {
				fmt.Fprintf(os.Stderr, "Detected unused import in %s at line %d\n", msg.pos.Filename, msg.pos.Line)
				fmt.Fprintf(os.Stderr, msg.text+"\n")
				foundUnusedImports = true
			}
		}
	}

	// error out only if we found unused imports
	if foundUnusedImports {
		return errorUnusedImports
	}
	return nil
}

// Simple tool that lints protos for unused imports
func main() {
	flag.Parse()

	inBazelRun := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if inBazelRun != "" {
		os.Chdir(inBazelRun)
	}

	if *protoFiles == "" {
		fmt.Fprintf(os.Stderr, "One or more proto files must be passed!\n")
		os.Exit(1)
	}

	protos := strings.Split(*protoFiles, ",")

	if err := parseAndFindUnusedProtos(protos); err == errorUnusedImports {
		os.Exit(1)
	}

}
