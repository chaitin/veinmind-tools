package sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const trainingPrefix = `你是一个云原生安全研究员,下面的json是云安全事件日志,请详细描述下列安全事件,指出可能会造成的风险,并给出解决方案`

func Dialogue(ctx context.Context, token, prefix string, content string) (openai.ChatCompletionResponse, error) {
	client := openai.NewClient(token)
	if prefix == "" {
		prefix = trainingPrefix
	}
	return client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: strings.Join([]string{prefix, content}, ":"),
				},
			},
		},
	)
}

func DialogueStream(ctx context.Context, token, prefix string, content string) (*openai.ChatCompletionStream, error) {
	client := openai.NewClient(token)
	if prefix == "" {
		prefix = trainingPrefix
	}
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: strings.Join([]string{prefix, content}, ":"),
			},
		},
		Stream: true,
	}
	return client.CreateChatCompletionStream(ctx, req)
}

func Read(stream *openai.ChatCompletionStream) error {
	defer func() {
		fmt.Println()
		stream.Close()
	}()
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf(resp.Choices[0].Delta.Content)
	}
}
