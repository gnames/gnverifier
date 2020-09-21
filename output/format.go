package output

type Format int

const (
	InvalidFormat Format = iota
	CSV
	CompactJSON
	PrettyJSON
)

func NewFormat(s string) Format {
	switch s {
	case "csv":
		return CSV
	case "compact":
		return CompactJSON
	case "pretty":
		return PrettyJSON
	default:
		return InvalidFormat
	}
}

