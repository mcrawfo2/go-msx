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
			modTime: time.Date(2021, 10, 21, 20, 25, 35, 789069243, time.UTC),
		},
		"/app": &vfsgen۰DirInfo{
			name:    "app",
			modTime: time.Date(2021, 10, 21, 20, 25, 35, 784785654, time.UTC),
		},
		"/app/defaults-app.properties": &vfsgen۰FileInfo{
			name:    "defaults-app.properties",
			modTime: time.Date(2021, 10, 6, 19, 26, 35, 190261628, time.UTC),
			content: []byte("\x70\x72\x6f\x66\x69\x6c\x65\x3d\x64\x65\x66\x61\x75\x6c\x74\x0a\x62\x61\x6e\x6e\x65\x72\x2e\x65\x6e\x61\x62\x6c\x65\x64\x3d\x66\x61\x6c\x73\x65\x0a"),
		},
		"/certificate": &vfsgen۰DirInfo{
			name:    "certificate",
			modTime: time.Date(2021, 8, 19, 20, 39, 15, 272677276, time.UTC),
		},
		"/certificate/defaults-certificate.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-certificate.properties",
			modTime:          time.Date(2021, 8, 19, 20, 38, 39, 630342884, time.UTC),
			uncompressedSize: 373,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\x08\x4e\x2d\x2a\x4b\x2d\x52\x70\x4e\x2d\x2a\xc9\x4c\xcb\x4c\x4e\x2c\x49\xe5\x4a\x46\xb0\xf5\x8a\xf3\x4b\x8b\x92\x53\xf5\x8a\xc1\xaa\xf4\x0a\x8a\xf2\xcb\x32\x53\x52\x8b\x14\x6c\x15\xd2\x32\x73\xf0\xa9\x04\xc9\xe8\x82\xd4\x28\xd8\x2a\xa8\x54\x43\x05\x4b\x72\x8a\x11\x12\xb5\x78\x74\x67\xa7\x56\x62\xd5\x0c\x13\xaf\xe5\xe2\x52\x56\x70\xce\xc9\x4c\xcd\x2b\x21\xe4\xf2\x64\xb0\x2a\x62\x5c\x0e\x55\x89\xea\xf2\x8c\x92\x92\x02\x0c\x19\xac\x4e\x87\x2a\x42\x71\x3a\xb2\x6e\x84\xdb\x01\x01\x00\x00\xff\xff\xb7\xa8\x3d\x54\x75\x01\x00\x00"),
		},
		"/consul": &vfsgen۰DirInfo{
			name:    "consul",
			modTime: time.Date(2021, 10, 18, 15, 13, 58, 351839537, time.UTC),
		},
		"/consul/defaults-consul.properties": &vfsgen۰FileInfo{
			name:    "defaults-consul.properties",
			modTime: time.Date(2021, 8, 19, 20, 38, 39, 663877492, time.UTC),
			content: []byte("\x73\x70\x72\x69\x6e\x67\x2e\x63\x6c\x6f\x75\x64\x2e\x63\x6f\x6e\x73\x75\x6c\x2e\x63\x6f\x6e\x66\x69\x67\x2e\x77\x61\x74\x63\x68\x2e\x77\x61\x69\x74\x2d\x74\x69\x6d\x65\x20\x3d\x20\x31\x35\x0a"),
		},
		"/discovery": &vfsgen۰DirInfo{
			name:    "discovery",
			modTime: time.Date(2021, 8, 19, 20, 39, 15, 278117913, time.UTC),
		},
		"/discovery/consulprovider": &vfsgen۰DirInfo{
			name:    "consulprovider",
			modTime: time.Date(2021, 8, 19, 20, 39, 40, 223868931, time.UTC),
		},
		"/discovery/consulprovider/defaults-discovery.properties": &vfsgen۰FileInfo{
			name:    "defaults-discovery.properties",
			modTime: time.Date(2020, 4, 3, 20, 22, 7, 0, time.UTC),
			content: []byte("\x73\x70\x72\x69\x6e\x67\x2e\x63\x6c\x6f\x75\x64\x2e\x63\x6f\x6e\x73\x75\x6c\x2e\x64\x69\x73\x63\x6f\x76\x65\x72\x79\x2e\x69\x6e\x73\x74\x61\x6e\x63\x65\x49\x64\x3d\x75\x75\x69\x64\x0a"),
		},
		"/fs": &vfsgen۰DirInfo{
			name:    "fs",
			modTime: time.Date(2021, 8, 19, 20, 39, 15, 281916234, time.UTC),
		},
		"/fs/defaults-fs.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-fs.properties",
			modTime:          time.Date(2021, 8, 19, 20, 38, 39, 667551255, time.UTC),
			uncompressedSize: 356,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x8f\x41\x0e\x84\x20\x0c\x45\xf7\x9e\x82\x85\x6b\x39\x81\x87\xa9\x88\xa4\x09\x14\x42\x71\x36\xc6\xbb\x4f\x0a\x44\x9c\xd5\xec\xda\xd7\xff\xdb\xdf\x83\x97\x1c\x63\x51\xab\xd2\x93\xd4\x96\xe3\x99\x8d\x65\x01\x1f\xc8\xda\xe3\xa6\xe7\x8b\x53\x46\x72\x0b\xa4\xe4\xd1\x40\xc1\x48\x0b\x41\xb0\xb7\x58\x4c\xa4\x03\x5d\x35\xd8\x62\xfe\x88\x37\x24\xc8\xd8\xd6\x9f\x9c\xf5\x86\x24\xd8\x47\x03\x5e\x58\x2d\x84\x70\x01\x87\xe4\x84\xed\xc8\x45\x4b\xc8\x76\x2d\x04\xa0\x5d\xb8\x09\xbb\x86\x94\xa6\xfe\x82\x84\xf7\x16\xd8\xaa\x55\xcd\x57\x87\xf7\x98\x8e\xc7\xea\xb4\xb7\x6f\xc1\x73\xf2\xb1\x0f\x59\xf3\x34\xc5\xcb\x33\xe2\xfc\x2c\xad\x4d\x9f\xdd\xd3\x37\x00\x00\xff\xff\xa6\x8a\xe1\x8d\x64\x01\x00\x00"),
		},
		"/integration": &vfsgen۰DirInfo{
			name:    "integration",
			modTime: time.Date(2021, 10, 21, 19, 54, 41, 131598364, time.UTC),
		},
		"/integration/defaults-integration.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-integration.properties",
			modTime:          time.Date(2021, 9, 17, 12, 26, 14, 891494228, time.UTC),
			uncompressedSize: 180,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x4a\xcd\xcd\x2f\x49\x2d\x4e\x2d\x2a\xcb\x4c\x4e\xd5\x2b\x2d\x4e\x2d\xca\x4d\xcc\x4b\x4c\x4f\xcd\x4d\xcd\x2b\x81\x89\x42\x69\x5b\xac\xb2\x5c\xa8\x26\x24\x96\x96\x64\x90\xa3\xaf\x38\x35\xb9\x28\xb5\xa4\x98\x28\xad\x80\x00\x00\x00\xff\xff\xd0\x7b\x4a\x40\xb4\x00\x00\x00"),
		},
		"/leader": &vfsgen۰DirInfo{
			name:    "leader",
			modTime: time.Date(2021, 8, 19, 20, 39, 40, 231599668, time.UTC),
		},
		"/leader/defaults-leader.yml": &vfsgen۰CompressedFileInfo{
			name:             "defaults-leader.yml",
			modTime:          time.Date(2020, 3, 18, 13, 9, 35, 0, time.UTC),
			uncompressedSize: 177,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xce\x41\x0a\xc2\x30\x10\x85\xe1\x7d\x4e\x31\x8b\x6e\x4d\xf7\xb9\x82\x08\x5e\x61\x4c\x5e\x25\x38\x26\x61\x26\x15\xa4\xf4\xee\x12\xba\x93\xee\xde\xea\xfb\x5f\xac\xc5\x56\xf1\x02\x4e\x50\x0f\x41\xec\xb9\x96\xe0\x88\x50\xf8\x21\x48\x81\x16\x16\x83\x23\x4a\x58\x78\x95\x7e\x63\xeb\xd0\x2b\xbe\x81\x0c\xfa\xc9\x11\xf3\xb4\x59\xd3\x5c\x9e\x9e\x5b\x93\x1c\x79\x10\xbe\xf0\x1b\xfb\x7c\xc0\x8e\xe8\x18\x77\xad\x0d\xda\x33\x6c\x24\x88\x2e\xf4\x1a\xd0\xb4\x9d\xff\xf0\xff\xcd\xdd\xfd\x02\x00\x00\xff\xff\x1f\x71\x63\xb6\xb1\x00\x00\x00"),
		},
		"/populate": &vfsgen۰DirInfo{
			name:    "populate",
			modTime: time.Date(2021, 8, 19, 20, 39, 15, 295006831, time.UTC),
		},
		"/populate/defaults-populate.properties": &vfsgen۰FileInfo{
			name:    "defaults-populate.properties",
			modTime: time.Date(2020, 9, 29, 13, 0, 4, 0, time.UTC),
			content: []byte("\x70\x6f\x70\x75\x6c\x61\x74\x65\x2e\x72\x6f\x6f\x74\x20\x3d\x20\x2f\x70\x6c\x61\x74\x66\x6f\x72\x6d\x2d\x63\x6f\x6d\x6d\x6f\x6e\x0a"),
		},
		"/security": &vfsgen۰DirInfo{
			name:    "security",
			modTime: time.Date(2021, 8, 19, 20, 39, 40, 247064430, time.UTC),
		},
		"/security/idmdetailsprovider": &vfsgen۰DirInfo{
			name:    "idmdetailsprovider",
			modTime: time.Date(2021, 8, 19, 20, 39, 40, 243315263, time.UTC),
		},
		"/security/idmdetailsprovider/defaults-security-idmdetailsprovider.properties": &vfsgen۰FileInfo{
			name:    "defaults-security-idmdetailsprovider.properties",
			modTime: time.Date(2021, 8, 19, 20, 39, 15, 307459485, time.UTC),
			content: []byte("\x73\x65\x63\x75\x72\x69\x74\x79\x2e\x74\x6f\x6b\x65\x6e\x2e\x64\x65\x74\x61\x69\x6c\x73\x2e\x61\x63\x74\x69\x76\x65\x2d\x63\x61\x63\x68\x65\x2e\x74\x74\x6c\x3a\x20\x35\x73\x0a"),
		},
		"/sqldb": &vfsgen۰DirInfo{
			name:    "sqldb",
			modTime: time.Date(2021, 10, 15, 13, 45, 4, 756034411, time.UTC),
		},
		"/sqldb/defaults-sqldb.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-sqldb.properties",
			modTime:          time.Date(2021, 8, 19, 20, 39, 15, 348888103, time.UTC),
			uncompressedSize: 370,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\xc1\x0a\x83\x30\x0c\x86\xef\x7d\x8a\x1e\xbc\x5a\x41\xd8\x06\x05\xd9\x5e\xa5\xb6\x41\x65\xd5\x74\x49\xdd\x0e\xe2\xbb\x8f\x2a\x1e\x86\xdd\xed\xa7\x5f\xd2\x2f\xfc\xae\x55\x16\xed\x93\xd0\xd8\x5e\xf5\xc8\x51\x36\xd2\xa3\x35\x3e\x65\xf1\x43\x03\x52\xa2\xf5\xb5\xbe\xdc\x84\xe0\x40\xc3\xd4\x29\x67\xa2\x61\x9c\xc9\x82\x72\x34\xbc\x81\xb4\x0c\xc8\xb1\x23\xe0\xcc\xc8\x64\x46\xd0\x99\xf7\x99\x81\x36\x26\x09\x31\x66\x06\x82\x61\xfe\x20\xb9\xdc\x72\x8a\xe5\x9e\xcb\xfd\x93\xe3\x82\x97\xd7\x55\x55\x2c\xff\x7d\xab\xce\xd1\x43\xb6\x3e\x8a\xe5\xd4\x4f\x5a\x39\xd5\xb2\x66\x2d\x9b\xe1\xce\xec\x47\x74\xd0\xb8\x81\x4d\xeb\x41\x7c\x03\x00\x00\xff\xff\x35\xea\x73\xd3\x72\x01\x00\x00"),
		},
		"/transit": &vfsgen۰DirInfo{
			name:    "transit",
			modTime: time.Date(2021, 10, 14, 18, 14, 30, 609782414, time.UTC),
		},
		"/transit/defaults-transit.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-transit.properties",
			modTime:          time.Date(2020, 3, 19, 18, 59, 25, 0, time.UTC),
			uncompressedSize: 268,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\xce\xb1\x0e\xc2\x30\x0c\x04\xd0\x9d\xaf\xc8\x0f\xb8\x03\x12\x95\x18\xfa\x31\x6e\x38\x50\x14\xe3\x58\x8e\x51\x9b\xbf\x47\x2c\x4c\x48\xc0\x7e\xef\xee\x0c\x4e\x01\x65\x0d\x82\x66\x1f\x16\xa5\xe9\x04\xe5\x55\x70\x49\x4b\xba\xb2\x74\x1c\x3e\xa7\x58\x36\x1e\x9d\xb2\x83\x03\x54\x31\xfa\x17\x50\x31\xc8\xbc\x19\x3c\x0a\xfa\x14\xc3\x90\x96\xc4\xe8\xc7\xd3\x4c\xb7\x7c\x3f\xcf\xbf\x41\xec\xd6\x3c\x5e\x17\xff\x1b\x64\x91\xb6\x91\x09\x17\x0d\xec\x41\x2b\xe7\xfa\xb0\x77\xc7\x33\x00\x00\xff\xff\x9f\x8e\xd6\xc3\x0c\x01\x00\x00"),
		},
		"/webservice": &vfsgen۰DirInfo{
			name:    "webservice",
			modTime: time.Date(2021, 10, 21, 19, 54, 41, 131742241, time.UTC),
		},
		"/webservice/defaults-security-actuator.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-security-actuator.properties",
			modTime:          time.Date(2021, 8, 19, 20, 39, 40, 283594360, time.UTC),
			uncompressedSize: 987,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\x41\xab\xda\x40\x14\x85\xf7\xf9\x15\x82\x5d\x15\x12\xec\xae\x04\xb2\x48\x4d\xc0\x80\x9a\xa2\x96\x42\x37\xe1\x9a\x39\x26\x03\x93\x19\x99\x3b\xa9\xca\xc3\xff\xfe\xe0\xe9\xf3\x89\x26\xa8\xbb\x61\xce\x3d\xdf\xb9\x87\x3b\x1c\x2c\x51\xb6\x56\xba\x83\xd7\x90\xa6\x0a\x0d\xb4\x0b\xf8\xfc\x17\x40\xd3\x5a\x41\xfc\x3a\x24\xd8\x50\xab\x5c\xe4\x6c\x8b\xce\x49\x0b\x36\xad\x2d\x91\x89\xe8\xdb\xdb\xdd\x2f\x5f\x5e\xbe\x14\xe1\xb1\x93\xb0\x85\x6d\x24\xb3\x34\x9a\xa3\x6c\x59\xc4\xbf\xb3\x22\x4e\x66\xd9\xbc\x3b\xce\x28\x70\xb4\xc8\xa7\x69\x31\x9e\x66\xe9\x7c\xd5\xb3\xbe\xd8\x1a\xa9\x5d\x40\x4a\xfe\xc7\x67\x9b\x68\x43\x8a\xbb\x5b\xdc\x18\x58\x2a\xe8\x12\xa2\xbf\xf6\xc5\x50\x83\x94\xab\x5f\x89\x38\x3b\x5e\xc8\x90\x7a\x63\x6e\x12\xae\x0d\x77\xe4\xda\xec\x12\x38\x92\x8a\xa3\xbf\x93\x74\x5e\xc4\x7f\x56\x93\x7c\x91\xfd\x4b\x13\xcf\x1b\x0e\xc6\xf9\x62\xd9\x65\xe7\x60\x87\x75\x50\x1a\xcb\x01\x29\x65\x76\x10\xb9\x95\x95\xd4\x7c\x7d\xd8\x6b\xd9\x37\x27\x3d\xfc\x7e\x7c\x8e\x37\x83\xab\x8d\xe8\xe7\x35\x27\xfd\x69\xde\x04\x24\x60\xfb\x79\xf5\x49\x7f\xcc\xc3\x7e\x6b\xb8\x9f\x77\x96\x2f\xbc\xa7\xd6\x1b\x5b\x08\x68\x27\x49\xf5\x2c\xe8\x97\x5f\x13\xe1\xc7\x55\x1f\x71\x1b\xda\xc7\x15\xee\x68\x0d\xed\x7d\xaa\x10\xfe\xf8\x39\x1a\xf1\xd1\x7b\x0f\x00\x00\xff\xff\xde\xef\x7a\xa4\xdb\x03\x00\x00"),
		},
		"/webservice/defaults-webservice.properties": &vfsgen۰CompressedFileInfo{
			name:             "defaults-webservice.properties",
			modTime:          time.Date(2021, 8, 19, 20, 39, 40, 284174759, time.UTC),
			uncompressedSize: 830,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x92\x5f\x6f\x9b\x30\x14\xc5\xdf\xf3\x29\x2c\x6d\x8f\x31\xa3\x64\xa9\x78\xe1\x81\xd2\x68\x54\xea\xb4\x6a\x44\xda\xf6\x64\xb9\xe6\x52\xac\xba\xb6\x65\x5f\xca\xa2\x29\xdf\x7d\xe2\x5f\xd7\xa0\x90\x89\x17\x7c\xce\xef\x9e\x7b\x2c\xf0\xe0\x5e\xc1\x05\x1e\x39\x4a\x41\x2d\xc7\x9a\x24\xe4\x53\xdb\xb6\xab\xd1\x41\xc7\x05\x50\xd0\xfc\x51\x41\x49\x12\x52\x71\xe5\x61\x32\x4b\x78\x6c\x9e\x96\x4c\x6b\x1c\x92\x84\xc4\x61\x1c\x4e\x52\x6d\x7c\x27\x85\x41\xff\x4c\xea\x7c\x7e\xf5\x81\x18\x55\x12\x54\x9e\x0c\x04\x11\x46\x57\xf2\xe9\xad\x92\xf2\x81\x00\x87\xb4\x92\x0a\x48\x32\x42\x81\x70\xf8\x9e\x78\x86\xc3\x0c\x78\x86\xc3\x49\x04\x9f\xfc\x8f\x7f\xce\x25\x1f\xbb\x22\x1a\xda\xbe\x88\x75\xe6\x55\x96\x8b\x55\x64\x25\x05\x47\xa0\xde\x34\x4e\xfc\x5b\xd9\x25\x08\x69\x6b\x70\xfe\x64\xa2\x97\xa8\x6f\x24\x82\x27\x09\xd9\xdf\x17\x6c\x97\xdd\xe6\x3b\xf6\xbd\x48\xd9\x8f\xbb\x7d\xce\xb2\x3c\xcd\xf2\x34\x0a\xd9\xc3\xb7\xfb\x5f\x57\x9b\x70\xbb\x3e\x03\xa5\xbb\x82\x5d\x45\x31\xfb\x92\x7d\x65\x45\x9e\x46\xdb\xeb\x25\x2a\xda\x5e\x4f\xd4\x26\xfe\x7c\x89\xca\x6e\xb2\x8e\xea\x91\xff\x45\x2c\x0d\xf7\xd7\x36\xce\x93\x1a\x78\xf9\xee\xee\x9d\x16\x88\xc6\xa3\x79\xa1\x5c\x29\xd3\x42\x49\x47\x64\xfc\x0a\xa2\x71\x12\x0f\x03\x38\x23\x8e\xeb\x9f\xf4\x66\x43\xf7\xdd\xff\x78\x57\x0e\x87\xc2\x72\xfd\xf6\xce\x5f\xac\x82\xf1\xf0\xc0\x1d\x68\x1c\xec\x73\xdb\xe1\xb7\x35\xfe\xe2\xf6\x19\x71\x5c\x67\x46\x23\x68\xa4\xb7\xd2\x5b\xe3\x25\x4a\xa3\x57\x7f\x03\x00\x00\xff\xff\x31\x57\x62\x11\x3e\x03\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/app"].(os.FileInfo),
		fs["/certificate"].(os.FileInfo),
		fs["/consul"].(os.FileInfo),
		fs["/discovery"].(os.FileInfo),
		fs["/fs"].(os.FileInfo),
		fs["/integration"].(os.FileInfo),
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
	fs["/integration"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/integration/defaults-integration.properties"].(os.FileInfo),
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
