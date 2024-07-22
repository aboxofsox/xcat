package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"
)

var (
	findLineNumberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb638"))
	lineNumberStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#677d8a"))
	helpTextStyle       = lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("#666666"))
)

var (
	viewMode    = flag.Bool("v", false, "view mode")
	lineNumbers = flag.Bool("n", false, "line numbers")
	startLine   = flag.Int("s", 1, "start line")
	endLine     = flag.Int("e", 0, "end line")
	findString  = flag.String("f", "", "find string")
	helpFlag    = flag.Bool("h", false, "help")
	initFlag    = flag.Bool("init", false, "initialize config")
)

var flagMap = map[string]string{
	"-init": "initialize config",
	"-v":    "view mode",
	"-n":    "line numbers",
	"-s":    "start line",
	"-e":    "end line",
	"-f":    "find string",
	"-h":    "help",
}

func resolveFileType(t string) string {
	switch t {
	case ".wsb":
		return ".xml"
	case ".tsx":
		return ".ts"
	case ".jsx":
		return ".js"
	default:
		return t
	}
}

func main() {
	flag.Parse()
	c, err := loadConfig()
	if err != nil {
		fmt.Println("error loading config")
		return
	}

	if *helpFlag {
		printHelp(flagMap)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("no file to read")
		return
	}

	filename := args[0]
	if strings.HasSuffix(filename, ".exe") {
		return
	}
	if strings.Contains(filename, "@") {
		parts := strings.Split(filename, "@")
		filename = parts[0]
		lines := strings.Split(parts[1], ":")
		if len(lines) > 1 {
			*startLine, _ = strconv.Atoi(lines[0])
			*endLine, _ = strconv.Atoi(lines[1])
		} else {
			*startLine, _ = strconv.Atoi(lines[0])
		}
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	defer file.Close()

	filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + resolveFileType(filepath.Ext(filename))
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)

	style := styles.Get(c.SyntaxHighlightingStyle)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file")
		return
	}

	if *endLine == 0 {
		*endLine = bytes.Count(contents, []byte("\n")) + 1
	}

	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		fmt.Println("error tokenizing file")
		return
	}

	w := bytes.Buffer{}
	err = formatter.Format(&w, style, iterator)
	if err != nil {
		fmt.Println("error formatting file")
		return
	}
	defer w.Reset()

	if *viewMode {
		Show(filename, w.String())
		clearScreen()
		return
	}
	if *lineNumbers {
		useLineNumbers(os.Stdout, &w, *startLine, *endLine)
		return
	}

	if *findString != "" {
		find(&w, *findString)
	}

	useDefault(w)

}

func printHelp(flags map[string]string) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for name, description := range flags {
		fmt.Fprintf(tw, "%s\t%s\n", helpTextStyle.Render(name), helpTextStyle.Render(description))
	}
	tw.Flush()
}

func find(r io.Reader, s string) {
	scanner := bufio.NewScanner(r)
	i := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(removeAsciiEscapeCodes(scanner.Text()), s) {
			fmt.Printf(
				"%s %s\n",
				findLineNumberStyle.Render(strconv.Itoa(i)),
				line,
			)
		} else {
			fmt.Printf("%s %s\n", strconv.Itoa(i), line)
		}
		i += 1
	}
}

func useLineNumbers(w io.Writer, r io.Reader, start, end int) {
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#677d8a"))
	scanner := bufio.NewScanner(r)
	i := 1
	for scanner.Scan() {
		if i >= start && (end == 0 || i <= end) {
			line := scanner.Text()
			s := strconv.Itoa(i)
			fmt.Fprintf(w, "%s %s\n", lineStyle.Render(s), line)
		}
		i += 1
	}
}

// useDefault prints the contents of a buffer to stdout.
// If stdout is a terminal (i.e. not piped) it will print the
// buffer as is. If stdout is piped, it will remove any ASCII
// escape codes from the buffer.
func useDefault(b bytes.Buffer) {
	if info, _ := os.Stdout.Stat(); (info.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println(b.String())
	} else {
		fmt.Println(removeAsciiEscapeCodes(b.String()))
	}
}

func removeAsciiEscapeCodes(s string) string {
	var inEscapeCode bool
	var result string
	for _, c := range s {
		if c == '\x1b' {
			inEscapeCode = true
		}
		if !inEscapeCode {
			result += string(c)
		}
		if c == 'm' {
			inEscapeCode = false
		}
	}
	return result
}
