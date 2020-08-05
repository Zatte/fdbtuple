## Foundation DB tuple layer in golang without C-GO dependencies

### FDB Tuple layer 
Is a cross language standarized binary packing for arbitrary data where packed bytes have predictable byte sort order. Meant to build keys in a key/value store. 

### Why
The packing is usefull in most byte ordered key value store to create key spaces. This fork only removes the dependency of c-go which foundation db bidings ships with. 

### Details
See official foundation db documentation

https://apple.github.io/foundationdb/developer-guide.html#namespace-management
https://godoc.org/github.com/apple/foundationdb/bindings/go/src/fdb/tuple
