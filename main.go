package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/swetjen/oto/parser"
	"github.com/swetjen/oto/render"
)

// Version is set during build.
var Version = "dev"

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(stdout io.Writer, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Println(args[0] + " " + Version + ` usage:
	oto [flags] paths [[path2] [path3]...]`)
		fmt.Println(`
flags:`)
		flags.PrintDefaults()
	}
	var (
		template      = flags.String("template", "", "plush template to render")
		outfile       = flags.String("out", "", "output file (default: stdout)")
		pkg           = flags.String("pkg", "", "explicit package name (default: inferred)")
		v             = flags.Bool("v", false, "verbose output")
		paramsStr     = flags.String("params", "", "list of parameters in the format: \"key:value,key:value\"")
		ignoreListIn  = flags.String("ignore", "", "comma separated list of interfaces to ignore")
		ignoreFieldIn = flags.String("fields", "", "comma separated list of fields to ignore")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if *template == "" {
		flags.PrintDefaults()
		return errors.New("missing template")
	}
	params, err := parseParams(*paramsStr)
	if err != nil {
		flags.PrintDefaults()
		return errors.Wrap(err, "params")
	}
	p := parser.New(flags.Args()...)
	ignoreItems := strings.Split(*ignoreListIn, ",")
	if ignoreItems[0] != "" {
		p.ExcludeInterfaces = ignoreItems
	}
	ignoreFields := strings.Split(*ignoreFieldIn, ",")
	if ignoreFields[0] != "" {
		p.ExcludeFields = ignoreFields
	}
	p.Verbose = *v
	if p.Verbose {
		fmt.Println("oto - github.com/pacedotdev/oto", Version)
	}
	def, err := p.Parse()
	if err != nil {
		return err
	}
	if *pkg != "" {
		def.PackageName = *pkg
	}
	b, err := ioutil.ReadFile(*template)
	if err != nil {
		return err
	}
	out, err := render.Render(string(b), def, params)
	if err != nil {
		return err
	}
	var w io.Writer = stdout
	if *outfile != "" {
		f, err := os.Create(*outfile)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	if _, err := io.WriteString(w, out); err != nil {
		return err
	}
	if p.Verbose {
		var methodsCount int
		for i := range def.Services {
			methodsCount += len(def.Services[i].Methods)
		}
		fmt.Println()
		fmt.Printf("\tTotal services: %d", len(def.Services))
		fmt.Printf("\tTotal Methods: %d", methodsCount)
		fmt.Printf("\tTotal Objects: %d\n", len(def.Objects))
		fmt.Printf("\tOutput size: %s\n", humanize.Bytes(uint64(len(out))))
	}
	return nil
}

// parseParams returns a map of data parsed from the params string.
func parseParams(s string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if s == "" {
		// empty map for an empty string
		return params, nil
	}
	pairs := strings.Split(s, ",")
	for i := range pairs {
		pair := strings.TrimSpace(pairs[i])
		segs := strings.Split(pair, ":")
		if len(segs) != 2 {
			return nil, errors.New("malformed params")
		}
		params[strings.TrimSpace(segs[0])] = strings.TrimSpace(segs[1])
	}
	return params, nil
}
