package dtos

// ArtefactSummaryDTO represents a summarized view of an artefact.
type ArtefactSummaryDTO struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	Material               string  `json:"material"`
	CollectionName         *string `json:"collectionName,omitempty"`
	ArchaeologistName      *string `json:"archaeologistName,omitempty"`
	ArchaeologicalSiteName *string `json:"archaeologicalSiteName,omitempty"`
	ShelfCode              *int    `json:"shelfCode,omitempty"`
	Level                  *int    `json:"level,omitempty"`
	Column                 *string `json:"column,omitempty"`
}
