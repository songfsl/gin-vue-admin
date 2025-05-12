package request

type GetLoginTimeByUsersIdReq struct {
	ID        uint   `json:"ID" form:"ID"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
