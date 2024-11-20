package interp

import (
	"fmt"
	"reflect"
	"strings"
)

func printFrameData(data []reflect.Value) {
	fmt.Printf("frame data:\n")
	for i, d := range data {
		fmt.Printf("data[%d]: %#v\n", i, d)
	}
}

func printNodeTree(n *node, depth int) {
	fmt.Printf("%s%s, n.ident: %q, n.findex: %d, n.level: %d\n", strings.Join(make([]string, depth+1), " "), kinds[n.kind], n.ident, n.findex, n.level)
	for _, c := range n.child {
		printNodeTree(c, depth+2)
	}
}
