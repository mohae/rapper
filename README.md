# rapper
[![GoDoc](https://godoc.org/github.com/mohae/rapper?status.svg)](https://godoc.org/github.com/mohae/wrap)[![Build Status](https://travis-ci.org/mohae/rapper.png)](https://travis-ci.org/mohae/rapper)  
cli for rewriting files in a path as wrapped text. The path can either be a single file or a directory. For directories, `wrap` will wrap the lines of all files within that directory but will not descend into any sub-directories.

Lines lengths are a maximum of 80 characters unless the `length` flag is passed with a non-zero, positive, integer.

Files can be filtered by file extension. The `ext` flag specifies the file extension. For multiple file extensions the `ext` flag should be repeated for each file extension. What should be done with the filtered file is determined by the `exclude` or `include` flag. If both `exclude` and `include` are _false_, nothing will be filtered; both cannot be true.

## Use
This assumes that you have Go installed and that you're `$GOPATH/bin` is in your path.

Get and compile `rapper`:  

    $ go install github.com/mohae/rapper

Run rapper with `verbose` output:

    $ rapper -v path/to/dir

Run rapper; only wrap `.txt` and `.md` files. The extension does not have to include the `.`:

    $ rapper -include -ext .txt -ext md path/to/dir

For usage output:

    $ rapper -h

## TODO
Automatically skip some common non-text extensions, e.g. .jpg, .jpeg, .gif, etc. This may be a box not worth opening and, instead, rely on the user to use the `exclude` and `include` flags properly.
