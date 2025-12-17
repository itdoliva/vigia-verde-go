package treeevent

// GeoLocation representa o campo "location" com latitude/longitude.
// Mais pra frente, isso vai ser preenchido a partir do JSON do POST.
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CreateTreeEventInput representa os dados necessários
// para criar um treeEvent no Firestore.
type CreateTreeEventInput struct {
	Location  GeoPoint `json:"location"`
	EventType string   `json:"eventType"`
	Title     string   `json:"title"`
	AuthorID  string   `json:"authorId"`
}
