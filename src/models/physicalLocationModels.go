package models

type LevelNumber int
const (
	Level1 LevelNumber = 1
	Level2 LevelNumber = 2
	Level3 LevelNumber = 3
	Level4 LevelNumber = 4
)

type ColumnLetter string
const (
	ColumnA ColumnLetter = "A"
	ColumnB ColumnLetter = "B"
	ColumnC ColumnLetter = "C"
	ColumnD ColumnLetter = "D"
)

type PhysicalLocationModel struct {
	ID      int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Level   LevelNumber  `json:"level" gorm:"column:level;type:int;not null"`
	Column  ColumnLetter `json:"column" gorm:"column:column;type:varchar(1);not null"`
	ShelfId int          `json:"shelfId" gorm:"column:shelf_id;type:int;not null"`
	Shelf   ShelfModel   `json:"shelf" gorm:"foreignKey:ShelfId;references:ID"`
}
