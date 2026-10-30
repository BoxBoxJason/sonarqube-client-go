---
title: Errors and pagination
weight: 50
---

## Errors

Options are validated before any request is sent: a missing or invalid parameter is returned as a [`*ValidationError`](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar#ValidationError), naming the field. When SonarQube answers with an error status, the error is a [`*ResponseError`](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar#ResponseError), carrying the status code, the message and the raw body.

{{% example "ExampleValidationError" %}}

To branch on a status code without unwrapping the error, use the helpers `IsNotFound`, `IsUnauthorized`, `IsForbidden`, `IsConflict`, `IsRateLimited` and `IsServerError`:

{{% example "ExampleIsNotFound" %}}

## Pagination

Paginated V1 endpoints take `PaginationArgs` (`Page`, starting at 1, and `PageSize`, up to 500) and return a `Paging` with the page index, the page size and the total:

{{% example "ExampleIssuesService_Search" %}}

Most of them also have an `...All` variant (`SearchAll`, `ListAll`) that walks through every page and returns the merged results:

{{% example "ExampleUsersService_SearchAll" %}}

{{< callout type="warning" >}}
SonarQube refuses to page past 10,000 results on some endpoints, such as the issues search. Narrow the search with filters rather than fetching everything on large instances.
{{< /callout >}}

In the CLI, the same behavior is available with `--all` (see [CLI pagination]({{< relref "cli#pagination" >}})).
