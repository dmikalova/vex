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
// back to replacedIn for whatever effect-tree node it holds.
type selfHouseResolvable interface {
	// houseReplaced returns a copy with each occurrence of house from replaced by to.
	// It returns any so every implementer shares one signature; each returns its
	// own concrete type.
	houseReplaced(from, to House) any
}

// resolveSelfHouse returns def with every SelfHouse sentinel it holds — in an
// ability's effect tree, a Target, a constant ability, a HouseLock, a play
// permission — replaced by the card's own house.
func resolveSelfHouse(def CardDefinition) CardDefinition {
	return replacedIn(def, SelfHouse, def.House)
}

// Rehouse returns def rehoused to house: its House and every reference to its
// printed house move to house. The self-house convention
// (TestNoCardHardcodesItsOwnHouse) guarantees a card never spells its own printed
// house out, so every reference to it came from card.House.Self — which is why a
// rehoused Maverick's ability names the house it is now printed in, not the one it
// left (The Flex rehoused to Mars chooses a Mars creature, not a Brobnar one).
func Rehouse(def CardDefinition, house House) CardDefinition {
	if def.House == house {
		return def
	}
	return replacedIn(def, def.House, house)
}

// replacedIn returns node with each occurrence of house from replaced by to. It is
// the typed entry point a selfHouseResolvable implementer uses for a nested node
// it holds; node must not be a nil interface.
func replacedIn[T any](node T, from, to House) T {
	return replaceHouse(reflect.ValueOf(node), from, to).Interface().(T)
}

// replaceHouse deep-copies v with each House equal to from replaced by to,
// descending through structs, slices, interfaces (the effect, condition, count,
// and chooser nodes), and pointers. Anything else comes back unchanged.
//
// A struct is copied wholesale before its exported fields are rewritten, so its
// unexported fields survive the copy even though reflection cannot set them; a
// struct that hides part of the definition there is selfHouseResolvable and
// replaces itself instead.
func replaceHouse(v reflect.Value, from, to House) reflect.Value {
	if v.Type() == houseType {
		if House(v.Uint()) == from {
			return reflect.ValueOf(to)
		}
		return v
	}
	if r, ok := v.Interface().(selfHouseResolvable); ok {
		return reflect.ValueOf(r.houseReplaced(from, to))
	}
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		out := reflect.New(v.Type()).Elem()
		out.Set(replaceHouse(v.Elem(), from, to))
		return out
	case reflect.Pointer:
		if v.IsNil() {
			return v
		}
		out := reflect.New(v.Type().Elem())
		out.Elem().Set(replaceHouse(v.Elem(), from, to))
		return out
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		out.Set(v)
		for i := range out.NumField() {
			if f := out.Field(i); f.CanSet() {
				f.Set(replaceHouse(f, from, to))
			}
		}
		return out
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := range v.Len() {
			out.Index(i).Set(replaceHouse(v.Index(i), from, to))
		}
		return out
	default:
		return v
	}
}
