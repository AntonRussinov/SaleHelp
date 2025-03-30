package domain

type Image struct {
	ID        int    // внутренний идентификатор изображения
	OfferID   int    // идентификатор предложения, к которому относится изображение
	ImagePath string // путь или URL к изображению
}
