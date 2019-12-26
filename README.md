# mtsql

Run SQL queries against CSV files.

## Usage

```
mtsql "SELECT * FROM cities"
```

```
mtsql "SELECT City, State FROM cities"
```

```
mtsql "PROFILE SELECT City, State FROM cities WHERE State = 'WA'"
```

## Development

```
go build .
go test ./...
```

