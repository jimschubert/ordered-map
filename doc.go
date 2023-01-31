/*
Package orderedmap defines a generic Ordered Map data structure with an API similar to that of [container/list].

The intent of an ordered map is to retain a map's key insertion order, similar to [LinkedHashMap] from Java.
It is *not* to constrain a map to *sorted* order.

# Usage

Go does not allow for extensions to built-in maps. As such there is no "range" or keyed assignments.

Construct a map with contents as follows:

	myMap := orderedmap.New[string, string].
		Set("First", "1st").
		Set("Second", "2nd").
		Set("Third", "3rd")

Iterate the map with the provided iterator:

	it := myMap.Iterator()
	for i := it.Next(); i != nil; i = it.Next() {
		fmt.Printf("Shorthand for %q is %q.\n", i.Key, i.Value)
	}

Or, with lengthier for loops:

	it := myMap.Iterator()
	var i *orderedmap.KeyValuePair[string, string]
	for {
		i = it.Next()
		if i == nil {
			break
		}
		fmt.Printf("Shorthand for %q is %q.\n", i.Key, i.Value)
	}

Get the value defined at some key:

	myValue, ok := myMap.Get("Second")

There are also some utility functions:

	myValue:= myMap.GetOrDefault("Tenth", "10th")
	first := myMap.First()
	last := myMap.Last()

There are functions to manipulate the map as well (including the order of elements):

	removed, ok := myMap.Remove("First")
	err := myMap.MoveToFront("Third")
	err := myMap.MoveToBack("Third")
	err := myMap.MoveAfter("Third", "Second")
	err := myMap.MoveBefore("Third", "Fourth")
	err := myMap.InsertAfter("Third", "3rd", "Second")
	err := myMap.InsertBefore("Third", "3rd", "Fourth")

[container/list]: https://pkg.go.dev/container/list
[LinkedHashMap]: https://docs.oracle.com/javase/8/docs/api/java/util/LinkedHashMap.html
*/
package orderedmap
