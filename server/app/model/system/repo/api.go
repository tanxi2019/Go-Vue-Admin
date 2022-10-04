package repo

import "server/app/model/system"

type ApiTreeResp struct {
	ID       int           `json:"ID"`
	Desc     string        `json:"desc"`
	Category string        `json:"category"`
	Children []*system.Api `json:"children"`
}
