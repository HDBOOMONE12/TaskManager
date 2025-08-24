CREATE TABLE IF NOT EXISTS tasks
(
    id          bigint generated always as identity primary key,
    user_id     bigint      not null references users (id) on delete cascade,
    title       varchar(50) not null,
    description varchar(50),
    status      varchar(50) not null default 'todo',

    due_date    timestamptz,
    priority    int         not null default 0,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),


    check (status in ('todo', 'doing', 'done'))

)