package dtos

// ArtefactSummaryDTO represents a summarized view of an artefact, including optional related entity names.
type ArtefactSummaryDTO struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	Material               string  `json:"material"`
	ArchaeologicalSiteName *string `json:"archaeologicalSiteName,omitempty"`
	ArchaeologistName      *string `json:"archaeologistName,omitempty"`
	CollectionName         *string `json:"collectionName,omitempty"`
	ShelfCode              *int    `json:"shelfCode,omitempty"`
	Column                 *string `json:"column,omitempty"`
	Level                  *int    `json:"level,omitempty"`
}
