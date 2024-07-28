-- Insert roles into roles table if they do not already exist
INSERT INTO roles (id, machine_name, name, description)
SELECT v.id, v.machine_name, v.name, v.description
FROM (
  VALUES
    (
      uuid_generate_v4(),
      'administrator',
      'Administrator',
      'Has broad access to system functions and data, but may have some restrictions. Can manage content, but cannot manage user or modify system-wide settings.'
    ),
    (
      uuid_generate_v4(),
      'user',
      'User',
      'Regular user with basic privileges. Can access and interact with the system based on predefined permissions, but cannot modify system settings or manage other users.'
    )
) AS v(id, machine_name, name, description)
WHERE NOT EXISTS (
  SELECT 1 FROM roles WHERE roles.machine_name = v.machine_name
);
