package files

type TreeNode struct {
	Name     string      `json:"name"`
	IsDir    bool        `json:"isDir"`
	Children []*TreeNode `json:"children,omitempty"`
}
