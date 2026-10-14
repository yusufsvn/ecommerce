alter table refreshtokens add constraint refreshtokens_family_id_key unique (family_id);
alter table refreshtokens add constraint refreshtokens_user_id_key unique (user_id);