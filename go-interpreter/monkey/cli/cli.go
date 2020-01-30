package cli

import (
	"github.com/spf13/cobra"

	lexerRepl "github.com/wreulicke/go-sandbox/go-interpreter/monkey/lexer/repl"
	parserRepl "github.com/wreulicke/go-sandbox/go-interpreter/monkey/parser/repl"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use: "monkey",
	}
	c.AddCommand(NewLexerCommand(), NewParserCommand())
	return c
}

func NewLexerCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "lexer",
		Short: "l",
		Run: func(cmd *cobra.Command, args []string) {
			lexerRepl.Start()
		},
	}
	return c
}

func NewParserCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "parser",
		Short: "p",
		Run: func(cmd *cobra.Command, args []string) {
			parserRepl.Start()
		},
	}
	return c
}
