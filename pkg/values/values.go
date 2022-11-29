package values

import (
	"html/template"
	"strings"
	"text/template/parse"

	"github.com/divolgin/helmsplain/pkg/log"
	"github.com/pkg/errors"
)

var nodeTypes = map[parse.NodeType]string{
	parse.NodeText:    "NodeText",
	parse.NodeAction:  "NodeAction",
	parse.NodeBool:    "NodeBool",
	parse.NodeChain:   "NodeChain",
	parse.NodeCommand: "NodeCommand",
	parse.NodeDot:     "NodeDot",

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

func GetFromFiles(filenames ...string) ([]string, error) {
	t, err := template.New("").Funcs(funcMap()).ParseFiles(filenames...)
	if err != nil {
		return nil, errors.Wrap(err, "parse files")
	}

	result := []string{}
	for _, tmpl := range t.Templates() {
		if tmpl.Tree == nil {
			continue
		}
		result = append(result, getFromTree(tmpl.Tree.Root)...)
	}

	return result, nil
}

func GetFromData(data string) []string {
	t := template.Must(template.New("t").Funcs(funcMap()).Parse(data))
	return getFromTree(t.Tree.Root)
}

func getFromTree(node parse.Node) []string {
	result := []string{}

	log.Debugf("node:%s -> %s\n", nodeTypes[node.Type()], node.String())

	if node.Type() == parse.NodeField {
		ref := node.String()
		if strings.HasPrefix(ref, ".Values.") {
			result = append(result, ref)
		}
	} else if node.Type() == parse.NodeCommand {
		ref := node.String()
		if strings.HasPrefix(ref, ".Values.") {
			result = append(result, ref)
		} else {
			for _, arg := range node.(*parse.CommandNode).Args {
				result = append(result, getFromTree(arg)...)
			}
		}
	} else if node.Type() == parse.NodePipe {
		for _, cmd := range node.(*parse.PipeNode).Cmds {
			result = append(result, getFromTree(cmd)...)
		}
	} else if node.Type() == parse.NodeAction {
		for _, cmd := range node.(*parse.ActionNode).Pipe.Cmds {
			result = append(result, getFromTree(cmd)...)
		}
	} else if node.Type() != parse.NodeText { // text nodes dump a lot in the output
		log.Debugf("%#v\n", node)
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			refs := getFromTree(n)
			result = append(result, refs...)
		}
	}

	return result
}
