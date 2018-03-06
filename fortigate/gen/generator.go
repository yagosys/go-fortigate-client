package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Nexinto/go-fortigate-client/fortigate"
)

func mkName(n string) string {
	parts := strings.Split(n, "-")
	for i, p := range parts {
		p = strings.Title(p)
		p = strings.Replace(p, ".", "", -1)
		p = strings.Replace(p, "+", "Plus", -1)
		p = strings.Replace(p, "/", "_", -1)
		p = strings.Replace(p, "(", "_", -1)
		p = strings.Replace(p, ")", "", -1)
		if strings.HasPrefix(p, "3") || strings.HasPrefix(p, "8") {
			p = "X" + p
		}
		parts[i] = p
	}
	return strings.Join(parts, "")
}

func main() {

	c, err := fortigate.NewWebClient(fortigate.WebClient{
		URL:      os.Getenv("FORTIGATE_URL"),
		User:     os.Getenv("FORTIGATE_USER"),
		Password: os.Getenv("FORTIGATE_PASSWORD"),
		ApiKey:   os.Getenv("FORTIGATE_API_KEY")})
	if err != nil {
		panic(err)
	}

	endpoints, err := c.Schema()
	if err != nil {
		panic(err)
	}

	var endpoints2 []fortigate.Endpoint

	for _, e := range endpoints {
		if strings.HasPrefix(e.Path, "diagnose__tree") {
			continue
		}

		if strings.HasPrefix(e.Path, "execute__tree") {
			continue
		}

		if e.Schema.Category == "complex" {
			// TODO
			continue
		}

		if e.Path != "firewall" && e.Path != "certificate" && e.Path != "vpn" {
			continue
		}

		if e.Path == "firewall" && e.Name == "vip" {
			// Create an alias "VIP" because compatibility
			e.Alias = "VIP"
			endpoints2 = append(endpoints2, e)
		}
		if e.Path == "firewall" && e.Name == "policy" {
			endpoints2 = append(endpoints2, e)
		}
	}

	tt := `// WARNING: This file was generated by generator.go

package fortigate

import (
  "fmt"
  "net/http"
  "strconv"
)

// A fortigate API client
type Client interface {
{{ range $e := . }}

  // List all {{ typeName $e }}s
  List{{ typeName $e }}s() ([]*{{ typeName $e }},error)

  // Get a {{ typeName $e }} by name
  Get{{ typeName $e }}(mkey {{ goType $e.Schema.MkeyType }}) (*{{ typeName $e }},error)

  // Create a new {{ typeName $e }}
  Create{{ typeName $e }}(*{{ typeName $e }}) ({{ goType $e.Schema.MkeyType }}, error)

  // Update a {{ typeName $e }}
  Update{{ typeName $e }}(*{{ typeName $e }}) error

  // Delete a {{ typeName $e }} by name
  Delete{{ typeName $e }}(mkey {{ goType $e.Schema.MkeyType }}) error
{{ end }}
}

// Fake Fortigate Client
type FakeClient struct {
{{ range $e := . }}
  {{ typeName $e }}s map[{{ goType $e.Schema.MkeyType }}]*{{ typeName $e }}
{{ if (eq "integer" $e.Schema.MkeyType) }}
  {{ typeName $e }}Counter int
{{ end -}}
{{ end }}
}

// Create a new fake client
func NewFakeClient() *FakeClient {
	return &FakeClient{
{{ range $e := . }}
	{{ typeName $e }}s: map[{{ goType $e.Schema.MkeyType }}]*{{ typeName $e }}{},
{{ if (eq "integer" $e.Schema.MkeyType) -}}
  {{ typeName $e }}Counter: 1,
{{ end -}}
{{ end }}
	}
}

{{ range $e := . }}
{{ range $f, $attr := $e.Schema.Children }}
{{ if (eq $attr.Type "option") }}
// {{ wrapcomment $attr.Help }}
type {{ typeName $e }}{{ mkName $f }} {{ goType $attr.Type }}

{{ end -}}{{/* if */}}

{{ if (eq $attr.Category "table") -}}
// {{ wrapcomment $attr.Help }}
type {{ typeName $e }}{{ mkName $f }} struct {
{{ range $tfield, $tattr := $attr.Children }}
  // {{ wrapcomment $tattr.Help }}
  {{ mkName $tfield }} {{ goType $tattr.Type }} {{ bq }}json:"{{ $tfield }},omitempty"{{ bq }}
{{ end -}}{{/* range $tfield */}}
}
{{ end -}}{{/* if */}}
{{ end -}}{{/* range $f */}}

{{ range $f, $attr := $e.Schema.Children }}
{{ if (eq $attr.Type "option") -}}
const (
{{- range $o := $attr.Options }}
  // {{ wrapcomment $o.Help }}
  {{ typeName $e }}{{ mkName $f }}{{ mkName $o.Name }} {{ typeName $e }}{{ mkName $f }} = "{{ $o.Name }}"
{{ end -}}{{/* range $o */}}
)
{{ end -}}{{/* if */}}
{{ end -}}{{/* range $f */}}

// {{ wrapcomment $e.Schema.Help }}
type {{ typeName $e }} struct {
{{ range $f, $attr := $e.Schema.Children }}
  // {{ wrapcomment $attr.Help }}
{{- if (eq $attr.Category "table") }}
  {{ mkName $f }} []{{ typeName $e }}{{ mkName $f }} {{ bq }}json:"{{ $f }},omitempty"{{ bq }}
{{ else if (eq $attr.Type "option") }}
  {{ mkName $f }} {{ typeName $e }}{{ mkName $f }} {{ bq }}json:"{{ $f }},omitempty"{{ bq }}
{{ else }}
  {{ mkName $f }} {{ goType $attr.Type }} {{ bq }}json:"{{ $f }},omitempty"{{ bq }}
{{ end -}}{{/* if */}}
{{ end -}}{{/* range $f */}}
}

// Returns the value that identifies a {{ typeName $e }}
func (x *{{ typeName $e }}) MKey(){{ goType $e.Schema.MkeyType }} {
  return x.{{ mkName $e.Schema.Mkey }}
}

// The results of a Get or List operation
type {{ typeName $e }}Results struct {
  Results []*{{ typeName $e }} {{ bq }}json:"results"{{ bq }}
  Mkey  {{ goType $e.Schema.MkeyType }} {{ bq }}json:"mkey"{{ bq }}
  Result
}

// List all {{ typeName $e }}s
func (c *WebClient) List{{ typeName $e }}s() (res []*{{ typeName $e }}, err error) {
  var results {{ typeName $e }}Results
   _, err = c.do(http.MethodGet, "{{ $e.Path }}/{{ $e.Name }}", nil, nil, &results)
	if err != nil {
    return []*{{ typeName $e }}{}, fmt.Errorf("error listing {{ typeName $e }}s: %s", err.Error())
  }
  res = results.Results
  return
}

// Get a {{ typeName $e }} by name
func (c *WebClient) Get{{ typeName $e }}(mkey {{ goType $e.Schema.MkeyType }}) (res *{{ typeName $e }}, err error) {
  var results {{ typeName $e }}Results
  _, err = c.do(http.MethodGet, "{{ $e.Path }}/{{ $e.Name }}/" + {{ bareMkeyAsString $e "mkey" }}, nil, nil, &results) 
	if err != nil {
    return &{{ typeName $e }}{}, fmt.Errorf("error getting {{ typeName $e }} '%s': %s", {{ bareMkeyAsString $e "mkey" }}, err.Error())
  }
  res = results.Results[0]
  return
}

// Create a new {{ typeName $e }}
func (c *WebClient) Create{{ typeName $e }}(obj *{{ typeName $e }}) (id {{ goType $e.Schema.MkeyType }}, err error) {
  _, err = c.do(http.MethodPost, "{{ $e.Path }}/{{ $e.Name }}", nil, obj, nil)   
	if err != nil {
    return {{ emptyliteralfor $e.Schema.MkeyType }}, fmt.Errorf("error creating {{ typeName $e }} '%s': %s", {{ mkeyAsString $e "obj" }}, err.Error())
  }
  return
}

// Update a {{ typeName $e }}
func (c *WebClient) Update{{ typeName $e }}(obj *{{ typeName $e }}) error {
  _, err := c.do(http.MethodPut, "{{ $e.Path }}/{{ $e.Name }}/" + {{ mkeyAsString $e "obj" }}, nil, obj, nil)
	if err != nil {
    return fmt.Errorf("error updating {{ typeName $e }} '%s': %s", {{ mkeyAsString $e "obj" }}, err.Error())
  }
  return err
}

// Delete a {{ typeName $e }} by name
func (c *WebClient) Delete{{ typeName $e }}(mkey {{ goType $e.Schema.MkeyType }}) error {
  _, err := c.do(http.MethodDelete, "{{ $e.Path }}/{{ $e.Name }}/" + {{ bareMkeyAsString $e "mkey" }}, nil, nil, nil)
	if err != nil {
    return fmt.Errorf("error deleting {{ typeName $e }} '%s': %s", {{ bareMkeyAsString $e "mkey" }}, err.Error())
  }
  return err
}

// List all {{ typeName $e }}s
func (c *FakeClient) List{{ typeName $e }}s() (res []*{{ typeName $e }}, err error) {
  for _, r := range c.{{ typeName $e }}s {
    res = append(res, r)
  }
  return
}

// Get a {{ typeName $e }} by name
func (c *FakeClient) Get{{ typeName $e }}(mkey {{ goType $e.Schema.MkeyType }}) (*{{ typeName $e }}, error) {
	if res, ok := c.{{ typeName $e }}s[mkey]; ok {
		return res, nil
	} else {
		return &{{ typeName $e }}{}, fmt.Errorf("error getting {{ typeName $e }} '%s': not found", {{ bareMkeyAsString $e "mkey" }})
	} 
}

// Create a new {{ typeName $e }}
func (c *FakeClient) Create{{ typeName $e }}(obj *{{ typeName $e }}) (id {{ goType $e.Schema.MkeyType }}, err error) {
{{ if (eq "integer" $e.Schema.MkeyType) -}}
  id = c.{{ typeName $e }}Counter
  c.{{ typeName $e }}Counter ++
{{ else -}}
  id = obj.{{ mkName $e.Schema.Mkey }}
{{ end -}}
	c.{{ typeName $e }}s[id] = obj
	return
}

// Update a {{ typeName $e }}
func (c *FakeClient) Update{{ typeName $e }}(obj *{{ typeName $e }}) (err error) {
	c.{{ typeName $e }}s[obj.{{ mkName $e.Schema.Mkey }}] = obj
	return nil
}

// Delete a {{ typeName $e }} by name
func (c *FakeClient) Delete{{ typeName $e }}(mkey {{ goType $e.Schema.MkeyType }}) (err error) {
	delete(c.{{ typeName $e }}s, mkey)
	return nil
}

{{ end -}}{{/* range $e */}}
`

	funcMap := template.FuncMap{
		"typeName": func(e fortigate.Endpoint) string {
			if e.Alias != "" {
				return e.Alias
			} else {
				return mkName(e.Path) + mkName(e.Name)
			}
		},
		"mkName": func(t string) string {
			return mkName(t)
		},
		"goType": func(t string) string {
			switch t {
			case "integer":
				return "int"
			default:
				return "string"
			}
		},
		"bq": func() string {
			return "`"
		},
		"emptyliteralfor": func(s string) string {
			switch s {
			case "integer":
				return "0"
			default:
				return `""`
			}
		},
		"wrapcomment": func(s string) string {
			return strings.Join(strings.Split(strings.Replace(s, "\t", "", -1), "\n"), "\n// ")
		},
		"mkeyAsString": func(e fortigate.Endpoint, vn string) string {
			switch e.Schema.MkeyType {
			case "integer":
				return "strconv.Itoa(" + vn + "." + mkName(e.Schema.Mkey) + ")"
			default:
				return vn + "." + mkName(e.Schema.Mkey)
			}
		},
		"bareMkeyAsString": func(e fortigate.Endpoint, vn string) string {
			switch e.Schema.MkeyType {
			case "integer":
				return "strconv.Itoa(" + vn + ")"
			default:
				return vn
			}
		},
	}

	t, err := template.New("type").Funcs(funcMap).Parse(tt)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer

	err = t.Execute(&buffer, endpoints2)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./fortigate/types.go", buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

}
