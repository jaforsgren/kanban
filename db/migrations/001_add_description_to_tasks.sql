DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'tasks'
        AND column_name = 'description'
    ) THEN
        ALTER TABLE tasks ADD COLUMN description TEXT DEFAULT '';
    END IF;
END $$;