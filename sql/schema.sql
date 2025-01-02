create table availability (
    id_availability INTEGER PRIMARY KEY,
    id_professional INTEGER NOT NULL,
    init_datetime DATETIME NOT NULL,
    end_datetime DATETIME NOT NULL,
    init_hour TEXT NOT NULL,
    end_hour TEXT NOT NULL,
    type_availability INTEGER NOT NULL,
    weekday_name TEXT NOT NULL,
    interval INTEGER NOT NULL,
    resting INTEGER NOT NULL DEFAULT 0,
    priority_entry INTEGER NOT NULL,
    is_deleted INTEGER NOT NULL,
    FOREIGN KEY(id_professional) REFERENCES professional(id_professional)
  );

create table slot (
    id_slot INTEGER PRIMARY KEY,
    inserted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id_availability INTEGER,
    id_professional INTEGER NOT NULL,
    slot DATETIME NOT NULL,
    weekday_name TEXT NOT NULL,
    interval INTEGER NOT NULL,
    priority_entry INTEGER NOT NULL,
    status_entry TEXT NOT NULL,
    external_id TEXT,
    owner TEXT,
    is_deleted INTEGER NOT NULL DEFAULT 0,
    deleted_at DATETIME,
    id_blocker INTEGER,
    FOREIGN KEY(id_blocker) REFERENCES blocker(id_blocker),
    FOREIGN KEY(id_availability) REFERENCES availability(id_availability),
    FOREIGN KEY(id_professional) REFERENCES professional(id_professional),
    CHECK (status_entry IN ('open', 'busy', 'block', 'reserved'))
);

create table blocker (
  id_blocker INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  id_professional INTEGER NOT NULL,
  init_datetime DATETIME NOT NULL,
  end_datetime DATETIME NOT NULL,
  is_deleted INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(id_professional) REFERENCES professional(id_professional)
);

create table professional (
  id_professional INTEGER PRIMARY KEY,
  reference_key TEXT NOT NULL UNIQUE,
  especialidade TEXT NOT NULL,
  nome TEXT NOT NULL
);

create table attribute (
  id_attribute INTEGER PRIMARY KEY,
  id_professional INTEGER NOT NULL,
  attribute TEXT NOT NULL,
  value TEXT NOT NULL,
  FOREIGN KEY(id_professional) REFERENCES professional(id_professional),
  UNIQUE (id_professional, attribute, value)
);
