CREATE TABLE IF NOT EXISTS tweets (
    id UUID PRIMARY KEY,
    text TEXT NOT NULL,
    hint TEXT,
    answer VARCHAR(10) CHECK (answer IN ('positive', 'negative', 'neutral'))
);

-- Check if the table is empty before loading data
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM tweets LIMIT 1) THEN
        -- Load CSV data into the tweets table
        COPY tweets(id, text, hint, answer)
        FROM '/docker-entrypoint-initdb.d/data.csv'
        DELIMITER ','
        CSV HEADER;
    END IF;
END $$;
