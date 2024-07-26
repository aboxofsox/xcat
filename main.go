package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"
)

var MaxFileSize int64 = 1024 * 1024

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
	initFlag    = flag.Bool("i", false, "initialize config")
	bigFile     = flag.Bool("b", false, "big file")
)

var flagMap = map[string]string{
	"-b": "big file",
	"-v": "view mode",
	"-n": "line numbers",
	"-s": "start line",
	"-e": "end line",
	"-f": "find string",
	"-h": "help",
	"-i": "initialize config",
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

	if c.LineNumberColor != "" {
		lineNumberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(c.LineNumberColor))
	}

	if c.MaxFileSize != 0 {
		MaxFileSize = c.MaxFileSize
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

	filename := resolve(args[0])

	ext := filepath.Ext(filename)

	// ignore disallowed file types
	// by default, binary executables and library files are ignored
	// add or remove types in the config file
	if slices.Contains(c.DisallowedFileTypes, ext) {
		return
	}

	// @ symbol is used to specify a specific line number range.
	// example: foo.txt@1:10 will display lines 1-10.
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

	info, err := os.Stat(filename)
	if err != nil {
		fmt.Println("error reading file")
		return
	}

	// sometimes you might not want to display a large file
	if info.Size() > MaxFileSize && !*bigFile {
		fmt.Println("file too large to display, use -big flag to display")
		return
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
		clearScreen()
		Show(filename, w.String(), *lineNumbers)
		return
	}
	if *lineNumbers {
		// useLineNumbers(os.Stdout, &w, *startLine, *endLine)
		out, err := useLineNumbers(&w, *startLine, *endLine)
		if err != nil {
			fmt.Println("error using line numbers")
			return
		}
		w.Write(out)
	}

	if *findString != "" {
		find(&w, *findString)
	}

	usePipe(w)

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
		if strings.Contains(escape(scanner.Text()), s) {
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

func useLineNumbers(r io.Reader, start, end int) ([]byte, error) {
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#677d8a"))
	scanner := bufio.NewScanner(r)
	buff := bytes.Buffer{}
	i := 1
	for scanner.Scan() {
		if i >= start && (end == 0 || i <= end) {
			line := scanner.Text()
			s := strconv.Itoa(i)
			// fmt.Fprintf(w, "%s %s\n", lineStyle.Render(s), line)
			buff.WriteString(fmt.Sprintf("%s %s\n", lineStyle.Render(s), line))
		}
		i += 1
	}
	return buff.Bytes(), nil
}

// usePipe prints the contents of a buffer to stdout.
// If stdout is a terminal (i.e. not piped) it will print the
// buffer as is. If stdout is piped, it will remove any ASCII
// escape codes from the buffer.
func usePipe(b bytes.Buffer) {
	if info, _ := os.Stdout.Stat(); (info.Mode() & os.ModeCharDevice) != 0 {
		// if err := writeString(b.String(), os.Stdout); err != nil {
		// 	panic(err)
		// }
		fmt.Println(b.String())
	} else {
		b = removeLineNumbers(b)
		// if err := writeString(escape(b.String()), os.Stdout); err != nil {
		// 	panic(err)
		// }
		fmt.Println(escape(b.String()))
	}
}

func writeString(s string, w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}

func resolve(filename string) string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(filename, "~/") {
		filename = strings.Replace(filename, "~", home, 1)
	}
	return filename
}

func removeLineNumbers(b bytes.Buffer) bytes.Buffer {
	var result bytes.Buffer
	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		result.WriteString(parts[1] + "\n")
	}
	return result

}

func escape(s string) string {
	var inEscapeCode bool
	var result strings.Builder
	for i := 0; i < len(s); {
		r, size := decodeRune(s[i:])
		if r == '\x1b' {
			inEscapeCode = true
			i += size
			continue
		}
		if inEscapeCode {
			if r == 'm' {
				inEscapeCode = false
			}
			i += size
			continue
		}
		result.WriteRune(r)
		i += size
	}
	return result.String()
}

func decodeRune(s string) (rune, int) {
	return utf8.DecodeRuneInString(s)
}
