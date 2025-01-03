package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"
)

type Filer interface {
	IsFile() bool
}

type TemplateDef struct {
	Engine   string
	Filepath string
}

type File struct {
	Name         string
	Content      string `json:"-"`
	*TempWrapper `json:",omitempty"`
}

type TempWrapper struct {
	TemplateDef `json:"template"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	reg := regexp.MustCompile(`\{.+\}`)
	if reg.Match(data) {
		s := string(data)[1 : len(data)-1]
		before, after, _ := strings.Cut(s, ":")

		before = strings.ReplaceAll(before, `"`, "")

		if path.Ext(before) != "" {
			f.Name = before
			var tmp TempWrapper
			err := json.Unmarshal([]byte(after), &tmp)
			if err != nil {
				return err
			}
			f.TempWrapper = &tmp
		} else {
			return fmt.Errorf("the data your unmarshaling is not a file")
		}

		return nil
	}
	s := string(data)[1 : len(data)-1]
	if path.Ext(s) != "" {
		f.Name = s
	} else {
		return fmt.Errorf("the data your unmarshaling is not a file")
	}

	return nil
}

func (f File) IsFile() bool {
	return true
}

type Folder struct {
	Name   string
	Filers FolderStruct
}

// cleanSpaces removes all spaces and newline.
// the result should be a []byte containing an inlined string without spaces.
func cleanSpaces(data []byte) []byte {
	s := string(data)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	data = []byte(s)
	return data
}

func (f *Folder) UnmarshalJSON(data []byte) error {
	data = cleanSpaces(data)
	reg := regexp.MustCompile(`\{.+\}`)
	if reg.Match(data) {
		s := string(data)[1 : len(data)-1]
		before, after, _ := strings.Cut(s, ":")

		before = strings.ReplaceAll(before, `"`, "")
		f.Name = before
		var filers FolderStruct
		err := json.Unmarshal([]byte(after), &filers)
		if err != nil {
			return err
		}
		f.Filers = filers

		return nil
	}
	s := string(data)[1 : len(data)-1]
	if path.Ext(s) == "" {
		f.Name = s
		f.Filers = nil
	}
	return nil
}

func (f Folder) IsFile() bool {
	return false
}

type FolderStruct []Filer

type FilerPlugin struct {
	IsFolder bool
	Name     string
	Content  string
	Children []FilerPlugin
}

type FolderStructPlugin []FilerPlugin

func (f FolderStruct) ToPluginInput() (string, error) {
	converted := convert(f)

	data, err := json.Marshal(converted)

	return string(data), err
}

func convert(f FolderStruct) []FilerPlugin {
	res := make([]FilerPlugin, len(f))
	for i, filer := range f {
		var fp FilerPlugin
		if filer.IsFile() {
			data := filer.(File)
			fp = FilerPlugin{
				IsFolder: !filer.IsFile(),
				Name:     data.Name,
				Content:  data.Content,
			}
		} else {
			data := filer.(Folder)
			fp = FilerPlugin{
				IsFolder: !filer.IsFile(),
				Name:     data.Name,
				Children: convert(data.Filers),
			}
		}

		res[i] = fp
	}
	return res
}

func (f *FolderStruct) UnmarshalJSON(data []byte) error {
	data = cleanSpaces(data)
	var (
		dec             = json.NewDecoder(strings.NewReader(string(data)))
		parts           = []string{}
		s               = ""
		lastOpenedDelim = ""
		countOdd        = 0
		countOpenDelim  = 0
	)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if dec.InputOffset() == 1 {
			continue
		}

		// the idea is to check on matching delimiter to parse object/array into `s`.
		// once `s` is complete (previously opened delim are closed) we append it to `parts` and reinit `s` to empty string
		stringed := fmt.Sprintf("%v", t)

		if _, ok := t.(json.Delim); ok {
			if stringed == "{" || stringed == "[" {
				if stringed == "{" {
					countOdd = 1
				}
				countOpenDelim++
				s += stringed
				lastOpenedDelim = stringed
			} else {
				countOpenDelim--
				s += stringed

				if countOpenDelim == 0 {
					parts = append(parts, s)
					s = ""
				}
			}
			continue
		}
		stringed = fmt.Sprintf("%q", stringed)
		if countOpenDelim == 0 {
			parts = append(parts, stringed)
		} else {
			if lastOpenedDelim == "[" && dec.More() {
				stringed += ","
			}
			if lastOpenedDelim == "{" && dec.More() {
				if countOdd%2 == 0 {
					stringed += ","
				} else {
					stringed += ":"
				}
				countOdd++
			}
			s += stringed
		}

	}

	var errs []error
	var file File
	var folder Folder

	for _, s := range parts {
		err := json.Unmarshal([]byte(s), &file)
		if err != nil {
			errs = append(errs, err)
			err = json.Unmarshal([]byte(s), &folder)
			if err != nil {
				errs = append(errs, err)
				return errors.Join(errs...)
			} else {
				*f = append(*f, folder)
			}
		} else {
			*f = append(*f, file)
		}
	}
	return nil
}

type folderStructType string

const (
	file   folderStructType = "file"
	folder folderStructType = "folder"
)

type jsonFile struct {
	jsonFiler
	Spec *struct {
		Content string `json:"content,omitempty"`
	} `json:"spec,omitempty"`
}

type jsonFolder struct {
	jsonFiler
	Spec *struct {
		Children FolderStruct `json:"children,omitempty"`
	} `json:"spec,omitempty"`
}

type jsonFiler struct {
	Name string           `json:"name"`
	Type folderStructType `json:"type"`
}

func (f File) MarshalJSON() ([]byte, error) {
	j := jsonFile{
		jsonFiler: jsonFiler{
			Name: f.Name,
			Type: file,
		},
	}

	if f.Content != "" {
		j.Spec = &struct {
			Content string "json:\"content,omitempty\""
		}{
			Content: f.Content,
		}
	}

	return json.Marshal(j)
}

func (f Folder) MarshalJSON() ([]byte, error) {
	j := jsonFolder{
		jsonFiler: jsonFiler{
			Name: f.Name,
			Type: folder,
		},
	}

	if f.Filers != nil {
		j.Spec = &struct {
			Children FolderStruct "json:\"children,omitempty\""
		}{
			Children: f.Filers,
		}

	}

	return json.Marshal(j)
}
