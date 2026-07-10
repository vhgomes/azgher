CREATE TABLE tags (
    id       SERIAL PRIMARY KEY,
    category TEXT NOT NULL,
    value    TEXT NOT NULL,
    CONSTRAINT uq_tags_category_value UNIQUE (category, value)
);
