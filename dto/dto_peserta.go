package dto

type PesertaBPJS struct {
	ID     string  `gorm:"primaryKey;type:uuid" json:"id"`
	PSTV01 string  `json:"pstv01"`
	PSTV02 string  `json:"pstv02"`
	PSTV03 string  `json:"pstv03"`
	PSTV04 int     `json:"pstv04"`
	PSTV05 int     `json:"pstv05"`
	PSTV06 int     `json:"pstv06"`
	PSTV07 string  `json:"pstv07"`
	PSTV08 int     `json:"pstv08"`
	PSTV09 string  `json:"pstv09"`
	PSTV10 string  `json:"pstv10"`
	PSTV11 string  `json:"pstv11"`
	PSTV12 int     `json:"pstv12"`
	PSTV13 string  `json:"pstv13"`
	PSTV14 string  `json:"pstv14"`
	PSTV15 float64 `json:"pstv15"`
	PSTV16 int     `json:"pstv16"`
	PSTV17 string  `json:"pstv17"`
	PSTV18 *string `json:"pstv18"`

	CreatedBy string `json:"created_by" gorm:"type:uuid"`
}
