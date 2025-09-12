# Trope Field: Analysis and Best Practices

## System-Level vs. User-Level
Trope is best used as a system-level attribute, setting the narrative context for a batch or series of images, similar to artistic or literary style. This keeps prompts modular and future-proofs the design for multi-image storytelling. User-level use is possible but can blur the line between subject and story context, making outputs less clear and harder to scale for multi-image arcs.

## Integration with Other Fields
Trope should provide a narrative “frame” or archetype, not dictate the literal content of the image. It informs the AI about the type of story being told, so the generated image fits within a recognizable narrative arc, but does not override the specifics of subject, action, or emotion. To avoid dominance, prompt templates should give equal weight to all system-level fields, with user-level fields remaining the primary drivers of visual content. In a multi-image context, trope could define the sequence or progression, while in single-image mode, it simply colors the interpretation.

## CSV Curation Guidance
Ensure tropes.csv is weighted toward universal, cross-genre story patterns (the kind found in Campbell, Booker, or Propp). Avoid genre-locked or overly specific tropes unless you want to support genre-specific series. Consider periodic review and curation to keep the list balanced and inspiring.

# Unified Prompting for DalleDress: System vs User Prompts

## Overview

This document proposes a refactor of the DalleDress image generation pipeline to:
- Combine the current two-step prompt enhancement and image generation into a single OpenAI API call using system/user prompt separation.
- Clearly separate system-level ("stage direction") and user-level ("creative content") attributes.
- Add a new system-level attribute: `trope` (the story pattern or narrative arc for the image or series).

## Motivation

- **Reduce API calls**: Eliminate the separate "enhance" step, halving OpenAI usage and latency.
- **Cleaner prompts**: Prevent technical or constraint details from leaking into the creative prompt.
- **Consistency**: Make image generation more reproducible for identical data.
- **Extensibility**: Support richer creative storytelling via the new `trope` attribute.

## Attribute Classification

### System-Level (Stage Direction)
- `Color` (e.g., "Use only the colors #xxxxxx and #yyyyyy")
- `Orientation` (e.g., "landscape", "portrait")
- `BackStyle` (background style)
- `ArtStyle` (artistic style, e.g., "cubism")
- `LitStyle` (literary style, if present)
- `Gaze` (where the subject is looking)
- `Quality`, `Size`, `Style` (if present in payload)
- Explicit constraints (e.g., "Do not put text in the image.")

### User-Level (Creative Content)
- `Noun` (main subject)
- `Emotion` (feeling to convey)
- `Adjective`, `Adverb` (descriptors)
- `Occupation` (if relevant)
- `Action` (what the subject is doing)
- `Trope` (the story pattern or narrative arc)

## Proposed Prompt Structure

- **System Prompt**: Encodes all system-level attributes and constraints. Example:
  > You are an expert visual artist. Always use the following constraints:
  > - Use only the colors {{Color1}} and {{Color2}}.
  > - Orientation: {{Orientation}}.
  > - Artistic style: {{ArtStyle}}.
  > - Literary style: {{LitStyle}} (if present).
  > - Gaze: {{Gaze}}.
  > - Background style: {{BackStyle}}.
  > - Do not put text in the image.
  > - Quality: {{Quality}}. Size: {{Size}}.

- **User Prompt**: Encodes creative content. Example:
  > Draw a {{Adverb}} {{Adjective}} {{Noun}} with human-like characteristics, feeling {{Emotion}}, working as a {{Occupation}}, performing {{Action}}. Narrative trope: {{Trope}}.

## Implementation Steps

1. Refactor prompt generation to split system/user prompts.
2. Update OpenAI API call to use both system and user messages in a single request.
3. Remove the separate EnhancePrompt step.
4. Add `trope` to the DalleDress data model and prompt templates.
5. Test for reproducibility and quality.

## Open Questions
- Should `trope` ever influence system-level constraints (e.g., style or color)?
- Are there other attributes that should be reclassified?

## Affected Go Files and Expected Changes

The following Go files are expected to require changes for this refactor:

- **dalle/prompt.go**
  - Refactor prompt generation to split system and user prompts.
  - Remove the EnhancePrompt and enhancePromptWithClient logic.
  - Update template logic to support the new `trope` attribute.

- **dalle/context.go**
  - Update DalleDress and related data structures to include `trope`.
  - Refactor prompt construction to generate both system and user messages.
  - Remove or update any logic that previously called EnhancePrompt.

- **dalle/image.go**
  - Update RequestImage to accept both system and user prompts and send them in a single OpenAI API call.
  - Remove any logic that expects an enhanced prompt as a separate step.

- **dalle/series.go** (if present)
  - Update series creation logic to support the new `trope` attribute.

- **Any test files**
  - Update or add tests to cover the new prompt structure and attribute handling.

Other files may require minor updates depending on how prompts and attributes are passed through the system.

## References
- [Prompt Engineering Guide](https://www.promptingguide.ai/introduction/elements)
- [Optimizing Prompts](https://www.promptingguide.ai/guides/optimizing-prompts)

---

**Author:** GitHub Copilot
**Date:** 2025-09-11
