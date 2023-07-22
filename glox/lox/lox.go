package lox

import (
	"bufio"
	"fmt"
	"glox/interpreter"
	"glox/lexer"
	"glox/parser"
	"io"
	"os"
)

const PROMPT = ">> "

type Runner interface {
	RunFile(path string) error
	StartREPL(stdin io.Reader)
}

func NewRunner(stdErr io.Writer, stdout io.Writer) Runner {
	return &Lox{stdErr: stdErr, stdout: stdout}
}

type Lox struct {
	stdErr io.Writer
	stdout io.Writer
}

func (r *Lox) RunFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	bytes := make([]byte, info.Size())
	if _, err = bufio.NewReader(file).Read(bytes); err != nil {
		return err
	}
	r.run(string(bytes))
	return nil
}

func (r *Lox) StartREPL(stdin io.Reader) {
	scanner := bufio.NewScanner(stdin)

	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		r.run(line)
	}
}

func (r *Lox) run(src string) {
	scnr := lexer.New(src)

	tokens := scnr.Tokenize()
	prsr := parser.New(tokens)
	exp, err := prsr.Parse()

	if err != nil {
		fmt.Fprintf(r.stdErr, "%v\n", err.Error())
		return
	}

	glox := interpreter.New()
	out := glox.Interpret(exp)
	if err, isErr := out.(error); isErr {
		fmt.Fprintf(r.stdErr, "%v\n", err.Error())
		return
	}

	fmt.Fprintf(r.stdout, "%v\n", out)

}
