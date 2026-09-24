package models

type Resolution string

func (r Resolution) String() string {
	return string(r)
}

const (
	P144  Resolution = "P144"
	P240  Resolution = "P240"
	P360  Resolution = "P360"
	P480  Resolution = "P480"
	P720  Resolution = "P720"
	P1080 Resolution = "P1080"
	P1440 Resolution = "P1440"
	P2160 Resolution = "P2160"
)

func (r Resolution) Valid() bool {
	switch r {
	case P144, P240, P360, P480, P720, P1080, P1440, P2160:
		return true
	}
	return false
}

type Video struct {
	Id                   int          `json:"id"`
	Title                string       `json:"title"`
	Author               string       `json:"author"`
	CanBeDownloaded      bool         `json:"canBeDownloaded"`
	MinAgeRestriction    *int         `json:"minAgeRestriction"`
	CreatedAt            string       `json:"createdAt"`
	PublicationDate      string       `json:"publicationDate"`
	AvailableResolutions []Resolution `json:"availableResolutions"`
}
