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
	"Makefile": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff|\x91;O+1\x10\x85\xeb̯\x98\xe2\x167\xc5\xe4\xf6+\xdd\x02\xc4S\x8a\x00\x05QP!\xaf=\xebX\xf8\xb1\xf2\x83\x90\u007f\x8fv\xe3\x046\x04\x1aK\xf6\xf9f|f\xce\xf9\xd3\xed\xf2\xe2r\x85\xffQ\a\x8c\xc5㟿\xf5\xe9\xe5jyv\xfd8G\xe9Կ\xb6\x18[υ\x0eH$\x83\xef\x8c\xfe\xa6m\x9d\xdd7\xa8\xe5\x00\x8b\x87\x9b\xfb\xbb\xe7\x063\xa7\x8c\xca\fG\x90\xaf\x1c\xb1/\xad5i\x8dҲ\xf0\xd8G\x96\xc19\x93\x01\x06\xb2\x81\xd9\xc1\xc9\x1cU\xd8x\x1b\x84\xa2A\"\xc5}\x9a\xc8\xfcβd\xa6\xe2M\x1e\x91\x040\xfc4m\xa2\xd9s\x14\x99i\xf4J\xc6wa\xa2\x1b\x9f\xb2\xb0\x96v\xddDk\xb9\x8e\x99Nb\x8a{\xf6\x8a\xbd\xdc\xfe\x8a\xa5\x8dК#\x153\x91w\x1e>\xbf\x02P\xdc\x16ݜ\x80Fa\x8a\x8e\vl`\xa6\x03\xba\xa0\xf0\x8d\xbd\n\xf1hc\x03\xb2\x1bu_@u\xe3MUO\x15\xf4%\xad\x01\xf6\xe0\x8f\xc4\xd7\xf7\nS\x9d\x98#9\xe1M\xc7)\x03\x8c\xd960\x8b\x0eiՍ\xf1\x1f.\xd55\x1c\x92?\x8a+P\xe72|\x04\x00\x00\xff\xff\x1cseɡ\x02\x00\x00",
		hash:  "c05a7b85677f4564bcb8583ba76069f71df9f98410910a1e244e4637971b2493",
		mime:  "",
		mtime: time.Unix(1584101878, 0),
		size:  673,
	},
	"README.md": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\x94Ao\xe36\x10\x85\xef\xfa\x15\x03\xa4\a\xbbH\xe4k\xb1(\n\xb4\xc16\x97\xdd\"h\xb6@\x81`\x01\x8fɑĆ\xe2\bC\xd2]\xb7\xe8\u007f/8\x94mi\x93\xc3^\f\x8bz\xef\x893\xf3\x917\xf0ݿ8Mm\xc0\x91\xfek\x9a\x9b\x1bx\xe2,\x86\xe0\x9e-5ͧ\x81\xa0c\xef\xf9o\x17z\xb0N\xc8$\x16G\x11P\b&ᣳd\xa1c\x81\x13g\x01Ö\x00\x83\x05ás}\x16L\x8eû\xa6\xf9\x1e\xf6f\xb4\xfb\x06\x00\xe0\x0e\xde\u007f!\x93\x13\x1e<\x01\x85$\xa7\x89]H\xb1\x05\x98\xdf\xff<M\xde\x19\xf5.\x04ࢦ\xecp\x9av#\xba\xd0\xf6\xbcog\xcb/\xd9y\xfb\x96\xf8P^\xd4߫\xbe\xecǲy!\xd9\xdf\xc2\xfe\xe5\x87x\xde\xd9=\x87\x84.\x90\xb8\u007f\xea\xe7Q\x92\xebФx5z6\xe8φ\x0f\xe5a].t\xceS\xa9\xe6\xfd\x17\xe3\xb3\xf6Gx\x04\xa1\x89\xa3K,\xa7\xb6i~e\x81.K\x1aH\xc0RB\xe7#p\x80\x9eKS\xff\"\x93 &\xc9&e\xa1[\x88D\xd0<?%\f\x16\xc5\xc2\x03\xc3\xe3,\xfa\x80'\xce\xe9\xf3fHi\x8a\xefv\xbbޥ!\x1fZ\xc3\xe3\xaeg\x8f\xa1\xbf\x8b\xb3+\xee\xe6\xe0;\xaf\x9em\xab\xd3ֶ5\xcd\x02\x02ȑ\"<\xfc\xf6\a|\xc4\x17\xd2Ѧ\x81@\x1b\b\xf1\x14\x13\x8d-\xc0\x1a\x8c\x84\xd2S\x8aM\x81\x02\x8f\xe8|\x19m\x9dz\xa2\x98\xf6\xb0\x9a;AY\x8c\x1a\x8d\xde\xc3\xc86\xfb\x82T\xb0\xd0S \xc1D\xda-I\x116\xcf5\xe2\xf3\xa6ݕ?ۭ\xce\xce-R\xeb\xe8˒\xb8C\xd6\t\x9cI=i(]q\xdb<Wo\x89+\u007f\xb6\xdb\xca\x02\x1dr\xffu`Y[x\x17\xcc|%\xd4Ep#\xf6*\x9a\xf2\xc1\xbb8\\T\x8f\xf5y\xa5\xd3}\xb9\x10\x13zO\x02#\x06\xd7QLzN<a8\xe3\xf5;\x8d|\xac\rS˛U:\x8a\xfa]!\xc3\xe3\xe8\xd2\xd9\xfc(4\x95\x91\xe8\xa9,\xdd6\x03\x99\x97;\x17\x9a\xe6\x13ϕ\x11\xe0<\xbeJi\x99u\x1c\xc8\xfb[\x90\x1c`?\x16\b~\xac\x8a\x9f\xf6\x15\x9a\x87\x82\x96-\xb7\x83\x8b\x17^\r\x068\x10\xf0D\x81Ji\xa0\x00\xda[\xdd\xf6\x80\x11\"\x1dIЃ\xa5\x0e\xb3O\x1a\xbf:7\xb1\x12\xb3dq\xa3gm{\xe9F\x0e\xba\xc1\xd1\x19\xe1Hrt\x86\x00{,\x8d\x04\x04\x15\x9f\xbbZ\x8f\"w\xf0\xf1\xe9\xcfW\xb1B#'\xfa\xd6ܪ~+x>\tU\xa0aPj\ndT\xb4\xbe\x16\\\xd4\x1dj\x92~\xee\xd9/.\x8fBd\xad\x16b>\\\x00\x9eC\xb9S\xc7\xdc\xec\xb6\x14\xa4\x93\xa9\xb0\x9cI\xbb/OK\xe1E\xa7Я\xa8]\xa8\xd6Tu\xec-\xc9\xd5Y\x91\u007f\xe5]\xd1\xdc\xfc\x1f\x00\x00\xff\xffp\xb6n\x80H\x06\x00\x00",
		hash:  "df84590813410be2d70f764cb7770a30e0edc359f81eb03869edf7015e604b04",
		mime:  "",
		mtime: time.Unix(1582578977, 0),
		size:  1608,
	},
	"cmd/app/bootstrap.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x84R\xc1\x8e\x9b0\x10\xbd\xf3\x15֪W\f\xe9\xad\xdc*\xf5\xd0S\xba\xd2\xe6\xd4\xdb`O\xd8Q\xc0\xb6f\x06\xbai\x95\u007f\xafL\xc8.I\x91\x8a8\xc0{o\x9eߌ\x87\xc21ZH\xa9)\x8c\t0`c>\xfd\x81\x94l\xfe\xbe\x14\xc6x\x14ǔ\x94b\xb8Q+(+@\x95\xa9\x1d\x15%{\x18\xe3IR\x0f\xe7r\xed\xb6`\xfb\xc5Ԙ\x04\x8cA\x1b\x93z\xd0c\xe4a\x06\xf5\x9cp\x05\x15s\xb6v\xa4\xdegg`\xa5#8}L8\v\xbe\x81\xe2\x81\xf2\x81O\x9f\xebݗ\xb2ޕ\xf5\xeeP\xd7\xcd\xfc\xdaz~~>\xdd\xe4\xfbqh\x91\x1b\xf3\xb2\xff\xfa\xfc\xf2\xfdǡ0\x86\xb1G\x90\xf7\xc4\x13\xb2,\rv\x1c\xc7\xd4\x18\x17\a\xebH\\\xb4\x83\xbc=L\xeba$Kqf?z\xb0\xcb\t\x97\xf2\x0e]\xe5\xb9\x14\x85 Oȹ\xdd\x14yn\xf5\x8a\xd8\xfc\x9b\xad]\f\x8aoZ&\xd0\xd7\x15\xbd\xc0Ϡ\xafY%\nJn\x11Y[Y[y\x12\xad8F\xad&ધ\xb6\xba\x1f#\x06h{\xf4\x8dQ\x1e\xb1($1\x85\xcezP\xb0\x0eD x\x86\x9c\xeb\x84gI\xe0\xb0\xdcږm\x13\xd7\xc7\xd1\xe7\x882\xf6\xcd?\xaa\xb9\xa7#u\xd7\xedy\xe4|\x9e\xf8\x84|ޢ\xef\x0f\x98`\xecu˿ÀL\xee\xff\x0e\xa2\x8c0\xd8\x13\x1cO`[\n\xfez\x15\xdb%\x8c\x9e\xe4\x83~g\x15T\xb6\xaa~A\xd7m\xfa)\x83í܌\xf9֯5\xb7}\xfbM\xe9D\xa1\xf8\x1b\x00\x00\xff\xff\xbe\xd4\xf0h\xb7\x03\x00\x00",
		hash:  "f42bdb5fd8c83e7bf15b3fb390242c008f6b243569eb7c3e30d56eadb1e358ab",
		mime:  "",
		mtime: time.Unix(1582591592, 0),
		size:  951,
	},
	"cmd/app/main.go": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\x92Ak\xdc0\x10\x85\xcf֯\x18L\x0f6x\xb5\xf4Z\xc8a\xb3\xf5Bi6\rNZ\xc8q\"k\x1dQk\xa4J\xe3eK\xf0\u007f/\xf2*\r\x85B\xf6\xe4\xf1\xf8\xcd\xe7\xa77\xf2\xa8~\xe2\xa0\xc1\xa2!!\x8c\xf5.0T\xa2(\x95#\xd6'.S\xc9n5\x18~\x9e\x9e\xa42Q9\xa9\x9c]\xdf\xee~\xac\xae\xbf\xaf\a\xb7\xb2\xf1\xb4F\xef/T*\x8c\x11\xa9\x0f\xb8\xb6f\b\xc8\xfa\xc29\xfe\xedu,E-\x84r\x14\x17\x8f\xe8\xfd-Z\rWP~xA\xef%\xa1\xd5\xf3\xa29L\xa4\xc0\x90᪆\x97E)\xbfQ{\xd4\xc4U\xaa\x97j\xeb\xacE\xea\x1bH\x9d\xfc\xb2?{j \x01*\xc5'\xc89\xc8\xed\xf9Y\x83\x0e\xc1\x85\x04\xfd?\xf5\x9e1\xf0\x99y\xf7\x8cQ\u007f\x9e\x82\xa1\xa1\x01\xec3\xdc8\x8a\xb5(\x8a\xa0y\n\x04dFQ̵\x98\xb3\xe9\u007ft\xef8\xb0H\xe6\xa0#ç+\xc8i\xca}\xee킳y\"Qj!^\xff\xb8$)\xdb\x04\xb91\x91\xd3I^Ar\xd3\xf7\xdb_\xe3='\xcf\u007fmT9ޣ\x0e\xd18\x9a\xe5ǲ\x81r\x1b4\xb2\x86\x83\t\x91\x81\xf1i\xd4K\xb7k7\x0f-<l\xaeo\xda\xfc\xad:\xe28iHN\xe0\xae\xfb\xb2\xdft\x8f\xf0\xb5}\xac˺\x11\xc5,wfd\x1d\xaa\xb7\x04\xd2U|[[7Q\x95\x17\x9d$\u007f\x02\x00\x00\xff\xff=\u007f\xa1\xfc\xb0\x02\x00\x00",
		hash:  "7c18deaf72c827a6632997894607bca2bddaa8c0fde6f1a420aec0f923695c5f",
		mime:  "",
		mtime: time.Unix(1581778705, 0),
		size:  688,
	},
	"cmd/app/profile.production.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff,\xce\xc1\r\xc3 \f\x85\xe1;S\xb0@\x9c;\x1bd\x8cԼ6$\x14\x10\xb6\x1b\xb1}E\xd5\xebӧ\xa7_\xc0֓\x0e\xba0\x84\xce[\x83\xf3\xfe\xc2X\xa4Zg\x04\xdf\xf0\xfe/m\xd7#\xf8uB\xad\x1d\xeby\xeb\xd2\xecqa\xd0DNZO\xe5E\x9c\xabE\xe2Z\xc42\xc5$\\?\xe8cަ\"\xba\x17\xc6\x16\x837KѹC\xb5\x11焢\xa4Y\xb6\xf2\xcbA\xf0\xcf=\v\xdc7\x00\x00\xff\xff\x1c\xc0\x18\x99\x9d\x00\x00\x00",
		hash:  "17f63f2fde4844eb88714014c05dd8e48e5ea23346c480d8c64eee5ed138aa4b",
		mime:  "",
		mtime: time.Unix(1584104387, 0),
		size:  157,
	},
	"cmd/build/build.go": {
		data:  "package main\n\nimport (\n\t\"cto-github.cisco.com/NFV-BU/go-msx/build\"\n)\n\nfunc main() {\n\tbuild.Run()\n}\n",
		hash:  "1afc43b37b6664090421f8a3bde52704b1365a691a132a3c5d3ccd24bec6cc18",
		mime:  "",
		mtime: time.Unix(1573950111, 0),
		size:  0,
	},
	"cmd/build/build.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x94\x8fAj+1\f@\xf7s\n\x13\xfe*\xe09\x80w\x1f\x9a\xe6\x1c\x8aF\x0e&\xb6d$\x8dI)\xb9{\x99t\xa0P\xbahw\xc6z\xe2\xe9ѝpu\xb8TJS\bؖ\x14\xa0\xf7\xed)\x9c\xcb\xf5\xb5T\xb2m\x12B\f\x17\x117W\xe8\xf3[\xab\xfb\u07ffw\xe8}fh\xf4\x98\xbbʲ\xa2\x17\xe1'05\xbbo\xabJ\x95\xc0(\xed\xec \xb5\"\xfc\x98B\xe8\x15<\x8b\xb6OA\a%\xf6\xff\xea%\x03\xfan\xdd\x1c(m\xc6b(\xf3h\x96F\xb3h\xa4\xa3 EsP'\xfd=\x19Q\x94\xfe\x80\xdf \xdf\xe0\a\x9e\xf3H\x9cG,\xectUؚ#\n\xdbZc%X\xf6\x93\xf6\xd4\xef\xe5\xf1\xf4r>=\x81\xc2Xׅ\xce*k\xb7\x14\x0e_\x82\xe3\xf10}\x04\x00\x00\xff\xff\xf5s\xbc\x95\x9c\x01\x00\x00",
		hash:  "4412e2b744fa333266741167bb3aab143e04a615f096ebf6a7db567cfa0ea5f6",
		mime:  "",
		mtime: time.Unix(1581778621, 0),
		size:  412,
	},
	"docker/Dockerfile": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffl\x90\xc1\x0f\xda \x18\xc5\xef\xfc\x15_\xd4+\xa0\xf3\xb0\xc5\xc4\x03\xb6蚩5Tݖe1\xb4E%\xb6\xa5\x01ڋ\xf1\u007f_\xdanqKv\xe0\xf0=x|\xbf\xf7\x98\xd8\xc0\x8a%\xfc\x12\xed؆/s\x93=\x94\xbd7)ɴ\xcb\f\xc9LI\xdb\xd2Ẑ\xfejl\x89s\xd5\xe2\xe1Q\xaf\xa7\xd2)\xec\xbcU>\xbb/\xe6\xe4\x13\x99\xe2\xf9\xc7)B\xe3\xf1\xb8;\xb0jt\x91Cf*/u\xa5l\xaf\xa3\xb5\x88wp3\x85\xacn\x8b\x19\x99}\xf8\xf3\x01\xb0\x04\xd2Π,\xea\xb9N\xd16\xbc\xac\xb7l\x93 \xbe?\x0f3\x17\x83\x02\xb84\xf9\xb2UUn,ba\b\x04\xa8\xack\x8a\xbe\xc6\xe2K\x18\x89~BⴇR>\x14\xe4\xda\xf97V\xa8\x9d\xb7:m\xbc6\x15\x04\xff\xa1\x9b<ߥ\xbc\xfa\xe5\t\x17g..ɑ\x1d\xa3\xe0\xc0\x8e\x9f\x81\xb6\xd2\xd2B\xa7t\xf2\x94uM*Y\xaa\x17\xe2\xdf\x0eq\xc2a\xf2tʶʒ\xdaX\xffz3u\x16\xdbT\u007f[(\nv!\xfc\x18\xd1\xc6Y\x9a\xea\u007f\xeeF?Q\x10\x1f\xbe\x03\xc6Wk\xca\xe5\xefr\x86\xa0]\"j\x8d\xf1\x14(\xfa\x15\x00\x00\xff\xff8&\xe7&\xc7\x01\x00\x00",
		hash:  "b5d2b788864ab730906a63f9c71128eec3da4840b5e6b13aae259a4ac2ac7a49",
		mime:  "",
		mtime: time.Unix(1581790314, 0),
		size:  455,
	},
	"gitignore": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff<\x8e\xc1N31\f\x84\xef~\x8aH\xbd\xfc\u007f\x84\xbc\x02\t\xc4\x15\x81\xb8r\xe8\x11\xa1U\x9ax\xb7\xae\xb2q\x88ݲ}{\x94\x16\xb8\x8c=\xe3\xb1\xf4mܳ,\x953%\xf7\xb6;P47q&\xbdq[\v\xc6х\x92\xdc˹\x84\x85\xa3˼S\xf7o\xbb\x0f\xed\xaf\xad\xff\xc1\xa3\x80\xc7\x00\x1eU\x006\xeeUr\xa2\xa60\xca\xee\x00\xa3\x91ZO\x9fZܳQ\xb4c#\xa7\x95\"O\x1c\x1d\xadFEY\x8a\x0e\xb5\xd1\xc4+)x|\xbf\u007fx<}~\xc0\xcfD9\x1a\x80\xc78\xcb-\xcer\xdd\xee0\xc2\x18g\x19\x13M\xc7\xf2kf\xb1s%\xed\xad\x8b\xa7\xb5J3\xf4p\x05Y\x02\x97~\x03\x8f\xb4\x12x\xbc\xd0y\xacM\xa6\x1erm]\x97\xdc\xf5K\x01\x90\x13\x85\x01\x12\xab\rp\xa2\x92\xa4\r\x90%\x86<@\xff\x1d\xe0;\x00\x00\xff\xff\xbd\xf5\x8b$@\x01\x00\x00",
		hash:  "5a6f301e0263ffcf31b195f9fd92ba67e1ef7e097bc9776a3218f024a44a51ba",
		mime:  "",
		mtime: time.Unix(1575838263, 0),
		size:  320,
	},
	"go.mod.tpl": {
		data:  "module cto-github.cisco.com/NFV-BU/${app.name}\n\ngo 1.12\n",
		hash:  "2c6648bda73dc5b5b1cda267ca61dbf56e5135ce61ca321c34588d01e87a6a33",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1576101799, 0),
		size:  0,
	},
	"idea/modules.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\x8fMK\xc40\x10\x86\xef\xf9\x15\xc3Ы\x19\x05\x0f\"I\xf6\xe0\a(,.K=Kh\xc7\x1a\xc9\x17i+\x82\xf8ߥ\xb5\xd8z\xdb\xdb\xf0\xce3\xcc\xfb\xa8\xddg\xf0\xf0\xc1\xa5w)j\xbc\x90\xe7\b\x1c\x9bԺ\xd8i|\xae\xefϮpg\x84\xca%\xbds3\xac\xe4%\x1a\x01\x00\xa0\x9a\x14r\x8a\x1c\a\x886\xb0\xc6\xc3/\xb9O\xed\xe8yo\xa3\xed\xb8,\xeċyѯ\xc9&\x85W\xe7y,^\xe34\\\x13U\x87\xe3\xd3\xe3\xddM\xfdr\xfbp\xacH\xba\x96-U_6g9=\xfb\x96.x\x9c\x8f\xb2\x1d\xde4\x9e\x84Ӧ\f\xfdk\xa3\xe8O\xc6\bE\x8b\xb3\x11?\x01\x00\x00\xff\xff\xfc\xfa\xf5T%\x01\x00\x00",
		hash:  "ea1f6cbe7e59c77ad3ab8717f32458cc1859ad197075132cd395c9cac348fcc6",
		mime:  "application/xml",
		mtime: time.Unix(1575824055, 0),
		size:  293,
	},
	"idea/project.iml.tpl": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff|\x90AK\xc3@\x10F\xef\xf9\x15\xc3ЫF\xc1\x83\xc8n\n\xd2T\x14\xa3PZ<\x965;\xa9\x8b\x9b\x992\xbbQ\xfb\xef\xa5i\xa9\x17\xe9\xfd͛\xc7g\xa6?}\x84/\xd2\x14\x84-^_^!\x10\xb7\xe2\x03o,\xae\x96\xf3\x8b[\x9cV\x85\xe9\xc5\x0f\x91 \xef\xb6d\xf1\xad\xbe_7\xaf\xb3\xd5s\x8d\u007f\xa77X\x15\x00\xa6\x95~+L\x9c\x81]O\x16\x1fd/t\uf47cŬ\x03!\x94\xff\x82/\xf4\u074cO\x16\"\xb9q\xec6\xa4\xa3rd9\xef\xc9A\xa3\xc5.D\xba+\xcbɡ`={\\L\x8eN\x00#\xeaIkκ;\xb6\x06\xfe \r\x99\xfc\x93\xff<\x83%\x19\xb4\xa5\xb9DO\x8aЉ.)\xe5d\xb1s1\x9d\x92\xcbSsU\x98\xf20IU\xfc\x06\x00\x00\xff\xff\xb2\xf9\r`C\x01\x00\x00",
		hash:  "b9b9d3416624fe97b81477ec39ab098014029985a43f2daaa239010da63305c4",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1575824154, 0),
		size:  323,
	},
	"idea/runConfigurations/make_clean.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffT\x90Ak\xb3@\x10\x86\xef\xf9\x15Ð\xbb\U0001dcc1\x90\x18\xf1kհ\x98\xb3\f:Z\xabΆuLɿ/\xb1\x86\xda\xcb\x1ev\xf7}\x9f\x87wW\xba\xe1\xe6\x84EAh`\x83\x17\xef>\xb9T;\xc9\xd1I\xdd6\x93'm\x9d$$\u0530\xc7\xfd\x06\x00`W\xaeߖ\xe4@\x1dC\xd93\t\x82>nl09\xbc\x85\xe7\xf8=,\xf2\x83\x8d¼\xb0״8f\xe99\x8e\xae\xf6\x90\xc7Y\x8aPS\xa9\xce?ҹ!\xa1\x8e\xeb\xb6g\x84\x8ak\x9az5XS?2΄x\x8cXؓreP\xfdċ\xcc,4,Qx\x1e?>ۋ\xcd\xfe\x87Ǽ8\xc5v\x1b\xfcv+\xf9\x86\xd5\xe0\xa2\xfa\xe5|\xd7Jsj=\xcf*\x06\x11\xc87\xd3\xc0\xa2\xa3\xc1\x15d\x06\xb1\xdcG\bV\xe4\xe0\x85^۰~\xb8\n\xee\x06\xff\xe1\xeb\xf3.\xf83\xda~\xf3\xbcX\xb6\xdfo\xbe\x03\x00\x00\xff\xff\x05씙\x88\x01\x00\x00",
		hash:  "576e1cf9fad57c85416a716dccb22d687f14a04f60b766a6dd35f396c9eac250",
		mime:  "application/xml",
		mtime: time.Unix(1575822514, 0),
		size:  392,
	},
	"idea/runConfigurations/make_dist.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffT\x90Ak\x83@\x10\x85\xef\xf9\x15Ð\xbb\xd0s6\x10\x8c\x11۪a1g\x19t\xb4[u7\xaccJ\xfe}\x891\xd4^\xf6\xb0\xbb\xef}\x1foW\xb9\xe1\xea,[\x01K\x03+<{\xf7͕\xe8Ɇ\xce6\xa6\x9d<\x89q6%K-{\xdco\x00\x00v\xd5\xfamI\x0e\xd41\xd4f\x14\x04\xb9_Yaz\xf8\x88N\xc9gT\x16\a\x1dGE\xa9/Y\x19\xe6\xd9)\x89/\xfaP$y\x86\xd0P%\xce߳\xb9 \xa5\x8e\x1b\xd33B\xcd\rM\xbd(l\xa8\x1f\x19g@2\xc6lٓp\xadP\xfcċ\xcb\xec3,Qx\x1cO\x9d\xedY\xe7\xefQX\x94\xc7Do\x83\xbfn!߲(|\x9a\xfe8\xdf\x19\xdb\x1e\x8d\xe7\xd9D!\x02\xf9v\x1a\xd8ʨpŘ9lo#\x04+p\xf0\"\xafeX\xbe\\\r7\x85o\xf8\xfa\xbc\v\xfeM\xb6\xdf<.\x96\xe5\xf7\x9b\xdf\x00\x00\x00\xff\xff\x9b\u007f_\xe0\x86\x01\x00\x00",
		hash:  "042b4a3b558a1245ef52aee76edb3582195592a4cb90641b03960ca76bb4d30a",
		mime:  "application/xml",
		mtime: time.Unix(1575822528, 0),
		size:  390,
	},
	"idea/runConfigurations/make_docker.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffT\x90\xc1n\xab0\x10E\xf7\xf9\x8a\xd1({\xa4\xb7\x8e\x91\x10!\x88\xd7\x02\x91E\xd6h\x04\x03\xa5\x80\x1d\x99!U\xfe\xbe\n\x05\x95n\xbc\xb0}\xef9\xba\xa7ʎwk\xd8\b\x18\x1aY\xe1\xd5\xd9O\xaeD\xcf&\xb4\xa6\xe9\xdaّt֤d\xa8e\x87\xfe\x01\x00\xe0T\xed\xdf\xd6\xe4H=Cm\xab\x9e\x1d\x82<\xef\xac0\rޢK\xf2\x1e\x95E\xa0\xe3\xa8(\xf5-+\xc3<\xbb$\xf1M\aE\x92g\b\rUb\xdd3[*R\xea\xb9\xe9\x06F\xa8\xb9\xa1y\x10\x85\r\r\x13\xe3\x82H\xa6\x98\r;\x12\xae\x15\x8a\x9by\xb5Y\x8c\xc65\n\xaf\xe3G\xe8x\xd5\xf9\xff(,\xcas\xa2\x8f\xdeo\xb7\x90kY\x14n\xae_\xd6\xf5\x9diϝ\xe3\xc5E!\x02\xb9v\x1e\xd9ȤpGYHl\x1e\x13x;\xb4\xb7\xb1\xf7:,\x1f\xb6\x86\x87\xc2\u007f\xb8}>y\u007ff\xf3\x0f\xaf\x8bu}\xff\xf0\x1d\x00\x00\xff\xffۢ\xec\xfd\x8a\x01\x00\x00",
		hash:  "3ed376aff730b1df7a0dae1c71cf0319dfc624e6fca7737dabecea437b234df8",
		mime:  "application/xml",
		mtime: time.Unix(1575822541, 0),
		size:  394,
	},
	"idea/runConfigurations/project__local_.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff|\x92Ak\x021\x10\x85\xef\xfe\x8a\x10<\xb4\a7г+\xb4V\xc5\xd2Z\x91\xb6W\x99&\xb315Ʉlb\x91\xd2\xff^\\\\p\x95\x9aS\x98\xf7>\xdec\x98\xa1$\x17ȣÕÒ/#}\xa1L\xab\xec\xc7\xe4+\xa3s\x84dȿ\x80\a\x8d\x91\x8fz\x8c16\x94\xa7ڑ\xec\xff@\b\xc5\xe1\xff\xcbn,I\xb0\xb7\x9c)\xac \xdbT\xf2\nl\x8d\x9c\xa5}\xc0\x92\xcf\xe8>\x04kd\x83\x9fG\xf1&\xe2\xecU \x13\xc5\xfd\xa2I\x9a\x11;\xe1y\x93?\xafg\xe81BBU\xf2\x143\x1e\xab6u\x1d\xa9l\xf1\xb2'\x17'\xa6o\x8a[\xe3\xf5Z\x99\x88M\x18ہ\xcd\a`\xb9z}\x9a\x8c\xdf֏\xf3U_H\xa7\x04\x84\xd0A5\xad\x03Dp\x980\xd6-60\x1dϥ\xa1#o\x8dW\xad0\x9d?O:be,.!m\xae6\x12\x0e\x8c/4\x9d\xa5\xca-hlA\x99h\xa0M\xda\xe4\xcfB\x9aZR!ɉ\xc5\xf4c\xf0\xf0.\xfe\xdb\xcb\xf5}t\xac\x0eӆ\x14ە\xfc\xae\x9d\x0fE\xe7VF\xbd\xc3\xe0xr\xa3\xde_\x00\x00\x00\xff\xffaT\x84\x80\u007f\x02\x00\x00",
		hash:  "e825a8bf92de1fb46b523a24a477952a1531d513693012d8c06f037292000da3",
		mime:  "application/xml",
		mtime: time.Unix(1575824540, 0),
		size:  639,
	},
	"idea/runConfigurations/project__remote_.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff|\x92Ao\xe20\x10\x85\xef\xfc\n\xcb\xe2\xb0{H\"\xed\x99 m) \xaa\x96\"\xd4\xf6\x8a\xa6\xce$\xb8\xd8\x1e˙P\xa1\xaa\xff\xbdJ \x12\xa6-9E\xf6\xfb\xe6=?\xcdH\x91\xf5\xe4бp`1\x97\xab@o\xa8xݸ\t\xb9RWM\x00\xd6\xe4\x1e\xc0A\x85A\x8e\aB\b1R\xe7w'r\xf8\x01ާ\xed\xff\xa7\xf8\x13\xd0\x12\xe3_)\n,\xa11\x9c\xcb\x12L\x8dR\xf0\xc1c.\xe7\xf4\xdf{\xa3U\xc7_z\xc9\xce\xe3\xe2+A1\x85ò\xb3\x9a\x938\xe3e\x17`Q\xcf\xd1a\x00\xc6\"\x97\x1c\x1a<e\xed\xf2Z*\x1a\x83߃\xca\xecL\xf4Na\xa7]\xb5)t\xc0\xceL\xec\xc14-\xb0Z?\xdeM'O\x9b\xdb\xc5z\x98)[d\xe0}\x84V\xb4\xf1\x10\xc0\"c\xa8{,ё\xe6\aAr\xec1\xf5\xc0[\x11\xbb\x18R`D\x92\xf8@\xa56(\x8e\x85F\x03w\xda\x15\xfd\xa8\xd9\xe2~\x1a]\xb6Ъ\x1d{\xed\r\x99\x05\xedҊ.r\xaa\x1dT\u0603\x8a)\xa94o\x9b\xd7T\xe9ZQ\xaa\xc8f\xcb\xd9Kr\xf3\x9c\xfd\xd6\xe4\xf5\x06#\xa9E\xdeR!\xf6\xb9\xfcן\x8f\xb2h\xbdƃ\xf6ഥ\xe3\xaf\x00\x00\x00\xff\xff\xd4s\x05\xf0\xb1\x02\x00\x00",
		hash:  "742c8c03944b07e4ee9a279b58d6a2b783b6a2f95dc23c2ebe7fb0fe4b2b687b",
		mime:  "application/xml",
		mtime: time.Unix(1575824558, 0),
		size:  689,
	},
	"idea/vcs.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffDα\n\xc20\x10\xc6\xf1\xbdOq\x1c]\xb5\n\x0e\x0eI:X\x15\x05QJu\x95\x92\x86\x121w!\tE\xdf^Ъ\xeb\xc7\x0f\xfe\x9f(\x1f\xee\x0e\x83\t\xd12I\x9cOg\b\x864w\x96z\x89\xe7f3Yb\xa92\xe1\x03ߌN\u007f\xb9@\x95\x01\x00\b\xcd\xce3\x19J@\xad3\x12/:V6\x18\x9d8<\x0f\xad\xf7\x96\xfa8ڷw\x9f\r\xba/\x92\x98\x9f\xea\xe3~\xbdj\xaeծ\xce\x11\x06\x1d%nmB(\xc6F\xf1\x8b\xa8L\x14\xe3\x17\xf5\n\x00\x00\xff\xff\\3\a\xb4\xbc\x00\x00\x00",
		hash:  "4d519e4e626e33220b5a4c104b8fa665ff6bd0bec841dc6be828cd645814b014",
		mime:  "application/xml",
		mtime: time.Unix(1575824180, 0),
		size:  188,
	},
	"idea/workspace.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\xd0AK\xc30\x14\a\xf0\xfb>E\b\xbbnQ\xf0\xe0!\xed\x90Y\xe7Dl\xa9\xd1k\xc9\xdag\x17\xc9\xde+\xe9\xeb\xa8\xdf^\xba\xba\t҃\xc7\xe4\xe5\xff\xff\x91\xa7W\xfd\xc1\x8b#\x84\xd6\x11F\xf2zy%\x05`I\x95\xc3:\x92o\xe6aq+W\xf1L7\x81>\xa1\xe4ߗ72\x9e\t\xa1K:4\x84\x80,\xd0\x1e \x92\x9b4OS#Ecy\x1fIյAy*\xadW5I\xa1\xa6#\xf4\xecv\xc1\x06\a\xed\xa9S\bM\r;\u009f\xb9\xc3\n\xfa\x04\xd9\x05\xd8Pfy/\xc5\xd1\xfa\x0e\"\xf9a}\v\xe7Zu\xe9\x9dV\x1c/_\x81\xd9a=\xcd\xe4\xc9:y1\xc5fk\x8a\xe1\vEvg\x1e/\xd2<\xcbӧdm\x8a\xfbm>\xff\xa7\xf8^S6n\xed\xec9d\xa8\x83\x1d\xd0\x05\xa0\xddy\xa8b\x0e\x1dh55\x193M\xa0\xfe+\xae\\\x80\x92\xb5\x1aO\u007f\xf1\xd3\xfd\x00\xc5\xdf\x01\x00\x00\xff\xff\xa9N\xa5\x18\xcf\x01\x00\x00",
		hash:  "5f0eb712ac71ed257fb6199f64e029bf05b6327af032bcd7461b057646e69be1",
		mime:  "application/xml",
		mtime: time.Unix(1575822164, 0),
		size:  463,
	},
	"k8s/kubernetes-deployment.yml.tpl": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xacVMO#9\x10\xbd\xe7WX\xd1H\x9c:Ͱ;he\x89\xc3\nؙ\xd5\x00\x89\b\xcc\x05\xa1\xc8qW\x12o\xfc\xb5vu\x0fQ6\xff}ew\x00w\xa7\xc3\x04\t\xb8t^\x95_=W\x95\xcbβ\xacǬ\xf8\x01\xce\v\xa3)a\xd6\xfa\xbc\xfa\xdc[\n]Pr\x01V\x9a\x95\x02\x8d=\x05\xc8\n\x86\x8c\xf6\b\xd1L\x81\xb7\x8c\x03%\x95\xf2[\x80\x92Okf\xed |oz\xde\x02\x0f\xae\x0e\xac\x14\x9cyJ\xd6kR\xbc\xd0M\x94)`\x02\xbazha\xffI\xf3\x13\xdc\xe3\xc3\xd1vᄛR\xe3\xd1\xe3\xc3Q\xc2~\xf4H6\x9b\x1e!\x1e$p4.\x04\"D1\xe4\x8b+6\x05\xe9k\x80\x84\xdd4e\xd5\xf0ܙ2\x1a\x96\xe5\x14\x9c\x06\x04?\x88س\x037ڗ2\x9b\x1b\uf165\x84Ii~\xf6\bAPV2\x84m\xbc$#\xe1o7\t5.\x1b\x92\xf6\x8a:@\xd6^a\x91Vk\x83\f\x85\xd1I\xac\x99,A\xe3T\xe0@\x98\xdc2\xe7\xc1Q27\xca?E\x97\xe7*m\xa9\x91\t\r.Y\x9em7U\x87}\x81\t\x11\x8a\xcd!ִ6M\"@6\x1b\xfa\nUuOեz݀RL\x174\x81B\x98\x9d\x00\x01d\xf3\xd0xM,\x9b\n]\x9c\x1d\x0f\xe2\u007f\xdbƥ\x00\x8d\xfb\xac\xa1V\x1c4\x82;{\x15Y\xf0\xa6\xbe\xe8\xe9\x00\xdd*\xfb\xc7\b}V\xbb\r<\xb8Jp\x18t\xea\x8c\xccY!\xdcY^\xdb\xf3\x00\xec\x883z&\xe6\r\xb7\x1aJ\x1c+#K\x05ס\xe5};E*\xa0#\x86\vJ\xf6\xae\u007fmC\xbb\x10;U\xec깴\x94\x89\xbdQ\xcf\x14\u007fOQ\xfby\xe9]>\x15:O\x18\xfa\xed\xbcd֙\x99\x90\xd0\u00ad3E\xc9C;'\x06)*\xd0\xe0\xfdș)4\xc3-\x10\xedW@\xdaʆ\x8d\t\xfb\xb4\x0e\x05\x04\x17\xea\x87\xf0\x14Ӹ\xc9Y\xa1\x84\xceY m\xaf2\x0e\x93U\xe1g\xb3I\x84\x16(\x98\xbc\x00\xc9Vc\xe0F\x17\x9e\x92ߎ\x9b=g\xc1\tS$\xe6\xc4\xea\x80\x15⣷\xb2\x00&q\xf1!{9}\xcfV\xbc)\x1d\x87V\xc7:\xf8\xb7\x04\xdf\xee\xe306\x95q+J\xfa\xa7\xbf_\x8b~\xcb\xc8mII\xffs\x13\x96B\x897xN\xbe\x9c\xee#:Ia\xd0U\xbbC\xebs1\x1e\xdd\xfe}\xf3ur~5\xbc\xbf\x98\x9c\x0fo\xc6\xf7W\x93o\xc3\xf1]\x8b\xb2b\xb2\x04J\xfa\xeb5\xa9X)q\xe2\xf9\x02T<$y.\rgra<\xf6\x0f\x0e1\x1a\xde\xee\v\xf1Ǘ\xe3\xe3\x03\x88~\xfcy\u007fu\xf7\x96Ԩ\xb35\xbc\x0e\xe6}K\xdf\xc9;\xf4\x8dϿ]^_\x1e\x9c̃y\xef\x86\xdf/o\xbah\xffrF\xb5\xbb%<\x16\xb8\x03\xfc\x0e\xab[\x98\xedZ\x9f\xe7\xa6\xf2OQO\x87\xc3\x12V\x94\xa0YB:\x91\xc2q\xf2\x1d\xb7Y}\x95\x8e\xde>{\x87\x8e\xfb%\xac<\x1a\a\x9d\x93\xbe\xd3\xd8\\\x0f\xc8s\xefe\xce\xc1\xa1\xcf9\xcb\u0087\x98\t\xce\xc2\x13\x83;\xec$\x8e\u07bdTi\xc7\xcb \xbdg\xe2\x1bb&\xe6\xd7\xcc6w\xb3\xffJ\xea\x10\x1f\x0eQ\xd4\xdd\xdb\x1d~\xf1b\xcd+\xe5_R\x92\xef\xbeU\x12ٿ\xe4\xdb\xc9̴ԅ\x84\x98\x93\xff\x03\x00\x00\xff\xff\xd7U\x9ec\x1b\v\x00\x00",
		hash:  "1f64ebc0171e357cd2d06c1a76583bc511b785a015d71c070359079584dd21de",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1580740271, 0),
		size:  2843,
	},
	"k8s/kubernetes-init.yml.tpl": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xac\x93ۋ\xda@\x14\xc6\xdf\xf3W\x1c\xa4 \x14\xa2\xbb\x85B\x99\xb7bm\xb7\xac5\xe2e_D\xc289k\a\xe7ֹ\xa4\x954\xff{I\xa2\xdd4\xd5\xed\x16\x9c'\xe7\xfbf~\xe7\xccwL\x1c\xc7\x115\xfc\x01\xad\xe3Z\x11\xc8o\xa3=W\x19\x81\x99\xce\"\x89\x9ef\xd4S\x12\x01(*\xd1\x19ʐ@\xaf(`\x1f\xb6h\x15zt\xe9o\aʲw<I\xe0UA\x8d\x19T\xbf\xcb\b@\xd0-\nWq\x00\xa81]\x1b`gu\xa8\xe5'\xf0\xa0\xd6\xca\xc8\x19d\xd5M\x8b\xceS\xebgZpv 0\xc5\x1cm\x04\xc0\xb4\xf2\x94+\xb4G||\xb6\x81jqIwH\xa0(\xda^Z\xabP\x96\xa4\xa3\xe7M$P\x9e\xae3-%U\x199n\xabJ\xc3\xe0\xecp\xcb\xd5\xf0\xefb\x95\x1d\xc7\xc6\xeaG.\xb0\xa5\x19\xab\xb3\xc0<ת%J\xbe\xb3ԟ\x8eYt:X\x86\uea54\xc5o\x01\x9do)\x00̄f\x16\xb7\xaf34B\x1f$*\x9fJ\x9da\x8a*_w\xb4\x9fB\u007fG\xbbY\xf7-\x1a\xc1\x19M\x99\x0e\xca\xf77\xeb~\xab\xf7\xfe\xe68\xc3j\xa1\xca\xdbomR]\xcc柧\x9f\xd2\xd1$Y}HG\xc9t\xb1\x9a\xa4w\xc9b\xd9\xea+\xa7\"T\xff\x12\xa6\x95\vb\xe0\xd0\xe6\x9c\xe1\xa0\xd9\xf6^D\x9c%\xf3s\xc4woon^\x06X\x8c\xee\xc6_\xc6g\x10E\x019\r§\x8e}E\x89\xad\xe7>\x0f\\&\xf7\xe3i\x97\xf7\xd1jٞ\b\x80Cf\xd1\xdf\xe3a\x8e\x8f\u007f:\xa7\xefB\xba\x1fM\x12\x1dw\x8f\a\x02^\xefQ=\xdf\xcf\xc3\xfb\xd5dy)\xf2\xfai\xff\x97xû\x14\xf8\x9b\u007f\x06\xdeܿ^\xde\r\xef\x9aq\xd7\xf5c\x9e\xc9ˉ\xff\n\x00\x00\xff\xff\x827~\xc6\x06\x05\x00\x00",
		hash:  "785dd276701797be8d684158a29b3fffb60ed1cef8eaed4128f5631aba41d378",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1575754811, 0),
		size:  1286,
	},
	"k8s/kubernetes-poddisruptionbudget.yml.tpl": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffl\x8e\xbdJ\x051\x10\x85\xfb<\xc5\x14\xb6Y\xd96\xa5XZ\xd8h+\x93dX\xc3\xe6gH&\x8b\xb2\xec\xbbKV\xf1\xde\v\xb7\x1b\xbe3\x87\xf3i\xad\x15rx\xa7\xdaB\xc9\x06\xb8\xc4\xe0\xbe\x1f\xb7ْ\xe0\xac\u0590\xbd\x81\xd7\xe2\x9fC\xab\x9d%\x94\xfc\xd4\xfdB\xa2\x12\tz\x144\n \xa2\xa5\xd8\xc6\x05\x80\xcc\x06\x1evd\x9e2&:N\xb8\xd4\xd2O\xbcvK5\x93P\x9bN6\xe2\xf1vS\xd1\xec\xed\x1fo\x8c\x8e\f\xec;\\\x9a\x1f\xff\x01\x1c\x87jLn,'\xfcz˸a\x88h#\x19\x98\x15@\xa3HNJ\xfd5K(\xee\xf3\xe5J\xf5\x8e\xecO\x00\x00\x00\xff\xffk\x814}\x0f\x01\x00\x00",
		hash:  "c805658a62f677b62df2e28c718972fd6af2171a164787d45ce21a1be731672a",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1575753121, 0),
		size:  271,
	},
	"local/profile.remote.json5": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xaa\xe6RPP*.(\xca\xccKW\xb2R\x00\xf1\x14\x14\x94\x92s\xf2KS\xe0\\\x90@~^qi\x0e\x92\x88\x82\x82RF~q\x89\x92\x95\x82\x92\xa1\xb9\x91\x9e\xa1\x05\b\x99\x1b(\xe9@\xe5ka\f\xa5\xb2\xc4Ҝ\x12\xd2u\x82\xa9Z\x1d\xaeZ.@\x00\x00\x00\xff\xff'\xe0\x96ڢ\x00\x00\x00",
		hash:  "a86fa658debf8f0c9c3fbd9cbab1ef9afb961d81413600a280d95100f5b5e74a",
		mime:  "",
		mtime: time.Unix(1575818187, 0),
		size:  162,
	},
	"manifest/assembly.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xac\x90?O\xc30\x10\xc5\xf7|\n+b\xaa\x14_[\xb6\xcaxb,S\x17V\u05fe\xb6FI\x1c\xf9.$\x05\xf1\xddQ\x13\xecV\xfc\xd9\xd8\xee\xfd\xfc\xdeӝ\x95!\xc2f_\x9f\xc5\xd8\xd4-=\x94'\xe6n\x03ИWl\xa5\xe9\x8c=\xa1\f\xf1\b]\xdd\x1f}K\xf3C\x95RՌ!iXɕ\\\x97\xa2\x10s\xdff$\x9f;\x87a\x90\xc3\xfdԶ^.W\xf0\xfc\xb4\xdd\xd9\x136\xa6\xf2-\xb1i-\x96\x97\x1c\xf9\rMx\x1b\xaca\x1f\xda\xffYJ\xfcU2\x92\xcb\xd6j\xb2ʑ\\\xa9\v!\x94w:\"\x85>Z$\x05\xdeM\xf0\x10bc\x98.sV\xfa\xcdw\n\xbe\xe6\x8b\tn\\\xea\xe0k\xdc\xe152\xcbY\t\xa1\x9c\x8fh9ĳ\xbe{\xefbxA\xcbro\b\x9d\x8f\x1f\xc0&\x1e\x91\x01G\x8e\xa6\xbaY\xe6\x1aJ5\xa1\xe7\xae\xe7\xc7\xccA\xc1w\x94\xac8ںwH\t\\\x91\xa6ha\xb1P\x90\xf4O\a#\xf1o\x96\fҝpsh\x16\xa4\v\x95\xff[\x17\x9f\x01\x00\x00\xff\xffhơ\x96\x81\x02\x00\x00",
		hash:  "87ce119ee5b2b1d41b508a7e93c2505538e54b87f3ac29d2c0660426b34908f3",
		mime:  "application/xml",
		mtime: time.Unix(1581777887, 0),
		size:  641,
	},
	"manifest/pom.xml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xacWM\x8f\xdb6\x13\xbe\xfbW\xe85\xf6*\xd1\xfb&\x87\xc2\xd0*h\x9a\x06\r\xb0\xdb\x16\xf9Bo\x01M\x8d%\xa6\")p(ٮ\xb1\xff\xbd \xa9oɻ\xdb\xcd\xfa$\xce<\xc3\x19\xce<C\x8e\xe37GQ\x045h\xe4Jެ\xaf\xa3\xcd:\x00\xc9T\xcaev\xb3\xfe\xf2\xf9}\xf8\xd3\xfaM\xb2\x8aK\xad\xbe\x033\xc1Q\x14\x12oֹ1\xe5\x96\x10Ak\x90\x11-)\xcb!R:#\u007f\xfeqG^G\x1b\xbb\x8bCn\x8f\xc8;\xf4\xe1p\x88\x0e\xaf\x1c\xee\xff\x9b\xcd5\xf9\xeb\xee\xf6\x13\xcbAАK4T2X\xaf\x82\xf6wD\xbeE\xa7\xbdU\x8c\x1a\x17ݣ^\x9d\xe5%\xd4\x11S/\f\x1d6:b\xbaN\x9c\xc3X\xa8\x14\x8a\xaf>\t\x89\xd3\xc6d$\xf30I\x05$WgZ\x96\x91\xfd\xbc\x0f\x05\x95|\x0fhb\xe2T\x1e\x95iU\x95\x1f҄)\x111\x8eLE\xb5\xc0\xa8\x85FW\xe7\xbd*R\xd0\xf71i\x91ގj\xc3\xf7\x94\x99\x0f\xe9\x05\x1f\x03\x80\xb7()\xfb\x9bf\\fI\xa9DL\xfa\xa5W7UM\xae\xce\xcd\xd7}xu\xdeU\xbcH\xbf\xc9J\xec\\\b-f\xd5\xec\xa8U\t\xdap\xc0\xa4-E\xec\xc3M\xdeZ\xc3\U0001387b\x02b\xd2\b\xbd\x15\x19\x9ay\x112\xd1(\x99\x92\x12\x98-`\x82Ll3n\xb6\xb6B\xb8%\x84\x19\x15f\xdc\xe4ծI\x14S\x82\xfc\xfe\xfek\xf8\xf6\v\xa9\x05\x86\x1cUA\r\xa4Q\xc6ML\x06\x1b\xf9\x9dS\xa8\xa1\xb0\x8e\u007fy)\x17K;z_\x86f\xc9o\xbf\xfe\xfc.&\xf6\xcb\xcb*]$\xcf\xf3c-\x9bԹD5Is\xb5\xe9JQT\x19\x97]\x1d\x9auW\x96\x11_<\xad5\xa0\xaa4\x03\f=vΘ\x11-^Eז\xe7\xf5\x90\xe2\x0e\x00G`\x95={O\x82\x91x(\r\x82\x98[\xaa\x97\xa7\xde}Lx:\xc1\x949EH\xb8\xe4\x86ӂ\xff\x031\xf1\x921*S\xb4\xc0\xb1\xac\x91\xce<8\xe1ؚ,\x98[\xf6\xedyVi:\x0f<\bbU\x99\xb22\xef\xb8\x06f\x94>%W\xe7\x1dEH\xb9\xbe'\x86\xea\f\f\x81\xa3\xd1t\xe8xj2ݲ\x83N5\x03\xdd\\eټ\x14\xc5\xc0qzɥoR^\x18ж\xf7\x8d\xael\u007fv\xeby\x18\xe4R\x1c\xbdf\x9aFr1\x8f1Y\xe0\xc5@\xd8\x13\x98\x8c\x19<'\xf4\xff°[|\xce9\x06\x1c\x83CNM\xe0/6\xc0\xc0\xe4\x10\xa4P\x16\xea$@\x9a`\xaft\xa0L\x0e\x1a\xa3\xce0\f\x1fj\x11\x8a\bbW\x9c~\xb0C.g#\x05d\x9a\x97F\xe9i\n{M\xd2F\x11\x1dEa/\x9dN1\xca\xeb\xe2V\x17+\xf1\x9c\xae\xd5@\r\x84T\xb3\x9c\xd7p\xb1k\x9b\xec\xff\xe7\x96E.3\xfbT<\xa9U\x9fM\xa3%\"\xb5\xef\xaa\xd2Y\xc4T\n9\xad0\x12껚<\xb9S\x8e\xb8\xfb\xd7?\x8d\xa1\xe7˓hr\x1d\xbd~\x81k\xd4g\xb7\xa6\x05O\xa9y\xc6\x05\xe9\x8a\xf9\xd2\xd9~\x90\xeb\x98+m>B\xcd\xed\xc1oAf&O61Y\x12\xaf\x1e#p_\xd6\xd5p\xd9\xcd\x13d\xf06\xc6)G\xa3\xf9΅zG%\xcd\xc0^\a\x83\x985\x94\n\xf9\xfc\xb2\xb4\xacǓ(s%\xedkR\x00\xc5%\u07bban\x8e\xebg\xbc\x0e9\x1c\x01@f\x8c7\xbc\x11\x14\r\xe8\xc1(\xd0\x12H\xe9\x13\x99\xef\xdc\xcd\x03\xed-<\x0f?FIK̕\xf9\xf8\x94\xb3\xb5\xe0G\x0f\xd7\x03_\xfct\xfd֓\xe3]:\x89}\xe3\x96\v\xbb\x1aWu8\x9e>\\\xec\x9akS\xd1\xe2\x9b\xc5,\xe4\xe2\xc7\x0e\xe8\xee\x1229ިVKC\x00H;A\xa7\xcd;ݮ&\x1d\xbb\xb0\xc1\x9c\x16\x03\x89\x9f\xbc\xdd n\xff\xa6%\xab\u007f\x03\x00\x00\xff\xff\x19V\xc6\xde\xd8\r\x00\x00",
		hash:  "e210ad97358ab12941b6d9ad9f771d9c96392a41a684fad784291246bf7b79ca",
		mime:  "application/xml",
		mtime: time.Unix(1581783160, 0),
		size:  3544,
	},
	"manifest/resources/manifest-images.yml": {
		data:  "images:\n   - image: dockerhub.cisco.com/vms-platform-dev-docker/${app.name}\n     tag: \"${version}-${build_number}\"\n",
		hash:  "7f3946b0bdbd6613544086ee1f704f6d1781e22766f5a32185a13f2f7b7ed55b",
		mime:  "",
		mtime: time.Unix(1581779149, 0),
		size:  0,
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
