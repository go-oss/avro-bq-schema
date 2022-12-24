# avro-bq-schema

Convert [Apache Avro](https://avro.apache.org/docs/1.11.1/specification/) Schema to [BigQuery Table Schema](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables#TableSchema).

```sh
go install github.com/go-oss/avro-bq-schema@latest
```

## Usage

```sh
avro-bq-schema schema.avsc > bq.json
```
