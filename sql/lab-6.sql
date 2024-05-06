Lab 6.1 
CREATE TABLE IF NOT EXISTS user (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL
);


INSERT INTO user (user_name, password, salt, created_by)
VALUES ('admin', '051504d937d4351e2c5f72a84bb9944f', 'tBwTPvHD', 'admin');

adminP@ssw0rdtBwTPvHD

INSERT INTO user (user_name, password, salt, created_by)
VALUES ('puppy', 'f2c6e6e5451121ee31d0bc4c9bcf1e3c', 'iOlOTtU5', 'admin');

puppyP@ssw0rdiOlOTtU5
f2c6e6e5451121ee31d0bc4c9bcf1e3c

-- Lab 6.2

CREATE TABLE IF NOT EXISTS role (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name TEXT NOT NULL,
	role_desc TEXT,
	is_super_admin BOOLEAN,
	created_at DATETIME,
	created_by TEXT NOT NULL,
	status_id INTEGER
);


CREATE TABLE IF NOT EXISTS user_roles (
    user_role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
	user_id INTEGER,
	created_at DATETIME,
	created_by TEXT NOT NULL
);

-- drop table role_permissions 

CREATE TABLE IF NOT EXISTS role_permissions (
    role_permission_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER,
    resource_type_id INTEGER, -- 1 GraphQLResolve, 2 GraphQLField, 3 API Function, 4 Web Menu, 5 Mobile Menu
	resource_name TEXT,
    can_execute BOOLEAN,
    can_read BOOLEAN,
    can_write BOOLEAN,
    can_delete BOOLEAN,
	created_at DATETIME,
	created_by TEXT NOT NULL
);


INSERT INTO role (role_name, role_desc, is_super_admin, created_at, created_by, status_id)
VALUES
    ('Admin', 'Administration role with full access', 1, '2023-05-05 10:00:00', 'system', 1),
    ('Manager', 'Manager role with limited access', 0, '2023-05-05 10:01:00', 'system', 1),
    ('User', 'Regular user role with read-only access', 0, '2023-05-05 10:02:00', 'system', 1);

INSERT INTO user_roles (role_id, user_id, created_at, created_by)
VALUES
    (1, 1, '2023-05-05 10:10:00', 'system'),
    (2, 2, '2023-05-05 10:11:00', 'system')
  

    INSERT INTO role_permissions (role_id, resource_type_id, resource_name, can_execute, can_read, can_write, can_delete, created_at, created_by)
    VALUES
   
    (2, 1, 'contacts.gets', 1, 1, 0, 0, '2023-05-05 10:23:00', 'system'),
    (2, 1, 'contacts.getById', 1, 1, 0, 0, '2023-05-05 10:24:00', 'system'),
    (2, 1, 'contacts.getPagination', 0, 1, 0, 0, '2023-05-05 10:26:00', 'system'),
    (2, 1, 'contactMutations.createContact', 1, 1, 0, 0, '2023-05-05 10:25:00', 'system'),
    (2, 1, 'contactMutations.updateContact', 0, 1, 0, 0, '2023-05-05 10:25:00', 'system'),
    (2, 1, 'contactMutations.deleteContact', 0, 1, 0, 0, '2023-05-05 10:25:00', 'system')
   
drop view vw_user_role_permissions
CREATE VIEW vw_user_role_permissions AS
SELECT 
    u.user_id,
    u.user_name,
    r.role_id,
    r.role_name,
    r.role_desc,
    r.is_super_admin,
    ur.user_role_id,
    ur.created_at AS user_role_created_at,
    rp.role_permission_id,
    rp.resource_type_id,
    rp.resource_name,
    rp.can_execute,
    rp.can_read,
    rp.can_write,
    rp.can_delete,
    rp.created_at AS role_permission_created_at
FROM
    user_roles ur
    JOIN user u ON ur.user_id = u.user_id
    JOIN role r ON ur.role_id = r.role_id
    LEFT JOIN role_permissions rp ON r.role_id = rp.role_id



select user_id, resource_name, sum(can_execute) as can_execute ,
sum(is_super_admin) as is_super_admin
from vw_user_role_permissions 
where user_id = 2 and resource_name = 'contacts.gets' and resource_type_id =1
group by user_id, resource_name




select user_id, resource_name, sum(can_execute) as can_execute ,
sum(is_super_admin) as is_super_admin
from vw_user_role_permissions 
where user_id = 1 and resource_name = 'contacts.gets' and resource_type_id =1
group by user_id, resource_name

select user_id,
sum(is_super_admin) as is_super_admin
from vw_user_role_permissions 
where user_id = 1

select * from user