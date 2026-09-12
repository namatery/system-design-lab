# URL Shortener

A service that creates compact links and redirects them to their original URLs.
This project focuses on a read-heavy design with fast redirects and durable
storage.

For example, a long URL like this:

```sh
https://example.com/articles/r/very-long-url
```

Becomes:

```sh
https://short.io/xyz12s
```

## Requirements

### Functional

- Generate a short URL from a long URL
- Redirect a short URL to its destination
- Redirect should be highly available
- Support optional expiration
- Collect basic analytics

### Non-Functional

#### Redirect

- Redirect must have low latency
- Redirect should remain available during individual service-instance failures
- URL mapping must be durable after successful creation

#### Link creation

- Link creation should prevent duplicate short codes
- Successfully created links must not be lost
- Link creation may tolerate higher latency than redirects

#### Analytics

- Analytics failures must not prevent redirects
- Analytics do not need to be immediately consistent

### Design assumptions

To make the design concrete, we assume the following targets.

#### Workload

| Metric | Assumption |
|---|---:|
| Redirects | 100 million/day |
| New links | 10 million/month |
| Peak traffic | 5× average |
| Traffic distribution | 80% of redirects target 20% of active links |
| Retention | 5 years for non-expired links |

Traffic is read-heavy and uneven: a small number of popular links may receive a
large share of redirects.

#### Data

- Original URLs average 300 bytes and are limited to 2,048 bytes.
- A stored link record is approximately 500 bytes before indexes and replication.
- Short codes are case-sensitive.
- Links are immutable after creation.
- Expired links are removed asynchronously.

#### Service targets

| Metric | Target |
|---|---:|
| Redirect latency | p95 server-side latency < 100 ms |
| Redirect availability | At least 99.99% |
| Link-creation latency | p95 server-side latency < 500 ms |
| Link-creation availability | At least 99.9% |
| Durability | An acknowledged link survives a process or node failure |
| Analytics consistency | Eventual; a delay of several minutes is acceptable |

The initial deployment runs in one region across multiple availability zones.
Redirects take priority over link creation and analytics during partial outages.

#### Derived capacity

| Resource | Average | Peak or five-year total |
|---|---:|---:|
| Redirect requests | ~1,157 requests/s | ~5,787 requests/s |
| Link-creation requests | ~4 requests/s | ~20 requests/s |
| Read/write ratio | ~300:1 | — |
| Stored links | — | 600 million records |
| Raw link data | — | ~300 GB |

The storage estimate excludes indexes, replicas, backups, cache entries, and
analytics events. A 99.99% redirect-availability target allows approximately
4.4 minutes of downtime per month.
