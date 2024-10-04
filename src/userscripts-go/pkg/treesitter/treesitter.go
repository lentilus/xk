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

// GenericEnvironment holds the environment and its begin-node (often containing arguments)
type GenericEnvironment struct {
	EnvironmentNode *sitter.Node   // The node for the environment
	ArgumentNodes   []*sitter.Node // The node for the environment's \begin
}

// FindGenericEnvironment searches for LaTeX environments and extracts their \begin nodes.
func FindGenericEnvironment(
	node *sitter.Node,
	source []byte,
	environment string,
) []GenericEnvironment {
	var foundEnvironments []GenericEnvironment

	if node == nil {
		return foundEnvironments
	}

	// Check if the current node is a generic_environment
	if node.Type() == "generic_environment" {
		// Get the \begin node
		beginNode := node.ChildByFieldName("begin")

		// Ensure the \begin node exists
		if beginNode != nil {
			// Inside the begin node, check for the curly_group after \begin that contains the environment name
			nameNode := beginNode.Child(0).NextSibling()
			name := nameNode.Content(source)

			// Ensure the curly_group exists
			if nameNode != nil && nameNode.Type() == "curly_group_text" {
				// Check if the curly_group text is the expected environment name (e.g., "{flashcard}")
				if name == "{"+environment+"}" {
					// Store the found environment and its begin node
					argNodes := []*sitter.Node{}

					// Append remaining siblings to argNodes
					for sibling := nameNode.NextSibling(); sibling != nil; sibling = sibling.NextSibling() {
						argNodes = append(argNodes, sibling)
					}

					foundEnvironments = append(foundEnvironments, GenericEnvironment{
						EnvironmentNode: node,
						ArgumentNodes:   argNodes,
					})
				}
			}
		}
	}

	// Recursively search through each child node
	for i := 0; i < int(node.ChildCount()); i++ {
		foundEnvironments = append(
			foundEnvironments,
			FindGenericEnvironment(node.Child(i), source, environment)...)
	}

	return foundEnvironments
}
