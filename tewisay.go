/* tewisay - cowsay clone

To the extent possible under law, the author(s) have dedicated all
copyright and related and neighboring rights to this software to the public
domain worldwide. This software is distributed without any warranty.

You should have received a copy of the CC0 Public Domain Dedication along
with this software.
If not, see <http://creativecommons.org/publicdomain/zero/1.0/>. */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	flag "github.com/neeee/pflag"
	"github.com/neeee/rwidth"
)

var (
	cowfile = flag.StringP("file", "f", "tes", "what cowfile to use")

	borderStyle = flag.StringP("border", "b", "unicode",
		"which border to use (use \"list\" to show all)")

	list = flag.BoolP("list", "l", false, "list cowfiles")

	tongue = flag.StringP("tongue", "T", "  ", "change tounge")
	eyes   = flag.StringP("eyes", "e", "oo", "change eyes")
)

type border [10]string

var borders = map[string]border{
	/* Format:
	top    left, middle,  right,
	middle left, padding, right,
	bottom left, middle,  right,
	line, */

	"say": {
		" ", "_", " ",
		"|", " ", "|",
		" ", "─", " ",
		"\\",
	},
	"classicish": {
		" ", "_", " ",
		"<", " ", ">",
		" ", "-", " ",
		"\\",
	},
	"think": {
		" ", "_", " ",
		"(", " ", ")",
		" ", "─", " ",
		"o",
	},
	"unicode": {
		"┌", "─", "┐",
		"│", " ", "│",
		"└", "─", "┘",
		"╲",
	},
	"thick": {
		"┏", "━", "┓",
		"┃", " ", "┃",
		"┗", "━", "┛",
		"╲",
	},
	"rounded": {
		"╭", "─", "╮",
		"│", " ", "│",
		"╰", "─", "╯",
		"╲",
	},
}

func countRunes(s string) int {
	n := 2
	s = rmAnsiEsc(s)
	for _, r := range s {
		if r == '\t' {
			n += 8 - (n % 8)
		} else {
			n += rwidth.Width(r)
		}
	}
	return n
}

func balloon(text string, b border) string {
	text = strings.TrimSuffix(text, "\x1b[m")
	text = strings.TrimRight(text, " \t")
	text = strings.Trim(text, "\n")

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
		topBorder    = strings.Repeat(b[1], maxlen)
		bottomBorder = strings.Repeat(b[7], maxlen)
	)

	var esc string
	for _, line := range lines {
		s := fmt.Sprintf("%s%s\x1b[0m%s", esc, line,
			strings.Repeat(" ", maxlen-countRunes(line)))
		middle = append(middle, b[3]+b[4]+s+b[4]+b[5])
		esc = lastEsc(line)
	}

	return fmt.Sprintf("%s%s%s\n"+
		"%s\n"+
		"%s%s%s",
		b[0], topBorder, b[2],
		strings.Join(middle, "\n"),
		b[6], bottomBorder, b[8])
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

	return replaceVar(cow, "thoughts", b[9])
}

func readCowfile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	out, err := ioutil.ReadAll(file)
	return string(out), err
}

func getCowfile(name string) (string, error) {
	if strings.Contains(name, "/") {
		out, err := readCowfile(name)
		if err == nil {
			return out, nil
		}
		return "", fmt.Errorf("Could not find %s cowfile!", name)
	}

	cowpaths := os.Getenv("COWPATH")
	if cowpaths == "" {
		cowpaths = "/usr/share/cows"
	}

	for _, cowpath := range strings.Split(cowpaths, ":") {
		name := cowpath + "/" + name + ".cow"
		out, err := readCowfile(name)
		if os.IsNotExist(err) {
			continue
		}
		if err == nil {
			return string(out), err
		}
	}
	return "", fmt.Errorf("Could not find %s cowfile!", name)
}

func listCowfiles() {
	cowpaths := os.Getenv("COWPATH")
	if cowpaths == "" {
		cowpaths = "/usr/share/cows"
	}
	for _, cowpath := range strings.Split(cowpaths, ":") {
		files, err := ioutil.ReadDir(cowpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		var cows []string
		for _, file := range files {
			if path.Ext(file.Name()) == ".cow" {
				cows = append(cows,
					strings.TrimSuffix(file.Name(), ".cow"))
			}
		}
		fmt.Printf("Cow files in %s:\n", cowpath)
		fmt.Println(strings.Join(cows, " "))
	}
}

func main() {
	flag.Parse()
	switch {
	case *list:
		listCowfiles()
		return
	case *borderStyle == "list":
		var l []string
		for k := range borders {
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
		nb, ok := borders[*borderStyle]
		if !ok {
			fmt.Printf("error: no border called \"%s\".\n"+
				"pass -lb to list borders\n", *borderStyle)
			return
		}
		b = nb
	}

	fmt.Printf("%s\n%s", balloon(tosay, b), prepare(cow, b))
}
