package objects

type Translation struct {
	Id          int
	Translation string
	Language    string
	ParentId    int
}

func NewTranslation(id int, translation string, language string, parentId int) *Translation {
	return &Translation{
		Id:          id,
		Translation: translation,
		Language:    language,
		ParentId:    parentId,
	}
}
