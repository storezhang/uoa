package uoa

type (
	cosStatement struct {
		// 操作
		Actions []string `json:"action,omitempty"`
		// 效力
		Effect string `json:"effect,omitempty"`
		// 资源
		Resources []string `json:"resource,omitempty"`
		// 生效条件
		Condition map[string]map[string]interface{} `json:"condition,omitempty"`
	}

	cosPolicy struct {
		// 版本
		Version string `json:"version,omitempty"`
		// 语句
		Statements []cosStatement `json:"statement,omitempty"`
	}
)
