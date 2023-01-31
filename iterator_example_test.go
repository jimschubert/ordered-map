package orderedmap_test

import (
	"fmt"

	orderedmap "github.com/jimschubert/ordered-map"
)

func ExampleIterator() {
	var animalSounds = orderedmap.New[string, string]().
		Set("Cat", "Meow").
		Set("Dog", "Woof").
		Set("Cow", "Moo").
		Set("Snake", "Ssss").
		Set("Fox", "Ring-ding-ding-ding-dingeringeding!")

	it := animalSounds.Iterator()
	for i := it.Next(); i != nil; i = it.Next() {
		fmt.Printf("The %s says %q.\n", i.Key, i.Value)
	}

	// Output:
	// The Cat says "Meow".
	// The Dog says "Woof".
	// The Cow says "Moo".
	// The Snake says "Ssss".
	// The Fox says "Ring-ding-ding-ding-dingeringeding!".
}

func ExampleIterator_alternate() {
	var animalSounds = orderedmap.New[string, string]().
		Set("Cat", "Meow").
		Set("Dog", "Woof").
		Set("Cow", "Moo").
		Set("Snake", "Ssss").
		Set("Fox", "Ring-ding-ding-ding-dingeringeding!")

	it := animalSounds.Iterator()
	var i *orderedmap.KeyValuePair[string, string]
	for {
		i = it.Next()
		if i == nil {
			break
		}
		fmt.Printf("The %s says %q.\n", i.Key, i.Value)
	}

	// Output:
	// The Cat says "Meow".
	// The Dog says "Woof".
	// The Cow says "Moo".
	// The Snake says "Ssss".
	// The Fox says "Ring-ding-ding-ding-dingeringeding!".
}
