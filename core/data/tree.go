package data

type (
	Tree[N comparable] = struct {
		Node
		Children []Tree
	}
)
