# Gabriel Novellia Assignment

## How to build
`docker build -t app`

## How to run
`docker run app`

## Endpoints

POST Import: `0.0.0.0:8080/import`

GET Records: `0.0.0.0:8080/records`

GET Records (with id): `0.0.0.0:8080/records/:id`

POST Transform: `0.0.0.0:8080/transform`

GET Analytics: `0.0.0.0:8080/analytics`

## Example Endpoint usage

All jsons match the examples given in PDF

```
curl --location --request POST '0.0.0.0/import' \
--header 'Content-Type: text/plain' \
--data-raw '{JSONL HERE}'
```

```
curl --location --request GET '0.0.0.0/records/{id}'
```
```
curl --location --request GET '0.0.0.0/records?subject=Patient/PT-007&resourceType=MedicationRequest
```
```
curl --location --request POST '0.0.0.0/transform' \
--header 'Content-Type: text/plain' \
--data-raw '{JSONL HERE}'
```
```
curl --location --request GET '0.0.0.0/analytics'
```
