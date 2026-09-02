package engine

import "reflect"

// A card whose ability names a house nearly always names its own: Battle Fleet,
// a Mars card, reveals Mars cards; Pitlord, a Dis card, locks you into Dis;
// Witch of the Wilds, an Untamed card, lets you play an Untamed card off-house.
// Writing the house out a second time lets the two drift, so the author writes
// the SelfHouse sentinel and the card's own house is filled in here, once, when
// the definition is built. Only a card that names a *different* house — one that
// really is about Sanctum rather than about itself — spells that house out.

var (
	houseType = reflect.TypeOf(House(0))
)

// selfHouseResolvable is implemented by the value types that keep part of a card
// definition in unexported fields (Target, KeyCostChange). Reflection can read
// but never write those, so such a type resolves its own sentinels, delegating
// back to resolvedIn for whatever effect-tree node it holds.
type selfHouseResolvable interface {
	// selfHouseResolved returns a copy with each SelfHouse in it replaced by house.
	// It returns any so every implementer shares one signature; each returns its
	// own concrete type.
	selfHouseResolved(house House) any
}

// resolveSelfHouse returns def with every SelfHouse sentinel it holds — in an
// ability's effect tree, a Target, a constant ability, a HouseLock, a play
// permission — replaced by the card's own house.
func resolveSelfHouse(def CardDefinition) CardDefinition {
	return resolvedIn(def, def.House)
}

// resolvedIn returns node with its SelfHouse sentinels resolved to house. It is
// the typed entry point a selfHouseResolvable implementer uses for a nested node
// it holds; node must not be a nil interface.
func resolvedIn[T any](node T, house House) T {
	return selfHouseResolved(reflect.ValueOf(node), house).Interface().(T)
}

// selfHouseResolved deep-copies v with each SelfHouse in it replaced by house,
// descending through structs, slices, interfaces (the effect, condition, count,
// and chooser nodes), and pointers. Anything else comes back unchanged.
//
// A struct is copied wholesale before its exported fields are rewritten, so its
// unexported fields survive the copy even though reflection cannot set them; a
// struct that hides part of the definition there is selfHouseResolvable and
// resolves itself instead.
func selfHouseResolved(v reflect.Value, house House) reflect.Value {
	if v.Type() == houseType {
		if House(v.Uint()) == SelfHouse {
			return reflect.ValueOf(house)
		}
		return v
	}
	if r, ok := v.Interface().(selfHouseResolvable); ok {
		return reflect.ValueOf(r.selfHouseResolved(house))
	}
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		out := reflect.New(v.Type()).Elem()
		out.Set(selfHouseResolved(v.Elem(), house))
		return out
	case reflect.Pointer:
		if v.IsNil() {
			return v
		}
		out := reflect.New(v.Type().Elem())
		out.Elem().Set(selfHouseResolved(v.Elem(), house))
		return out
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		out.Set(v)
		for i := range out.NumField() {
			if f := out.Field(i); f.CanSet() {
				f.Set(selfHouseResolved(f, house))
			}
		}
		return out
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := range v.Len() {
			out.Index(i).Set(selfHouseResolved(v.Index(i), house))
		}
		return out
	default:
		return v
	}
}
