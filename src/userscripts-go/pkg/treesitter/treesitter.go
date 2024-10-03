package treesitter

// #cgo CFLAGS: -std=c11 -fPIC -I./src
// #include "parser.c"
// #include "scanner.c"
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

// Get the tree-sitter Language for this grammar.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_latex())
}

// WalkNodes recursively walks through the syntax tree nodes.
func WalkNodes(node *sitter.Node, source []byte, indent int) {
	if node == nil {
		return
	}

	// Print the node type and its text
	fmt.Printf("%s%s: %s\n",
		strings.Repeat("..", indent),
		node.Type(),
		string(source[node.StartByte():node.EndByte()]))

	// Recursively walk through each child node
	for i := 0; i < int(node.ChildCount()); i++ {
		WalkNodes(node.Child(i), source, indent+2)
	}
}

// GenericCommand holds the command and its argument nodes.
type GenericCommand struct {
	CommandNode  *sitter.Node // The node for the LaTeX command
	ArgumentNode *sitter.Node // The node for the command's argument
}

// FindGenericCommand searches for a specific LaTeX command and extracts its arguments.
func FindGenericCommand(node *sitter.Node, source []byte, command string) []GenericCommand {
	var foundCommands []GenericCommand

	if node == nil {
		return foundCommands
	}

	// Check if the current node is a generic_command
	if node.Type() == "generic_command" {
		// Extract the command name node
		commandNode := node.Child(0)

		// Check if the command matches the specified command
		if string(source[commandNode.StartByte():commandNode.EndByte()]) == "\\"+command {
			// The argument is in the curly_group child node (if it exists)
			if node.ChildCount() > 1 && node.Child(1).Type() == "curly_group" {
				argumentNode := node.Child(1) // Get the argument node
				foundCommands = append(foundCommands, GenericCommand{
					CommandNode:  commandNode,
					ArgumentNode: argumentNode,
				})
			}
		}
	}

	// Recursively search through each child node
	for i := 0; i < int(node.ChildCount()); i++ {
		foundCommands = append(foundCommands, FindGenericCommand(node.Child(i), source, command)...)
	}

	return foundCommands
}
