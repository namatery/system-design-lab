CREATE KEYSPACE IF NOT EXISTS url_shortener
WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'datacenter1': 1
};

CREATE TABLE IF NOT EXISTS url_shortener.links (
    alias text PRIMARY KEY,
    long_url text,
    created_at timestamp,
    expires_at timestamp
);

