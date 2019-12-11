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
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffl\x91=O\xc30\x10\x86\xe7ޯ\xb8\x81\x81\x0e.\xbb%\x06\x10\x9fR\x05\xa8\x88\x81\t9\xf6ŵ\xb0}\x91?(\xfc{\x147*J\xe9rR\xf2>z/\xf7\xe4\xfa\xedq}s\xbb\xc1K\xb4\x8c\xa9F<;\x9f^}ܭ\xaf\xee_\x97\xa8\x83\xb9\xe8\xaa\xf3\xd3\\YF!4\xc7\xde\xd9\u007f\xd9O\xf0\x00\xab\x97\x87\xe7\xa7w\x89\x85rA\xe3\xc6\xc1\xfa\x93\x12\x0e\xb5\xf3.oQ{R\x11\x87D\x9aCp\x05`$%,\x0e\xab\x97hx\x17=+#\xc6H\x18\x1a\xf2,\xa6oҵ\x90\xa8ѕ\x86d\x80qӼ\xc4R\xa4\xa4\n\x89\xf6q\xc2Şg\xb9\x8b\xb9(\xefžMu\x9e\xa6\xbb\xf2I,\uf535\x94Du\xb3x_\xfe\xd7\x01`\xa8\xabV\x9e\x80Z0G\x9b\x19\t\v\xcb\x18\xd8\xe0\x17E\xc3\xe9Hň\xeco\x00\x98\x1c\xcaS\xc8P\xf3\x16\xa0ٕ\xb0H\x01Ŧo?\xe0\xf00\xd5\xc3\xc1\xfd\x910\x16}(\xf0\x1b\x00\x00\xff\xffL_2\xec\x14\x02\x00\x00",
		hash:  "dafa5720a7876cccac37f0e313ffa67b51cea1cc41f428748fbbd4ec46cbfb57",
		mime:  "",
		mtime: time.Unix(1575843965, 0),
		size:  532,
	},
	"README.md": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\x94Ao\xe36\x10\x85\xef\xfa\x15\x03\xa4\a\xbbH\xe4k\xb1(\n\xb4\xc16\x97\xdd\"h\xb6@\x81`\x01\x8fȑ4\r\xc5\x11\x86\xa4\xbbn\xd1\xff^\x90\x94my\x93\xc3^\f\x93z\xf3ę\xf7Q7\xf0ݿ8ϭǉ\xfek\x9a\x9b\x1bx\x92\xa4\x86\xe0^,5ͧ\x91\xa0\x17\xe7\xe4o\xf6\x03XV2Q\x94)\x00*\xc1\xacr`K\x16zQ8JR0b\t\xd0[0\xe2{\x1e\x92bd\xf1\xef\x9a\xe6{؛\xc9\xee\x1b\x00\x80;x\xff\x85L\x8a\xd89\x02\xf2Q\x8f\xb3\xb0\x8f\xa1\x05X\x9e\xff<ώM\xa9]\t\x80Cq\xd9\xe1<\xef&d\xdf\x0e\xb2o\x97\x92_\x12;\xfb\x96\xb8\xcb\x0f\xea\xefE\x9f\xcfcż\x90\xeeoa\xff\xf2C8\x9d\xec^|D\xf6\xa4\xfcO}=j\xe4\x1eM\f\x97B'\x06ݩ\xe0C^\\\xb7\v=;\xcaݼ\xffb\\*\xf3Q\x99@i\x96\xc0Q\xf4\xd86ͯ\xa2\xd0'\x8d#)X\x8a\xc8.\x80x\x18$\x0f\xf5/2\x11B\xd4dbR\xba\x85@\x04\xcd\xf3SDoQ-<\b<.\xa2\x0fx\x94\x14?o\xc6\x18\xe7\xf0n\xb7\x1b8\x8e\xa9k\x8dL\xbbA\x1c\xfa\xe1.,Ua\xb7\x18߹R\xb3mK\xdaelM\xb3\x82\x00R\xa0\x00\x0f\xbf\xfd\x01\x1f\xf1\x85J\xb4q$(\x03\x84p\f\x91\xa6\x16\xe0\x1a\x8c\x88:P\fM\x86\x02\x0f\xc8.G[S\x8f\x14\xe2\x1e\xaer'ț\xa1X\xa3s0\x89M.#\xe5-\f\xe4I1R\x99\x96\xc6\x00\x9b\xe7j\xf1y\xd3\xee\xf2\x9f\xed\xb6d\xc7+\xd7\x1a}\xdeR\xeeRI\xe0D걘\xd2\x05\xb7\xcds\xad\xcdv\xf9\xcfv[Y\xa0.\r_\x1b\xe6\xbdU튙\xaf\x84e\x13x¡\x88\xe6\xd49\x0e\xe3Y\xf5X\u05eft\xc6\x11\xfa\x13G\xbf\xd3$\x87:\x99r\xe67\xdba\n\xe5\x05JF\xa6\x89\xe3\xa9\xf8Qiγ/\xd7/\x8fՌd^\xee\xd87\xcd'YZ \xc0%\xa7\x8ac\x0e5\x8c\xe4\xdc-h\xf2\xb0\x9fr\xda?V\xc5O\xfbJ\xc7Cf\xc8\xe6\xcf\x00\x873\x98\x06=t\x042\x93'\v\x9c\xa1ͪ\xdbr\xec\x11\x03\x04:\x90\xa2\x03K=&\x17\x8b\xfd\xd5\x05\t\x15\x8d5t\x9br\xa9\xb6\xe7i$_\x0e8\xb1Q\t\xa4\a6\x048 \xfb<\x1e(b\xc8\vt\xae\xde9\xe9\xe1\xe3ӟ\xafl\x95&\x89\xf4\xad\xbeU\xfd\x96\xf1\x82|\x15\x143\xc8=y2Et}\xff9\x94\x13\x16\xa7\xf2\xbag\xb7\xfaJd\xf4j\xb7\x10Rw&u1\x95\xbeT,\xc3nsC%\x99\n\xcb\t\xa9\xfb\xbcZ\vϺB\xf7\x15\x9e+\xd55U\xbd8Kz\xa9\xacl\xbf\xaa]c\v\xcd\xff\x01\x00\x00\xff\xff\xef\xf57\xc21\x06\x00\x00",
		hash:  "94a5227df4b2ba4a26c9c9548f3eda1d23b63a8b0cdfc0395b058f29181b2f69",
		mime:  "",
		mtime: time.Unix(1575915399, 0),
		size:  1585,
	},
	"cmd/app/bootstrap.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x84R\xc1\x8e\x9b0\x10\xbd\xfb+\xacU\xaf\x18\xd2^Zn\x95z\xe8)]is\xeam\xb0\av\x14\xb0\xad\xf1@7\xad\xf2\xef\x95\tiH\x8a\xb4\x88\x03\xbc\xf7\xe6\xf9\xf9\xd9\xe4\xdb` \xc6Zi\xeda\xc0Z\u007f\xf8\x031\x9a\xfc}VZ;L\x96)\n\x05\u007f\xa5VPV\x80\bS3\n\xa6졵\xa3\x14{8\x15k\xb7\x05\xdb/\xa6ZG`\xf4R\xeb\u0603\xb4\x81\x87\x19\x94S\xc4\x15\xa4\xe6l\xcdH\xbd\xcb\xce\xc0B-XyL8\v\xbe\x81\xe0\x81\xf2\x82O\x1f\xabݗ\xa2\xda\x15\xd5\xeePU\xf5\xfc\x9aj~~>]\xe5\xfbqh\x90k\xfd\xb2\xff\xfa\xfc\xf2\xfd\xc7Ai\xcd\xd8#$\xac\xf5'\xf3\xd9TJ\xeb\x8e\xc3\x18km\xc3`,%\x1b̐\xde\x1eJzhbBNKO\xb7\xe8f1>\x17w\xe8*\xc6Y\xa9\x84<!\xe7]\xc6\xc0\xf3\x0e/\x88ɿ\xd9\xda\x06/\xf8&E\x04y]\xd1\v\xfc\f\xf2\x9aUI@\xc8.\"cJcJGIJ\x0eA\xca\t\xb8\xec\xa9)\xef\xdbC\x0fM\x8f\xae\xd6\xc2#*\x95\"\x93\xef\x8c\x03\x01c!%\xf0\x8e!\xe7:\xe2)E\xb0\xb8\x1ck\n\x03\xe6\fdoC\xb6\x0f\xa3ˑ\xd2\xd8\xd7\xffY\xcf{h\xa9\xbb\\\x92G\xce\xe5\x86'\xe4\xd3\x16}\xbf\xc0\x04c/[\xfe\x1dzd\xb2\xef;$a\x84\xc1\x1c\xa1=\x82iȻK\xf5\xdb#\x8c\x8eҍ\xfe\xc7\nHښ\xfa\x05]\xb7\xe9'\f\x16\xb7r3\xe6S\xbe\xcc\\\xef\xd7o\x8aG\xf2\xeao\x00\x00\x00\xff\xff;\xd5Y\xfe\x9e\x03\x00\x00",
		hash:  "be222746081af149a99bbaaf90d2e6e599df6a8ffc5036916eba6c0cca6657fb",
		mime:  "",
		mtime: time.Unix(1575762766, 0),
		size:  926,
	},
	"cmd/app/main.go": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\x92Qk\xdb0\x10ǟ\xadOq\x98=X\xe0(\x1b{\x19\x83>\xa4\x99\x03cMW\xd2n\xd0ǫ\xac\xb8b\xd6I\x93\xce%\xa3\xe4\xbb\x0f9\xea\xca`\xb0<\xf9|\xfa\xdfO\xff\xbbS@\xfd\x03\a\x03\x0e-\ta]\xf0\x91\xa1\x11U\xad=\xb19p\x9dC\xf6\x8b\xc1\xf2\xe3\xf4\xa0\xb4M\xda+\xed\xdd\xf2z\xf3}q\xf9m9\xf8\x85K\x87%\x86p\xa6RcJH}ĥ\xb3CD6g\xd6\xf1\xaf`R-\xa4\x10\xdaS\x9a=b\b\xd7\xe8\f\\@\xfd\xe6\x19CP\x84\xce\x1cg\xcd~\"\r\x96,7\x12\x9eg\xa5\xfaJݓ!nr<Gk\xef\x1cR\xdfBΔ\x9f\xed\xc9S\v\x19\xd0h>@\x99\x83Z\x9f\xbe\x12L\x8c>f迩\xb7\x8c\x91ƠGL\xe6\xd3\x14-\r-`_\xe0\xd6S\x92\xa2\xaa\xa2\xe1)\x12\x90\x1dEu\x94\xe2XL\xff\xa5\xfb\x8f\x03\x87d\xf7&1|\xbc\x802M\xb5-\xb9M\xf4\xaeTd\x8a\x14\xe2\xe5\xc6y\x92\xaaː+\x9b8w\xf2\x02R\xab\xbe_\xff\x1co9{\xfec\xa3\xa9߫\x0f\xea\xadzW\xb7P\xaf\xa3A6\xb0\xb7110>\x8cf\xce\xee\xba\xd5]\aw\xab˫\xae\x9c5O8N\x06\xb2\x01\xb8\xd9}ޮv\xf7𥻗\xb5lEuT\x1b;\xb2\x89\xcdk\xe3\xf9\x05\xbenk7QS\xf6\x9b%\xbf\x03\x00\x00\xff\xff\x01\x8bOs\xa7\x02\x00\x00",
		hash:  "ced38a8bc08ad65f9f8a1b20f5e59da797c54e77b3e6be27719c13191f4143fe",
		mime:  "",
		mtime: time.Unix(1575736744, 0),
		size:  679,
	},
	"cmd/app/profile.production.yml": {
		data:  "security.jwt.keys:\n  key-source: pem\n  key-path: /keystore/jwt-pubkey.pem\n",
		hash:  "06dfef26cc039f6e97c2b4246df15c3781e5d950200709e01c9d61edca64998d",
		mime:  "",
		mtime: time.Unix(1575761467, 0),
		size:  0,
	},
	"cmd/build/build.go": {
		data:  "package main\n\nimport (\n\t\"cto-github.cisco.com/NFV-BU/go-msx/build\"\n)\n\nfunc main() {\n\tbuild.Run()\n}\n",
		hash:  "1afc43b37b6664090421f8a3bde52704b1365a691a132a3c5d3ccd24bec6cc18",
		mime:  "",
		mtime: time.Unix(1573950111, 0),
		size:  0,
	},
	"cmd/build/build.yml": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x94\x8f\xc1J\x041\f\x86\xef\xf3\x14e\xf1\xb4\xd0\"x\x91\xde\x04\xd7}\x8el&]ʶMI2eE|w\x99q\xc0\x8b\a\xbd\x95\xbf_\xbe\xe4\xa7;\xe1bp)\x14'\xe7\xb0\xce\xd1A\xef\xeb\x93[\xca\u05f7\\H\xd7\x1f缻0\x9b\x9a@\x0f\xef\xb5\xec\xd9\xc3\a\xf4\x1e\x1aT\xfa\f]x^\xd02\xb7\r\x98\xaa\xde\xd7Q\xa1B\xa0\x14\xddSx\x0e\x8f\x93s\xbd\x80%\x96\xfa\xed\xed \xd4\xecE,'@ۗ\xadj\xe4\x1a0+r\x18U\xe3\xa8\xea\x95dd$\xaf\x06b$\u007f'=\xb2\xd0?\xf0\x1b\xa4\x1b\xfc·4bK\xc3\xe7ft\x15X\xabz\xe4\xa6K\xf1\x85`\xdeO\x1a$\x9a\xb9\xed\x85\xfd\xe9\xf5|\xda\xf2ܰ,3\x9d\x85\x97\xae\xd1\x1d~\xbc\xc7\xe3a\xfa\n\x00\x00\xff\xff5sh\xa3\x8a\x01\x00\x00",
		hash:  "7f65287d82a8e5e1b0e96ba72b02b6f7c43449db992136757ef71f80ab461b9f",
		mime:  "",
		mtime: time.Unix(1575762638, 0),
		size:  394,
	},
	"docker/Dockerfile": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffl\x90\xc1n\xe20\x10\x86\xef~\x8a\x11\xe1j{Y\x0e\xabE\xe2\x10\x92Т\x02A\t\xd0VU\x85\x9cĀEbGc'\x17ĻW\x81\xa2\xb6R\x8f\xf3\xcf\xfc\xd27\x9f\xe7y\xc4\xf3<\x984\xaa, 7\xda\t\xa5%v\x19!\xd3$^\xc0\xc1\x94B\x1fF\x036\xf8K\xadC\xe9\xf2#\xf8)d]A\"\x89\x96[\x98lf\xf30Jvӹ\xff\x90\x02\xadL1n\xa5.\f\x12?\f\x81\x01\x17u\xcd\xc9s\x9c<\x85\xb3\xe4:\x91d\xb3\x84J\x9c$\x14\xca:B\xee\x18\xa1\xb2\x0eU\xd68e4\x04\xbf\xd0\x14&?I<6\x19˕\xcd\r\xcbM\xc5\xdb\xcaҺ\x14no\xb0\xa2\x85l\xe9\xed\xe8\x9ag\xc2\xca;\xf6h\xc8\xfe\xb1?t8\xfc\u007f\xa5N\xa3d\x1b%\xbbt\xed\xafg\xc1\xca_?\x02o\x05\xf2Re\xbc\u007f\x16uʹ\xa8\xe4\x85D/\xab8\x8d\xa0\u007f\xb6\x12[\x89\xac6\xe8._\xcft\x15l\xf4\xf7\n'\xc1\"\x84\xb7\x1eo,\xf2L\xfd\xd8\xf5\xdeI\x10\xaf^\x81\xd2=\x9aj\xfci\xf1f\xa8S\xc1\xd1\x18ǁ\x93\x8f\x00\x00\x00\xff\xff\xbb\xe2\xf5\x02\x99\x01\x00\x00",
		hash:  "84ddaa2d3a4191723b63e9400169f22fe470eb191adcfddf717eff6566229831",
		mime:  "",
		mtime: time.Unix(1575761491, 0),
		size:  409,
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
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xacV\xdfo\xe28\x10~\xcf_a\xa1\x95\xfa\x14\xd2\xed\xddV'K}8\xb5\xbd\xddӖ\x82J\xbb/U\x15\x19g >\xfc\xebl'[\xc4\xf1\xbf\x9f\xec\xd0℀\xa8\xb4\xed\v|3\xfe\xe6\xf3\xccx\x864M\x13\xa2\xd9\x0f0\x96)\x89\x11\xd1\xdaf\xf5\xe7d\xc9d\x81\xd1\rh\xaeV\x02\xa4K\x048R\x10Gp\x82\x90$\x02\xac&\x140\xaa\x85\xdd\x02\x18}Z\x13\xad\x87\xfe\xf3&\xb1\x1a\xa8w5\xa09\xa3\xc4b\xb4^\xa3\xe2\x9d.\x17\xaa\x80\x1cd\xfd\xdc\xc1\xfe\xe3\xea'\x98\x97\xe7\xb3\xed\xc1\x9c\xaaJ\xba\xb3\x97糈\xfd\xec\x05m6\tB\x168P\xa7\x8c\x0f\x84\x90 \x8e\x96wd\x06\xdc6\x00\xf2\xb7i\xcbj\xe0\x85QU0,\xab\x19\x18\t\x0e\xec0`o\x0eTI[\xf1t\xa1\xace\x1a#¹\xfa\x99 \xe4@hN\x1cl\xe3E\x19\xf1\u007f\xfbIhpޒtP\xd4\t\xb2\x0e\n\v\xb4R*G\x1cS2\x8a5\xe7\x15H7cn\xc8T\xa6\x89\xb1`0Z(a_\x83\xcb[\x95\xb6Ԏ0\t&:\x9en/Մ}\x87\x11b\x82, Դ1\xe5\x01@\x9b\r\xdeAu\xd3SM\xa9v\x17\x10\x82\xc8\x02G\x90\x0f\xb3\x17\xc0\x83d\xe1\x1b\xaf\x8d\xa53&\x8b\xab\xf3a\xf8\xef\xda(g \xdd!\xab\xaf\x15\x05\xe9\xc0\\\xedD\x16\xb4\xad/x\x1apf\x95\xfe\xa3\x98\xbcj܆\x16L\xcd(\f{u\x06\xe6\xb4`\xe6*k\xec\x99\a\xf6\xc4)9g\x8b\x96[\x03E\x8e\xb5╀\x91oy\xdbM\x91\xf0脸\x12\xa3\x83\xe7wm\xa8K\xb6Wž\x9e\x8bK\x19\xd9[\xf5\x8c\xf1\x8f\x14u\x90U\xd6d3&\xb3\x88a\xd0\xcdK\xaa\x8d\x9a3\x0e\x1d\\\x1bUTԷsd\xe0\xac\x06\t\xd6N\x8c\x9aA;\\\xe9\x9c\xfe\n\x0ew\xb2\xa1C\xc2>\xad}\x01\xc1\xf8\xfa9x\ri\xdcd\xa4\x10LfēvO)\xe3\xa2S\xfek\xbbI\x98d\x8e\x11~\x03\x9c\xac\xa6@\x95,,F\xbf\x9d\xb7{N\x83a\xaa\x88̑\xd5\x00)د\xbeJ\t\x84\xbb\xf2\x97\xdc\xe5\xf2#W\xb1\xaa2\x14:\x1dk\xe0\xdf\nl\xb7\x8f\xfd\xd8\x14ʬ0\x1a\\\xfe>b\x83\x8e\x91\xea\n\xa3\xc1\xe76̙`Gx.\xbe\\\x1e\"\xba\x88a\x90u\xb7C\x9bw1\x9d<\xfc}\xff5\xbf\xbe\x1b?\xdd\xe4\xd7\xe3\xfb\xe9\xd3]\xfem<}\xecPքW\x80\xd1`\xbdF5\xa9\xb8\xcb--A\x84G\x92e\\Q\xc2Ke\xdd\xe0\xe4\x10\x93\xf1á\x10\u007f|9??\x81\xe8ǟOw\x8fǤ\x06\x9d\x9d\xe1u2\xef1}\x17\x1f\xd07\xbd\xfev;\xba=9\x99'\xf3>\x8e\xbf\xdf\xde\xf7\xd1\xfee\x94\xe8v\x8b\xff\xb1@\r\xb8\xef\xb0z\x80\xf9\xbe\xf5mn\n\xfb\x1a\xf4\xf48,a\x85\x91SK\x88'\x92\u007fN\xb6g\x9b5\xabtr\xfc\xed\x9d:\ue5f0\xb2N\x19\xe8\x9d\xf4\x1dc\xc3ٳ\xc3\xe3\x8d\x10\xb6\xfd\x9c-FD\xb7\xe3\x1e^\x1e=\x1a|\xbb\a\x85\xc9\xfe\x98\n+0\xab\x85}\x17\x9f%\xff\a\x00\x00\xff\xffBA~\x95l\n\x00\x00",
		hash:  "97e069e85f363dcc65aba6099e80981fad5e4928db9fc21b90206e49085faaa6",
		mime:  "application/vnd.groove-tool-template",
		mtime: time.Unix(1575750935, 0),
		size:  2668,
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
