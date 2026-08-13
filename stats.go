package math3d

// Stats contains one-pass summary statistics for a sequence of values.
// Count is zero for an empty sequence, in which case Min, Max, and Sum are
// their zero values.
type Stats[T any] struct {
	Count         int
	Min, Max, Sum T
}

// NewStats constructs summary statistics from their stored representation.
func NewStats[T any](count int, minimum, maximum, sum T) Stats[T] {
	return Stats[T]{Count: count, Min: minimum, Max: maximum, Sum: sum}
}

type statsVector interface {
	Vec2 | Vec3 | Vec4 | DVec2 | DVec3 | DVec4
}

// StatsOf summarizes vectors component by component. Unlike upstream, the
// first value initializes Min and Max, so an all-positive or all-negative
// sequence is not biased toward zero.
func StatsOf[T statsVector](values []T) Stats[T] {
	if len(values) == 0 {
		return Stats[T]{}
	}

	result := Stats[T]{Count: 1, Min: values[0], Max: values[0], Sum: values[0]}
	for _, value := range values[1:] {
		result.Count++
		result.Min = statsMin(result.Min, value)
		result.Max = statsMax(result.Max, value)
		result.Sum = statsAdd(result.Sum, value)
	}
	return result
}

// StatsAverage returns the component-wise average. It fails when Count is not
// positive.
func StatsAverage[T statsVector](stats Stats[T]) (T, bool) {
	if stats.Count <= 0 {
		var zero T
		return zero, false
	}
	return statsScale(stats.Sum, 1/float64(stats.Count)), true
}

// StatsExtents returns Max-Min component by component.
func StatsExtents[T statsVector](stats Stats[T]) T { return statsSub(stats.Max, stats.Min) }

// StatsMiddle returns the midpoint between Min and Max.
func StatsMiddle[T statsVector](stats Stats[T]) T {
	return statsAdd(stats.Min, statsScale(StatsExtents(stats), 0.5))
}

func statsAdd[T statsVector](a, b T) T {
	switch a := any(a).(type) {
	case Vec2:
		return any(a.Add(any(b).(Vec2))).(T)
	case Vec3:
		return any(a.Add(any(b).(Vec3))).(T)
	case Vec4:
		return any(a.Add(any(b).(Vec4))).(T)
	case DVec2:
		return any(a.Add(any(b).(DVec2))).(T)
	case DVec3:
		return any(a.Add(any(b).(DVec3))).(T)
	case DVec4:
		return any(a.Add(any(b).(DVec4))).(T)
	}
	panic("unreachable")
}

func statsSub[T statsVector](a, b T) T {
	switch a := any(a).(type) {
	case Vec2:
		return any(a.Sub(any(b).(Vec2))).(T)
	case Vec3:
		return any(a.Sub(any(b).(Vec3))).(T)
	case Vec4:
		return any(a.Sub(any(b).(Vec4))).(T)
	case DVec2:
		return any(a.Sub(any(b).(DVec2))).(T)
	case DVec3:
		return any(a.Sub(any(b).(DVec3))).(T)
	case DVec4:
		return any(a.Sub(any(b).(DVec4))).(T)
	}
	panic("unreachable")
}

func statsMin[T statsVector](a, b T) T {
	switch a := any(a).(type) {
	case Vec2:
		return any(a.Min(any(b).(Vec2))).(T)
	case Vec3:
		return any(a.Min(any(b).(Vec3))).(T)
	case Vec4:
		return any(a.Min(any(b).(Vec4))).(T)
	case DVec2:
		return any(a.Min(any(b).(DVec2))).(T)
	case DVec3:
		return any(a.Min(any(b).(DVec3))).(T)
	case DVec4:
		return any(a.Min(any(b).(DVec4))).(T)
	}
	panic("unreachable")
}

func statsMax[T statsVector](a, b T) T {
	switch a := any(a).(type) {
	case Vec2:
		return any(a.Max(any(b).(Vec2))).(T)
	case Vec3:
		return any(a.Max(any(b).(Vec3))).(T)
	case Vec4:
		return any(a.Max(any(b).(Vec4))).(T)
	case DVec2:
		return any(a.Max(any(b).(DVec2))).(T)
	case DVec3:
		return any(a.Max(any(b).(DVec3))).(T)
	case DVec4:
		return any(a.Max(any(b).(DVec4))).(T)
	}
	panic("unreachable")
}

func statsScale[T statsVector](value T, scale float64) T {
	switch value := any(value).(type) {
	case Vec2:
		return any(value.Scale(float32(scale))).(T)
	case Vec3:
		return any(value.Scale(float32(scale))).(T)
	case Vec4:
		return any(value.Scale(float32(scale))).(T)
	case DVec2:
		return any(value.Scale(scale)).(T)
	case DVec3:
		return any(value.Scale(scale)).(T)
	case DVec4:
		return any(value.Scale(scale)).(T)
	}
	panic("unreachable")
}
