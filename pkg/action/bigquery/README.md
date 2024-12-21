# BigQuery

Actions for BigQuery data integration.

## Prerequisites

- **Google Cloud account**: You need to have a Google Cloud account and the necessary permissions to create and modify BigQuery datasets and tables.
- **BigQuery dataset**: You need to have a BigQuery dataset to insert the alert into. You can find instructions on how to create a dataset [here](https://cloud.google.com/bigquery/docs/datasets).

## `bigquery.insert_alert`

This action inserts an alert into BigQuery.

### Arguments

Example policy:

```rego
run contains res if {
  res := {
    id: "your-action",
    uses: "bigquery.insert_alert",
    args: {
      "project_id": "my-project",
      "dataset_id": "my-dataset",
      "table_id": "my-table",
    },
  },
}
```

- `project_id` (string, required): Specifies the ID of the Google Project to insert the alert into.
- `dataset_id` (string, required): Specifies the ID of the BigQuery dataset to insert the alert into.
- `table_id` (string, required): Specifies the ID of the BigQuery table to insert the alert into. If the table does not exist, it will be created.

### Table schema

| Field name | Type      | Mode     |
|------------|-----------|----------|
| id         | STRING    | REQUIRED |
| schema     | STRING    | REQUIRED |
| created_at | TIMESTAMP | REQUIRED |
| title      | STRING    | REQUIRED |
| description| STRING    | REQUIRED |
| source     | STRING    | REQUIRED |
| namespace  | STRING    | REQUIRED |
| attrs      | RECORD    | REPEATED |
| └─ id      | STRING    | REQUIRED |
| └─ key     | STRING    | REQUIRED |
| └─ value   | STRING    | REQUIRED |
| └─ type    | STRING    | REQUIRED |
| └─ ttl     | INTEGER   | REQUIRED |
| └─ global  | BOOLEAN   | REQUIRED |
| refs       | RECORD    | REPEATED |
| └─ Title   | STRING    | REQUIRED |
| └─ URL     | STRING    | REQUIRED |
| data       | JSON      | REQUIRED |


## `bigquery.insert_data`

This action inserts any data into BigQuery.

### Arguments

Example policy:

```rego
run contains res if {
  res := {
    id: "your-action",
    uses: "bigquery.insert_data",
    args: {
      "project_id": "my-project",
      "dataset_id": "my-dataset",
      "table_id": "my-table",
      "data": {
        "name": "John Doe",
        "age": 42,
      },
    },
  },
}
```

- `project_id` (string, required): Specifies the ID of the Google Project to insert the data into.
- `dataset_id` (string, required): Specifies the ID of the BigQuery dataset to insert the data into.
- `table_id` (string, required): Specifies the ID of the BigQuery table to insert the data into. If the table does not exist, it will be created.
- `tags` (array of string, optional): Specifies the tags to insert into the table. This field will be REPEATED STRING field type in BigQuery.
- `data` (object, required): Specifies the data to insert into the table. This field will be JSON field type in BigQuery.

### Table schema

| Field name | Type      | Mode     |
|------------|-----------|----------|
| id         | STRING    | REQUIRED |
| alert_id   | STRING    | REQUIRED |
| created_at | TIMESTAMP | REQUIRED |
| tags       | STRING    | REPEATED |
| data       | JSON      | REQUIRED |
