package resolve

// Merge combines multiple environment maps into one.
// Maps are applied left-to-right; later maps override earlier ones.
func Merge(maps ...map[string]string) map[string]string {
	out := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// Filter returns a new map containing only the key/value pairs for which
// keep returns true.
func Filter(env map[string]string, keep func(key, value string) bool) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if keep(k, v) {
			out[k] = v
		}
	}
	return out
}

// Rename applies a renaming function to every key in env.
// If rename returns an empty string the key is dropped.
func Rename(env map[string]string, rename func(key string) string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if nk := rename(k); nk != "" {
			out[nk] = v
		}
	}
	return out
}
