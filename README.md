# rep

Simple CLI tool for text replacement.

## Installation

```sh
go install github.com/nek023/rep@latest
```

## Usage

You can replace text by starting rep and typing a command.  
The command `%s/foo/bar` will replace `foo` with `bar`.

```sh
# Replace text by passing data through a pipe
echo example | rep 

# Replace the content of a file
rep hello.txt
```
