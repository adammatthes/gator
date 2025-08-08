#!/bin/bash
goose -dir sql/schema postgres "$DATABASE_URL" down
