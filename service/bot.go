package service

import (
	"encoding/json"
	"fmt"
)

var botUrl string = "https://api.red-pill.ai/v1/chat/completions"
var botHeaders map[string]string = map[string]string{
	"Authorization": "Bearer sk-5x3E9Bk8az6t3JD1BZfWJlsPv0cyw2iJ1VJ3WqlWagjW41J0",
	"Content-Type":  "application/json; charset=utf-8",
}

type BotMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type BotContent struct {
	Messages []BotMessage `json:"messages"`
	Model    string       `json:"model"`
	Agent    string       `json:"agent"`
	Stream   bool         `json:"stream"`
}

type BotResponse struct {
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

func GetBotResponse(content string) (string, error) {
	c := BotContent{
		Messages: []BotMessage{
			{
				Content: content,
				Role:    "user",
			},
		},
		Model:  "gpt-4o",
		Agent:  "vitalikbuterin",
		Stream: false,
	}
	res, err := HttpPost(botUrl, c, botHeaders)
	if err != nil {
		return "", err
	}

	var bs BotResponse
	err = json.Unmarshal(res, &bs)
	if len(bs.Error.Message) > 0 {
		return "", fmt.Errorf(bs.Error.Message)
	}
	if len(bs.Choices) < 1 {
		return "", err
	}
	return bs.Choices[0].Message.Content, err
}

/*
curl -X "POST" "https://api.red-pill.ai/v1/chat/completions" \
     -H 'Authorization: Bearer sk-5x3E9Bk8az6t3JD1BZfWJlsPv0cyw2iJ1VJ3WqlWagjW41J0' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d '{
    "messages": [
        {
            "content": "Hello everyone! limit 3 words",
            "role": "user"
        }
    ],
    "model": "gpt-4o",
    "agent": "vitalikbuterin",
    "stream": false
}'

{
    "id": "chatcmpl-AHSAf7dwnwTpXuV2mmzX4Llc4Yiyf",
    "object": "chat.completion",
    "created": 1728723229,
    "model": "gpt-4o-2024-05-13",
    "choices": [
        {
            "index": 0,
            "message": {
                "role": "assistant",
                "content": "Hi there! Ask away!",
                "refusal": null
            },
            "logprobs": null,
            "finish_reason": "stop"
        }
    ],
    "usage": {
        "prompt_tokens": 5000,
        "completion_tokens": 6,
        "total_tokens": 5006,
        "prompt_tokens_details": {
            "cached_tokens": 0
        },
        "completion_tokens_details": {
            "reasoning_tokens": 0
        }
    },
    "system_fingerprint": "fp_51ac368d59"
}

{
  "error": {
    "message": "Invalid value: 'an'. Supported values are: 'system', 'assistant', 'user', 'function', and 'tool'.",
    "type": "invalid_request_error",
    "param": "messages[1].role",
    "code": "invalid_value"
  }
}
*/
