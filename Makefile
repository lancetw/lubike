test:
	@govendor test -cover +local

deploy:
	@git push heroku master

run:
	@go install .
	@heroku local
