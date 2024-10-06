package models

type Stats struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	CPU       int   `json:"cpu"`
	Memory    int   `json:"memory"`
	Storage   int   `json:"storage"`
	Network   int   `json:"network"`
	Uptime    int   `json:"uptime"`
	CreatedAt int64 `json:"createdAt" gorm:"autoCreateTime"`
}
