package lox

import (
	"bufio"
	"fmt"
	"glox/interpreter"
	"glox/lexer"
	"glox/parser"
	"glox/utils"
	"io"
	"os"
)

const PROMPT = ">> "

type Runner interface {
	RunFile(path string) error
	StartREPL(stdin io.Reader)
}

func NewRunner(stdErr io.Writer) Runner {
	return &Lox{stdErr: stdErr}
}

type Lox struct {
	stdErr io.Writer
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
	exp, errs := prsr.Parse()

	if len(errs) > 0 {
		messages := utils.Map[error, string](errs, func(err error) string { return fmt.Sprintf("%s\n", err.Error()) })
		fmt.Fprintf(r.stdErr, "%v", messages)
	}

	glox := interpreter.New()

	fmt.Println(glox.Interpret(exp))

}
