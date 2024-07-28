-- Insert role into roles table if it not already exists
INSERT INTO roles (id, machine_name, name, description)
SELECT uuid_generate_v4(), 'super_admin', 'Super Admin', 'Has full access and control over all system functions and data.'
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE machine_name = 'super_admin');

-- Insert users into users table if they do not already exist
INSERT INTO users (id, name, work_email, oid)
SELECT v.id, v.name, v.work_email, v.oid::uuid
FROM (
  VALUES
  (uuid_generate_v4(), 'Larry Nguyen (TECHVIFY.ITS)', 'linhnd@techvify.com.vn', '93571118-ebcf-4910-8f33-e638d95cae99')
) AS v(id, name, work_email, oid)
WHERE NOT EXISTS (
  SELECT 1 FROM users
  WHERE work_email IN ('linhnd@techvify.com.vn')
);

-- Insert users with super admin user role if they do not already exist
INSERT INTO user_roles (id, user_id, role_id)
SELECT uuid_generate_v4(), users.id, roles.id
FROM users
CROSS JOIN roles
WHERE users.work_email IN ('linhnd@techvify.com.vn')
  AND roles.machine_name = 'super_admin'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles
    WHERE user_roles.user_id = users.id
    AND user_roles.role_id = roles.id
  );
