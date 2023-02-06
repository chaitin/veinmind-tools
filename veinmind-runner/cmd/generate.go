//go:generate go run ./ generate doc --path ../docs/
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docPath     string
	docType     string
	generateCmd = &cmd.Command{
		Use:   "generate",
		Short: "Generate relevant information",
	}
	generateDocumentCmd = &cmd.Command{
		Use: "doc",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Stat(docPath)
			if os.IsNotExist(err) {
				err := os.MkdirAll(docPath, 0666)
				if err != nil {
					return err
				}
			}

			switch strings.ToLower(docType) {
			case "markdown", "md":
				return doc.GenMarkdownTree(cmd.Parent().Parent(), docPath)
			case "manpage", "man":
				return doc.GenManTree(cmd.Parent().Parent(), &doc.GenManHeader{}, docPath)
			case "yaml", "yml":
				return doc.GenYamlTree(cmd.Parent().Parent(), docPath)
			case "rest", "rst":
				return doc.GenReSTTree(cmd.Parent().Parent(), docPath)
			default:
				return errors.New(fmt.Sprintf("cmd: doc type doesn't match for %s", docType))
			}
		},
	}
)

func init() {
	generateCmd.AddCommand(generateDocumentCmd)
	generateDocumentCmd.Flags().StringVarP(&docPath, "path", "p", "docs", "document path")
	generateDocumentCmd.Flags().StringVarP(&docType, "type", "t", "markdown", "document type, markdown/manpage/yaml/rest")
}
