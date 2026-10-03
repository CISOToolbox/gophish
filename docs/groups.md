# Groups & targets

A **group** is a named list of **targets** (recipients). Campaigns are sent to one or more groups.

## Target fields

Each target has:

| Field | Description |
|---|---|
| **First Name** | Used by `{{.FirstName}}` |
| **Last Name** | Used by `{{.LastName}}` |
| **Email** | The recipient address (required) |
| **Position** | Used by `{{.Position}}` (e.g. job title) |

These values feed the [template variables](email-templates.md#template-variables), so personalization is only as good as the data in your groups.

## Creating a group

Add targets one by one in the console, or **bulk import** from CSV.

### CSV format

The importer expects a header row with these columns (names are matched case-insensitively):

```csv
First Name,Last Name,Email,Position
Ada,Lovelace,ada@example.com,Analyst
Alan,Turing,alan@example.com,Engineer
```

Only **Email** is strictly required; the other columns may be blank.

## Deduplication

When a campaign is created, recipients are de-duplicated by email address across the selected groups, so a person present in two groups receives a single email.

## API

See [REST API](api.md). Endpoints live under `/api/groups/`; CSV import is `POST /api/import/group`.
