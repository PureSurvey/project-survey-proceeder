package objects

type Template struct {
	Id            int
	Code          string
	DefaultParams string
}

func NewTemplate(id int, code string, defaultParams string) *Template {
	return &Template{
		Id:            id,
		Code:          code,
		DefaultParams: defaultParams,
	}
}
