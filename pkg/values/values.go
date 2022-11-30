package values

import (
	"html/template"
	"strings"
	"text/template/parse"

	"github.com/divolgin/helmsplain/pkg/log"
	"github.com/pkg/errors"
)

// this is for nice output only
var nodeTypes = map[parse.NodeType]string{
	parse.NodeText:       "NodeText",
	parse.NodeAction:     "NodeAction",
	parse.NodeBool:       "NodeBool",
	parse.NodeChain:      "NodeChain",
	parse.NodeCommand:    "NodeCommand",
	parse.NodeDot:        "NodeDot",
	parse.NodeField:      "NodeField",
	parse.NodeIdentifier: "NodeIdentifier",
	parse.NodeIf:         "NodeIf",
	parse.NodeList:       "NodeList",
	parse.NodeNil:        "NodeNil",
	parse.NodeNumber:     "NodeNumber",
	parse.NodePipe:       "NodePipe",
	parse.NodeRange:      "NodeRange",
	parse.NodeString:     "NodeString",
	parse.NodeTemplate:   "NodeTemplate",
	parse.NodeVariable:   "NodeVariable",
	parse.NodeWith:       "NodeWith",
	parse.NodeComment:    "NodeComment",
	parse.NodeBreak:      "NodeBreak",
	parse.NodeContinue:   "NodeContinue",
}

type Value struct {
	Key string
	Pos parse.Pos
}

func GetFromFiles(filenames ...string) ([]Value, error) {
	t, err := template.New("").Funcs(funcMap()).ParseFiles(filenames...)
	if err != nil {
		return nil, errors.Wrap(err, "parse files")
	}

	result := []Value{}
	for _, tmpl := range t.Templates() {
		if tmpl.Tree == nil {
			continue
		}
		result = append(result, getFromNode(tmpl.Tree.Root, "")...)
	}

	return result, nil
}

func GetFromData(data string) []Value {
	t := template.Must(template.New("t").Funcs(funcMap()).Parse(data))
	return getFromNode(t.Tree.Root, "")
}

func getFromNode(node parse.Node, withPrefix string) []Value {
	result := []Value{}

	log.Debugf("node:%s -> %s\n", nodeTypes[node.Type()], node.String())

	switch node.Type() {
	case parse.NodeField:
		return getFromFieldNode(node.(*parse.FieldNode), withPrefix)

	case parse.NodeCommand:
		return getFromCommandNode(node.(*parse.CommandNode), withPrefix)

	case parse.NodePipe:
		for _, cmd := range node.(*parse.PipeNode).Cmds {
			result = append(result, getFromNode(cmd, withPrefix)...)
		}

	case parse.NodeAction:
		for _, cmd := range node.(*parse.ActionNode).Pipe.Cmds {
			result = append(result, getFromNode(cmd, withPrefix)...)
		}

	case parse.NodeIf:
		return getFromIfNode(node.(*parse.IfNode), withPrefix)

	case parse.NodeWith:
		return getFromWithNode(node.(*parse.WithNode), withPrefix)

	case parse.NodeList:
		for _, n := range node.(*parse.ListNode).Nodes {
			refs := getFromNode(n, withPrefix)
			result = append(result, refs...)
		}

	case parse.NodeText:
		// no-op for now

	default:
		// text nodes dump a lot in the output
		log.Debugf("%#v\n", node)
	}

	return result
}

func getFromFieldNode(node *parse.FieldNode, withPrefix string) []Value {
	result := []Value{}

	ref := node.String()
	if strings.HasPrefix(ref, ".") {
		ref = withPrefix + ref
		if strings.HasPrefix(ref, ".Values.") {
			result = append(result, Value{
				Key: ref,
				Pos: node.Position(),
			})
		}
	}

	return result
}

func getFromCommandNode(node *parse.CommandNode, withPrefix string) []Value {
	result := []Value{}

	ref := node.String()
	if strings.HasPrefix(ref, ".") {
		ref = withPrefix + ref
		if strings.HasPrefix(ref, ".Values.") {
			result = append(result, Value{
				Key: ref,
				Pos: node.Position(),
			})
		}
	} else {
		for _, arg := range node.Args {
			result = append(result, getFromNode(arg, withPrefix)...)
		}
	}

	return result
}

func getFromIfNode(node *parse.IfNode, withPrefix string) []Value {
	result := []Value{}

	for _, cmd := range node.Pipe.Cmds {
		result = append(result, getFromNode(cmd, withPrefix)...)
	}
	if node.List != nil {
		for _, n := range node.List.Nodes {
			result = append(result, getFromNode(n, withPrefix)...)
		}
	}
	if node.ElseList != nil {
		for _, n := range node.ElseList.Nodes {
			result = append(result, getFromNode(n, withPrefix)...)
		}
	}

	return result
}

func getFromWithNode(node *parse.WithNode, withPrefix string) []Value {
	result := []Value{}

	for _, cmd := range node.Pipe.Cmds {
		cmdString := cmd.String()
		if strings.HasPrefix(cmdString, ".") {
			withPrefix = withPrefix + cmdString
		} else {
			result = append(result, getFromNode(cmd, withPrefix)...)
		}
	}
	if node.List != nil {
		for _, n := range node.List.Nodes {
			result = append(result, getFromNode(n, withPrefix)...)
		}
	}
	if node.ElseList != nil {
		for _, n := range node.ElseList.Nodes {
			result = append(result, getFromNode(n, withPrefix)...)
		}
	}

	return result
}
