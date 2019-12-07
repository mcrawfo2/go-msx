package skel

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type staticFilesFile struct {
	data  string
	mime  string
	mtime time.Time
	// size is the size before compression. If 0, it means the data is uncompressed
	size int
	// hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
	hash string
}

var staticFiles = map[string]*staticFilesFile{
	"cmd/app/bootstrap.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffT\x8eAN\xc50\f\x05\xf79\xc5;A*\xc1\xeag\xcd\x16\x84\xc4\tL\xbeK-5\xb1\xe5\x18Do\x8f\xd2.\x80\xed\xd83z\xd2W\xcddV\x12Щq\xc1\xd0ƃ\xfdK*'\xe0Σ\xbaX\x88\xf6\x82'\xe5\x817m\x1c\x9b\xf4\x8f\x04P\x84\xcb\xfbg\xf0\x98>p\x97a;\x1d/gh>\xe2Y\xaa\xebo\x0e0r\xeeQ`;Ū\xdeN\x18\x87\xf1\x1f\x94\xa6\xc0>\x9b\xa6\x1e\x05\xb7\x87\xdbc\x02\xaa\xf6\xe0\xefx\xa5\xd8\n\x96\xffCGPH\xbdN9/9/\x17I?\x01\x00\x00\xff\xff\xa4V\\\xbb\xe2\x00\x00\x00",
		hash:  "100a6a992fc21f005e7fbfe4194f69daf18f86b6fcea3fd7ea529c3e549c5961",
		mime:  "",
		mtime: time.Unix(1574385796, 0),
		size:  226,
	},
	"cmd/app/main.go": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\x92Qk\xdb0\x10ǟ\xadOq\x98=X\xe0(\x1b{\x19\x83>\xa4\x99\x03cMW\xd2n\xd0ǫ\xac\xb8b\xd6I\x93\xce%\xa3\xe4\xbb\x0f9\xea\xca`\xb0<\xf9|\xfa\xdfO\xff\xbbS@\xfd\x03\a\x03\x0e-\ta]\xf0\x91\xa1\x11U\xad=\xb19p\x9dC\xf6\x8b\xc1\xf2\xe3\xf4\xa0\xb4M\xda+\xed\xdd\xf2z\xf3}q\xf9m9\xf8\x85K\x87%\x86p\xa6RcJH}ĥ\xb3CD6g\xd6\xf1\xaf`R-\xa4\x10\xdaS\x9a=b\b\xd7\xe8\f\\@\xfd\xe6\x19CP\x84\xce\x1cg\xcd~\"\r\x96,7\x12\x9eg\xa5\xfaJݓ!nr<Gk\xef\x1cR\xdfBΔ\x9f\xed\xc9S\v\x19\xd0h>@\x99\x83Z\x9f\xbe\x12L\x8c>f迩\xb7\x8c\x91ƠGL\xe6\xd3\x14-\r-`_\xe0\xd6S\x92\xa2\xaa\xa2\xe1)\x12\x90\x1dEu\x94\xe2XL\xff\xa5\xfb\x8f\x03\x87d\xf7&1|\xbc\x802M\xb5-\xb9M\xf4\xaeTd\x8a\x14\xe2\xe5\xc6y\x92\xaaː+\x9b8w\xf2\x02R\xab\xbe_\xff\x1co9{\xfec\xa3\xa9߫\x0f\xea\xadzW\xb7P\xaf\xa3A6\xb0\xb7110>\x8cf\xce\xee\xba\xd5]\aw\xab˫\xae\x9c5O8N\x06\xb2\x01\xb8\xd9}ޮv\xf7𥻗\xb5lEuT\x1b;\xb2\x89\xcdk\xe3\xf9\x05\xbenk7QS\xf6\x9b%\xbf\x03\x00\x00\xff\xff\x01\x8bOs\xa7\x02\x00\x00",
		hash:  "ced38a8bc08ad65f9f8a1b20f5e59da797c54e77b3e6be27719c13191f4143fe",
		mime:  "",
		mtime: time.Unix(1575736744, 0),
		size:  679,
	},
	"cmd/build/build.go": {
		data:  "package main\n\nimport (\n\t\"cto-github.cisco.com/NFV-BU/go-msx/build\"\n)\n\nfunc main() {\n\tbuild.Run()\n}\n",
		hash:  "1afc43b37b6664090421f8a3bde52704b1365a691a132a3c5d3ccd24bec6cc18",
		mime:  "",
		mtime: time.Unix(1573950111, 0),
		size:  0,
	},
	"cmd/build/build.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x94\x8f\xb1j\x031\f\x86\xf7{\n\x911`S\xe8R\xbc\x15\x9a\xe69\x14E\x0e&\xb6e$\x9dI\u07fe\\{ХC\xbb\t\xfd\x9f>\xf4\xf3\x83iu\xbcTN\v\x00\xb5k\x02\x1cc\x1b\xa5\xe7r{/\x95mK\x00\x02\\D\xdc\\qďV\x97\xa5\xd9cK\x94+\xa3q\x82\xe7\xf8\x12\x9f\x16\x80Qѳh\xfb>\x1b\xa8\xdc\xfdU\xbdd$\xdf]\x9b\x8d\xa4E*F\x12g\xb34\x9b\x05c\x9d\x858\x98\xa3:\xeb\xdf\xc9@\xa2\xfc\x0f\xfc\x8e\xf9\x8e\xbf\xf0=\xcf\xd4\xf3\f\xa5;\xdf\x14\xbdH\x0f$\xdd\xd6\x1a*\xe3u\u007fi\xb2Z\x91\xbe\x17\x0e\xa7\xb7\xf3\xe9k_:\xd5\xf5\xcag\x95uX\x82Ï\xf7x<,\x9f\x01\x00\x00\xff\xffUoT\xb0i\x01\x00\x00",
		hash:  "40c49ccd56b64c9795a3c31f0dd99841004745689f0f0650e02a070481cee89f",
		mime:  "",
		mtime: time.Unix(1575735740, 0),
		size:  361,
	},
	"docker/Dockerfile": {
		data:  "FROM dockerhub.cisco.com/vms-platform-dev-docker/vms-base:latest\n\nEXPOSE ${server.port}\n",
		hash:  "fd0630e562f4a727270500ae7d7e0cd287608497a5f5dead34f5c65c178a3719",
		mime:  "",
		mtime: time.Unix(1575740441, 0),
		size:  0,
	},
	"go.mod": {
		data:  "module cto-github.cisco.com/NFV-BU/${app.name}\n\ngo 1.12\n",
		hash:  "2c6648bda73dc5b5b1cda267ca61dbf56e5135ce61ca321c34588d01e87a6a33",
		mime:  "",
		mtime: time.Unix(1575740932, 0),
		size:  0,
	},
	"k8s/kubernetes-deployment.yml.tpl": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xacV[o\xe28\x14~ϯ\xb0\xd0H}\n\xe9tw\xaa\x95\xa5>\x8c(;\x9d\xed\x05\xc4eV\xa3\xaa\x8a\x8cs\n\x1e|[\xdbI\x8bX\xfe\xfb\xc8\x0e\x1d\x92\x10\x10#\x15^\xc2w\x8e\xbf\xf3\x9dK|\x88\xe38\"\x9a}\x03c\x99\x92\x18\x11\xadmR|\x8c\x96Lf\x18]\x83\xe6j%@\xbaH\x80#\x19q\x04G\bI\"\xc0jB\x01\xa3B\xd8-\x80ч5Ѻ\xeb\x9f7\x91\xd5@\xbd\xab\x01\xcd\x19%\x16\xa3\xf5\x1ae\xbf\xe8R\xa12HA\x16\x8f\r\xec\u007f\xae^\xc0<=\x9em\x0f\xa6T\xe5ҝ==\x9eU\xd8Ϟ\xd0f\x13!d\x81\x03u\xca\xf8@\b\t\xe2\xe8\xe2\x8è\xdb\x12@>\x9b\xba\xac\x12\x9e\x1b\x95\a\xc32\x9f\x81\x91\xe0\xc0v\x03\xf6\xe6@\x95\xb49\x8f\xe7\xcaZ\xa61\"\x9c\xab\x97\b!\aBs\xe2`\x1b\xafR\x11\xff\xd9/B\x89\U000da903\xa2N\x90uPX\xa0\x95R9☒\x95X\xcf<\a\xe9f\xccu\x99J41\x16\fFs%\xeckpy\xebҖ\xda\x11&\xc1T\x8e\xc7ۤʰ\xbf`\x84\x98 s\b=-Mi\x00\xd0f\x83wPQ\xceT٪]\x02B\x10\x99\xe1\n\xe4\xc3\xec\x05\xf0 \x99\xfb\xc1\xabc\xf1\x8c\xc9\xec\xea\xbc\x1b\xbeM\x1b\xe5\f\xa4;d\xf5\xbd\xa2 \x1d\x98\xab\x9dȌ\xd6\xf5\x05O\x03ά\xe2\x1f\x8aɫҭk\xc1\x14\x8cB\xb7Ug`\x8e3f\xae\x92Ҟx`O\x9c\x92\xcfl^s+\xa1\x8ac\xa1x.\xe0ޏ\xbcm\x96HxtH\xdc\x02\xa3\x83\xe7wc\xa8\x17l\xaf\x8bm3Wme\xc5^\xebg\x15om*g\x05H\xb0vh\xd4\f\xea\xba\x17\xce\xe9/\xe0pC\xa4\x0ey|X\xfb\xba\x82\xf1eu\xf0\x1a\xb2\xdb$$\x13L&ē6O)\xe3*\xa7\xfc\xcfz\xef\x98d\x8e\x11~\r\x9c\xac\xc6@\x95\xcc,F\u007f\x9c\xd7GA\x83a*\xab\x98+V\x03$c\xef\x9d\xca\x02\bw\x8bw\xc9\xe5\xf2wR\xb1*7\x14\x1a\x83d\xe0\xbf\x1cls\xbc\xfcm&\x94YaԹ\xfc\xf3\x9eu\x1aF\xaas\x8c:\x1f\xeb0g\x82\x1d\xe1\xb9\xf8ty\x88\xe8\xa2\n\x83,\x9a\xa3^\x8e\xebx8\xfa\xfa\xf0%\xed\xdd\r\xa6\xd7io\xf00\x9eޥ7\x83\xf1\xa4AY\x10\x9e\x03F\x9d\xf5\x1a\x15$\xe7.\xb5t\x01\"\xccn\x92pE\t_(\xeb:'\x87\x18\x0eF\x87B\xfc\xf5\xe9\xfc\xfc\x04\xa2o\x9f\xa7w\x93cR\x83\xceƝr2\xef1}\x17\xbf\xa1oܻ\xe9\xdf\xf7O.\xe6ɼ\x93\xc1m\xff\xa1\x8d\xf6o\xa3DsZ\xfc\x0e\xa7\x06\xdc-\xacF\xf0\xbco}\xbb΄}\rzZ\x1c\x96\xb0\xc2ȩ%\xc8v\x85\xfd\xdet\xf4u\xf2=\xbd\xed\u007f\x1f\xa7\xff\xfc;\t\x0f\x83\xe9\xa8w w\r\xe2t\xa2\xe1\xe7\xc9M;M\xb2\x84\x95u\xca@\xf2\xe3\xc5\xc5:\x9f-aխs\xfbw\u07b6l\xc2r\r\x0f\x8f_\x10\xa7\xae\x8a7\x19\xad[\xa2a,9[\xf6\u007fu\x9b\x84\u007f\n\xcfl~Ot=\xee\xe1\xc5Ӣ\xc1\xbf\x93Aa\xb4\u007f\x97\x86\xf5\x99\x14\xc2\xeej\x18\xfd\f\x00\x00\xff\xffHJ\xf9\x1d\xa8\n\x00\x00",
		hash:  "144104093fb1c086e185ef0553a4d1329353e9de49976264b5cc2c78f590e90f",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1575733217, 0),
		size:  2728,
	},
}

// NotFound is called when no asset is found.
// It defaults to http.NotFound but can be overwritten
var NotFound = http.NotFound

// ServeHTTP serves a request, attempting to reply with an embedded file.
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	f, ok := staticFiles[path]
	if !ok {
		if path != "" && !strings.HasSuffix(path, "/") {
			NotFound(rw, req)
			return
		}
		f, ok = staticFiles[path+"index.html"]
		if !ok {
			NotFound(rw, req)
			return
		}
	}
	header := rw.Header()
	if f.hash != "" {
		if hash := req.Header.Get("If-None-Match"); hash == f.hash {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("ETag", f.hash)
	}
	if !f.mtime.IsZero() {
		if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && f.mtime.Before(t.Add(1*time.Second)) {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("Last-Modified", f.mtime.UTC().Format(http.TimeFormat))
	}
	header.Set("Content-Type", f.mime)

	// Check if the asset is compressed in the binary
	if f.size == 0 {
		header.Set("Content-Length", strconv.Itoa(len(f.data)))
		io.WriteString(rw, f.data)
	} else {
		if header.Get("Content-Encoding") == "" && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			header.Set("Content-Encoding", "gzip")
			header.Set("Content-Length", strconv.Itoa(len(f.data)))
			io.WriteString(rw, f.data)
		} else {
			header.Set("Content-Length", strconv.Itoa(f.size))
			reader, _ := gzip.NewReader(strings.NewReader(f.data))
			io.Copy(rw, reader)
			reader.Close()
		}
	}
}

// Server is simply ServeHTTP but wrapped in http.HandlerFunc so it can be passed into net/http functions directly.
var Server http.Handler = http.HandlerFunc(ServeHTTP)

// Open allows you to read an embedded file directly. It will return a decompressing Reader if the file is embedded in compressed format.
// You should close the Reader after you're done with it.
func Open(name string) (io.ReadCloser, error) {
	f, ok := staticFiles[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	if f.size == 0 {
		return ioutil.NopCloser(strings.NewReader(f.data)), nil
	}
	return gzip.NewReader(strings.NewReader(f.data))
}

// ModTime returns the modification time of the original file.
// Useful for caching purposes
// Returns zero time if the file is not in the bundle
func ModTime(file string) (t time.Time) {
	if f, ok := staticFiles[file]; ok {
		t = f.mtime
	}
	return
}

// Hash returns the hex-encoded SHA256 hash of the original file
// Used for the Etag, and useful for caching
// Returns an empty string if the file is not in the bundle
func Hash(file string) (s string) {
	if f, ok := staticFiles[file]; ok {
		s = f.hash
	}
	return
}
