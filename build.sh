docker build -t telegram-bot .
docker run --env-file=token.env --rm telegram-bot
