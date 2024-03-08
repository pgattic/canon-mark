
# Marks - Annotation Manager

## Usage

Marks can received piped information through the standard input, informing it where to search for annotations.

Standard input commands are as follows:

- `@` - Specify path to the directory within `~/.marks/[profile]/` where the annotations will be queried
- `@@` - Specify the name of the file within the directory where the annotations will be sourced
- `@@@` - Specify the line number within the file where the annotations will be sourced

Marks will search for your saved annotations in the `~/.marks/` directory, under a specific profile. If no profile is specified, all profiles will be searched. Thus, in this case, marks will search in all immediate subdirectories under `~/.marks/` for `nt-kjv-canon/John/3.json`, applying annotations/highlights that overlap with verses 16 or 17, and apply them to the standard output. 

### Displaying annotated text

A good example would be from [Canon](https://github.com/pgattic/canon)'s output of `canon "John 3:16-17" -v` (-v for "verbose"):

```
@nt-kjv-canon/John/
@@3
@@@16 For God so loved the world, that he gave his only begotten Son, that whosoever believeth in him should not perish, but have everlasting life.
@@@17 For God sent not his Son into the world to condemn the world; but that the world through him might be saved.
```

if this were to be piped into `marks`, the text would be printed out with the appropriate highlights applied, such as `canon "John 3:16-17 -v | marks -n"`.

If being used frequently with Canon, a few bash functions may be useful to simplify the process of displaying annotated text, such as:

```
can() { canon "$@" -v | marks -n; }

canl() { can "$@" | less --wordwrap -R; }
```

### Storing annotations

While marks is reading input text, it keeps track of which paragraphs it is given, and can be told to store that range in an annotation in the user's directory.

For example: telling Marks to store an annotation for John 3:16-17 would be as follows:

```
canon "John 3:16-17" -v | marks add --ul --bg=blue
```

As of yet, the command line options are as follows:

- `--ul` - Underline the text
- `--bg=color` - Set the background color
- `--fg=color` - Set the foreground color

More options are planned, such as tags and links to other references.

## Development Environment

In the process of developing this software, I used:

- Neovim (Code editor)
- Golang (Programming langauge)

## Useful Websites

- [Shakiba Mosihiri's answer on StackOverflow](https://stackoverflow.com/a/28938235/23116882)
- [ANSI Escape Codes Generator](https://ansi.gabebanks.net/)
- [ANSI Escape Code - Wikipedia](https://en.wikipedia.org/wiki/ANSI_escape_code)

## Roadmap

- [ ] Highlight words within paragraphs
- [x] Support adding a new annotation from the command line (nobody wants to do that manually!)
- [ ] Support multiple profiles
- [ ] Support filtering by tags
- [ ] Hold onto the date that each annotation was added
- [ ] List all annotations that overlap the input range
- [ ] Delete annotations

