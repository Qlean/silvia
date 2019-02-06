alter table adjust.events
  ADD COLUMN error_type varchar(100),
  ADD COLUMN error_text varchar(65535),
  ADD COLUMN error_event varchar(65535);
