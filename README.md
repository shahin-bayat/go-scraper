How to Start?

1. Create a Postgres db

```bash
	 docker run --name scraperdb \
	          -e POSTGRES_PASSWORD=goscraper \
	          -e POSTGRES_USER=shahin \
	          -e POSTGRES_DB=scraperdb \
	          -p 5432:5432 \
	          -d postgres
```

2. Run the server

```bash
  // go run main.go
```
