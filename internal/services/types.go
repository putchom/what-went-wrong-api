package services

type Entitlements struct {
	MaxGoals               int  `json:"maxGoals"`
	LogRetentionDays       *int `json:"logRetentionDays"` // nil = unlimited
	CanUseAiExcuse         bool `json:"canUseAiExcuse"`
	CanUsePremiumTemplates bool `json:"canUsePremiumTemplates"`
}
