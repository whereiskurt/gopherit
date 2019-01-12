// Code generated by vfsgen; DO NOT EDIT.

// +build release

package config

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// TemplateFolder statically implements the virtual filesystem provided to vfsgen.
var TemplateFolder = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2019, 1, 12, 21, 43, 27, 576602205, time.UTC),
		},
		"/default.gophercli.yaml": &vfsgen۰CompressedFileInfo{
			name:             "default.gophercli.yaml",
			modTime:          time.Date(2018, 12, 15, 15, 23, 48, 183169106, time.UTC),
			uncompressedSize: 340,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xcf\x4d\x4b\xc3\x40\x10\xc6\xf1\xfb\x7e\x8a\x07\x22\xf6\x52\xf2\x52\xdf\xf7\xa6\x55\x51\xaa\x08\x09\x7a\x5f\x37\x93\xec\xea\x9a\x09\xbb\x63\x4b\xbe\xbd\x34\x14\xbd\x7b\x7d\xe0\xff\x1b\x26\xcb\xb0\xe6\x71\x82\x38\x82\x0d\x9e\x06\x41\x22\x2b\x9e\x07\x18\x1c\x3d\xbc\x3c\xdf\xdd\x3e\xd6\xe8\x7c\x20\x2c\xf2\x9e\x47\x47\x31\x9f\xcc\x57\x58\xc0\x0c\x2d\x78\x4b\x31\xfa\x96\x54\x96\x61\xe7\xc5\xcd\xcc\xb5\xb5\x94\xd2\x86\xa6\xe3\x86\x6c\x24\xd9\xd0\xa4\xd6\x33\xad\x15\x70\x63\x12\xbd\xd6\x4f\x1a\x4e\x64\xd4\x45\x11\xd8\x9a\xe0\x38\x89\xae\xca\xaa\xac\x14\xfe\x7a\x0d\xf3\x6e\xab\xd5\xc9\xb2\x77\xfe\xe2\xf2\x4a\x01\xbf\xa0\x46\x4b\xdd\xe9\xd9\xf9\xf2\xe3\x33\x94\xd5\x4a\xa9\x86\xe2\x96\xa2\xfe\x67\x0e\xd4\xcc\x72\xcf\xa1\xa5\xa8\x91\x17\x96\x87\xce\xf7\x45\xcb\x36\x32\x4b\xa1\xd4\x1b\xc5\xe4\x79\xd8\xfb\x8d\xe3\x5d\x43\xfb\x63\x90\xf8\x4d\x87\xe5\xf0\xe0\x3c\xfd\x04\x00\x00\xff\xff\x0e\xd9\x61\xb1\x54\x01\x00\x00"),
		},
		"/default.test.gophercli.yaml": &vfsgen۰CompressedFileInfo{
			name:             "default.test.gophercli.yaml",
			modTime:          time.Date(2018, 12, 15, 15, 23, 48, 0, time.UTC),
			uncompressedSize: 340,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xcf\x4d\x4b\xc3\x40\x10\xc6\xf1\xfb\x7e\x8a\x07\x22\xf6\x52\xf2\x52\xdf\xf7\xa6\x55\x51\xaa\x08\x09\x7a\x5f\x37\x93\xec\xea\x9a\x09\xbb\x63\x4b\xbe\xbd\x34\x14\xbd\x7b\x7d\xe0\xff\x1b\x26\xcb\xb0\xe6\x71\x82\x38\x82\x0d\x9e\x06\x41\x22\x2b\x9e\x07\x18\x1c\x3d\xbc\x3c\xdf\xdd\x3e\xd6\xe8\x7c\x20\x2c\xf2\x9e\x47\x47\x31\x9f\xcc\x57\x58\xc0\x0c\x2d\x78\x4b\x31\xfa\x96\x54\x96\x61\xe7\xc5\xcd\xcc\xb5\xb5\x94\xd2\x86\xa6\xe3\x86\x6c\x24\xd9\xd0\xa4\xd6\x33\xad\x15\x70\x63\x12\xbd\xd6\x4f\x1a\x4e\x64\xd4\x45\x11\xd8\x9a\xe0\x38\x89\xae\xca\xaa\xac\x14\xfe\x7a\x0d\xf3\x6e\xab\xd5\xc9\xb2\x77\xfe\xe2\xf2\x4a\x01\xbf\xa0\x46\x4b\xdd\xe9\xd9\xf9\xf2\xe3\x33\x94\xd5\x4a\xa9\x86\xe2\x96\xa2\xfe\x67\x0e\xd4\xcc\x72\xcf\xa1\xa5\xa8\x91\x17\x96\x87\xce\xf7\x45\xcb\x36\x32\x4b\xa1\xd4\x1b\xc5\xe4\x79\xd8\xfb\x8d\xe3\x5d\x43\xfb\x63\x90\xf8\x4d\x87\xe5\xf0\xe0\x3c\xfd\x04\x00\x00\xff\xff\x0e\xd9\x61\xb1\x54\x01\x00\x00"),
		},
		"/template": &vfsgen۰DirInfo{
			name:    "template",
			modTime: time.Date(2019, 1, 12, 15, 33, 2, 868011399, time.UTC),
		},
		"/template/client": &vfsgen۰DirInfo{
			name:    "client",
			modTime: time.Date(2019, 1, 12, 14, 38, 20, 372064868, time.UTC),
		},
		"/template/client/table.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "table.tmpl",
			modTime:          time.Date(2019, 1, 8, 4, 46, 51, 644563666, time.UTC),
			uncompressedSize: 3130,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x96\x5f\x6f\xda\x30\x10\xc0\xdf\xf9\x14\x27\x0f\xfa\x47\x25\x71\xd9\xd6\x97\x4a\x3c\x4c\x42\x9a\x78\xd9\x53\x1f\x91\xae\x59\x71\x42\xd6\x36\x89\x12\x8b\x69\xb2\xf3\xdd\xa7\xb3\x43\xa1\xc5\x0e\x21\x50\x07\x39\xc1\xf7\xc7\x67\x72\xbf\x3b\x38\x87\x95\x94\xc5\x3d\xe7\x45\x24\xf3\xf2\xcf\x73\xf8\x94\xbf\xf2\x2a\x8f\xe5\xdf\xa8\x14\x5c\x46\x51\xc2\xbf\x14\xd3\x65\x5a\x15\x2f\xd1\xbf\x8b\x78\x3a\x5f\x97\xa9\xbc\x90\xd3\x24\x2f\x56\xa2\x1c\xdd\xfe\x18\x28\xb5\x14\x71\x9a\x09\x60\x3f\xcd\xda\x43\xf4\xfb\x45\xb0\xba\x1e\xdc\xf4\x1a\x27\x98\x81\x63\xe0\xde\x8a\xcb\x0c\x91\x14\x11\x69\xa2\x2f\x1a\x34\x3d\x9b\x15\xb3\xe0\x34\xe3\x80\x8f\xa0\x39\x20\x2c\x34\x5c\xbe\xcd\x60\x56\xe0\x12\x51\xbb\xcc\x34\x5c\xa1\x36\xf3\xb5\xd9\xc7\xcc\xe6\x02\x44\x4e\x37\xf7\x6e\x0b\xc4\x31\xe8\x05\x22\x72\x0d\x21\xcd\xe4\x06\xcd\x0a\x3d\x7a\xcc\x34\x19\xd8\x93\x6f\xb4\x3a\xfd\x24\x1d\xc6\x07\xb3\xa0\xc7\xb8\x19\x28\x05\x45\x99\x66\x32\x06\xa6\x61\x14\x7c\xaf\xf4\x28\x98\xdc\xd2\xfc\xed\xae\xd2\xc0\x80\xcd\x67\x0c\xd8\xaf\xe8\x55\x30\x60\x33\x51\x3d\x95\x69\x21\xd3\x3c\x63\x40\x39\xd6\x77\xd3\x32\xca\x12\x01\xc3\xe7\x31\x0c\xd7\x70\x3f\x85\x10\x82\xba\x1e\xd0\xa9\xda\x03\x5a\x64\x0c\xc2\xf9\x0c\x42\x0a\x08\xc2\x9d\x78\x28\x1c\xa5\x02\x10\xd9\xd2\xf8\xea\x1d\x9a\xc8\x96\x75\x3d\xd8\xe1\xea\x61\x95\x66\x49\x4f\xac\x7a\xa9\x5b\x6c\xd0\x09\x90\x4b\xdd\xe4\xa1\x41\xe6\x0a\xaf\x1b\x5e\x2c\x52\x1e\x75\x42\xc3\x02\xb3\x45\xe6\x11\x5c\xb8\x6c\xbd\xbf\xbb\x2c\x42\x4e\x75\x02\xa5\x41\x03\xb7\x88\x8c\xbd\xea\xfb\x63\x83\x4c\x47\xf5\x33\xe0\xe0\xc2\x00\x6c\xda\xc1\x28\xf8\x7a\x57\xe9\x16\x0e\x8e\xcd\xb5\xdd\xf4\x4f\x28\xff\x93\x03\x00\x80\x7e\x17\xca\xa7\x20\xe0\x4c\xfc\xa6\xa1\x50\xfa\x57\x27\xb6\x15\xff\xf8\x1c\x87\xae\x04\xf9\x88\x92\x17\x32\x8f\xc3\x43\x5d\xca\x8f\xa1\xdb\x61\x87\xfe\xe5\x03\xd5\xe9\xf0\x70\x67\xf3\xa3\xec\x8e\xf0\x70\xcf\xf3\xc2\xee\x76\xd8\xa5\x1b\x7a\xca\x41\xd7\xb7\xdc\x79\xec\x39\x3c\x06\xe2\x6e\x44\x79\x3b\x99\xa5\xb9\xa9\x32\x13\x23\x98\xb4\x16\x99\x46\x62\x58\x0c\xcf\xd9\x8c\xdb\x0f\xe0\xab\x53\x5b\x89\x24\x89\x34\x12\x5b\x28\x1a\xf9\x91\x47\xa7\xa2\x36\x5c\x27\x54\xd7\xe8\x66\x4a\x1b\x3d\xec\x56\xb7\xe1\x5a\x5a\xb9\xdc\xc8\x65\x5b\xf5\x53\x2a\x8d\x21\xcb\xe5\x26\xb0\xba\x56\xea\xfc\x41\x31\xd6\x7c\xc8\xbd\xad\xa0\x27\xff\x0b\x69\x7f\x27\x76\x93\xff\x01\x00\x00\xff\xff\x7d\xd8\xe2\xa8\x3a\x0c\x00\x00"),
		},
		"/template/cmd": &vfsgen۰DirInfo{
			name:    "cmd",
			modTime: time.Date(2018, 12, 23, 21, 6, 24, 800610086, time.UTC),
		},
		"/template/cmd/client.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "client.tmpl",
			modTime:          time.Date(2018, 12, 23, 20, 34, 40, 860286587, time.UTC),
			uncompressedSize: 1649,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x94\xc1\x4e\xdb\x40\x10\x86\xef\x7e\x8a\x91\xa9\x94\x20\xd9\x46\xd0\x1b\x92\x0f\x25\x50\x8a\x54\x08\x6a\x82\x38\x20\x0e\x6b\xef\xd8\xde\x6a\xbd\x9b\x7a\xc7\x88\x28\xf2\xbb\x57\x1e\x3b\x21\x05\x3b\x12\xbd\xa0\xc5\x33\xf9\xe6\x9f\xdd\xf9\x67\xb3\x91\x98\x29\x83\xe0\xcf\xb4\x42\x43\x0f\x4e\xe4\xe8\x43\xd3\x78\x33\x5b\x96\xc2\xc8\x73\x0f\x00\x20\xb7\xab\x02\xab\x54\x2b\x48\x39\x0d\x9e\x16\x0f\x17\xb3\xf9\xed\xed\xb7\xbb\xcb\x67\x78\x9a\xdf\x2f\x6f\xe6\x77\x8b\x67\xef\x5d\xf2\x48\x56\x9b\x34\x5d\x16\x08\x12\x33\x51\x6b\x82\xb4\xab\x05\x93\x8e\x3e\x81\x54\x18\x48\x10\x6c\xa9\x88\x50\x46\xc7\x9e\xb7\xa8\x93\x3e\xcb\x75\x92\xb4\x72\xc4\x87\x7a\x25\x05\x21\x1f\x25\x6a\x24\xf4\xbc\xf9\x8a\x94\x35\x7d\x22\x84\x65\x00\x10\xb6\x7f\xc3\xb0\xb4\x12\xf9\xa3\xa3\x4a\x99\x9c\x8f\xb0\x40\x02\xb2\x40\x22\xd1\x18\xa4\xee\x25\xf8\xed\xac\x09\x5e\x4b\x0d\x4f\xbd\xc2\x73\x8e\x3d\xf7\xbc\x9c\x79\x79\xcb\xeb\x7a\x1d\xe2\x15\x08\xd7\x1c\xbc\xb9\x84\xa9\x65\x41\x42\x47\xc7\x3d\x82\x18\x41\x2d\x82\x8a\xfe\x87\x43\x88\x65\x1b\x1c\x22\x30\x25\xcc\x0d\x8b\x30\xa2\xc4\x21\x42\x27\x00\x38\xfc\x22\x74\x8d\x90\xd9\x0a\x84\x94\xca\xe4\x41\x77\x71\x7c\x72\x28\xaa\xb4\xad\xf4\x0f\x5b\x32\x5b\xa2\x4b\x0f\xb0\xdb\x70\xa5\x58\xdc\x67\xe0\xc4\xc2\x69\x4c\x38\xb7\xfd\x9f\xba\x89\x75\xd3\x98\xee\x0e\xfd\x59\xd9\x61\xca\x0f\x96\x8a\xb4\xc0\x5d\x2d\xa0\xaa\xc6\x93\x4c\x68\x87\x00\xbf\xf0\x4f\x8d\x8e\xa0\x4d\x69\x2b\xd8\x8c\x1f\x30\xb1\x72\x0d\xd3\x1f\xcb\xe5\x3d\x5c\x5f\x2d\x1d\x58\xa3\xd7\x7b\x4f\x08\x61\x9a\xb5\x72\xd3\xcc\x6a\xd9\x0e\xd2\xbe\xdc\x19\x57\xeb\x23\x7c\x19\xd3\xed\x3c\xc2\x24\x3a\x89\x58\xcd\xc9\xe4\xd8\xf3\xae\x5e\x45\xb9\xd2\xe8\xce\x37\x1b\xc2\x72\xa5\x05\xed\x1c\xfd\x53\x39\xf2\x9b\x06\x3e\x46\x1e\xd8\x38\xc3\xb1\x4b\x76\x92\xdf\x34\xde\x66\x83\x46\x36\x8d\xe7\xbd\x5f\x15\x0c\x6e\x37\x45\x27\xf6\xe8\x08\x86\x2c\xad\xdc\x9b\xab\xfb\xff\xb7\xa3\x0c\x09\x66\xb6\x42\x70\x6f\xd6\x8e\x7a\xd8\x97\x8f\x1b\x87\xfd\x3e\x10\xde\xed\x81\x8f\x9f\x77\xfe\x8c\x4f\x83\xb3\xf1\x1c\x36\x60\x3c\x9e\x90\x9f\x42\x48\x67\xc1\xd7\x03\x55\x4c\xdc\x7d\x3b\xe5\xa1\x8e\x33\x44\x1a\xbf\xb8\xfe\xde\xdf\xae\x6e\x9f\xd8\x6d\xb3\x3d\xe5\x5b\x7f\xc7\xfe\x85\x4d\x92\x35\x2c\x4a\x45\x85\xbf\x75\x66\xec\x7f\xaf\xab\x6a\x1d\xf9\x87\x50\x7d\x83\x3b\xc3\xf9\xb7\xd6\xe0\x1a\x1e\x85\xd6\x48\xfe\xd6\x2c\xb1\xff\x58\x60\x85\x3c\xb5\x25\x27\x28\x17\xf9\xe3\x6d\xf4\x23\x32\xdc\x46\xb7\x89\xf7\xda\xd8\x72\xfe\x06\x00\x00\xff\xff\xde\x95\x9f\xf9\x71\x06\x00\x00"),
		},
		"/template/cmd/gophercli.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "gophercli.tmpl",
			modTime:          time.Date(2018, 12, 23, 20, 58, 7, 786942517, time.UTC),
			uncompressedSize: 1261,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\x51\x6f\xda\x30\x10\x7e\xcf\xaf\xb8\xf1\xc2\x8a\x62\xd2\xad\xdb\x4b\xa5\x49\x74\xac\x9b\x90\x4a\xa9\x04\xed\x4b\xa0\x56\x20\x47\xb0\xe6\xc4\xa9\x7d\x49\x87\x20\xfc\xf6\x29\x0e\x24\xb4\x63\x6a\x2f\x0f\xf1\x7d\xf7\x7d\x9f\xcf\xf6\x6d\x36\x21\x2e\x45\x82\xd0\xfa\xa5\xd2\x15\xea\xfe\xcd\xe0\xde\x04\x11\xb6\x8a\xc2\x81\x17\xe1\x72\xc6\xd8\xae\x0c\xc6\x18\xeb\xf2\xba\xcc\x5d\x97\xbb\x9d\x47\xce\x39\xdf\x03\x65\xb4\xdb\x9d\xa8\x33\x6d\x75\xdc\x9a\xe7\x01\xe7\x1e\x78\x6d\x9b\x3c\x76\x4b\xc0\xc6\x14\x1e\x7b\x4f\x00\xb0\xac\x99\x3e\x40\x6f\x09\x5b\xe8\x9d\x9d\x95\xe9\x16\xea\xb5\x04\x38\x07\xee\xd5\xcc\x69\xbb\x34\x99\xee\xec\xe6\xe5\x06\x30\xb5\xbb\x5b\xeb\x69\x4d\xdb\x1e\x9d\x83\x4b\xce\x25\x3f\x02\x06\x35\xad\x68\x40\xdf\xda\xf0\xd9\x49\xde\x31\x0a\x5b\xfb\x1d\x03\xa7\x79\xb0\x83\x1d\x9c\xe6\xbd\x90\xbf\x8e\x03\xcd\xef\xad\x88\x52\x73\xe9\x79\x91\x30\xd4\x8d\x04\xad\xb2\x79\x77\xa1\x62\x6f\x8e\x72\xae\x62\x8c\xd5\xcc\x71\xaa\x57\x1c\x4c\x3e\x40\x66\xd0\xc0\x55\x7f\x78\xfd\x71\x32\x3c\x83\xab\xbb\x01\x90\x02\x8d\xb9\xc0\x67\x08\x92\x10\x62\x15\x8a\xe5\x1a\x2a\x81\xb1\x10\xad\x50\x68\x98\xac\x44\x12\x99\xae\xe3\xfc\x14\x96\xa6\x11\x44\xb2\x54\x3a\x0e\x48\xa8\x04\x02\xba\xb4\x1d\x35\xdd\xd4\x8d\x3c\xaf\x50\xa3\x30\xbf\x33\x4d\x5e\x64\x7d\x05\x79\xe7\xe7\x2c\xc1\xe7\x20\x4d\x19\x61\x9c\xca\x80\xd0\x73\x1c\x3b\x63\x95\x4f\x45\x5c\x48\x01\x7e\x7f\x34\x1c\x5e\xdd\xfe\x98\x81\x3f\xbe\xff\xde\x24\xa3\xbb\xc9\x60\x74\x3b\x2e\x4f\x27\xd5\x3c\x90\x30\x4a\xcb\x4e\x4c\xa5\x7f\x40\x3d\x57\x46\xd0\xfa\x72\x7f\x51\x8c\x19\x21\x31\x21\x17\x80\x19\x8b\x8c\x91\x40\xaa\x28\x12\x49\xe4\xa9\x8c\xd2\x8c\x40\x62\x8e\x12\x7c\xfb\xfb\x34\xab\x95\x4f\x99\xc0\x52\x08\xec\xe9\x3d\xca\xcf\x8d\xb2\xbc\x22\xb7\x5a\xe6\xef\x50\x5e\xb0\x10\x97\x41\x26\xa9\x71\x08\x71\x9e\x45\x6e\xe5\x90\xbf\xed\xf0\xa5\x51\x92\x0e\x16\x78\x50\xe6\x6f\x2a\xbf\x36\x4a\x9b\x7f\xbb\xa8\x87\x6d\x8c\x64\xca\x31\x80\xbd\x2c\x3f\x5c\xee\xde\x20\xc9\x62\xd4\x62\x11\x48\xb9\x06\xbf\x3e\x81\x73\xfd\x27\x88\x53\x89\xfb\x17\xd9\x6c\xe0\xf0\xd4\xd0\xea\x4b\x81\x09\xdd\x08\x43\xad\xa2\x38\x55\xba\x4f\xc3\x80\xb0\x05\xaf\xab\x63\xd4\x39\xea\x31\x05\x9a\xfe\x5f\x54\xe9\xbf\xb5\x07\xd4\x46\xa8\x64\xff\x2b\xeb\xce\x66\x83\x49\x58\x14\x7f\x03\x00\x00\xff\xff\x9b\x7b\x60\x67\xed\x04\x00\x00"),
		},
		"/template/cmd/server.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "server.tmpl",
			modTime:          time.Date(2018, 12, 22, 14, 59, 53, 259735282, time.UTC),
			uncompressedSize: 476,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x90\x4f\x4b\xc3\x40\x10\xc5\xef\xf9\x14\x8f\xe2\xd1\xf4\xdf\xb1\xd0\x83\x56\x0f\x1e\xda\x08\xb1\xa7\xd2\x43\x9a\x4c\xd3\x40\x76\x67\xd9\x9d\x8a\xb0\xe4\xbb\x4b\x36\x71\x11\xab\x5e\x42\x76\xe6\xbd\x99\xf7\x1b\xef\x2b\x3a\x37\x9a\x30\xc9\xc9\xbe\x93\xdd\xbb\xa2\xa6\x09\xba\x2e\x01\x80\x0d\x2b\x55\xe8\x6a\x15\x1e\x40\xcd\xe6\x42\xb6\x6c\x1b\xb8\x20\xc6\x21\xdf\x3f\x6e\xb2\xed\xf6\x61\xf7\x74\xc4\x21\x7b\x7d\x7b\xc9\x76\xf9\x31\x09\xf2\xfc\x7a\x2a\x07\xbb\xfb\xf2\x3b\x29\xac\xc4\x7f\x36\x83\x30\x33\xd2\xb0\x8e\xa2\x34\x35\xf7\xfd\x87\xad\x40\x5f\xd5\x89\x6c\x98\x46\x02\xb9\x10\x42\x59\x18\x6d\xe3\x84\x34\x58\x4f\x87\x21\xcf\x1f\x85\x32\x2d\xb9\x15\xbc\x87\x90\x32\x6d\x21\x11\x2a\xef\xf7\xf6\x50\xbf\x37\xd9\x04\x60\xef\x49\x57\x5d\x97\x24\x3f\x6f\x12\xed\x63\xc2\xbb\xdb\x43\x7c\x27\xfb\xab\x3d\x42\xad\x17\xf3\xe5\x62\xbe\x44\x9a\x56\x5c\x5a\x66\x59\x4f\x67\x25\xeb\x73\x53\xcf\xc6\x42\x4c\x72\x1b\x64\x8c\xfa\xcf\x22\x36\xd1\xfe\x19\x00\x00\xff\xff\x70\xac\xa9\x14\xdc\x01\x00\x00"),
		},
		"/template/cmd/version.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "version.tmpl",
			modTime:          time.Date(2018, 12, 23, 21, 6, 24, 796610039, time.UTC),
			uncompressedSize: 435,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x90\x5f\x4b\xc3\x30\x14\xc5\xdf\xfb\x29\x0e\xc5\x07\x05\xf7\x05\x0a\x3e\x4d\x91\x81\x58\xc1\x3f\x2f\xe2\x43\x68\xef\xd6\x40\x9a\x84\xdc\xdb\xa9\x84\x7e\x77\x69\xd2\x3a\x65\x83\x41\x48\x08\xe7\xdc\x1f\xe7\x9e\x18\x5b\xda\x6a\x4b\x28\xdf\x28\xb0\x76\xf6\x95\xd5\x8e\x4a\x8c\x63\x01\x00\x6b\xd7\xf7\xca\xb6\x55\xfa\x00\x3b\xe7\x3b\x0a\x8d\xd1\xd8\x67\x37\xde\xeb\xa7\x97\x4d\xfd\xf8\xfc\x51\x24\x4b\x8c\x42\xbd\x37\x4a\x08\xe5\x7d\x32\xaf\x1f\x36\x07\xe4\x29\xe6\x4c\xca\x5a\xed\x45\x3b\xcb\x8b\xb6\x5a\x71\x73\x3d\xdd\x9d\xfb\x6c\x8c\x26\x2b\x37\x12\x06\x02\x70\xab\xd9\x1b\xf5\x0d\xe9\x08\x59\xf9\x8d\x34\x9d\x41\xfc\x20\xb8\x6c\x69\xab\x06\x23\xd5\x34\x74\x75\x60\xf2\xc2\x64\x0a\x7b\x0a\xa7\x98\x59\x39\xcf\x4c\xd0\xbb\x2f\xd5\x7b\x43\x5c\xfd\xdb\x7f\x2e\x74\x7e\xca\x71\x2c\x62\x24\xdb\x4e\x35\x1c\xb5\xfe\xc7\x94\x53\x5e\x1c\x77\xbd\x8c\xff\x04\x00\x00\xff\xff\x06\x16\xdf\x23\xb3\x01\x00\x00"),
		},
		"/template/server": &vfsgen۰DirInfo{
			name:    "server",
			modTime: time.Date(2019, 1, 7, 23, 18, 1, 233030457, time.UTC),
		},
		"/template/server/html": &vfsgen۰DirInfo{
			name:    "html",
			modTime: time.Date(2019, 1, 7, 23, 18, 1, 233030457, time.UTC),
		},
		"/vfsgen_templates.go": &vfsgen۰CompressedFileInfo{
			name:             "vfsgen_templates.go",
			modTime:          time.Date(2019, 1, 12, 21, 43, 27, 572602187, time.UTC),
			uncompressedSize: 1153,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x52\x4b\x6f\xdb\x38\x10\x3e\x8b\xbf\x62\x56\x17\xdb\x58\x2f\x79\x37\x90\xcb\xe6\x85\x00\x45\xd2\x83\xdb\x1e\x03\x8a\x1a\x51\xac\x29\x52\x20\x47\x31\x8c\x22\xff\xbd\x18\x4a\x76\x60\xfb\x64\x71\x1e\xdf\x63\x3e\xa5\xe0\x97\xf3\x1e\x2c\x06\x4c\x9a\x10\x34\xac\x08\x87\xd1\x6b\xc2\xfc\x7e\x7e\x95\x36\xae\xe0\xe8\xa8\x07\xed\x3d\xc4\x0e\xa8\x47\xe8\x9c\xc7\x0c\x53\x68\x31\x01\xf5\x2e\x43\x17\x7d\x8b\x49\x28\x05\x7b\xfe\x74\x19\x02\x1a\xcc\x59\xa7\x13\x34\x68\xf4\x94\x79\x7d\xe3\x02\x3f\x8c\x29\xda\xa4\x07\x30\x3a\x40\x83\x90\xa6\x00\x5d\x8a\x03\xe8\x70\x3a\xf6\x98\x10\x62\xf8\x42\x39\x65\x42\x2e\xb5\xbc\x7c\xd0\x27\x08\x91\xa0\xd7\x1f\xbc\x2f\xa1\xd7\xe4\x3e\x70\x81\x87\x95\x54\x26\x86\xce\x59\x75\xd6\xa1\x56\x12\xe0\x47\x76\xc1\xc2\x47\x97\x2d\x06\x38\x22\x98\x84\xb3\xdc\x4c\x9a\x9c\x01\x1b\x0b\x14\x03\x14\xa1\x26\x06\xc2\x40\xf9\xac\xf6\x62\x0a\xe0\xd0\x60\xdb\x62\x2b\xe1\xa2\xb3\x8d\x01\xe7\xb1\x66\x72\xbe\x05\xd2\x36\x4b\x31\x6a\x73\xd0\x16\x61\xd0\x2e\x08\xa1\x94\x8d\xbb\x8b\xcd\x36\x16\xc9\x33\x9f\xf7\xcb\x72\x69\xa3\x10\x6e\x18\x63\x22\x58\x8b\xaa\xb6\x8e\xfa\xa9\x91\x26\x0e\x2a\xf7\x53\x32\x31\x7e\x53\xf3\x4c\x7d\x53\x75\x69\x1a\x33\x06\xe5\xa3\x4d\x53\xe6\x6a\x40\x52\x3d\xd1\xc8\xff\x63\x79\xc9\x94\x5c\xb0\xb9\x16\x1b\x21\xba\x29\x98\xc2\x6c\xbd\x81\x3f\xa2\x8a\x13\x8d\x13\x3d\x39\x8f\x41\x0f\x08\xbb\x3b\xa8\xc7\x83\xbd\x75\xf2\x2a\x11\xb5\x10\x95\x52\xf0\xfa\xb6\x7f\xdc\xc1\x4b\xc7\xa6\x5e\xae\xc8\x8e\xbd\x3c\x3c\x2e\x99\x81\xe4\x6c\x4f\xff\x19\xef\xcc\x01\xe2\x94\xc0\x1c\xdb\xf5\x86\x8d\x73\x21\xbb\x16\xd9\xe3\x5b\x24\x51\x99\x63\xbb\x85\x77\xa6\x12\xb3\x7c\x46\xe2\x19\x51\x9d\x1b\x9e\xe6\x6b\x33\xd1\x65\x74\xe1\x73\xdf\xa3\x39\x80\x63\x3e\xab\x54\x28\x05\xbe\xfc\x82\xc4\xc4\x6e\xa0\x96\xe0\x6c\x39\x5f\xa0\xdb\xdf\x53\xa6\xaf\x50\x8d\x9a\xfa\x2c\x45\xe5\x3a\x58\xdc\x93\xf7\x31\x90\x76\x21\xaf\x0b\xc1\xfa\x92\xb2\xba\x18\x79\x4b\xf0\x0e\x6a\xa9\x6a\x51\xdd\x3a\xcc\xef\x52\xd5\xf0\x2f\x5c\x17\x44\xf5\x29\x44\x85\xa9\x48\x9b\x4f\x2d\x9f\x17\xd3\xd7\x7c\x4e\xf9\xe0\xd2\xfa\x1a\x64\xb3\x3d\x77\xbe\x8d\xe4\x62\xc8\xcc\xe3\xbc\x70\x07\xfc\xbb\x06\xd9\x8a\xaa\xfa\x3e\xa7\xf3\x75\x6e\x59\x4c\xac\xb9\xf2\x3f\x67\x78\xaf\x6d\x2e\xa3\x75\x42\x8f\x3a\x63\x29\xfd\xd4\xc9\xe9\xc6\x2f\x53\xf5\xfe\x8a\x06\x77\x7c\x6e\x44\x31\x8b\x05\xfc\x73\x07\xc1\xf9\x62\xca\x9c\x4a\xf9\xa4\x49\x7b\x1f\xd6\x98\xd2\xa6\x08\xfd\x14\x7f\x03\x00\x00\xff\xff\x9e\x8f\x06\x74\x81\x04\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/default.gophercli.yaml"].(os.FileInfo),
		fs["/default.test.gophercli.yaml"].(os.FileInfo),
		fs["/template"].(os.FileInfo),
		fs["/vfsgen_templates.go"].(os.FileInfo),
	}
	fs["/template"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/client"].(os.FileInfo),
		fs["/template/cmd"].(os.FileInfo),
		fs["/template/server"].(os.FileInfo),
	}
	fs["/template/client"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/client/table.tmpl"].(os.FileInfo),
	}
	fs["/template/cmd"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/cmd/client.tmpl"].(os.FileInfo),
		fs["/template/cmd/gophercli.tmpl"].(os.FileInfo),
		fs["/template/cmd/server.tmpl"].(os.FileInfo),
		fs["/template/cmd/version.tmpl"].(os.FileInfo),
	}
	fs["/template/server"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/server/html"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
