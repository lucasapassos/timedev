create table availability (
    id_availability INTEGER PRIMARY KEY,
    id_professional INTEGER NOT NULL,
    init_datetime TEXT NOT NULL,
    end_datetime TEXT NOT NULL,
    init_hour TEXT NOT NULL,
    end_hour TEXT NOT NULL,
    type_availability INTEGER DEFAULT 1,
    weekday_name TEXT NOT NULL,
    interval INTEGER NOT NULL
); 