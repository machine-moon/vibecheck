-- Create the tweets table
CREATE TABLE IF NOT EXISTS tweets (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    hint TEXT,
    answer VARCHAR(10) CHECK (answer IN ('positive', 'negative', 'neutral'))
);

-- Load CSV data into the tweets table
COPY tweets(text, hint, answer)
FROM '/docker-entrypoint-initdb.d/create_tables.csv'
DELIMITER ','
CSV HEADER;