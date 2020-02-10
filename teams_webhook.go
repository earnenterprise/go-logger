package logger

// TeamsWebhook ...
type TeamsWebhook struct {
	Context          string                `json:"@context"`
	Type             string                `json:"@type"`
	ThemeColor       string                `json:"themeColor"`
	Title            string                `json:"title,omitempty"`
	Text             string                `json:"text,omitempty"`
	Summary          string                `json:"summary,omitempty"`
	Sections         []TeamsWebhookSection `json:"sections,omitempty"`
	PotentialActions []TeamsWebhookAction  `json:"potentialAction,omitempty"`
}

// TeamsWebhookSection ...
type TeamsWebhookSection struct {
	Title            string               `json:"title,omitempty"`
	StartGroup       bool                 `json:"startGroup"`
	ActivityImage    *TeamsWebhookImage   `json:"activityImage,omitempty"`
	ActivityTitle    string               `json:"activityTitle,omitempty"`
	ActivitySubtitle string               `json:"activitySubtitle,omitempty"`
	ActivityText     string               `json:"activityText,omitempty"`
	HeroImage        string               `json:"heroImage,omitempty"`
	Text             string               `json:"text,omitempty"`
	Facts            []TeamsWebhookFact   `json:"facts,omitempty"`
	Images           []TeamsWebhookImage  `json:"images,omitempty"`
	PotentialActions []TeamsWebhookAction `json:"potentialAction,omitempty"`
}

// TeamsWebhookImage ...
type TeamsWebhookImage struct {
	Image string `json:"image,omitempty"`
	Title string `json:"title,omitempty"`
}

// TeamsWebhookFact ...
type TeamsWebhookFact struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// TeamsWebhookAction ...
type TeamsWebhookAction struct {
}

// CreateMessageCard ...
func CreateMessageCard(summary string) *TeamsWebhook {
	webhook := TeamsWebhook{Type: "MessageCard", Context: "https://schema.org/extensions", Sections: []TeamsWebhookSection{}, Summary: summary, ThemeColor: "0075FF"}
	return &webhook
}

// AddSectionWithText ...
func (webhook *TeamsWebhook) AddSectionWithText(title string, startGroup bool, text string) *TeamsWebhookSection {
	webhookSection := TeamsWebhookSection{Title: title, StartGroup: startGroup, Text: text}
	webhook.Sections = append(webhook.Sections, webhookSection)
	return &webhook.Sections[len(webhook.Sections)-1]
}

// AddSectionWithFacts ...
func (webhook *TeamsWebhook) AddSectionWithFacts(title string, startGroup bool, facts map[string]string) *TeamsWebhookSection {
	webhookSection := TeamsWebhookSection{Title: title, StartGroup: startGroup, Facts: []TeamsWebhookFact{}}
	for k, v := range facts {
		webhookSection.Facts = append(webhookSection.Facts, TeamsWebhookFact{Name: k, Value: v})
	}
	webhook.Sections = append(webhook.Sections, webhookSection)
	return &webhook.Sections[len(webhook.Sections)-1]
}
