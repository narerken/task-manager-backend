create table attendance_sessions (
                                     id serial primary key,
                                     department_id int not null,
                                     date date not null,
                                     state text not null default 'draft'
                                         check (state in ('draft', 'published')),
                                     created_by int,
                                     updated_by int,
                                     created_at timestamp default now(),
                                     updated_at timestamp default now(),
                                     deleted_at timestamp,
                                     constraint attendance_sessions_department_date_unique unique (department_id, date)
);

create table attendance_entries (
                                    id serial primary key,
                                    session_id int not null references attendance_sessions(id) on delete cascade,
                                    student_id int not null,
                                    status text not null
                                        check (status in ('present', 'absent', 'late', 'excused')),
                                    comment text,
                                    marked_by int,
                                    created_at timestamp default now(),
                                    updated_at timestamp default now(),
                                    deleted_at timestamp,
                                    constraint attendance_entries_session_student_unique unique (session_id, student_id)
);

create index idx_attendance_sessions_department_date on attendance_sessions (department_id, date);
create index idx_attendance_entries_session on attendance_entries (session_id);
