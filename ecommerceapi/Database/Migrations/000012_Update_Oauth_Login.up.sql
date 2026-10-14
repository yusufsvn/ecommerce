-- Add new columns for OAuth login
Alter table users add column oauth_provider varchar(255) default null;
Alter table users add column oauth_id varchar(255) default null;

-- Change the constraints
Alter table users alter column phone drop not null;
alter table users alter column password drop not null;