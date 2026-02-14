package export

import (
	"encoding/json"

	"github.com/ellebam/hireme/api/internal/domain"
)

// ParsePersonal parses personal section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParsePersonal(raw json.RawMessage) domain.PersonalContent {
	var content domain.PersonalContent
	_ = json.Unmarshal(raw, &content)
	return content
}

// ParseSummary parses summary section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseSummary(raw json.RawMessage) domain.SummaryContent {
	var content domain.SummaryContent
	_ = json.Unmarshal(raw, &content)
	return content
}

// ParseExperience parses experience section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseExperience(raw json.RawMessage) domain.ExperienceContent {
	var content domain.ExperienceContent
	_ = json.Unmarshal(raw, &content)
	return content
}

// ParseEducation parses education section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseEducation(raw json.RawMessage) domain.EducationContent {
	var content domain.EducationContent
	_ = json.Unmarshal(raw, &content)
	return content
}

// ParseSkills parses skills section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseSkills(raw json.RawMessage) domain.SkillsContent {
	var content domain.SkillsContent
	_ = json.Unmarshal(raw, &content)
	return content
}

// ParseLanguages parses languages section content from raw JSON.
// Returns zero-value struct on parse failure.
func ParseLanguages(raw json.RawMessage) domain.LanguagesContent {
	var content domain.LanguagesContent
	_ = json.Unmarshal(raw, &content)
	return content
}
