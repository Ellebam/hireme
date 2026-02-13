package export

import (
	"encoding/json"

	"github.com/ellebam/hireme/api/internal/domain"
)

// ParsePersonal parses personal section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParsePersonal(raw json.RawMessage) domain.PersonalContent {
	var c domain.PersonalContent
	_ = json.Unmarshal(raw, &c)
	return c
}

// ParseSummary parses summary section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseSummary(raw json.RawMessage) domain.SummaryContent {
	var c domain.SummaryContent
	_ = json.Unmarshal(raw, &c)
	return c
}

// ParseExperience parses experience section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseExperience(raw json.RawMessage) domain.ExperienceContent {
	var c domain.ExperienceContent
	_ = json.Unmarshal(raw, &c)
	return c
}

// ParseEducation parses education section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseEducation(raw json.RawMessage) domain.EducationContent {
	var c domain.EducationContent
	_ = json.Unmarshal(raw, &c)
	return c
}

// ParseSkills parses skills section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseSkills(raw json.RawMessage) domain.SkillsContent {
	var c domain.SkillsContent
	_ = json.Unmarshal(raw, &c)
	return c
}

// ParseLanguages parses languages section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseLanguages(raw json.RawMessage) domain.LanguagesContent {
	var c domain.LanguagesContent
	_ = json.Unmarshal(raw, &c)
	return c
}
