package model

type Stats struct {
	UserID         string `gorm:"primaryKey" json:"-"`
	MaxStorage     int64  `json:"maxStorage"`
	UsedStorage    int64  `json:"usedStorage"`
	UploadedFiles  int    `json:"uploadedFiles"`
	TotalViews     int    `json:"totalViews"`
	TotalWatchtime int64  `json:"totalWatchTime"`
}
