package attrfmt

import (
	"golang.org/x/exp/slog"
)

// HideSensitiveValueOnKeys will hide all value if key of the attribute
// is match with the keys specified in the argument.
func HideSensitiveValueOnKeys(keys map[string]struct{}) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		v := attr.Value
		for v.Kind() == slog.KindLogValuer {
			v = v.LogValuer().LogValue()
		}

		if v.Kind() == slog.KindGroup {
			groupValue := make([]slog.Attr, 0)
			for _, a := range v.Group() {
				// add attribute key as the group list,
				// because these attributes is defined under slog.Group
				groupValue = append(groupValue, HideSensitiveValueOnKeys(keys)(append(groups, attr.Key), a))
			}

			attr.Value = slog.GroupValue(groupValue...)
			return attr
		}

		if _, exist := keys[attr.Key]; exist {
			attr.Value = slog.StringValue("[FILTERED]")
			return attr
		}

		for _, group := range groups {
			if _, exist := keys[group]; exist {
				attr.Value = slog.StringValue("[FILTERED]")
				return attr
			}
		}

		return attr
	}
}
