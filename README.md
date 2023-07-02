# Implementation of the [Lox Programming language](https://craftinginterpreters.com/the-lox-language.html)

[WIP]
---

## Quick Start
You start playing with the REPL by running the main file located in `crafting-interpreters/glox`
```sh
go run main.go
```

## Grammar

Production rules:

```txt
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
