UPDATE organization_settings
SET name = 'Garago'
WHERE id = 1 AND BTRIM(name) = '';
