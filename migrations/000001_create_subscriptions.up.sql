CREATE TABLE subscriptions
(
    id           UUID PRIMARY KEY,
    service_name TEXT NOT NULL ,
    price        INTEGER NOT NULL ,
    user_id      UUID NOT NULL ,
    start_date   DATE NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    end_date DATE
);