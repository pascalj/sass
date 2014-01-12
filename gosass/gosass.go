package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "flag"
  "errors"
  "github.com/pascalj/sass"
)

var stdin = flag.Bool("s", false, "Read from stdin.")
var outputType = flag.String("t","nested", "Output style. Can be nested, compact, compressed, or expanded.")
var inPath = flag.String("i","", "Input. May be a path to a file or directory.")
var outPath = flag.String("o","", "Output. May be a path to a file or directory.")
var sourceComments = flag.Bool("c",false, "Add line number comments.")

func main() {
  if len(os.Args) < 2 {
    usage()
    return
  }

  options := sass.NewOptions()
  parseOptions(options)

  inStat, inErr := os.Stat(*inPath)
  outStat, outErr := os.Stat(*inPath)
  output := ""
  var compileError error;

  switch {
    case *stdin:
      inputBytes, _ := ioutil.ReadAll(os.Stdin)
      output, compileError = sass.Compile(string(inputBytes), options)
    case inErr == nil && inStat.Mode() & os.ModeType == 0:
      output, compileError = sass.CompileFile(*inPath, options)
    case inErr == nil && outErr == nil && inStat.IsDir() && outStat.IsDir():
      compileError = sass.CompileFolder(*inPath, *outPath, options)
    default:
      usage()
      return
  }
  if outErr != nil || outStat.Mode() & os.ModeType == 0 {
    ioutil.WriteFile(*outPath, []byte(output), 0644)
  }


  if compileError != nil {
    fmt.Println("Compiler error:", compileError)
    return
  }

  if *outPath == "" {
    fmt.Println(output)
  }
}

func usage() {
  fmt.Println("Usage:", os.Args[0], "[options]\r\n")
  fmt.Println("Description:")
  fmt.Println("  Converts SCSS to CSS.\r\n")
  fmt.Println("Options:\r\n")
  flag.PrintDefaults()
}

func parseOptions(options *sass.SassOptions) error {
  flag.Parse()
  switch *outputType {
    case "nested":
      options.OutputStyle = sass.SASS_STYLE_NESTED
    case "compact":
      options.OutputStyle = sass.SASS_STYLE_COMPACT
    case "compressed":
      options.OutputStyle = sass.SASS_STYLE_COMPRESSED
    case "expanded":
      options.OutputStyle = sass.SASS_STYLE_NESTED
    default:
      usage()
      return errors.New("Invalid type argument.")
  }
  options.SourceComments = *sourceComments
  return nil
}