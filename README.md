# Fayth

**Fayth** is a Go library for talking to AI models. Inspired by the mystical Fayth from *Final Fantasy X*.

---

##  What is Fayth?

Fayth provides a unified interface for working with different AI providers.

Currently, it supports **OpenAI**, with more providers planned in future updates.

---

##  Quick Start

Here's how to get up and running with OpenAI:

```go
import (
    "context"
    "log"
    "nyxze/fayth/model/openai"
)

// Initialize the model
model, err := openai.New()
if err != nil {
    log.Fatal(err)
}

// Send a message
msg := model.NewTextMessage(model.User, "Hello, world!")
generation, err := model.Generate(context.Background(), []model.Message{msg})
if err != nil {
    log.Fatal(err)
}

// Output the AI response
log.Println(generation.Text)
```

---

##  Installation

```bash
go get github.com/nyxze/fayth
```

---

##  Documentation

Coming soon. For now, check out the [examples](./examples) folder for working code snippets.

---

##  Roadmap

* [x] OpenAI integration
* [ ] Anthropic / Claude support
* [ ] Local model adapters (e.g., Ollama, LM Studio)
* [ ] Streaming support
* [ ] More robust message history helpers

---

## License

MIT — do what you want, just follow the teachings. 
