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
  -v        view mode
  -n        line numbers
  -s        start line
  -e        end line
  -f        find string
  -h        help
```

## Config
When you run the tool for the first time, a `config.json` file is created in `~/.config/xcat`.
```json
{
        "syntax_highlighting_style": "monokai"
}
```

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
You can pipe the output of `xcat`. It doesn't work if you use line numbers, unfortunately, but soon™️. When piping, the output of `xcat` is trimmed of all ASCII escape characters.
```powershell
xcat main.go | clip
```
```powershell
xcat main.go | % { $_ }
```

