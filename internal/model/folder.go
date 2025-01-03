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

// Filer defines the family of types in folder_struct definition.
// Filer types are either File or Folder.
type Filer interface {
	IsFile() bool
}

type TemplateDef struct {
	Engine   string // template engine to run
	Filepath string // path of the template to apply
}

// A File is one of the two types that can be defined in folder_struct.
// It has a Name, and it can have a template definition.
// The content is retrieved at runtime from the template definition.
type File struct {
	Name         string              // name of the file
	Content      string              `json:"-"` // content of the file
	*TempWrapper `json:",omitempty"` // template definition of the file
}

type TempWrapper struct {
	TemplateDef `json:"template"`
}

// UnmarshalJSON implementes encoding/json.Unmarshaller on File type
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

// IsFile implements the Filer interface on File type.
func (f File) IsFile() bool {
	return true
}

// A Folder is one of the two types that can be defined in folder_struct.
// It has a Name, and it can have children as an array of Filer.
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

// UnmarshalJSON implementes encoding/json.Unmarshaller on Folder type
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

// IsFile implements the Filer interface on Folder type.
func (f Folder) IsFile() bool {
	return false
}

// FolderStruct is a user-defined folder_struct. It's just an array of Filer.
type FolderStruct []Filer

// UnmarshalJSON implementes encoding/json.Unmarshaller on FolderStruct type
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
