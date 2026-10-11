package detector

import (
	"regexp"
	"strings"
)

type sig struct {
	cat string
	re  *regexp.Regexp
}

var sigs []sig

func add(cat, pat string) {
	sigs = append(sigs, sig{cat: cat, re: regexp.MustCompile(pat)})
}

func init() {
	add("sqli", `union(\s|/\*.*?\*/|%20|\+)*select`+
		`|select.{1,80}from.{1,80}(where|information_schema)`+
		`|information_schema\.(tables|columns)`+
		`|sysobjects|syscolumns|sys\.tables|pg_catalog|sqlite_master`+
		`|\bor\b\s*\d+\s*=\s*\d+|\band\b\s*\d+\s*=\s*\d+`+
		`|\bor\b\s*'[^']*'\s*=\s*'[^']*'|\band\b\s*'[^']*'\s*=\s*'[^']*'`+
		`|'\s*(or|and)\s*'|"\s*(or|and)\s*"`+
		`|'\s*or\s*1\s*=\s*1|"\s*or\s*1\s*=\s*1`+
		`|sleep\s*\(\s*\d|benchmark\s*\(|pg_sleep\s*\(|waitfor\s+delay`+
		`|;\s*(drop|alter|truncate)\s+table|;\s*insert\s+into|;\s*update\s+\S+\s+set|;\s*delete\s+from`+
		`|xp_cmdshell|sp_executesql|exec(\s+|\()master`+
		`|extractvalue\s*\(|updatexml\s*\(|exp\s*\(~`+
		`|having\s+\d+\s*=\s*\d+|group\s+by.{0,40}having`+
		`|cast\s*\(.{0,40}as\s+(int|char|signed)|convert\s*\(.{0,40},`+
		`|0x[0-9a-f]{6,}|concat\s*\(.{0,60}select`+
		`|order\s+by\s+\d{1,3}\s*(--|#|;|$)`+
		`|procedure\s+analyse\s*\(|--\s|#\s*$|'\s*--|'\s*#`)

	add("xss", `<script\b|</script|<svg\b|<iframe\b|<object\b|<embed\b`+
		`|<img[^>]*\bon\w+\s*=|<body[^>]*\bonload\s*=|<details[^>]*\bontoggle\s*=`+
		`|<marquee[^>]*\bonstart\s*=|<input[^>]*\bonfocus\s*=|<[a-z]+[^>]*\bon(error|load|mouseover|click|focus|toggle|start)\s*=`+
		`|javascript\s*:|vbscript\s*:|data\s*:\s*text/html`+
		`|document\.(cookie|domain|write|location)|window\.location`+
		`|eval\s*\(|settimeout\s*\(\s*['"]|setinterval\s*\(\s*['"]`+
		`|alert\s*\(|prompt\s*\(|confirm\s*\(|fetch\s*\(.{0,40}document\.cookie`+
		`|expression\s*\(|@import|<style[^>]*>.*expression`+
		`|%3cscript|%3c%2fscript|&#x?3c;?\s*script|&lt;\s*script`)

	add("traversal", `\.\./|\.\.\\|\.\.%2f|\.\.%5c|%2e%2e(%2f|%5c|/|\\)`+
		`|%252e%252e|\.\.%c0%af|\.\.%c1%9c|\.{4,}//|\.{4,}\\\\`+
		`|/etc/passwd|/etc/shadow|/etc/hosts|/proc/self/(environ|cwd|status)|/proc/\d+/`+
		`|c:\\\\?windows|c:/windows|boot\.ini|win\.ini|win32|system32`+
		`|php://(filter|input|output)|zip://|phar://|file:///|expect://`+
		`|\.\.\x5c|%00`)

	add("cmdi", `;\s*(cat|ls|id|whoami|uname|wget|curl|nc|bash|sh|ping|rm|chmod|nohup|python\d?|perl|telnet)\b`+
		`|\|\s*(cat|ls|id|whoami|nc|bash|sh|whoami)\b|\|\|\s*(cat|ls|id|whoami)\b`+
		`|&&\s*(cat|ls|id|whoami|uname|wget|curl)\b|&\s*(cat|ls|id|whoami)\b`+
		`|\$\([^)]{0,80}\)|`+"`[^`]{0,80}`"+
		`|/bin/(sh|bash|dash|zsh)\b|%0a\s*(cat|ls|id|whoami)`+
		`|powershell(\.exe)?\s|cmd(\.exe)?\s*/c|nc\s+-e\s|bash\s+-i\s`+
		`|wget\s+http|curl\s+http.{0,10}(-o|--output)|nohup\s.{0,40}&`)
}

const quickBytes = `'"<>;|&$` + "`" + `.%\`

func prescreen(s string) (bool, string) {
	if !strings.ContainsAny(s, quickBytes) {
		return false, ""
	}

	for _, sig := range sigs {
		if sig.re.MatchString(s) {
			return true, sig.cat
		}
	}

	return false, ""
}
