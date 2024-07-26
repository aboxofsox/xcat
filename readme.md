# xcat
A better, but also worse, version of `cat`.

## Benchmarks
`xcat` is about twice as slow as `cat` because the syntax highlighter parses each file and uses ASCII escape characters.  
```
cat avgerage execution time: 17.569848ms
xcat average execution time: 32.499113ms
```

## Usage
```
  -b     big file
  -v     view mode
  -n     line numbers
  -s     start line
  -e     end line
  -f     find string
  -h     help
  -i     initialize config
```

## Config
When you run the tool for the first time, a `config.json` file is created in `~/.config/xcat`.

You can also create a new config with the below default values with `xcat -i`.

```json
{
        "syntax_highlighting_style": "monokai",
        "line_number_color": "#677d8a",
        "disallowed_filet_ypes": ["exe", "dll", "so", "dylib", "bin", "o", "a", "lib"],
        "size": 0
}
```
`Size` is the max allowed file size. Use `0` if you don't want to set a limit.

### Examples

#### Line numbers and range
```
xcat -n main.go@212:227
```
```go
214 func removeAsciiEscapeCodes(s string) string {
215     var inEscapeCode bool
216     var result string
217     for _, c := range s {
218             if c == '\x1b' {
219                     inEscapeCode = true
220             }
221             if !inEscapeCode {
222                     result += string(c)
223             }
224             if c == 'm' {
225                     inEscapeCode = false
226             }
227     }
```

#### Piping
You can pipe the output of `xcat`. When piping, the output of `xcat` is trimmed of all ASCII escape characters.
```powershell
xcat main.go | clip
```
```powershell
xcat main.go | % { $_ }
```

