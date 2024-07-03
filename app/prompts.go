package app

var promptTemplate = `{{.LitPrompt false}}Here's the prompt:

Draw a {{.Adverb false}} {{.Adjective false}} {{.Noun true}} with human-like
characteristics feeling {{.Emotion false}}{{.Occupation false}}.

Noun: {{.Noun false}} with human-like characteristics.
Emotion: {{.Emotion false}}.
Occupation: {{.Occupation false}}.
Action: {{.Action false}}.
Artistic style: {{.ArtStyle false 1}}.
{{if .HasLitStyle}}Literary Style: {{.LitStyle false}}.
{{end}}Use only the colors {{.Color true 1}} and {{.Color true 2}}.
{{.Orientation false}}.
{{.BackStyle false}}.

Emphasize the emotional aspect of the image. Look deeply into and expand upon the
many connotative meanings of "{{.Noun true}}," "{{.Emotion true}}," "{{.Adjective true}},"
and "{{.Adverb true}}." Find the representation that most closely matches all the data.

Focus on the emotion, the noun, and the styles.`

var dataTemplate = `Orig:             {{.Orig}}
Seed:             {{.Seed}}
Adverb:           {{.Adverb false}}
AdverbShort:      {{.Adverb true}}
Adjective:        {{.Adjective false}}
AdjectiveShort:   {{.Adjective true}}
Noun:             {{.Noun false}}
NounShort:        {{.Noun true}}
Emotion:          {{.Emotion false}}
EmotionShort:     {{.Emotion true}}
Occupation:       {{.Occupation false}}
OccupationShort:  {{.Occupation true}}
Action:	          {{.Action false}}
ActionShort:	  {{.Action true}}
ArtStyle 1:       {{.ArtStyle false 1}}
ArtStyleShort 1:  {{.ArtStyle true 1}}
ArtStyle 2:       {{.ArtStyle false 2}}
ArtStyleShort 2:  {{.ArtStyle true 2}}
{{if .HasLitStyle}}LitStyle:         {{.LitStyle true}}
LitStyleShort:    {{.LitStyle false}}
{{end}}Color 1:          {{.Color false 1}}
Color 2:          {{.Color false 2}}
Color 3:          {{.Color false 3}}
Orientation:      {{.Orientation false}}
OrientationShort: {{.Orientation true}}
Gaze:             {{.Gaze false}}
GazeShort:        {{.Gaze true}}
BackStyle:        {{.BackStyle false}}
BackStyleShort:   {{.BackStyle true}}`

var terseTemplate = `{{.Adverb false}} {{.Adjective false}} {{.Noun true}} with human-like characteristics feeling {{.Emotion false}}{{.Occupation false}} in the style of {{.ArtStyle true 1}}`

var authorTemplate = `{{if .HasLitStyle}}You are an award winning author who writes in the literary
style called {{.LitStyle true}}. Take on the persona of such an author.
{{.LitStyle true}} is a genre or literary style that {{.LitStyleDescr}}.{{end}}`
