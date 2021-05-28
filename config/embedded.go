// Code generated by vfsgen; DO NOT EDIT.

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

// EmbeddedDefaultsFileSystem statically implements the virtual filesystem provided to vfsgen.
var EmbeddedDefaultsFileSystem = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2021, 5, 28, 17, 29, 50, 422178212, time.UTC),
		},
		"/app": &vfsgen۰DirInfo{
			name:    "app",
			modTime: time.Date(2021, 5, 28, 17, 29, 50, 570841745, time.UTC),
		},
		"/app/defaults-app.properties": &vfsgen۰FileInfo{
			name:    "defaults-app.properties",
			modTime: time.Date(2020, 3, 16, 22, 37, 55, 58928482, time.UTC),
			content: []byte("\x70\x72\x6f\x66\x69\x6c\x65\x3d\x64\x65\x66\x61\x75\x6c\x74\x0a"),
		},
		"/certificate": &vfsgen۰DirInfo{
			name:    "certificate",
			modTime: time.Date(2021, 5, 21, 12, 30, 56, 303822517, time.UTC),
		},
		"/certificate/defaults-certificate.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-certificate.properties",
			modTime:          time.Date(2020, 12, 16, 17, 2, 0, 358624584, time.UTC),
			uncompressedSize: 373,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\x08\x4e\x2d\x2a\x4b\x2d\x52\x70\x4e\x2d\x2a\xc9\x4c\xcb\x4c\x4e\x2c\x49\xe5\x4a\x46\xb0\xf5\x8a\xf3\x4b\x8b\x92\x53\xf5\x8a\xc1\xaa\xf4\x0a\x8a\xf2\xcb\x32\x53\x52\x8b\x14\x6c\x15\xd2\x32\x73\xf0\xa9\x04\xc9\xe8\x82\xd4\x28\xd8\x2a\xa8\x54\x43\x05\x4b\x72\x8a\x11\x12\xb5\x78\x74\x67\xa7\x56\x62\xd5\x0c\x13\xaf\xe5\xe2\x52\x56\x70\xce\xc9\x4c\xcd\x2b\x21\xe4\xf2\x64\xb0\x2a\x62\x5c\x0e\x55\x89\xea\xf2\x8c\x92\x92\x02\x0c\x19\xac\x4e\x87\x2a\x42\x71\x3a\xb2\x6e\x84\xdb\x01\x01\x00\x00\xff\xff\xb7\xa8\x3d\x54\x75\x01\x00\x00"),
		},
		"/consul": &vfsgen۰DirInfo{
			name:    "consul",
			modTime: time.Date(2021, 4, 8, 13, 52, 53, 21906239, time.UTC),
		},
		"/consul/defaults-consul.properties": &vfsgen۰FileInfo{
			name:    "defaults-consul.properties",
			modTime: time.Date(2021, 1, 13, 18, 39, 4, 452446588, time.UTC),
			content: []byte("\x73\x70\x72\x69\x6e\x67\x2e\x63\x6c\x6f\x75\x64\x2e\x63\x6f\x6e\x73\x75\x6c\x2e\x63\x6f\x6e\x66\x69\x67\x2e\x77\x61\x74\x63\x68\x2e\x77\x61\x69\x74\x2d\x74\x69\x6d\x65\x20\x3d\x20\x31\x35\x0a"),
		},
		"/discovery": &vfsgen۰DirInfo{
			name:    "discovery",
			modTime: time.Date(2021, 4, 8, 13, 52, 53, 23114761, time.UTC),
		},
		"/discovery/consulprovider": &vfsgen۰DirInfo{
			name:    "consulprovider",
			modTime: time.Date(2021, 5, 28, 17, 29, 50, 421769413, time.UTC),
		},
		"/discovery/consulprovider/defaults-discovery.properties": &vfsgen۰FileInfo{
			name:    "defaults-discovery.properties",
			modTime: time.Date(2020, 4, 3, 20, 22, 7, 409821879, time.UTC),
			content: []byte("\x73\x70\x72\x69\x6e\x67\x2e\x63\x6c\x6f\x75\x64\x2e\x63\x6f\x6e\x73\x75\x6c\x2e\x64\x69\x73\x63\x6f\x76\x65\x72\x79\x2e\x69\x6e\x73\x74\x61\x6e\x63\x65\x49\x64\x3d\x75\x75\x69\x64\x0a"),
		},
		"/fs": &vfsgen۰DirInfo{
			name:    "fs",
			modTime: time.Date(2021, 4, 8, 13, 52, 53, 27739748, time.UTC),
		},
		"/fs/defaults-fs.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-fs.properties",
			modTime:          time.Date(2021, 1, 13, 18, 39, 4, 453903911, time.UTC),
			uncompressedSize: 356,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x8f\x41\x0e\x84\x20\x0c\x45\xf7\x9e\x82\x85\x6b\x39\x81\x87\xa9\x88\xa4\x09\x14\x42\x71\x36\xc6\xbb\x4f\x0a\x44\x9c\xd5\xec\xda\xd7\xff\xdb\xdf\x83\x97\x1c\x63\x51\xab\xd2\x93\xd4\x96\xe3\x99\x8d\x65\x01\x1f\xc8\xda\xe3\xa6\xe7\x8b\x53\x46\x72\x0b\xa4\xe4\xd1\x40\xc1\x48\x0b\x41\xb0\xb7\x58\x4c\xa4\x03\x5d\x35\xd8\x62\xfe\x88\x37\x24\xc8\xd8\xd6\x9f\x9c\xf5\x86\x24\xd8\x47\x03\x5e\x58\x2d\x84\x70\x01\x87\xe4\x84\xed\xc8\x45\x4b\xc8\x76\x2d\x04\xa0\x5d\xb8\x09\xbb\x86\x94\xa6\xfe\x82\x84\xf7\x16\xd8\xaa\x55\xcd\x57\x87\xf7\x98\x8e\xc7\xea\xb4\xb7\x6f\xc1\x73\xf2\xb1\x0f\x59\xf3\x34\xc5\xcb\x33\xe2\xfc\x2c\xad\x4d\x9f\xdd\xd3\x37\x00\x00\xff\xff\xa6\x8a\xe1\x8d\x64\x01\x00\x00"),
		},
		"/leader": &vfsgen۰DirInfo{
			name:    "leader",
			modTime: time.Date(2021, 4, 8, 13, 52, 53, 45082656, time.UTC),
		},
		"/leader/defaults-leader.yml": &vfsgen۰CompressedFileInfo{
			name:             "defaults-leader.yml",
			modTime:          time.Date(2020, 3, 18, 13, 9, 35, 358637102, time.UTC),
			uncompressedSize: 177,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xce\x41\x0a\xc2\x30\x10\x85\xe1\x7d\x4e\x31\x8b\x6e\x4d\xf7\xb9\x82\x08\x5e\x61\x4c\x5e\x25\x38\x26\x61\x26\x15\xa4\xf4\xee\x12\xba\x93\xee\xde\xea\xfb\x5f\xac\xc5\x56\xf1\x02\x4e\x50\x0f\x41\xec\xb9\x96\xe0\x88\x50\xf8\x21\x48\x81\x16\x16\x83\x23\x4a\x58\x78\x95\x7e\x63\xeb\xd0\x2b\xbe\x81\x0c\xfa\xc9\x11\xf3\xb4\x59\xd3\x5c\x9e\x9e\x5b\x93\x1c\x79\x10\xbe\xf0\x1b\xfb\x7c\xc0\x8e\xe8\x18\x77\xad\x0d\xda\x33\x6c\x24\x88\x2e\xf4\x1a\xd0\xb4\x9d\xff\xf0\xff\xcd\xdd\xfd\x02\x00\x00\xff\xff\x1f\x71\x63\xb6\xb1\x00\x00\x00"),
		},
		"/populate": &vfsgen۰DirInfo{
			name:    "populate",
			modTime: time.Date(2021, 5, 21, 12, 30, 56, 310254062, time.UTC),
		},
		"/populate/defaults-populate.properties": &vfsgen۰FileInfo{
			name:    "defaults-populate.properties",
			modTime: time.Date(2020, 9, 29, 13, 0, 4, 731036249, time.UTC),
			content: []byte("\x70\x6f\x70\x75\x6c\x61\x74\x65\x2e\x72\x6f\x6f\x74\x20\x3d\x20\x2f\x70\x6c\x61\x74\x66\x6f\x72\x6d\x2d\x63\x6f\x6d\x6d\x6f\x6e\x0a"),
		},
		"/security": &vfsgen۰DirInfo{
			name:    "security",
			modTime: time.Date(2021, 5, 28, 17, 29, 50, 423616119, time.UTC),
		},
		"/security/idmdetailsprovider": &vfsgen۰DirInfo{
			name:    "idmdetailsprovider",
			modTime: time.Date(2021, 4, 13, 18, 0, 19, 499225393, time.UTC),
		},
		"/security/idmdetailsprovider/defaults-security-idmdetailsprovider.properties": &vfsgen۰FileInfo{
			name:    "defaults-security-idmdetailsprovider.properties",
			modTime: time.Date(2021, 4, 13, 18, 0, 19, 499332101, time.UTC),
			content: []byte("\x73\x65\x63\x75\x72\x69\x74\x79\x2e\x74\x6f\x6b\x65\x6e\x2e\x64\x65\x74\x61\x69\x6c\x73\x2e\x61\x63\x74\x69\x76\x65\x2d\x63\x61\x63\x68\x65\x2e\x74\x74\x6c\x3a\x20\x35\x73\x0a"),
		},
		"/sqldb": &vfsgen۰DirInfo{
			name:    "sqldb",
			modTime: time.Date(2021, 5, 28, 17, 29, 50, 427057634, time.UTC),
		},
		"/sqldb/defaults-sqldb.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-sqldb.properties",
			modTime:          time.Date(2021, 4, 8, 13, 52, 53, 101751411, time.UTC),
			uncompressedSize: 370,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\xc1\x0a\x83\x30\x0c\x86\xef\x7d\x8a\x1e\xbc\x5a\x41\xd8\x06\x05\xd9\x5e\xa5\xb6\x41\x65\xd5\x74\x49\xdd\x0e\xe2\xbb\x8f\x2a\x1e\x86\xdd\xed\xa7\x5f\xd2\x2f\xfc\xae\x55\x16\xed\x93\xd0\xd8\x5e\xf5\xc8\x51\x36\xd2\xa3\x35\x3e\x65\xf1\x43\x03\x52\xa2\xf5\xb5\xbe\xdc\x84\xe0\x40\xc3\xd4\x29\x67\xa2\x61\x9c\xc9\x82\x72\x34\xbc\x81\xb4\x0c\xc8\xb1\x23\xe0\xcc\xc8\x64\x46\xd0\x99\xf7\x99\x81\x36\x26\x09\x31\x66\x06\x82\x61\xfe\x20\xb9\xdc\x72\x8a\xe5\x9e\xcb\xfd\x93\xe3\x82\x97\xd7\x55\x55\x2c\xff\x7d\xab\xce\xd1\x43\xb6\x3e\x8a\xe5\xd4\x4f\x5a\x39\xd5\xb2\x66\x2d\x9b\xe1\xce\xec\x47\x74\xd0\xb8\x81\x4d\xeb\x41\x7c\x03\x00\x00\xff\xff\x35\xea\x73\xd3\x72\x01\x00\x00"),
		},
		"/transit": &vfsgen۰DirInfo{
			name:    "transit",
			modTime: time.Date(2021, 5, 26, 13, 13, 54, 213443855, time.UTC),
		},
		"/transit/defaults-transit.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-transit.properties",
			modTime:          time.Date(2020, 3, 19, 18, 59, 25, 223165396, time.UTC),
			uncompressedSize: 268,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\xce\xb1\x0e\xc2\x30\x0c\x04\xd0\x9d\xaf\xc8\x0f\xb8\x03\x12\x95\x18\xfa\x31\x6e\x38\x50\x14\xe3\x58\x8e\x51\x9b\xbf\x47\x2c\x4c\x48\xc0\x7e\xef\xee\x0c\x4e\x01\x65\x0d\x82\x66\x1f\x16\xa5\xe9\x04\xe5\x55\x70\x49\x4b\xba\xb2\x74\x1c\x3e\xa7\x58\x36\x1e\x9d\xb2\x83\x03\x54\x31\xfa\x17\x50\x31\xc8\xbc\x19\x3c\x0a\xfa\x14\xc3\x90\x96\xc4\xe8\xc7\xd3\x4c\xb7\x7c\x3f\xcf\xbf\x41\xec\xd6\x3c\x5e\x17\xff\x1b\x64\x91\xb6\x91\x09\x17\x0d\xec\x41\x2b\xe7\xfa\xb0\x77\xc7\x33\x00\x00\xff\xff\x9f\x8e\xd6\xc3\x0c\x01\x00\x00"),
		},
		"/webservice": &vfsgen۰DirInfo{
			name:    "webservice",
			modTime: time.Date(2021, 5, 28, 17, 29, 50, 573064784, time.UTC),
		},
		"/webservice/defaults-security-actuator.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-security-actuator.properties",
			modTime:          time.Date(2021, 5, 28, 17, 29, 50, 572353150, time.UTC),
			uncompressedSize: 987,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\x41\xab\xda\x40\x14\x85\xf7\xf9\x15\x82\x5d\x15\x12\xec\xae\x04\xb2\x48\x4d\xc0\x80\x9a\xa2\x96\x42\x37\xe1\x9a\x39\x26\x03\x93\x19\x99\x3b\xa9\xca\xc3\xff\xfe\xe0\xe9\xf3\x89\x26\xa8\xbb\x61\xce\x3d\xdf\xb9\x87\x3b\x1c\x2c\x51\xb6\x56\xba\x83\xd7\x90\xa6\x0a\x0d\xb4\x0b\xf8\xfc\x17\x40\xd3\x5a\x41\xfc\x3a\x24\xd8\x50\xab\x5c\xe4\x6c\x8b\xce\x49\x0b\x36\xad\x2d\x91\x89\xe8\xdb\xdb\xdd\x2f\x5f\x5e\xbe\x14\xe1\xb1\x93\xb0\x85\x6d\x24\xb3\x34\x9a\xa3\x6c\x59\xc4\xbf\xb3\x22\x4e\x66\xd9\xbc\x3b\xce\x28\x70\xb4\xc8\xa7\x69\x31\x9e\x66\xe9\x7c\xd5\xb3\xbe\xd8\x1a\xa9\x5d\x40\x4a\xfe\xc7\x67\x9b\x68\x43\x8a\xbb\x5b\xdc\x18\x58\x2a\xe8\x12\xa2\xbf\xf6\xc5\x50\x83\x94\xab\x5f\x89\x38\x3b\x5e\xc8\x90\x7a\x63\x6e\x12\xae\x0d\x77\xe4\xda\xec\x12\x38\x92\x8a\xa3\xbf\x93\x74\x5e\xc4\x7f\x56\x93\x7c\x91\xfd\x4b\x13\xcf\x1b\x0e\xc6\xf9\x62\xd9\x65\xe7\x60\x87\x75\x50\x1a\xcb\x01\x29\x65\x76\x10\xb9\x95\x95\xd4\x7c\x7d\xd8\x6b\xd9\x37\x27\x3d\xfc\x7e\x7c\x8e\x37\x83\xab\x8d\xe8\xe7\x35\x27\xfd\x69\xde\x04\x24\x60\xfb\x79\xf5\x49\x7f\xcc\xc3\x7e\x6b\xb8\x9f\x77\x96\x2f\xbc\xa7\xd6\x1b\x5b\x08\x68\x27\x49\xf5\x2c\xe8\x97\x5f\x13\xe1\xc7\x55\x1f\x71\x1b\xda\xc7\x15\xee\x68\x0d\xed\x7d\xaa\x10\xfe\xf8\x39\x1a\xf1\xd1\x7b\x0f\x00\x00\xff\xff\xde\xef\x7a\xa4\xdb\x03\x00\x00"),
		},
		"/webservice/defaults-webservice.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-webservice.properties",
			modTime:          time.Date(2021, 1, 7, 15, 34, 17, 154063815, time.UTC),
			uncompressedSize: 598,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x8e\x41\x4f\x83\x40\x14\x84\xef\xfc\x8a\x97\xe8\x51\x90\x52\x69\xb8\x70\xc0\x95\x88\x49\x8d\x46\x9a\x18\x4f\x9b\xed\xf2\x28\x9b\x6e\x58\xb2\xbb\x94\x34\xc6\xff\x6e\xc0\xa2\xb6\x29\x9a\x77\x9b\xf9\x66\xde\x18\xd4\x3b\xd4\x9e\xb1\xcc\x0a\xee\x36\xcc\x56\x10\xc3\x75\xd7\x75\xce\xc1\xb1\x9a\x71\x74\xb1\x66\x6b\x89\x05\xc4\x50\x32\x69\x70\x34\x0b\x5c\xb7\x9b\x29\xb3\x51\xda\x42\x0c\x91\x1f\xf9\xa3\x54\x29\xd3\x4b\xbe\x37\xdc\xa8\x9e\xe6\x9d\x0b\x50\xb2\x00\x2b\x0d\x7c\x11\xc0\x55\x5d\x8a\xcd\xf7\x24\x69\x3c\x8e\xda\xba\xa5\x90\x08\xf1\x01\xf2\xb8\xb6\xbf\x89\x2d\xee\x4f\x80\x2d\xee\x8f\x2a\xd8\xe8\x5f\xbe\x9f\x6b\xfe\xe8\x87\xd4\xd8\x0d\x43\x1a\xad\x76\xa2\x98\x9c\x22\x4a\xc1\x99\x45\xd7\xa8\x56\xf3\x9f\x97\x7d\x03\x17\x4d\x85\xda\x1c\x25\x06\xc9\x35\xad\xb0\x68\x20\x86\xd5\x32\xa7\x29\xb9\xcb\x52\xfa\x92\x27\xf4\xf5\x61\x95\x51\x92\x25\x24\x4b\x02\x9f\x3e\x3f\x2d\xdf\x66\x73\x3f\xbc\x3a\x03\x25\x69\x4e\x67\x41\x44\xef\xc9\x23\xcd\xb3\x24\x08\x17\x53\x54\x10\x2e\x46\x6a\x1e\xdd\xfc\x45\x91\x5b\xd2\x53\x03\xf2\x5f\xc5\x54\xd8\xf9\x0c\x00\x00\xff\xff\x99\x36\x4a\x1e\x56\x02\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/app"].(os.FileInfo),
		fs["/certificate"].(os.FileInfo),
		fs["/consul"].(os.FileInfo),
		fs["/discovery"].(os.FileInfo),
		fs["/fs"].(os.FileInfo),
		fs["/leader"].(os.FileInfo),
		fs["/populate"].(os.FileInfo),
		fs["/security"].(os.FileInfo),
		fs["/sqldb"].(os.FileInfo),
		fs["/transit"].(os.FileInfo),
		fs["/webservice"].(os.FileInfo),
	}
	fs["/app"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/app/defaults-app.properties"].(os.FileInfo),
	}
	fs["/certificate"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/certificate/defaults-certificate.properties"].(os.FileInfo),
	}
	fs["/consul"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/consul/defaults-consul.properties"].(os.FileInfo),
	}
	fs["/discovery"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/discovery/consulprovider"].(os.FileInfo),
	}
	fs["/discovery/consulprovider"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/discovery/consulprovider/defaults-discovery.properties"].(os.FileInfo),
	}
	fs["/fs"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/fs/defaults-fs.properties"].(os.FileInfo),
	}
	fs["/leader"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/leader/defaults-leader.yml"].(os.FileInfo),
	}
	fs["/populate"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/populate/defaults-populate.properties"].(os.FileInfo),
	}
	fs["/security"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/security/idmdetailsprovider"].(os.FileInfo),
	}
	fs["/security/idmdetailsprovider"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/security/idmdetailsprovider/defaults-security-idmdetailsprovider.properties"].(os.FileInfo),
	}
	fs["/sqldb"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/sqldb/defaults-sqldb.properties"].(os.FileInfo),
	}
	fs["/transit"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/transit/defaults-transit.properties"].(os.FileInfo),
	}
	fs["/webservice"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/webservice/defaults-security-actuator.properties"].(os.FileInfo),
		fs["/webservice/defaults-webservice.properties"].(os.FileInfo),
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
	case *vfsgen۰FileInfo:
		return &vfsgen۰File{
			vfsgen۰FileInfo: f,
			Reader:          bytes.NewReader(f.content),
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

// vfsgen۰FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type vfsgen۰FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen۰FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen۰FileInfo) Name() string       { return f.name }
func (f *vfsgen۰FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen۰FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰FileInfo) IsDir() bool        { return false }
func (f *vfsgen۰FileInfo) Sys() interface{}   { return nil }

// vfsgen۰File is an opened file instance.
type vfsgen۰File struct {
	*vfsgen۰FileInfo
	*bytes.Reader
}

func (f *vfsgen۰File) Close() error {
	return nil
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
