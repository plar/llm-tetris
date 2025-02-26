# LLM Tetris Console Game

## Context and Motivation

This is a classic Tetris game developed in Go, designed to run in a console environment.
The project took about *three* hours to complete from scratch. This project has ~50% test coverage and can be used as a good starting point for further LLM exploration. It has own issues and flows that can be improved be senior developers. :)
It was created as a fun project to showcase the capabilities of modern LLMs in generating code.

I’ve been using this approach for about a year (with my own prompts and tools), and recently, [Harper Reed explained](https://harper.blog/2025/02/16/my-llm-codegen-workflow-atm/) it quite well. I decided to showcase his method(prompts) with this project to demonstrate its potential to a wider audience.

### Why Tetris?

Tetris has long been a popular choice for developers as a test project when exploring new programming languages, offering a perfect blend of complexity and enjoyment. Throughout my career, I've created approximately 20 different versions of Tetris, utilizing various languages. It is time of ChatGPT to create a new version of Tetris. :)

![LLM-Tetris](tetris.svg)

## Used tools

- [Go](https://go.dev/)
- VSCode with [continue.dev](https://continue.dev/) extension + ChatGPT4o (paid)
- ChatGPT for initial prompts

All Chat sessions can be found in the [chats](chats) folder. This folder contains all the prompts and responses from ChatGPT and continue.dev.

- [chats/session-01-chatgpt.md](chats/session-01-chatgpt.md) - Harper's initial Idea prompt
- [chats/session-02-chatgpt.md](chats/session-02-chatgpt.md) - Harper's TDD prompt
- [chats/session-03-continuedev.md](chats/session-03-continuedev.md) - Sessions from VSCode/continue.dev with prompts from PROMPT-PLAN.md
- [chats/session-04-continuedev.md](chats/session-04-continuedev.md) - Sessions from VSCode/continue.dev with prompts from PROMPT-PLAN.md
- [chats/session-05-continuedev.md](chats/session-05-continuedev.md) - Sessions from VSCode/continue.dev with prompts from PROMPT-PLAN.md

## Installation

### Run from Source

You need to have Go installed on your system. You can download it from the [Go website](https://golang.org/dl/).

Clone the repository: `$ git clone https://github.com/plar/llm-tetris.git`

Navigate to the project directory: `$ cd llm-tetris`

Run the game: `$ go run .`

### Build from Source

#### Compile for Different Platforms

Run the following commands in your terminal to build binaries for each target platform:

# For Linux
```
GOOS=linux GOARCH=amd64 go build -o tetris
```

# For Windows

```
GOOS=windows GOARCH=amd64 go build -o tetris.exe
```

# For macOS

```
GOOS=darwin GOARCH=amd64 go build -o tetris
```

## Usage

Run the game by executing the binary:

### On Windows

```shell
tetris.exe
```

### On Linux/macOS

```shell
./tetris
```

### Controls

- **Move Left**: Left Arrow
- **Move Right**: Right Arrow
- **Soft Drop**: Down Arrow
- **Rotate Clockwise**: Up Arrow
- **Hard Drop**: Spacebar
- **Pause**: P

## Configuration

The game supports configuration through a JSON file `config.json`:

- **Key remapping**: Adjust the keys for game controls.

Upon the first run, a default `config.json` file will be generated if one doesn't exist.

## Persistence

High scores are stored in `config.json` and persist across game sessions.

## License

This project is licensed under the MIT License.
