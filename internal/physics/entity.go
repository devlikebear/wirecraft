package physics

import "sort"

type EntityID string

type EntityType string

const (
	EntityTypePiston EntityType = "piston"
	EntityTypeMotor  EntityType = "motor"
)

type DynamicEntity struct {
	ID        EntityID
	Type      EntityType
	Transform Transform
	Velocity  Vec3
	Target    Transform
}

type EntitySet struct {
	entities []DynamicEntity
}

func NewEntitySet(entities []DynamicEntity) EntitySet {
	return EntitySet{entities: SortEntities(entities)}
}

func (s EntitySet) Entities() []DynamicEntity {
	return append([]DynamicEntity(nil), s.entities...)
}

func SortEntities(entities []DynamicEntity) []DynamicEntity {
	sorted := append([]DynamicEntity(nil), entities...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}
