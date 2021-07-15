package uoa

// Path 一个可以生成文件键的对象
type Path interface {
	// Paths 对应路径，按路径划分
	// 如果不需要目录路径，可以只返回一个数据的数组
	Paths() []string
}
