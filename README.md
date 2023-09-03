<samp>

# Implementation of the [Lox Programming language](https://craftinginterpreters.com/the-lox-language.html)

## [WIP]

## Quick Start

> Source files are not yet supported... will handle this later. Currently the safest way to play with the interpreter is using the REPL.

---

You have two options to run the REPL, you can either install the latest released binary, or you can clone the repo locally and run the `main.go` file.

### Using Release Binary

The script below will download the latest released binary in your system and request for execution permission, once granted, you can start the REPL by executing the command `./glox` in the directory where you ran the script.

```sh
curl -L -s https://api.github.com/repos/silverhairs/crafting-interpreters/releases/latest \
| grep "browser_download_url.*glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m | sed 's/x86_64/amd64/')" \
| cut -d '"' -f 4 \
| wget -qi - \
&& chmod +x glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m | sed 's/x86_64/amd64/') \
&& mv glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m | sed 's/x86_64/amd64/') $PWD/glox
```

**About the released binary**:

- You need to have [cURL](https://curl.se/) installed in your system for the above script to work.
- I don't have a specific release schedule, I just cut a new release whenever a major feature has been implemented.
- Currently the latest release only holds binaries for macos and linux, if you have windows... idk what to tell yah (seek help i guess). I will include windows later in the release workflow though, I just need to get a windows computer to test.

### Using source code

1. Clone the repo
2. Navigate to `crafting-interpreters/glox` (command: `cd ./glox`)
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
                | returnStmt
                | printStmt
                | blockStmt ;
    exprStmt   -> expression ";" ;
    ifStmt     -> "if" "(" expression ")" statement
                ("else" statement)? ;
    printStmt  -> "print" expression ";" ;
    blockStmt  -> "{" declaration* "}" ;
    returnStmt -> "return" expression? ";" ;

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
