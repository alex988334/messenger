package data

type ILink interface {
}

type Link struct {
	ParentModel  string
	ParentKey    string
	CurrentModel string
	CurrentKey   string
	Weight       int
}

func NewLink(parentModel string, parentKey string, currentModel string, currentKey string, weight int) *Link {

	return &Link{
		ParentModel:  parentModel,
		ParentKey:    parentKey,
		CurrentModel: currentModel,
		CurrentKey:   currentKey,
		Weight:       weight,
	}
}
