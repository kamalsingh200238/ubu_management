package utils

type AlertVariant string

const (
	AlertVariantPrimary AlertVariant = "primary"
	AlertVariantSuccess AlertVariant = "success"
	AlertVariantNeutral AlertVariant = "neutral"
	AlertVariantWarning AlertVariant = "warning"
	AlertVariantDanger  AlertVariant = "danger"
)

type AlertDetails struct {
	Closable bool         `json:"closable"`
	Heading  string       `json:"heading"`
	Messages []string     `json:"messages"`
	Variant  AlertVariant `json:"variant"`
	Duration int          `json:"duration"`
}
