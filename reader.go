package configfile

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"bytes"
)

// Path separator
var sep string = string(os.PathSeparator)

type ConfigReader struct {
	Subfolder            string
	ExcludeCurrentFolder bool
	ExcludeHomedir       bool
	IncludeEtc           bool
}

// Returns current user's home folder
func GetHomedir() (string, error) {
	if runtime.GOOS != "windows" {
		// Found using environment variable
		// Useful for cross-compile to darwin
		if home := os.Getenv("HOME"); home != "" {
			return home, nil
		}
	}

	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return user.HomeDir, nil
}

// Returns list of folders, where configuration files will be searched
func (r ConfigReader) ListFolders() ([]string, error) {
	folders := []string{}

	// Current folder
	if !r.ExcludeCurrentFolder {
		path, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		folders = append(folders, path)
	}

	// Home folder
	if !r.ExcludeHomedir {
		path, err := GetHomedir()
		if err != nil {
			return nil, err
		}

		folders = append(folders, path)
	}

	if r.IncludeEtc && runtime.GOOS != "windows" && sep == "/" {
		folders = append(folders, "/etc")
	}

	// Appending subfolder
	if r.Subfolder != "" {
		for i := 0; i < len(folders); i++ {
			folders[i] = folders[i] + sep + r.Subfolder
		}
	}

	return folders, nil
}

// Returns fullname for configuration file
func (r ConfigReader) Resolve(name string) (string, error) {
	ff, err := r.ListFolders()
	if err != nil {
		return "", err
	}
	for _, path := range ff {
		filename := path + sep + name
		_, err = os.Stat(filename)
		if err == nil {
			return filename, nil
		}
	}

	return "", fmt.Errorf("File %s not found", name)
}

// Searches for congfiguration file in all available folders
// and returns it or error
func (r ConfigReader) GetFile(name string) (*os.File, error) {
	fullname, err := r.Resolve(name)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fullname)
	if err == nil {
		return file, nil
	}

	return nil, fmt.Errorf("Unable to read configuration file %s. Not exisits or not readable", name)
}

func (r ConfigReader) ReadFile(name string) ([]byte, error) {
	file, err := r.GetFile(name)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(file);
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Reads configuration file and unmarshalls it data using JSON unmarshaller
func (r ConfigReader) ReadJson(name string, target interface{}) error {
	file, err := r.GetFile(name)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := json.NewDecoder(file)
	err = reader.Decode(&target)
	if err != nil {
		return err
	}

	return nil
}
