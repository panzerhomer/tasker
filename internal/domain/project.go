package domain

import "errors"

type Project struct {
	ID          int64  `json:omitempty`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p *Project) Validate() error {
	if len(p.Name) < 1 {
		return errors.New("project name must be more than 0 char")
	}

	return nil
}
