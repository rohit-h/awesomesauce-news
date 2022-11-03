build:
	@go build .

run: build
	@TELEGRAM_BOT_TOKEN=dummy TELEGRAM_CHAT_ID=dummy ./awesomesauce-news

digest: build
	@TELEGRAM_BOT_TOKEN=172144934:AAEECulMQGvyGGcOzKNVMIGfMOlZBcZ02hU TELEGRAM_CHAT_ID=@top_hackernews ./awesomesauce-news -digest

clean: awesomesauce-news
	@rm -v awesomesauce-news

clean-db: awesomesauce.db
	@rm -v awesomesauce.db
