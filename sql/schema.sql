create table availability (
    id_availability INTEGER PRIMARY KEY,
    id_professional INTEGER NOT NULL,
    init_datetime TEXT NOT NULL,
    end_datetime TEXT NOT NULL,
    init_hour TEXT NOT NULL,
    end_hour TEXT NOT NULL,
    type_availability INTEGER NOT NULL,
    weekday_name TEXT NOT NULL,
    interval INTEGER NOT NULL,
    priority_entry INTEGER NOT NULL
); 

create table slot (
    id_slot INTEGER PRIMARY KEY,
    id_availability INTEGER NOT NULL,
    id_professional INTEGER NOT NULL,
    slot DATETIME NOT NULL,
    weekday_name TEXT NOT NULL,
    interval INTEGER NOT NULL,
    priority_entry INTEGER NOT NULL,
    status_entry TEXT NOT NULL,
    FOREIGN KEY(id_availability) REFERENCES availability(id_availability)
)