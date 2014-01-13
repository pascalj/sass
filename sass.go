package sass

// #cgo LDFLAGS: -lsass
// #include <stdlib.h>
// #include <sass_interface.h>
import "C"
import (
	"errors"
	"unsafe"
)

type SassOuputStyle int

const (
	SASS_STYLE_NESTED     SassOuputStyle = iota
	SASS_STYLE_EXPANDED                  = iota
	SASS_STYLE_COMPACT                   = iota
	SASS_STYLE_COMPRESSED                = iota
)

type SassOptions struct {
	OutputStyle    SassOuputStyle
	SourceComments bool
	ImagePath      string
	IncludePaths   string
}

func NewOptions() *SassOptions {
	return &SassOptions{
		OutputStyle:    SASS_STYLE_NESTED,
		SourceComments: false,
		ImagePath:      "",
		IncludePaths:   "",
	}
}

func Compile(input string, options *SassOptions) (compiled string, err error) {
	ctx := newSassCtx()
	ctx.source_string = C.CString(input)
	ctx.options = newCSassOptions(options)
	defer C.free(unsafe.Pointer(ctx.source_string))

	C.sass_compile(ctx)

	if ctx.error_status > 0 {
		return "", errors.New(C.GoString(ctx.error_message))
	}
	return C.GoString(ctx.output_string), nil
}

func CompileFile(infilePath string, options *SassOptions) (compiled string, err error) {
	ctx := newSassFileCtx()
	ctx.input_path = C.CString(infilePath)
	ctx.options = newCSassOptions(options)
	defer C.free(unsafe.Pointer(ctx.input_path))

	C.sass_compile_file(ctx)

	if ctx.error_status > 0 {
		return "", errors.New(C.GoString(ctx.error_message))
	}
	return C.GoString(ctx.output_string), nil
}

func CompileFolder(searchPath string, outputPath string, options *SassOptions) error {
	ctx := newSassFolderCtx()
	ctx.search_path = C.CString(searchPath)
	ctx.output_path = C.CString(outputPath)
	ctx.options = newCSassOptions(options)
	defer C.free(unsafe.Pointer(ctx.search_path))
	defer C.free(unsafe.Pointer(ctx.output_path))

	C.sass_compile_folder(ctx)

	if ctx.error_status > 0 {
		return errors.New(C.GoString(ctx.error_message))
	}

	return nil
}

// private

func newSassCtx() *C.struct_sass_context {
	return C.sass_new_context()
}

func newSassFileCtx() *C.struct_sass_file_context {
	return C.sass_new_file_context()
}

func newSassFolderCtx() *C.struct_sass_folder_context {
	return C.sass_new_folder_context()
}

func newCSassOptions(options *SassOptions) C.struct_sass_options {
	var sassOptions C.struct_sass_options
	sassOptions.output_style = C.int(options.OutputStyle)
	if options.SourceComments {
		sassOptions.source_comments = 1
	}
	sassOptions.image_path = C.CString(options.ImagePath)
	sassOptions.include_paths = C.CString(options.IncludePaths)
	return sassOptions
}
