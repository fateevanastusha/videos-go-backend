package validators

import "github.com/fateevanastusha/videos-go-backend/internal/models"

func ValidateResolutions(resolutions []models.Resolution) []ErrorMessage {
	for _, res := range resolutions {
		if !res.Valid() {
			return []ErrorMessage{{Message: "must be one of P144, P240, P360, P480, P720, P1080, P1440, P2160", Field: "availableResolutions"}}
		}
	}
	return nil
}
