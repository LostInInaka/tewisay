/* tewisay - cowsay clone

To the extent possible under law, the author(s) have dedicated all
copyright and related and neighboring rights to this software to the public
domain worldwide. This software is distributed without any warranty.

You should have received a copy of the CC0 Public Domain Dedication along
with this software.
If not, see <http://creativecommons.org/publicdomain/zero/1.0/>. */

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
	_border = flag.String("b", "unicode", "which border to use")

	list  = flag.Bool("l", false, "list cowfiles")
	listb = flag.Bool("lb", false, "list borders")

	tongue = flag.String("T", "  ", "change tounges")
	eyes   = flag.String("e", "oo", "change tounges")
)

var escRxp = regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)

func countRunes(s string) int {
	n := 2
	s = escRxp.ReplaceAllString(s, "")
	for _, r := range s {
		if r == '\t' {
			n += 8 - (n % 8)
		} else if unicode.IsGraphic(r) && !(unicode.IsMark(r)) {
			n++
		}
	}
	return n
}

type border [9]string

var borders = map[string]border{
	/* Format:
	top    left, middle, right,
	middle left,         right,
	bottom left, middle, right,
	line, */

	"say": {
		" ", "_", " ",
		"| ", " |",
		" ", "─", " ",
		"\\",
	},
	"classicish": {
		" ", "_", " ",
		"< ", " >",
		" ", "-", " ",
		"\\",
	},
	"think": {
		" ", "_", " ",
		"( ", " )",
		" ", "─", " ",
		"o",
	},
	"unicode": {
		"┌", "─", "┐",
		"│ ", " │",
		"└", "─", "┘",
		"╲",
	},
	"thick": {
		"┏", "━", "┓",
		"┃ ", " ┃",
		"┗", "━", "┛",
		"╲",
	},
	"rounded": {
		"╭", "─", "╮",
		"│ ", " │",
		"╰", "─", "╯",
		"╲",
	},
}

func balloon(text string, b border) string {
	text = strings.Trim(text, "\n")
	text = strings.TrimSuffix(text, "\n\x1b[0m")

	var (
		maxlen int
		middle []string
	)

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if newlen := countRunes(line); newlen > maxlen {
			maxlen = newlen
		}
	}

	var (
		up   = strings.Repeat(b[1], maxlen)
		down = strings.Repeat(b[6], maxlen)
	)

	var lastEscs []string
	for _, line := range lines {
		s := fmt.Sprintf("%s%s%s\x1b[0m%s%s", b[3],
			strings.Join(lastEscs, ""), line,
			strings.Repeat(" ", maxlen-countRunes(line)), b[4])

		middle = append(middle, s)
		lastEscs = escRxp.FindAllString(line, -1)
	}

	return fmt.Sprintf("%s%s%s\n"+
		"%s\n%s%s%s",
		b[0], up, b[2],
		strings.Join(middle, "\n"),
		b[5], down, b[7])
}

func replaceVar(s string, v string, r string) string {
	s = strings.Replace(s, "${"+v+"}", r, -1)
	s = strings.Replace(s, "$"+v, r, -1)
	return s
}

func prepare(cow string, b border) string {
	// :c
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
	cow = strings.Join(ncow, "\n")

	// oh god
	cow = strings.Replace(cow, "\\\\", "\\", -1)
	cow = strings.Replace(cow, "\\@", "@", -1)
	cow = replaceVar(cow, "eyes", *eyes)
	cow = replaceVar(cow, "tongue", *tongue)

	return replaceVar(cow, "thoughts", b[8])
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
	switch {
	case *list:
		listCowfiles()
		return
	case *listb:
		var l []string
		for k, _ := range borders {
			l = append(l, k)
		}
		fmt.Println("Availible borders:\n",
			strings.Join(l, " "))
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

	var b border
	if path.Base(os.Args[0]) == "tewithink" {
		b = borders["think"]
	} else {
		nb, ok := borders[*_border]
		if !ok {
			fmt.Printf("error: no border called \"%s\".\n"+
				"pass -lb to list borders\n", *_border)
			return
		}
		b = nb
	}

	fmt.Printf("%s\n%s", balloon(tosay, b), prepare(cow, b))
}
