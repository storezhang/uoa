package uoa

// 内部接口封装
// 使用模板方法设计模式
type keyTemplate struct {
	key         Key
	environment string
}

func (t *keyTemplate) Paths() (paths []string) {
	paths = t.key.Paths()
	if "" != t.environment {
		paths = append([]string{t.environment}, paths...)
	}

	return
}
