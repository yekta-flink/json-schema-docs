package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/xeipuuv/gojsonpointer"
)

// resolveSchema recursively resolves schemas.
func resolveSchema(schem *schema, dir string, root *simplejson.Json) (*schema, error) {
	for k, prop := range schem.Properties {
		if prop.Ref != "" {
			tmp, err := resolveReference(prop.Ref, dir, root)
			if err != nil {
				return nil, err
			}
			*prop = *tmp
		}
		foo, err := resolveSchema(prop, dir, root)
		if err != nil {
			return nil, err
		}
		*prop = *foo
		prop.Key = k

		prop.Parent = schem
		if (schem.Key == "then") || (schem.Key == "else") {
			prop.Parent = schem.Parent
		}

	}

	// Handle conditionals.
	conditionals := []struct {
		Name      string
		Attribute *schema
	}{
		{
			Name:      "Then",
			Attribute: schem.Then,
		},
		{
			Name:      "Else",
			Attribute: schem.Else,
		},
	}

	for _, cnd := range conditionals {
		if cnd.Attribute != nil {
			// Resolve reference if there is any and replace schema.
			if cnd.Attribute.Ref != "" {
				tmp, err := resolveReference(cnd.Attribute.Ref, dir, root)
				if err != nil {
					return nil, err
				}
				*cnd.Attribute = *tmp
			}

			cnd.Attribute.Parent = schem
			cnd.Attribute.Key = strings.ToLower(cnd.Name)

			foo, err := resolveSchema(cnd.Attribute, dir, root)
			if err != nil {
				return nil, err
			}
			*cnd.Attribute = *foo
		}
	}

	if schem.Items != nil {
		if schem.Items.Ref != "" {
			tmp, err := resolveReference(schem.Items.Ref, dir, root)
			if err != nil {
				return nil, err
			}
			*schem.Items = *tmp
		}
		foo, err := resolveSchema(schem.Items, dir, root)
		if err != nil {
			return nil, err
		}
		*schem.Items = *foo
	}

	return schem, nil
}

// resolveReference loads a schema from a $ref.
//
// If ref contains a hashtag (#), the part before represents a cross-schema
// reference, and the part after represents a in-schema reference.
//
// If ref is missing a hashtag, the whole schema is being referenced.
func resolveReference(ref string, dir string, root *simplejson.Json) (*schema, error) {
	i := strings.Index(ref, "#")

	if i < 0 {
		// cross-schema reference to another schema
		schema, err := loadSchemaFromPath(filepath.Join(dir, ref))
		if err != nil {
			return nil, err
		}
		return schema, nil
	} else if i == 0 {
		// in-schema reference
		return resolveInSchemaReference(ref[i+1:], root)
	} else {

		schema, err := loadSchemaFromPath(filepath.Join(dir, ref[:i]))
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(schema)
		if err != nil {
			return nil, err
		}

		newRoot, err := simplejson.NewJson(b)
		if err != nil {
			return nil, err
		}
		return resolveInSchemaReference(ref[i+1:], newRoot)
	}

}

func resolveInSchemaReference(path string, root *simplejson.Json) (*schema, error) {
	// in-schema reference
	pointer, err := gojsonpointer.NewJsonPointer(path)
	if err != nil {
		return nil, err
	}

	v, _, err := pointer.Get(root.MustMap())
	if err != nil {
		return nil, err
	}

	var sch schema
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &sch); err != nil {
		return nil, err
	}

	return &sch, nil
}

// loadSchemaFromPath returns a schema at a given path.
func loadSchemaFromPath(path string) (*schema, error) {
	rc, err := openFileOrURL(path)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	return newSchema(rc, filepath.Dir(path))
}

// openFileOrURL opens a file from a URL or local path.
func openFileOrURL(path string) (io.ReadCloser, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}
	return os.Open(path)
}
