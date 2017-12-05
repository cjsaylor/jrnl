//package frontmatter provide a Marshaler/Unmarshaler for frontmatter files.
//
// A frontmatter file is a file with any textual content, and a yaml frontmatter block for metadata,
// as defined by Jekyll http://jekyllrb.com/docs/frontmatter/
//
// The frontmatter File must have the following format:
//
//      `---\n`
//      <yaml content>
//      `\n---\n`
//      <text content>
//
// Where, the 'yaml content' is handled by http://gopkg.in/yaml.v2
//
// And where, the text content is 'content' field's value.
//
// The 'content' field must:
//
//      exist
//      tagged `fm:"content"`
//      be exported
//      be of the correct type.
//
// A correct type is:
//
//     - string, *string
//     - any convertible to the above two.
//
// See example for details.
//
package frontmatter

import (
	"errors"
	yaml "gopkg.in/yaml.v2"
	"strings"
)

const (
	Tag       = "fm"      //Tag used to find the content field
	Content   = "content" // value used to identify the content field
	Header    = "---\n"   // front matter file header
	Separator = "\n---\n" // front matter metadata/content separator
)

var (
	ErrMissingSeparator      = errors.New("found a heading '---' without separator '---'")
	ErrUnexported            = errors.New("cannot set an unexported content field")
	ErrNoContentField        = errors.New("missing content field")
	ErrWrongContentFieldType = errors.New("No content field with the right type")
)

//Marshal any object.
func Marshal(v interface{}) ([]byte, error) {
	//first delegate the metadata marshalling, and errors are reported
	yamlBytes, yerr := yaml.Marshal(v)
	if yerr != nil {
		return nil, yerr
	}
	// the do the frontmatter reading
	content, cerr := ReadString(v, Tag, Content)
	if cerr != nil {
		return nil, cerr
	}
	// we trim space (useless in yaml) to have a more stable output format
	yamlContent := strings.TrimSpace(string(yamlBytes))
	//the result is just a simple join
	result := strings.Join([]string{Header, yamlContent, Separator, content}, "")
	return ([]byte)(result), nil
}

//Unmarshal any object from 'data'.
func Unmarshal(data []byte, v interface{}) (err error) {
	//always working as string
	txt := string(data)

	content := txt //by default

	if strings.HasPrefix(txt, Header) { // there is a header, therefore there MUST be a front matter

		//we remove the prefix
		txt = strings.TrimPrefix(txt, Header)
		// nice trick: we split using the separator
		// and we hope the get: metadata (valid yaml) and the content
		// all the rest is check
		splitted := strings.SplitN(txt, Separator, 2)

		if len(splitted) != 2 {
			return ErrMissingSeparator
		}

		metadata := splitted[0]
		// now this is supposed to be a valid yaml.
		yamlerr := yaml.Unmarshal(([]byte)(metadata), v)
		if yamlerr != nil {
			//there have been a yaml error, we report it, but we don't fail
			err = yamlerr
		} else {
			// on success we only keep the content
			// because the metadata content is assumed to have been read (and therefore not lost)
			// on error, the metadata content would be lost
			content = splitted[1]
		}
	}
	contenterr := WriteString(v, Tag, Content, content)
	if contenterr != nil {
		err = contenterr
	}
	return
}
