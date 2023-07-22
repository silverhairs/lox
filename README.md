# Implementation of the [Lox Programming language](https://craftinginterpreters.com/the-lox-language.html)

## [WIP]

## Quick Start

### From Binary (only available on macos and linux):

Paste the script below in your terminal and press **Enter**.

The script will download the latest `glox` binary release in your system and give it executable permission. After this script, you can run `./glox` and it should start the REPL.

> Make sure you have [curl](https://curl.se/) installed in your system before running the script

```sh
curl -L -s https://api.github.com/repos/silverhairs/crafting-interpreters/releases/latest \
| grep "browser_download_url.*glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m)" \
| cut -d '"' -f 4 \
| wget -qi - \
&& chmod +x glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m) \
&& sudo mv glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m) $PWD/glox
```

### From Source

Cloning the repository and runing `main.go` should start the REPL. `main.go` is located in `crafting-interpreters/glox`. You need to have Golang installed in your system for this option.

```sh
go run main.go
```

## Grammar

Production rules:

```txt
    program    -> statement* EOF ;
    statement  -> exprStmt
                | printStmt ;
    exprStmt   -> expression ";" ;
    printStmt  -> "print" expression ";" ;

    expression -> literal
                | unary
                | binary
                | grouping
                | ternary ;

    literal    -> NUMBER | STRING | boolean | "nil" ;
    unary      -> ( "-" | "!" ) expression ;
    grouping   -> "(" expression ")" ;
    binary     -> expression operator expression ;
    operator   -> "==" | "!=" | "<" | ">" | "<=" | ">="
                | "+" | "-" | "*" | "/" ;
    ternary -> expression "?" expression ":" expression;
```
