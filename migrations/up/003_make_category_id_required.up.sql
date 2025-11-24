-- Make category_id required in todos table
-- This migration updates existing todos and adds NOT NULL constraint

-- Step 1: Ensure at least one category exists
-- If no categories exist, create a default one
INSERT INTO categories (name, color, created_at, updated_at)
SELECT 'Default', '#3B82F6', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM categories LIMIT 1);

-- Step 2: Update todos with NULL category_id to use the first category
UPDATE todos
SET category_id = (SELECT id FROM categories ORDER BY id LIMIT 1)
WHERE category_id IS NULL;

-- Step 3: Drop existing foreign key constraint if exists
ALTER TABLE todos
DROP CONSTRAINT IF EXISTS todos_category_id_fkey;

-- Step 4: Add NOT NULL constraint
ALTER TABLE todos
ALTER COLUMN category_id SET NOT NULL;

-- Step 5: Re-add foreign key constraint with RESTRICT (prevent delete if used)
ALTER TABLE todos
ADD CONSTRAINT todos_category_id_fkey
FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT;

