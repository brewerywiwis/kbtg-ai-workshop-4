# Transfer API OpenAPI Specification

This directory contains the OpenAPI/Swagger specifications for the LBK Points Transfer API.

## Files

- `transfer.yaml` - Transfer API specification (referenced from external source)

## External Reference

The full transfer API specification can be found at:
https://github.com/mikelopster/kbtg-ai-workshop-nov/blob/main/workshop-4/specs/transfer.yml

## Key Endpoints

- `POST /transfers` - Create transfer
- `GET /transfers/{id}` - Get transfer by idempotency key
- `GET /transfers?userId={id}` - List user transfers (paginated)

## Authentication

The current implementation is public (no authentication required) as specified in the original API spec.
