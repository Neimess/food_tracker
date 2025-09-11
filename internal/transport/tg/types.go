package tg

type TogglePayload struct {
	Index        int   `json:"index"`
	IngredientID int64 `json:"indredient_id"`
}

type DefPayload struct {
	ID int64 `json:"id"`
}