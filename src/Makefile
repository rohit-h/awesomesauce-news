build:
	@go build .

run:
	@TELEGRAM_BOT_TOKEN=dummy TELEGRAM_CHAT_ID=dummy ./awesomesauce-news

clean: awesomesauce-news
	@rm -v awesomesauce-news

clean-db: awesomesauce.db
	@rm -v awesomesauce.db
