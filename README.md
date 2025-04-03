# SpamGoose Bot

SpamGoose is a Telegram bot inspired by [Jack Rhysider's](https://open.spotify.com/episode/67cV8KrMIH2S5NN9IoU8xr) experiment. It allows users to post messages directly to a specific Telegram channel. Messages are filtered to ensure they do not contain links or files, and valid messages are tagged with `#frombot`.

## Features

- Accepts text messages from users.
- Posts valid messages to the `@ragoose_dumps` Telegram channel.
- Filters out messages containing links, files, or `@` mentions.

## Requirements

- Go 1.20 or later
- A Telegram bot token (set as an environment variable `TELEGRAM_BOT_TOKEN`)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/daveragos/spamGoose.git
   cd spamGoose

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.