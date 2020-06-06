# pwgen - a password generator for memorable passwords

This is a simple password generator for memorable passwords.

This is the default dictionary: https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt

## Usage

```
pwgen is a password generator for memorable passwords
If no dictionary is supplied a default dictionary is downloaded.

Usage:
  pwgen [flags]

Flags:
  -c, --count uint32       the number of passwords to generate (default 10)
  -d, --dict string        the file to read/save the dicionary from/to (default "dict.txt")
  -h, --help               help for pwgen
      --max-length uint8   the maximum length of each word (default 10)
      --min-length uint8   the minimum length of each word (default 3)
  -s, --separator string   the separator to use between each word (default "-")
  -w, --words uint16       the number of words to use per password (default 5)
```
