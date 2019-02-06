alter table adjust.events
  ADD error_type varchar(256) default NULL;

alter table adjust.events
  ADD error_text varchar(65535) default NULL;

alter table adjust.events
  ADD error_event varchar(65535) default NULL;
