package elastic

type OnError struct {
	Error struct {
		RootCause []struct {
			Type         string `json:"type"`
			Reason       string `json:"reason"`
			ResourceType string `json:"resource.type"`
			ResourceID   string `json:"resource.id"`
			IndexUUID    string `json:"index_uuid"`
			Index        string `json:"index"`
		} `json:"root_cause"`
		Type         string `json:"type"`
		Reason       string `json:"reason"`
		ResourceType string `json:"resource.type"`
		ResourceID   string `json:"resource.id"`
		IndexUUID    string `json:"index_uuid"`
		Index        string `json:"index"`
	} `json:"error"`
	Status int `json:"status"`
}

type OnErrorBulkOperation struct {
	Error struct {
		Type      string `json:"type"`
		Reason    string `json:"reason"`
		IndexUUID string `json:"index_uuid"`
		Shard     string `json:"shard"`
		Index     string `json:"index"`
	} `json:"error"`
}
