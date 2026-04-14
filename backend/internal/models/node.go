package models

import "time"

// NodeType represents the type of node in a travel plan
type NodeType string

// Node type constants
const (
	NodeTypeAttraction NodeType = "attraction"
	NodeTypeTransition NodeType = "transition"
)

// String returns the string representation of a NodeType
func (n NodeType) String() string {
	return string(n)
}

// CheckNodeType validates if the provided type is a valid NodeType
func CheckNodeType(nodeType NodeType) bool {
	validTypes := map[NodeType]bool{
		NodeTypeAttraction: true,
		NodeTypeTransition: true,
	}
	return validTypes[nodeType]
}

// NodeDetail is an interface for node detail structures (Attraction or Transition)
// Allows generic operations on both detail types
type NodeDetail interface {
	// GetNodeType returns the type of node this detail belongs to
	GetNodeType() NodeType
	// SetNodeID sets the node ID for this detail
	SetNodeID(id string)
}

// Node represents a base point in a travel plan (attraction or transition)
// Uses single-table inheritance with type discriminator
// Embeds type-specific detail structures for convenient access
type Node struct {
	// id (UUID, primary key)
	// Unique identifier for the node
	ID string `gorm:"primaryKey;type:TEXT" json:"id"`

	// type (enum: attraction | transition, discriminator for polymorphism)
	// Determines which detail table contains additional data
	Type string `gorm:"type:TEXT;not null;index" json:"type"`

	// created_by (UUID, foreign key to User, creator for moderation)
	// User who created this node
	CreatedBy string `gorm:"type:TEXT;not null;index" json:"created_by"`

	// is_approved (boolean, default: true for system nodes, false for user-created)
	// Whether node is approved by admin (user-created nodes start as false)
	IsApproved bool `gorm:"type:BOOLEAN;default:true;index" json:"is_approved"`

	// created_at (timestamp, UTC, immutable)
	// When the node was created
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// updated_at (timestamp, UTC, nullable, updated only by admin approval)
	// When the node status was last updated by admin
	UpdatedAt *time.Time `gorm:"type:TIMESTAMP" json:"updated_at"`

	// AttractionNodeDetail (optional pointer, embedded via foreign key)
	// Contains attraction-specific details if type == "attraction"
	// Use GORM Preload("AttractionNodeDetail") to load
	AttractionNodeDetail *AttractionNodeDetail `gorm:"foreignKey:NodeID;references:ID" json:"attraction,omitempty"`

	// TransitionNodeDetail (optional pointer, embedded via foreign key)
	// Contains transition-specific details if type == "transition"
	// Use GORM Preload("TransitionNodeDetail") to load
	TransitionNodeDetail *TransitionNodeDetail `gorm:"foreignKey:NodeID;references:ID" json:"transition,omitempty"`
}

// TableName specifies the database table name for Node
func (Node) TableName() string {
	return "nodes"
}

// GetDetail returns the type-specific detail (AttractionNodeDetail or TransitionNodeDetail)
// based on the node's type discriminator. Returns nil if detail is not loaded or type is invalid.
func (n *Node) GetDetail() interface{} {
	switch NodeType(n.Type) {
	case NodeTypeAttraction:
		return n.AttractionNodeDetail
	case NodeTypeTransition:
		return n.TransitionNodeDetail
	default:
		return nil
	}
}

// GetAttractionDetail returns the attraction detail if this is an attraction node, nil otherwise
func (n *Node) GetAttractionDetail() *AttractionNodeDetail {
	if n.Type == string(NodeTypeAttraction) {
		return n.AttractionNodeDetail
	}
	return nil
}

// GetTransitionDetail returns the transition detail if this is a transition node, nil otherwise
func (n *Node) GetTransitionDetail() *TransitionNodeDetail {
	if n.Type == string(NodeTypeTransition) {
		return n.TransitionNodeDetail
	}
	return nil
}

// AttractionNodeDetail contains details specific to attraction nodes
type AttractionNodeDetail struct {
	// node_id (UUID, primary key, foreign key)
	// References the parent node
	NodeID string `gorm:"primaryKey;type:TEXT" json:"node_id"`

	// name (string, max 200 chars, required)
	// Name of the attraction
	Name string `gorm:"type:TEXT;not null;size:200;index" json:"name"`

	// category (enum: tourist_attraction, restaurant, hotel, museum, park, shopping, entertainment, other)
	// Type of attraction for filtering
	Category string `gorm:"type:TEXT;not null;size:50" json:"category"`

	// location (string, max 300 chars, required)
	// Physical location (address or coordinates)
	Location string `gorm:"type:TEXT;not null;size:300" json:"location"`

	// description (string, max 1000 chars, optional)
	// Detailed description of the attraction
	Description string `gorm:"type:TEXT;size:1000" json:"description"`

	// contact_info (string, max 200 chars, optional)
	// Phone, email, or website
	ContactInfo string `gorm:"type:TEXT;size:200" json:"contact_info"`

	// hours_of_operation (string, max 200 chars, optional)
	// Operating hours (e.g., "09:00-18:00" or "Mon-Sun: 09:00-18:00")
	HoursOfOperation string `gorm:"type:TEXT;size:200" json:"hours_of_operation"`

	// created_at (timestamp, UTC, immutable)
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}

// TableName specifies the database table name for AttractionNodeDetail
func (AttractionNodeDetail) TableName() string {
	return "attraction_node_details"
}

// GetNodeType returns the type of node this detail belongs to (attraction)
func (a *AttractionNodeDetail) GetNodeType() NodeType {
	return NodeTypeAttraction
}

// SetNodeID sets the node ID for this detail
func (a *AttractionNodeDetail) SetNodeID(id string) {
	a.NodeID = id
}

// AttractionCategory represents valid attraction categories
type AttractionCategory string

const (
	CategoryTouristAttraction AttractionCategory = "tourist_attraction"
	CategoryRestaurant        AttractionCategory = "restaurant"
	CategoryHotel             AttractionCategory = "hotel"
	CategoryMuseum            AttractionCategory = "museum"
	CategoryPark              AttractionCategory = "park"
	CategoryShopping          AttractionCategory = "shopping"
	CategoryEntertainment     AttractionCategory = "entertainment"
	CategoryOther             AttractionCategory = "other"
)

// String returns the string representation
func (c AttractionCategory) String() string {
	return string(c)
}

// ValidAttractionCategories returns all valid categories
func ValidAttractionCategories() []string {
	return []string{
		CategoryTouristAttraction.String(),
		CategoryRestaurant.String(),
		CategoryHotel.String(),
		CategoryMuseum.String(),
		CategoryPark.String(),
		CategoryShopping.String(),
		CategoryEntertainment.String(),
		CategoryOther.String(),
	}
}

// TransitionNodeDetail contains details specific to transition nodes
type TransitionNodeDetail struct {
	// node_id (UUID, primary key, foreign key)
	// References the parent node
	NodeID string `gorm:"primaryKey;type:TEXT" json:"node_id"`

	// title (string, max 200 chars, required)
	// Service/line identifier (e.g., "Bus Line 5", "M1 Train", "Walking")
	// Immutable identifier for the transition service
	Title string `gorm:"type:TEXT;not null;size:200" json:"title"`

	// mode (enum: walking, car, bus, train, bike, taxi, flight, other)
	// How the traveler moves between attractions
	Mode string `gorm:"type:TEXT;not null;size:50" json:"mode"`

	// description (string, max 1000 chars, optional)
	// General description of the transition journey, independent of plan
	Description string `gorm:"type:TEXT;size:1000" json:"description"`

	// hours_of_operation (string, max 200 chars, optional)
	// Operating hours/availability of this service
	// e.g., "Mon-Fri 6:00-23:00, Sat-Sun 7:00-22:00"
	// NULL for modes without fixed hours (walking, driving)
	HoursOfOperation *string `gorm:"type:TEXT;size:200" json:"hours_of_operation"`

	// route_notes (string, max 500 chars, optional)
	// Additional directions or notes for the journey
	RouteNotes string `gorm:"type:TEXT;size:500" json:"route_notes"`

	// created_at (timestamp, UTC, immutable)
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}

// TableName specifies the database table name for TransitionNodeDetail
func (TransitionNodeDetail) TableName() string {
	return "transition_node_details"
}

// GetNodeType returns the type of node this detail belongs to (transition)
func (t *TransitionNodeDetail) GetNodeType() NodeType {
	return NodeTypeTransition
}

// SetNodeID sets the node ID for this detail
func (t *TransitionNodeDetail) SetNodeID(id string) {
	t.NodeID = id
}

// TransitionMode represents valid transition modes
type TransitionMode string

const (
	ModeWalking TransitionMode = "walking"
	ModeCar     TransitionMode = "car"
	ModeBus     TransitionMode = "bus"
	ModeTrain   TransitionMode = "train"
	ModeBike    TransitionMode = "bike"
	ModeTaxi    TransitionMode = "taxi"
	ModeFlight  TransitionMode = "flight"
	ModeOther   TransitionMode = "other"
)

// String returns the string representation
func (m TransitionMode) String() string {
	return string(m)
}

// ValidTransitionModes returns all valid modes
func ValidTransitionModes() []string {
	return []string{
		ModeWalking.String(),
		ModeCar.String(),
		ModeBus.String(),
		ModeTrain.String(),
		ModeBike.String(),
		ModeTaxi.String(),
		ModeFlight.String(),
		ModeOther.String(),
	}
}

// Validation methods for models

// ValidateAttractionNodeDetail validates attraction node details
func (a *AttractionNodeDetail) Validate() bool {
	if len(a.Name) < 1 || len(a.Name) > 200 {
		return false
	}
	if len(a.Location) < 1 || len(a.Location) > 300 {
		return false
	}
	if len(a.Description) > 1000 || len(a.ContactInfo) > 200 || len(a.HoursOfOperation) > 200 {
		return false
	}
	return true
}

// ValidateTransitionNodeDetail validates transition node details
func (t *TransitionNodeDetail) Validate() bool {

	if len(t.RouteNotes) > 500 {
		return false
	}
	return true
}
