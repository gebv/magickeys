-- DATABASE magickeys

-- \df
create function sort_text_array(text[]) returns text[][] as $$
    select array_agg(n) from (select n from unnest($1) as t(n) order by n) as a;
$$ language sql immutable;

CREATE TABLE values (
    value_id uuid NOT NULL PRIMARY KEY,
    keys text[],
    value jsonb NOT NULL DEFAULT '{}',
    is_removed boolean DEFAULT false,
    created_at timestamp,
    updated_at timestamp DEFAULT now()
);
CREATE INDEX values_keys_gin_idx on values USING GIN (keys);
CREATE UNIQUE INDEX values_uniq_ifcontainuniqkey_idx on values (sort_text_array(keys))
    WHERE keys @> '{uniq}';