package configfile

import (
	"fmt"
	"os"
	"os/user"
	"io/ioutil"
	"encoding/json"
)

type ConfigFile struct {
	Filename string
	FullPath string
	PathSeparator string
	File *os.File
}


/// Creates new ConfigFile instance
/// If config file not found, returns error
func NewConfigFile(filename string, etcLookup bool) (*ConfigFile, error) {
	cnf := ConfigFile{Filename: filename}

	var stat bool
	var err error

	cnf.PathSeparator = string(os.PathSeparator)

	// Read from current folder
	current, err := os.Getwd()
	if stat, err = cnf.readFrom(current); err != nil {
		return nil, err
	} else if stat {
		return &cnf, nil
	}

	// Read from homedir
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	if stat, err = cnf.readFrom(user.HomeDir); err != nil {
		return nil, err
	} else if stat {
		return &cnf, nil
	}

	// Read from /etc
	if etcLookup {
		if stat, err = cnf.readFrom("/etc/"); err != nil {
			return nil, err
		} else if stat {
			return &cnf, nil
		}
	}

	return nil, fmt.Errorf("Unable to find configuration file %s", filename)
}

func (cnf *ConfigFile) readFrom(folder string) (bool, error) {
	if len(folder) > 0 && folder[len(folder)-1:] != cnf.PathSeparator {
		folder = folder + cnf.PathSeparator
	}

	name := folder + cnf.Filename

	if _, err := os.Stat(name); err == nil {
		// File found, opening
		cnf.File, err = os.Open(name)

		if err != nil {
			return false, err
		}

		cnf.FullPath = name
		return true, nil
	} else {
		// File not found
		return false, nil
	}
}

/// Reads all bytes from configuration file
func (cnf *ConfigFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(cnf.File)
}

func (cnf *ConfigFile) DecodeJson(strct interface{}) error {
	bytes, err := cnf.ReadAll()

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &strct)
}