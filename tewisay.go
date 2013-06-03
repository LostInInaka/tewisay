package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

var (
	cowfile = flag.String("f", "", "what cowfile to use")
	think   = path.Base(os.Args[0]) == "tewithink"
)

const (
	upper = "_"
	lower = "─"

	line  = "\\"
	left  = "|"
	right = "|"

	tline  = "o"
	tleft  = "("
	tright = ")"
)

func prepare(cow string) string {
	for old, neu := range map[string]string{
		"$the_cow = <<EOC;\n": "",
		"\nEOC":               "",
	} {
		cow = strings.Replace(cow, old, neu, -1)
	}
	if think {
		return strings.Replace(cow, "$thoughts", tline, -1)
	}
	return strings.Replace(cow, "$thoughts", line, -1)
}

func getCowfile(name string) (string, error) {
	if !strings.Contains(name, "/") {
		cowpath := os.Getenv("COWPATH")
		if cowpath == "" {
			cowpath = "/usr/share/cows"
		}
		name = cowpath + "/" + name + ".cow"
	}

	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return "", err
	}
	out, err := ioutil.ReadAll(file)
	return string(out), err
}

func balloon(text string) string {
	var (
		l    = utf8.RuneCountInString(text)
		up   = strings.Repeat(upper, l+2)
		down = strings.Repeat(lower, l+2)
	)

	if think {
		return fmt.Sprintf(" %s\n%s %s %s\n %s",
			up, tleft, text, tright, down)
	} else {
		return fmt.Sprintf(" %s\n%s %s %s\n %s",
			up, left, text, right, down)
	}
}

func main() {
	flag.Parse()
	tosay := strings.Join(flag.Args(), " ")
	cow, err := getCowfile(*cowfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n%s", balloon(tosay), prepare(cow))
}
