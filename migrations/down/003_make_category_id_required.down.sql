-- Rollback: Make category_id optional again

-- Step 1: Drop foreign key constraint
ALTER TABLE todos
DROP CONSTRAINT IF EXISTS todos_category_id_fkey;

-- Step 2: Remove NOT NULL constraint
ALTER TABLE todos
ALTER COLUMN category_id DROP NOT NULL;

-- Step 3: Re-add foreign key constraint with SET NULL (allow null on delete)
ALTER TABLE todos
ADD CONSTRAINT todos_category_id_fkey
FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;

