package meta

type Paging struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`

	BeforeID uint `json:"before_id"`
	AfterID  uint `json:"after_id"`
}

func (this *Paging) SetPage(page uint) {
	this.Offset = page * this.Limit
}
