package bot

import (
	"fmt"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"

	"sea_battle/internal/domain"
)

type AIBot struct {
	field      *domain.Field
	OpenAIClient openai.Client
	api_key    string
}

func NewAIBot(field *domain.Field) *AIBot {
	key, _ := os.LookupEnv("AI_KEY")
	return &AIBot{
		field: field,
		OpenAIClient: openai.NewClient(option.WithAPIKey(key)),
		api_key: key,
	}
}

func (ab *AIBot) translate_field(field *domain.Field) string {
	var res strings.Builder

	for i := 0; i < domain.Size; i++ {
		for j := 0; j < domain.Size; j++ {
			res.WriteString(strconv.Itoa(field.Matrix[i][j]) + " ")
		}
		res.WriteString("\n")
	}

	return res.String()
}

func (ab *AIBot) Shoot() domain.Pair {
	curr_field := ab.translate_field(ab.field)
	var f, s int

	for {
		req := fmt.Sprintf("You recieve a current field positon: \n%s\n, choose where you want to shoot, Reply with two integers only (e.g. \"1 2\")", curr_field)

		resp, err := ab.OpenAIClient.Responses.New(context.Background(), responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(req)},
			Model: openai.ChatModelChatgpt4oLatest,
		}, )

		if err != nil {
			return domain.Pair{X: 11, Y: 11}
		}

		p := strings.Split(resp.OutputText(), " ")
		f, _ := strconv.Atoi(p[0])
		s, _ := strconv.Atoi(p[1])
		if f >= 0 && f < domain.Size && s >= 0 && s < domain.Size {
			break
		}
	}
	
	return domain.Pair{X: f, Y: s}
}

func (ab *AIBot) Getter() *domain.Field {
	return ab.field
}

func (ab *AIBot) SetResult(shotRes domain.ShotResult) {}

