package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

var cowfile = flag.String("f", "", "what cowfile to use")

const (
	upper = "_"
	left  = "|"
	right = "|"
	lower = "─"
)

var toReplace = map[string]string{
	"$the_cow = <<EOC;\n": "",
	"\nEOC":               "",
	"$thoughts":           "\\",
}

func balloon(text string) string {
	l := utf8.RuneCountInString(text)
	return fmt.Sprintf(""+
		" %s\n"+
		"%s %s %s\n"+
		" %s",
		strings.Repeat(upper, l+2),
		left, text, right,
		strings.Repeat(lower, l+2))
}

func prepare(cow string) string {
	for old, neu := range toReplace {
		cow = strings.Replace(cow, old, neu, -1)
	}
	return cow
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
	if err != nil {
		return "", err
	}
	defer file.Close()
	out, err := ioutil.ReadAll(file)
	return string(out), err
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
