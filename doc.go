// Package math3d provides dependency-free geometric value types and operations.
//
// # Conventions
//
// The package uses right-handed Cartesian cross products. View and world
// helpers treat local forward as the negative Z direction; the package does
// not prescribe a world-up axis.
//
// Matrices use row-vector semantics. A point is transformed as
// [x y z 1] * M, while a direction or normal is transformed as
// [x y z 0] * M. Translation therefore occupies M41, M42, and M43. In a
// product A * B, A is applied first and B second.
//
// Rotation APIs accept radians unless their documentation explicitly says
// otherwise. Positive rotations follow the right-hand rule. Coordinate record
// fields whose upstream declarations do not specify units document their units
// individually rather than inheriting this rule.
//
// Primary floating-point values use float32, matching the upstream C# float
// representation. Types prefixed with D use float64. AxisAngle is also
// float64-based to match its upstream representation. Floating-point values
// retain normal Go IEEE 754 equality behavior, including its behavior for NaN.
// Approximate comparisons use an explicit absolute tolerance; conveniences
// may use the source default of 1e-7. No mutable global tolerance is used.
//
// Values have public fields and may be constructed and changed directly by
// callers. Package operations use value semantics: they do not mutate value
// arguments or receiver values, and return changed values instead. A zero value
// is safe to store and pass, but it is not necessarily geometrically meaningful;
// operations document any stronger preconditions and report ordinary failures
// with an ok result.
//
// Struct layout is an API representation, not a binary wire format. Any binary
// encoding supplied by this package defines component order and endianness
// explicitly instead of serializing raw struct memory.
package math3d
