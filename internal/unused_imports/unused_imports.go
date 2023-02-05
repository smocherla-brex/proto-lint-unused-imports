package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/reporter"
)

const (
	protoExtension = ".proto"
)

var (
	protosDir *string
	importPaths *string
	errorUnusedImports = errors.New("found unused imports")
)

func init () {
	protosDir = flag.String("protos-dir", "", "directory containing proto files")
	importPaths = flag.String("import-paths", "", "A list of import paths for the protos to be passed to the compiler, in a comma-separated string.")
}


func listProtos() ([]string , error) {
	if _, err := os.Stat(*protosDir); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("Directory %v does not exist", *protosDir))
	}

	files, err := ioutil.ReadDir(*protosDir)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Directory %v does not have any files!", *protosDir))
	}

	var protoFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".proto" && !file.IsDir() {
			protoFiles = append(protoFiles, path.Join(*protosDir, file.Name()))
		}
	}

	fmt.Fprintf(os.Stdout, "Found %d protos\n", len(protoFiles))
	return protoFiles, nil
	
}

func parseAndFindUnusedProtos(protoFiles []string) error {
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
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: impPaths,
		}),
		Reporter: reporter.NewReporter(nil, rep),
	}

	ctx := context.Background()

	// Parser will use protoc to compile and discover any validation errors
	_, err := c.Compile(ctx, protoFiles...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile protos: %v", err)
		return err
	}
	
	foundUnusedImports := false
	if len(msgs) > 0 {
		for _, msg := range msgs {
			if strings.HasPrefix(msg.text, "import") && strings.HasSuffix(msg.text, "not used") {
				fmt.Fprintf(os.Stderr, "Detected unused import in %s at line %d\n", msg.pos.Filename, msg.pos.Line)
				fmt.Fprintf(os.Stderr, msg.text + "\n")
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

	if *protosDir == "" {
		fmt.Fprintf(os.Stderr, "A directory must be passed!\n")
		os.Exit(1)
	}


	protoFiles, err := listProtos()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred listing proto files: %v\n", err)
		os.Exit(1)
	}

	if len(protoFiles) == 0 {
		fmt.Fprintf(os.Stderr, "No proto files found\n")
		os.Exit(1)
	}

	if err = parseAndFindUnusedProtos(protoFiles); err == errorUnusedImports {
		os.Exit(1)
	}

}