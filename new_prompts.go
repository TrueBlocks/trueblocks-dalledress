package main

var promptTemplate = `Draw a {{.Adverb true}} {{.Adjective true}} {{.Noun true}} who works as a
{{.Occupation true}} and is {{.Action false}} and feeling {{.Emotion false}}.
Noun: {{.Noun false}}.
Emotion: {{.Emotion false}}.
Occupation: {{.Occupation false}}.
Action: {{.Action false}}.
Artistic style: {{.ArtStyle false 1}}.
Use only the colors {{.Color true 1}} and {{.Color true 2}}.
{{.Orientation false}}.
{{.BackStyle false}}.
Expand upon the most relevant connotative meanings of {{.Noun true}}, {{.Emotion true}}, {{.Adjective true}}, and {{.Adverb true}}.
Find the representation that most closely matches the description.
Focus on the noun, the occupation, the emotion, and literary style.
Rewrite the prompt using the literary style: {{.LitStyle false}}`

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
LitStyle:         {{.LitStyle true}}
LitStyleShort:    {{.LitStyle false}}
Color 1:          {{.Color false 1}}
Color 2:          {{.Color false 2}}
Color 3:          {{.Color false 3}}
Orientation:      {{.Orientation false}}
OrientationShort: {{.Orientation true}}
Gaze:             {{.Gaze false}}
GazeShort:        {{.Gaze true}}
BackStyle:        {{.BackStyle false}}
BackStyleShort:   {{.BackStyle true}}`

var terseTemplate = `{{.Adverb false}} {{.Adjective false}} {{.Noun true}} {{.Occupation true}} {{.Action false}} and feeling {{.Emotion true}} in the style of {{.ArtStyle true 1}}`
