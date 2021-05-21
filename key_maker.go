package uoa

// 内部接口封装
// 使用模板方法设计模式
type keyMaker struct {
	key         Key
	environment string
}

func (k *keyMaker) Paths() (paths []string) {
	paths = k.key.Paths()
	if "" != k.environment {
		paths = append([]string{k.environment}, paths...)
	}

	return
}
