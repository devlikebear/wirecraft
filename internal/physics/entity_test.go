package physics

import "testing"

func TestSortEntitiesOrdersByID(t *testing.T) {
	entities := []DynamicEntity{
		{ID: EntityID("piston:2"), Type: EntityTypePiston},
		{ID: EntityID("motor:1"), Type: EntityTypeMotor},
		{ID: EntityID("piston:1"), Type: EntityTypePiston},
	}

	sorted := SortEntities(entities)

	wantIDs := []EntityID{
		EntityID("motor:1"),
		EntityID("piston:1"),
		EntityID("piston:2"),
	}
	if len(sorted) != len(wantIDs) {
		t.Fatalf("len(sorted) = %d, want %d", len(sorted), len(wantIDs))
	}
	for i, wantID := range wantIDs {
		if sorted[i].ID != wantID {
			t.Fatalf("sorted[%d].ID = %q, want %q", i, sorted[i].ID, wantID)
		}
	}
	if entities[0].ID != EntityID("piston:2") {
		t.Fatalf("SortEntities mutated input: first ID = %q", entities[0].ID)
	}
}

func TestEntitySetEntitiesReturnsSortedCopies(t *testing.T) {
	set := NewEntitySet([]DynamicEntity{
		{
			ID:   EntityID("piston:2"),
			Type: EntityTypePiston,
			Transform: Transform{
				Position: Vec3{X: 2, Y: 0, Z: 0},
				Rotation: IdentityQuat(),
				Scale:    UnitVec3(),
			},
		},
		{
			ID:   EntityID("piston:1"),
			Type: EntityTypePiston,
			Transform: Transform{
				Position: Vec3{X: 1, Y: 0, Z: 0},
				Rotation: IdentityQuat(),
				Scale:    UnitVec3(),
			},
		},
	})

	first := set.Entities()
	if first[0].ID != EntityID("piston:1") {
		t.Fatalf("first[0].ID = %q, want piston:1", first[0].ID)
	}
	first[0].ID = EntityID("mutated")

	second := set.Entities()
	if second[0].ID != EntityID("piston:1") {
		t.Fatalf("EntitySet exposed internal storage: second[0].ID = %q", second[0].ID)
	}
}

func TestDefaultTransformUsesIdentityRotationAndUnitScale(t *testing.T) {
	transform := DefaultTransform(Vec3{X: 3, Y: 4, Z: 5})

	if transform.Position != (Vec3{X: 3, Y: 4, Z: 5}) {
		t.Fatalf("Position = %+v, want requested position", transform.Position)
	}
	if transform.Rotation != IdentityQuat() {
		t.Fatalf("Rotation = %+v, want identity", transform.Rotation)
	}
	if transform.Scale != UnitVec3() {
		t.Fatalf("Scale = %+v, want unit scale", transform.Scale)
	}
}
