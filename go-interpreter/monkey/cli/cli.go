package cli

import (
	"github.com/spf13/cobra"

	lexerRepl "github.com/wreulicke/go-sandbox/go-interpreter/monkey/lexer/repl"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use: "monkey",
	}
	c.AddCommand(NewLexerCommand())
	return c
}

func NewLexerCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "lexer",
		Short: "s",
		Run: func(cmd *cobra.Command, args []string) {
			lexerRepl.Start()
		},
	}
	return c
}
