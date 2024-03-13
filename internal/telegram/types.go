package telegram

type (
	String string
)

func (s *String) Parse(value string) error {
	*s = String(value)
	return nil
}
