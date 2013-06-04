package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
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

var escRxp = regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)

func countRunes(s string) (n int) {
	s = escRxp.ReplaceAllString(s, "")
	for _, r := range s {
		if unicode.IsGraphic(r) && !(unicode.IsMark(r)) {
			n++
		}
	}
	return n
}

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

func balloon(text string) string {
	text = strings.Replace(text, "\t", "    ", -1)
	text = strings.Trim(text, "\n")

	var (
		length = 0
		middle []string
		r      = right
		l      = left
	)

	if think {
		l = tleft
		r = tright
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if newlen := countRunes(line); newlen > length {
			length = newlen
		}
	}

	var (
		up   = strings.Repeat(upper, length+2)
		down = strings.Repeat(lower, length+2)
	)

	var lastEscs []string
	for _, line := range lines {
		s := fmt.Sprintf("%s %s%s \x1b[0m%s%s", l,
			strings.Join(lastEscs, ""), line,
			strings.Repeat(" ", length-countRunes(line)), r)

		middle = append(middle, s)
		lastEscs = escRxp.FindAllString(line, -1)
	}

	return fmt.Sprintf(" %s\n%s\n %s",
		up, strings.Join(middle, "\n"), down)
}

func replaceVar(s string, v string, r string) string {
	s = strings.Replace(s, "${"+v+"}", r, -1)
	s = strings.Replace(s, "$"+v, r, -1)
	return s
}

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

	var ncow []string
	for _, line := range strings.Split(cow, "\n") {
		switch {
		case strings.HasPrefix(line, "$the_cow"):
		case strings.HasPrefix(line, "EOC"):
		case strings.HasPrefix(line, "#"):
		default:
			ncow = append(ncow, line)
		}
	}

	// oh god
	cow = strings.Join(ncow, "\n")
	cow = strings.Replace(cow, "\\\\", "\\", -1)
	cow = strings.Replace(cow, "\\@", "@", -1)
	cow = replaceVar(cow, "eyes", *eyes)
	cow = replaceVar(cow, "tongue", *tongue)

	if think {
		return replaceVar(cow, "thoughts", tline)
	}
	return replaceVar(cow, "thoughts", line)
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
		tosay = string(out)
	}
	cow, err := getCowfile(*cowfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n%s", balloon(tosay), prepare(cow))
}
