create table professional (
  id_professional SERIAL PRIMARY KEY,
  inserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  reference_key TEXT NOT NULL UNIQUE,
  especialidade TEXT NOT NULL,
  nome TEXT NOT NULL
);
create index idx_professional_insert on professional (inserted_at);
create index idx_professional_update on professional (updated_at);
create index idx_professional_reference_key on professional (reference_key);
create index idx_professional_especialidade on professional (especialidade);
create index idx_professional_nome on professional (nome);

create table attribute (
  id_attribute SERIAL PRIMARY KEY,
  inserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  id_professional INTEGER NOT NULL,
  attribute TEXT NOT NULL,
  value TEXT NOT NULL,
  FOREIGN KEY(id_professional) REFERENCES professional(id_professional),
  UNIQUE (id_professional, attribute, value)
);
create index idx_attribute_insert on attribute (inserted_at);
create index idx_attribute_update on attribute (updated_at);
create index idx_attribute_professional on attribute (id_professional);
create index idx_attribute_attribute on attribute (attribute);
create index idx_attribute_value on attribute (value);

create table blocker (
  id_blocker SERIAL PRIMARY KEY,
  inserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  title TEXT NOT NULL,
  description TEXT,
  id_professional INTEGER NOT NULL,
  init_datetime TIMESTAMP NOT NULL,
  end_datetime TIMESTAMP NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMP,
  FOREIGN KEY(id_professional) REFERENCES professional(id_professional)
);
create index idx_blocker_delete_at on blocker (deleted_at);
create index idx_blocker_insert on blocker (inserted_at);
create index idx_blocker_update on blocker (updated_at);
create index idx_blocker_professional on blocker (id_professional);
create index idx_blocker_init_datetime on blocker (init_datetime);
create index idx_blocker_end_datetime on blocker (end_datetime);
create index idx_blocker_deleted on blocker (is_deleted);

create table availability (
  id_availability SERIAL PRIMARY KEY,
  inserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  id_professional INTEGER NOT NULL,
  init_datetime TIMESTAMP NOT NULL,
  end_datetime TIMESTAMP NOT NULL,
  init_hour TEXT NOT NULL,
  end_hour TEXT NOT NULL,
  type_availability INTEGER NOT NULL,
  weekday_name TEXT NOT NULL,
  interval INTEGER NOT NULL,
  resting INTEGER NOT NULL DEFAULT 0,
  priority_entry INTEGER NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMP,
  FOREIGN KEY(id_professional) REFERENCES professional(id_professional)
);
create index idx_availability_insert on availability (inserted_at);
create index idx_availability_update on availability (updated_at);
create index idx_availability_delete on availability (deleted_at);
create index idx_availability_professional on availability (id_professional);
create index idx_availability_type on availability (type_availability);
create index idx_availability_priority on availability (priority_entry);
create index idx_availability_deleted on availability (is_deleted);

create table slot (
  id_slot SERIAL PRIMARY KEY,
  inserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  id_availability INTEGER,
  id_professional INTEGER NOT NULL,
  slot TIMESTAMP NOT NULL,
  weekday_name TEXT NOT NULL,
  interval INTEGER NOT NULL,
  priority_entry INTEGER NOT NULL,
  status_entry TEXT NOT NULL,
  external_id TEXT,
  owner TEXT,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMP,
  id_blocker INTEGER,
  FOREIGN KEY(id_blocker) REFERENCES blocker(id_blocker),
  FOREIGN KEY(id_availability) REFERENCES availability(id_availability),
  FOREIGN KEY(id_professional) REFERENCES professional(id_professional),
  CHECK (status_entry IN ('open', 'busy', 'block', 'reserved'))
);
create index idx_slot_availability on slot (id_availability);
create index idx_slot_professional on slot (id_professional);
create index idx_slot_slot on slot (slot);
create index idx_slot_interval on slot (interval);
create index idx_slot_priority on slot (priority_entry);
create index idx_slot_external_id on slot (external_id);
create index idx_slot_deleted on slot (is_deleted);
create index idx_slot_blocker on slot (id_blocker);


