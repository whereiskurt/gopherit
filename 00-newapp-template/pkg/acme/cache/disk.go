package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Disk struct {
	UseCrypto   bool
	CacheKey    []byte
	CacheFolder string
}

func NewDisk(folder string, key string, crypto bool) (d *Disk) {
	d = new(Disk)
	d.UseCrypto = crypto
	d.CacheKey = []byte(key)
	d.CacheFolder = strings.TrimSuffix(folder, "/")
	return
}

func (d *Disk) Fetch(filename string) (bb []byte, err error) {
	if _, stat := os.Stat(filename); os.IsNotExist(stat) {
		// File doesn't exist return no error
		return
	}

	bb, err = ioutil.ReadFile(filename)
	if err != nil {
		// Error! Failed to read file!
		return
	}
	if d.UseCrypto {
		var dd []byte
		dd, err = Decrypt(bb, d.CacheKey)
		if err == nil {
			// Successfully decrypted!
			bb = dd
		} else {
			err = errors.New(fmt.Sprintf("cache failed to decrypt file '%s' : %s", filename, err))
		}
	}
	return
}
func (d *Disk) Store(filename string, bb []byte) (err error) {

	// if d.jqpath != "" {
	// 	var pretty bytes.Buffer
	// 	raw := bb
	// 	cmd := exec.Command(d.jqpath, ".")
	// 	cmd.Stdin = strings.NewReader(string(raw))
	// 	cmd.Stdout = &pretty
	// 	err = cmd.Run()
	// 	if err == nil {
	// 		bb = []byte(pretty.String())
	// 	}
	// }

	if d.UseCrypto && len(d.CacheKey) > 0 {
		bb, err = Encrypt(bb, d.CacheKey)
		if err != nil {
			return
		}
	}

	err = os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, bb, 0644)

	return
}
func (d *Disk) Clear(filename string) {
	os.Remove(filename) // delete the cache file.
}