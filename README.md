<samp>

# Implementation of the [Lox Programming language](https://craftinginterpreters.com/the-lox-language.html)

## [WIP]

## Quick Start

> Source files are not yet supported... will handle this later. Currently the safest way to play with the interpreter is using the REPL.

---

You have two options to run the REPL, you can either install the latest released binary, or you can clone the repo locally and run the `main.go` file.

### Using Release Binary


The command below will download the latest released binary in your system, once downloaded, you can start the REPL by executing the downloaded binary in the directory where you ran the command.

> Disclaimer: The command does not work on windows, if you use windows, you can manually open the [latest release page](https://github.com/silverhairs/lox/releases/latest) and download the `glox-windows-amd64.exe` file.

```sh
curl -fsSL https://raw.githubusercontent.com/silverhairs/lox/main/install.sh | sh
```

**About the released binary**:

- You need to have [cURL](https://curl.se/) installed in your system for the above command to work.
- I don't have a specific release schedule, I just cut a new release whenever a major feature has been implemented.

### Using source code

1. Clone the repo
2. Navigate to `lox/glox` (command: `cd ./glox`)
3. Run the main.go file (command: `go run main.go`)

**Requirements**:

- You need to have Go installed in your system

## Grammar

Production rules:

```bnf
    program    -> declaration* EOF ;
    declaration-> funDecl
                | letDecl
                | statement ;
    funDecl    -> "fun" function ;
    function   -> IDENTIFIER "(" parameters ")" block ;
    parameters -> IDENTIFIER ("," IDENTIER)* :
    letDecl    -> ("var" | "let") IDENTIFIER ("=" expression) ? ";" ;

    statement  -> exprStmt
                | ifStmt
                | printStmt
                | blockStmt ;
    exprStmt   -> expression ";" ;
    ifStmt     -> "if" "(" expression ")" statement
                ("else" statement)? ;
    printStmt  -> "print" expression ";" ;
    blockStmt  -> "{" declaration* "}" ;

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

## Example code snippet

> For now the interpreter does not handle source files, so this is to be done in the REPL (line by line).

```js
let age = 15;
let is_adult = age >= 18;

if (is_adult){
    print "You are an adult";
  } else{
    print "You cannot drink";
  }

 {
   let age = age+3;
   print age;
   let message = age >= 18 ? "is an adult" : "cannot drink";
   print message;
 }
print age;
```

</samp>
