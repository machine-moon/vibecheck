CREATE TABLE IF NOT EXISTS tweets (
    id UUID PRIMARY KEY,
    text TEXT NOT NULL,
    hint TEXT,
    answer VARCHAR(10) CHECK (answer IN ('positive', 'negative', 'neutral'))
);

-- Load CSV data into the tweets table
COPY tweets(id, text, hint, answer)
FROM '/docker-entrypoint-initdb.d/data.csv'
DELIMITER ','
CSV HEADER

