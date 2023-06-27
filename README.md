# Implementation of the [Lox Programming language](https://craftinginterpreters.com/the-lox-language.html)

[WIP]

## Grammar

Production rules:

```txt
    expression -> literal
                | unary
                | binary
                | grouping ;

    literal    -> NUMBER | STRING | boolean | "nil" ;
    unary      -> ( "-" | "!" ) expression ;
    grouping   -> "(" expression ")" ;
    binary     -> expression operator expression ;
    operator   -> "==" | "!=" | "<" | ">" | "<=" | ">="
                | "+" | "-" | "*" | "/" ;
```
