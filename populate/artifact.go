package populate

// Standard manifest artifact (just points to a file for all of the data)
type Artifact struct {
	TemplateFileName string `json:"templateFileName"`
}

// Standard manifest
type Manifest map[string][]Artifact
