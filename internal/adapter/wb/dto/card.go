package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type GetCardsRequest struct {
	Settings Settings `json:"settings"`
}

type Settings struct {
	Cursor Cursor `json:"cursor"`
	Filter Filter `json:"filter"`
	Sort   Sort   `json:"sort"`
}

type Cursor struct {
	Limit     int    `json:"limit"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	NmID      int    `json:"nmID,omitempty"`
}

type Sort struct {
	Ascending bool `json:"ascending,omitempty"`
}

type Filter struct {
	TextSearch            string   `json:"textSearch"`                           // строка поиска
	AllowedCategoriesOnly bool     `json:"allowedCategoriesOnly" default:"true"` // только разрешённые категории
	TagIDs                []int    `json:"tagIDs"`                               // фильтрация по тегам
	ObjectIDs             []int    `json:"objectIDs"`                            // фильтр по объектам
	Brands                []string `json:"brands"`                               // список брендов
	ImtID                 int      `json:"imtID"`                                // конкретный имтID
	WithPhoto             int      `json:"withPhoto"`                            // -1: любые, 1: только с фото
}

// --- Ответ ---
type GetCardsResponse struct {
	Cards  []Card         `json:"cards"`
	Cursor ResponseCursor `json:"cursor"`
}

// Курсор для пагинации
type ResponseCursor struct {
	UpdatedAt string `json:"updatedAt"`
	NmID      int    `json:"nmID"`
	Total     int    `json:"total"`
}

// Одна карточка товара
type Card struct {
	NmID            int              `json:"nmID"`
	ImtID           int              `json:"imtID"`
	NmUUID          string           `json:"nmUUID"`
	SubjectID       int              `json:"subjectID"`
	SubjectName     string           `json:"subjectName"`
	VendorCode      string           `json:"vendorCode"`
	Brand           string           `json:"brand"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	NeedKiz         bool             `json:"needKiz"`
	Photos          []Photo          `json:"photos"`
	Video           string           `json:"video"`
	Wholesale       Wholesale        `json:"wholesale"`
	Dimensions      Dimensions       `json:"dimensions"`
	Characteristics []Characteristic `json:"characteristics"`
	Sizes           []Size           `json:"sizes"`
	Tags            []Tag            `json:"tags"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       string           `json:"updatedAt"`
}

// Фотографии
type Photo struct {
	Big      string `json:"big"`
	C246x328 string `json:"c246x328"`
	C516x688 string `json:"c516x688"`
	Square   string `json:"square"`
	TM       string `json:"tm"`
}

// Опт
type Wholesale struct {
	Enabled bool `json:"enabled"`
	Quantum int  `json:"quantum"`
}

// Габариты
type Dimensions struct {
	Length       int     `json:"length"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	WeightBrutto float64 `json:"weightBrutto"`
	IsValid      bool    `json:"isValid"`
}

// CharacteristicValue — value в ответе WB может быть числом или массивом строк.
type CharacteristicValue []string

func (v *CharacteristicValue) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*v = nil
		return nil
	}
	switch data[0] {
	case '[':
		var s []string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*v = s
		return nil
	case '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*v = []string{s}
		return nil
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var n json.Number
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*v = []string{n.String()}
		return nil
	default:
		return fmt.Errorf("characteristic value: unexpected type (first byte %q)", data[0])
	}
}

// Характеристики (атрибуты)
type Characteristic struct {
	ID    int                `json:"id"`
	Name  string             `json:"name"`
	Value CharacteristicValue `json:"value"`
}

// Размеры / вариации
type Size struct {
	ChrtID   int      `json:"chrtID"`
	TechSize string   `json:"techSize"`
	Skus     []string `json:"skus"`
}

// Теги
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
