-- Daily partitions bound the size of each link's event history.
CREATE TABLE IF NOT EXISTS url_shortener.click_events (
    alias text,
    day date,
    event_id timeuuid,
    clicked_at timestamp,
    referrer text,
    user_agent text,
    PRIMARY KEY ((alias, day), event_id)
) WITH CLUSTERING ORDER BY (event_id DESC);
