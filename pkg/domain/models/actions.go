package models

import (
	"time"
)

type PollingCommand struct {
	Polling   *RequestGoogleAction `json:"actions"`
	Type      *string              `json:"type,omitempty"`
	Timestamp *time.Time           `json:"timestamp,omitempty"`
}

type RequestGoogleAction struct {
	ActionID       string `json:"actionid" binding:"required,uuid"`
	RequestID      string `json:"requestid" binding:"required"`
	Pollmode       string `json:"pollmode" binding:"required,oneof=none 1m 5m"`
	Selectdocument string `json:"selectdocument" binding:"required,oneof=byuri byid"`
	Document       string `json:"document" binding:"required,url"`
	NameDocument   string `json:"namedocument" binding:"omitempty,max=255"`
	ResourceID     string `json:"resourceid" binding:"omitempty"`
	Operation      string `json:"operation" binding:"required,oneof=getallcontent insertrow"`
	Data           string `json:"data" binding:"omitempty"`
	CredentialID   string `json:"credentialid" binding:"required"`
	Sub            string `json:"sub" binding:"required,numeric"`
	Type           string `json:"type" binding:"required,oneof=googlesheets"`
	WorkflowID     string `json:"workflowid" binding:"required,uuid"`
	NodeID         string `json:"nodeid" binding:"required,max=255"`
	RedirectURL    string `json:"redirecturl" binding:"required"`
	Status         string `json:"status" binding:"omitempty,oneof=pending success failed"`
	CreatedAt      string `json:"createdat" binding:"required,datetime=2006-01-02T15:04:05Z"`
	Testmode       bool   `json:"testmode"`
}

type ResponseGetGoogleSheetByID struct {
	Error  string `json:"error"`
	Data   string `json:"data"`
	Status int    `json:"status"`
}
