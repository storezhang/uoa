package uoa

// Key 一个可以生成文件键的对象
type Key interface {
	// Paths 对应路径，按路径划分
	Paths() []string
}
