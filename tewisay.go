package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

var (
	cowfile = flag.String("f", "default", "what cowfile to use")
	list    = flag.Bool("l", false, "list cowfiles")
	think   = path.Base(os.Args[0]) == "tewithink"

	tongue = flag.String("T", "  ", "change tounges")
	eyes   = flag.String("e", "oo", "change tounges")

	borg     = flag.Bool("b", false, "borg")
	dead     = flag.Bool("d", false, "dead")
	greedy   = flag.Bool("g", false, "greedy")
	paranoid = flag.Bool("p", false, "noided")
	stoned   = flag.Bool("s", false, "stoned")
	tired    = flag.Bool("t", false, "tired")
	wired    = flag.Bool("w", false, "wired")
	young    = flag.Bool("y", false, "young")
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
	// fuck.
	switch {
	case *borg:
		*eyes = "=="
	case *dead:
		*eyes = "xx"
		*tongue = "U "
	case *greedy:
		*eyes = "$$"
	case *paranoid:
		*eyes = "@@"
	case *stoned:
		*eyes = "**"
		*tongue = "U "
	case *tired:
		*eyes = "--"
	case *wired:
		*eyes = "OO"
	case *young:
		*eyes = ".."
	}

	for old, neu := range map[string]string{
		"$the_cow = <<EOC;\n":     "",
		"$the_cow = <<\"EOC\";\n": "",
		"\nEOC":                   "",
		"$eyes":                   *eyes,
		"$tongue":                 *tongue,
	} {
		cow = strings.Replace(cow, old, neu, -1)
	}
	if think {
		return strings.Replace(cow, "$thoughts", tline, -1)
	}
	return strings.Replace(cow, "$thoughts", line, -1)
}

func countRunes(s string) (n int) {
	fmt.Println(s)
	for _, r := range s {
		if unicode.IsGraphic(r) && !(unicode.IsMark(r)) {
			n++
		}
	}
	return n
}

func balloon(text string) string {
	text = strings.Replace(text, "\t", "    ", -1)

	var (
		r      = right
		l      = left
		length = 0
		middle []string
	)

	if think {
		l = tleft
		r = tright
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if newlen := countRunes(line); newlen > length {
			fmt.Println(newlen)
			length = newlen
		}
	}

	var (
		up   = strings.Repeat(upper, length+2)
		down = strings.Repeat(lower, length+2)
	)

	for _, line := range lines {
		middle = append(middle,
			fmt.Sprintf("%s %s %s%s", l, line,
				strings.Repeat(" ", length-countRunes(line)), r))
	}

	return fmt.Sprintf(" %s\n%s\n %s",
		up, strings.Join(middle, "\n"), down)
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
	if os.IsNotExist(err) {
		return "", fmt.Errorf("couldn't find cowfile at %s!", name)
	}
	if err != nil {
		return "", err
	}
	defer file.Close()
	out, err := ioutil.ReadAll(file)
	return string(out), err
}

func listCowfiles() {
	cowpath := os.Getenv("COWPATH")
	if cowpath == "" {
		cowpath = "/usr/share/cows"
	}
	files, err := ioutil.ReadDir(cowpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	var cows []string
	for _, file := range files {
		if !(path.Ext(file.Name()) == ".cow") {
			continue
		}
		cows = append(cows,
			strings.TrimSuffix(file.Name(), ".cow"))
	}
	fmt.Printf("Cow files in %s:\n", cowpath)
	fmt.Println(strings.Join(cows, " "))
}

func main() {
	flag.Parse()
	if *list {
		listCowfiles()
		return
	}

	var tosay string
	if args := flag.Args(); len(args) != 0 {
		tosay = strings.Join(flag.Args(), " ")
	} else {
		out, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		tosay = strings.TrimSpace(string(out))
	}
	cow, err := getCowfile(*cowfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n%s", balloon(tosay), prepare(cow))
}
