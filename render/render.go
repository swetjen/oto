package render

import (
	"bytes"
	"encoding/json"
	"go/doc"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
	"github.com/pacedotdev/oto/parser"
	"github.com/pkg/errors"
)

// Render renders the template using the Definition.
func Render(tpl string, def parser.Definition, params map[string]interface{}) (string, error) {
	if err := validateTypeMappings(tpl, def); err != nil {
		return "", err
	}
	funcs := template.FuncMap{
		"camelize_down":       camelizeDown,
		"camelize_up":         camelizeUp,
		"camelize_up_field":   camelizeUpField,
		"json":                toJSONHelper,
		"json_inline":         toJSONInlineHelper,
		"format_comment_line": formatCommentLine,
		"format_comment_text": formatCommentText,
		"format_comment_html": formatCommentHTML,
		"format_tags":         formatTags,
		"object_golang":       ObjectGolang,
		"smart_prefix":        smartPrefix,
	}
	tmpl, err := template.New("oto").Funcs(funcs).Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	data := map[string]interface{}{
		"def":    def,
		"params": params,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toJSONHelper(v interface{}, prefix, indent string) (string, error) {
	if indent == "" {
		indent = "\t"
	}
	b, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toJSONInlineHelper(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func formatCommentLine(s string) string {
	var buf bytes.Buffer
	doc.ToText(&buf, s, "", "", 2000)
	s = strings.TrimSpace(buf.String())
	return s
}

func formatCommentText(s string) string {
	var buf bytes.Buffer
	doc.ToText(&buf, s, "// ", "", 80)
	return buf.String()
}

func formatCommentHTML(s string) string {
	var buf bytes.Buffer
	doc.ToHTML(&buf, s, nil)
	return buf.String()
}

// formatTags formats a list of struct tag strings into one.
// Will return an error if any of the tag strings are invalid.
func formatTags(tags ...string) (string, error) {
	alltags := &structtag.Tags{}
	for _, tag := range tags {
		theseTags, err := structtag.Parse(tag)
		if err != nil {
			return "", errors.Wrapf(err, "parse tags: `%s`", tag)
		}
		for _, t := range theseTags.Tags() {
			alltags.Set(t)
		}
	}
	tagsStr := alltags.String()
	if tagsStr == "" {
		return "", nil
	}
	tagsStr = "`" + tagsStr + "`"
	return tagsStr, nil
}

// smartPrefix prepends a string before s, allowing for the specific use
// case of pointers to objects. If the s begins with * (as in, *Object), the
// result will be *prefixObject to preserve its original meaning.
func smartPrefix(prefix, s string) string {
	if strings.HasPrefix(s, "*") {
		return "*" + prefix + s[1:]
	}
	return prefix + s
}
