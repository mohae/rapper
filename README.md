# rapper
[![GoDoc](https://godoc.org/github.com/mohae/rapper?status.svg)](https://godoc.org/github.com/mohae/wrap)[![Build Status](https://travis-ci.org/mohae/rapper.png)](https://travis-ci.org/mohae/rapper)  
Rapper is a cli for rewriting files in a path, or set of paths, as wrapped text. The path can either be a single file or a directory or any combination thereof. For directories, `wrap` will wrap the lines of all files within that directory but will not descend into any sub-directories.

Lines lengths are a maximum of 80 characters unless the `length` flag is passed with a non-zero, positive, integer.

Files can be filtered by file extension. The `ext` flag specifies the file extension. For multiple file extensions the `ext` flag should be repeated for each file extension. What should be done with the filtered file is determined by the `exclude` or `include` flag. If both `exclude` and `include` are _false_, nothing will be filtered; both cannot be true.

An extension is the suffix of the last '.' in the filename. If the filename does not have a dot, like 'README', it will not have a extension.

Files can also be wrapped as comments. C, C++, and Shell style comments are supported and specified with the `comment` flag.

Rapper will attempt to wrap all files in the target(s), including non-text files. Use the `ext` and either the `include` or `exclude` flags to filter files.

## Use
This assumes that you have Go installed and that you're `$GOPATH/bin` is in your path.

Get and compile `rapper`:  

    $ go install github.com/mohae/rapper

    Run rapper with `verbose` output:

        $ rapper -v path/to/dir

    Run rapper against multiple paths:

        $ rapper path/to/dir path/to/another/dir

Run rapper, only wrapping `.txt` and `.stuff` files. The extension does not have to include the `.`:

    $ rapper -include -ext .txt -ext stuff path/to/dir

For usage output:

    $ rapper -h
