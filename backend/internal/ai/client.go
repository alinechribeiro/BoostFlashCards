package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	model  string
}

func New(apiKey, model string) *Client {
	if apiKey == "" {
		return nil
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &Client{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

type Flashcard struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type Topic struct {
	Name       string      `json:"name"`
	Flashcards []Flashcard `json:"flashcards"`
}

type Subject struct {
	Name   string  `json:"name"`
	Topics []Topic `json:"topics"`
}

type SubjectFlashcardsResponse struct {
	Topics []Topic `json:"topics"`
}

type NewSubjectResponse struct {
	Subject Subject `json:"subject"`
}

type Insight struct {
	Text string `json:"text"`
}

type InsightFlashcardSpec struct {
	Text  string `json:"text"`
	Style string `json:"style"` // "qa" or "true_false"
}

type InsightFlashcard struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

// cleanJSON makes the AI response more likely to be valid JSON by
// stripping markdown fences and fixing common trailing-comma patterns.
func cleanJSON(s string) string {
	s = strings.TrimSpace(s)

	// Strip ```json ... ``` or ``` ... ``` wrappers if present.
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i != -1 {
			s = s[i+1:]
		}
		if j := strings.LastIndex(s, "```"); j != -1 {
			s = s[:j]
		}
	}

	// Remove some common trailing comma constructs.
	replacements := []struct {
		old string
		new string
	}{
		{",\n}", "\n}"},
		{",\r\n}", "\r\n}"},
		{", }", " }"},
		{",\n]", "\n]"},
		{",\r\n]", "\r\n]"},
		{", ]", " ]"},
	}
	for _, r := range replacements {
		s = strings.ReplaceAll(s, r.old, r.new)
	}

	s = strings.TrimSpace(s)

	// If the model still wrapped the JSON in prose, try to slice out the JSON-ish part.
	firstBrace := strings.IndexAny(s, "{[")
	lastBrace := strings.LastIndexAny(s, "}]")
	if firstBrace != -1 && lastBrace != -1 && lastBrace > firstBrace {
		s = s[firstBrace : lastBrace+1]
	}

	return strings.TrimSpace(s)
}

type GradingResult struct {
	Score        float64         `json:"score"`
	MaxScore     float64         `json:"max_score"`
	Grade        string          `json:"grade"`
	Feedback     string          `json:"feedback"`
	GradeBand    string          `json:"grade_band,omitempty"`
	Explanation  string          `json:"explanation,omitempty"`
	Strengths    json.RawMessage `json:"strengths,omitempty"`
	Improvements json.RawMessage `json:"improvements,omitempty"`
}

func (g *GradingResult) StrengthsText() string {
	return rawTextOrBulletList(g.Strengths)
}

func (g *GradingResult) ImprovementsText() string {
	return rawTextOrBulletList(g.Improvements)
}

// rawTextOrBulletList tries to decode the AI field as either:
// - simple string
// - array of strings (joined as bullet list)
// and falls back to the raw JSON text if needed.
func rawTextOrBulletList(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return ""
	}

	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		return single
	}

	var list []string
	if err := json.Unmarshal(raw, &list); err == nil && len(list) > 0 {
		for i, item := range list {
			list[i] = strings.TrimSpace(item)
		}
		return "• " + strings.Join(list, "\n• ")
	}

	// Fallback: return best-effort plain text
	return s
}

func (c *Client) GeneratePracticeQuestionForSubject(ctx context.Context, subjectName string) (string, error) {
	if c == nil || c.client == nil {
		return "", errors.New("AI client not configured")
	}

	system := "You are an expert GCSE UK exam setter. Respond ONLY with a single exam-style question, no JSON, no explanation."
	user := fmt.Sprintf(`Give ONE past-paper-style GCSE %s question that a student can answer in a few paragraphs.

The question should be self-contained and clearly state what the student must do.`, subjectName)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.5,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned from OpenAI")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func (c *Client) GradeAnswerForSubject(ctx context.Context, subjectName, question, answer string) (*GradingResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("AI client not configured")
	}

	system := "You are an experienced GCSE UK examiner. Respond ONLY with strict JSON, no markdown, no commentary."
	user := fmt.Sprintf(`You are marking a GCSE %s answer.

QUESTION:
%s

STUDENT ANSWER:
%s

Mark the answer out of 30 and estimate the GCSE grade (1-9) that this performance corresponds to for this subject.

Return ONLY JSON in this exact shape:
{
  "score": number,          // marks awarded, 0-30
  "max_score": number,      // always 30
  "grade": "string",        // e.g. "7"
  "grade_band": "string",   // e.g. "strong pass", "grade 8-9 standard"
  "feedback": "string",     // short paragraph of feedback to the student
  "strengths": "string",    // bullet-style text of what went well
  "improvements": "string"  // bullet-style text of what to improve
}`, subjectName, question, answer)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.3,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}

	content := cleanJSON(resp.Choices[0].Message.Content)
	var out GradingResult
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, err
	}
	if out.MaxScore == 0 {
		out.MaxScore = 30
	}
	return &out, nil
}

func (c *Client) GenerateFlashcardsForSubject(ctx context.Context, subjectName string, numCards int) (*SubjectFlashcardsResponse, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("AI client not configured")
	}
	if numCards <= 0 || numCards > 40 {
		numCards = 10
	}

	system := "You are an expert GCSE UK tutor. Respond ONLY with pure JSON, no markdown, no comments."
	user := fmt.Sprintf(`For the GCSE subject %q, generate approximately %d high‑quality flashcards.
Return JSON matching:
{"topics":[{"name":"string","flashcards":[{"front":"string","back":"string"}]}]}
Make front a concise question and back a clear, exam‑relevant answer.`, subjectName, numCards)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.3,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}
	content := cleanJSON(resp.Choices[0].Message.Content)

	var out SubjectFlashcardsResponse
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GenerateSubjectFromDescription(ctx context.Context, description string) (*NewSubjectResponse, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("AI client not configured")
	}

	system := "You are an expert GCSE UK curriculum designer. Respond ONLY with pure JSON, no markdown, no comments."
	user := fmt.Sprintf(`Based on the following description, propose ONE GCSE subject with a few key topics and flashcards:
%s

Return JSON matching exactly:
{"subject":{"name":"string","topics":[{"name":"string","flashcards":[{"front":"string","back":"string"}]}]}}

Use concise, clear GCSE‑appropriate wording.`, description)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.4,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}
	content := cleanJSON(resp.Choices[0].Message.Content)

	var out NewSubjectResponse
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExtractInsightsFromText finds key examinable ideas in a passage.
func (c *Client) ExtractInsightsFromText(ctx context.Context, subjectName, text string, maxInsights int) ([]Insight, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("AI client not configured")
	}
	if maxInsights <= 0 || maxInsights > 30 {
		maxInsights = 12
	}

	system := "You are an expert GCSE UK tutor. Respond ONLY with pure JSON, no markdown, no comments."
	user := fmt.Sprintf(`Read the passage below and extract up to %d distinct, exam-relevant insights for the subject %s.

Each insight should be a short, self-contained statement that could become a flashcard.

Return ONLY JSON in this shape:
{"insights":[{"text":"string"}]}

PASSAGE:
%s`, maxInsights, subjectName, text)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.3,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}
	content := cleanJSON(resp.Choices[0].Message.Content)

	var out struct {
		Insights []Insight `json:"insights"`
	}
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, err
	}
	return out.Insights, nil
}

// GenerateFlashcardsFromInsights turns user-approved insights into ready-to-save flashcards.
func (c *Client) GenerateFlashcardsFromInsights(ctx context.Context, subjectName string, specs []InsightFlashcardSpec) ([]InsightFlashcard, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("AI client not configured")
	}
	if len(specs) == 0 {
		return nil, errors.New("no insights provided")
	}

	// Encode specs as JSON to feed into the prompt in a structured way.
	data, err := json.Marshal(specs)
	if err != nil {
		return nil, err
	}

	system := "You are an expert GCSE UK tutor. Respond ONLY with pure JSON, no markdown, no comments."
	user := fmt.Sprintf(`You are creating flashcards for GCSE %s.

You will receive a JSON array of items, each with:
- "text": an insight the student should learn
- "style": either "qa" for question/answer or "true_false" for a true/false card.

For each item, create ONE flashcard object with:
- "front": the prompt shown to the student
- "back": the ideal answer (or "True"/"False" plus a short justification).

Rules:
- If style is "qa": make "front" a clear exam-style question and "back" a concise, mark-scheme-like answer.
- If style is "true_false": make "front" a statement starting with "True or false: ..." and "back" start with either "True" or "False", followed by a short explanation.

Return ONLY JSON in this exact shape:
{"flashcards":[{"front":"string","back":"string"}]}

INSIGHTS JSON:
%s`, subjectName, string(data))

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.4,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}
	content := cleanJSON(resp.Choices[0].Message.Content)

	var out struct {
		Flashcards []InsightFlashcard `json:"flashcards"`
	}
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, err
	}
	return out.Flashcards, nil
}

