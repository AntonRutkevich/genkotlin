package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/AntonRutkevich/astparser"
	"github.com/mkorolyov/genkotlin/kotlin"
)

var (
	inputDir         = flag.String("in", "", "directory with go files to be parsed")
	excludeRegexpStr = flag.String("e", "", "exclude regexp to skip files")
	includeRegexpStr = flag.String("i", "", "include regexp to limit input files")
	outputDir        = flag.String("o", "", "output directory for generated files, without package")
	outPackage       = flag.String("package", "", "package for generated kt files")
	useJsonNames     = flag.Bool("json_names", false, "use json tag for field names, eg com.myapp")
)

func main() {
	flag.Parse()

	// load golang sources
	cfg := astparser.Config{InputDir: *inputDir}
	if *excludeRegexpStr != "" {
		cfg.ExcludeRegexp = *excludeRegexpStr
	}
	if *includeRegexpStr != "" {
		cfg.IncludeRegexp = *includeRegexpStr
	}
	sources, err := astparser.Load(cfg)
	if err != nil {
		log.Fatalf("failed to load sources from %s excluding %s: %v", *inputDir, *excludeRegexpStr, err)
	}

	config := kotlin.Config{UseJsonTagNames: *useJsonNames}
	if *outPackage != "" {
		config.OutPackage = *outPackage
	}

	generator := kotlin.NewGenerator(config)
	// generate kotlin classes
	if *outputDir != "" {
		kotlinFiles := generator.Generate(sources)
		// save
		for f, body := range kotlinFiles {
			outPackagePath := ""
			if *outPackage != "" {
				outPackagePath = strings.Join(strings.Split(*outPackage, "."), "/")
			}

			// 0666 is enough, but https://github.com/golang/go/issues/22323
			packagePath := *outputDir + "/" + outPackagePath
			if err := os.MkdirAll(packagePath, 0777); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create directory %s: %v\n", packagePath, err)
			}

			filePath := packagePath + "/" + f + ".kt"

			if err := ioutil.WriteFile(filePath, body, 0777); err != nil {
				fmt.Fprintf(os.Stderr, "failed to save generated kotlin file %s: %v\n", filePath, err)
			}
		}
	}
}
