package resolver

type (
	Dependency struct {
		Resource string
		Locator  Locator
		Status   *string
	}

	Locator struct {
		Namespace *string
		Name      string
	}
)

// To returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

func (d *Dependency) String() string {
	result := d.Resource + ":"
	if d.Locator.Namespace != nil {
		result += *d.Locator.Namespace + "/"
	}
	result += d.Locator.Name
	if d.Status != nil {
		result += "?" + *d.Status
	}

	return result
}
