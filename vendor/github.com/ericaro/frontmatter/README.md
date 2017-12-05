[![Build Status](https://travis-ci.org/ericaro/frontmatter.png?branch=master)](https://travis-ci.org/ericaro/ringbuffer) [![GoDoc](https://godoc.org/github.com/ericaro/frontmatter?status.svg)](https://godoc.org/github.com/ericaro/frontmatter)

# Frontmatter 

frontmatter is a golang package that provides a Marshaler/Unmarshaler for frontmatter files.

 A frontmatter file is a file with any textual content, and a yaml frontmatter block for metadata,
 as defined by [Jekyll](http://jekyllrb.com/docs/frontmatter/)

 The frontmatter File must have the following format:

      `---\n`
      <yaml content>
      `\n---\n`
      <text content>

 Where, the 'yaml content' is handled by http://gopkg.in/yaml.v2

 And where, the text content is 'content' field's value.

 The 'content' field must:

      exist
      tagged `fm:"content"`
      be exported
      be of the correct type.

 A correct type is:

     - string, *string
     - any convertible to the above two.

 see [go doc](https://godoc.org/github.com/ericaro/frontmatter) for details.

# Installation

first get [go](http://golang.org)
then `go get github.com/ericaro/frontmatter`

This package depends on http://gopkg.in/yaml.v2

# License

frontmatter is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).


# Branches


master: [![Build Status](https://travis-ci.org/ericaro/frontmatter.png?branch=master)](https://travis-ci.org/ericaro/frontmatter) against go versions:

  - 1.2
  - 1.3
  - 1.4
  - tip
