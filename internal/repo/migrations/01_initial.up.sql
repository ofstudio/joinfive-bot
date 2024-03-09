CREATE TABLE "updates"
(
    "id"                INTEGER PRIMARY KEY,
    "created_at"        TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    "chat_id"           INTEGER   NOT NULL,
    "chat_type"         TEXT      NOT NULL,
    "chat_title"        TEXT      NOT NULL DEFAULT '',
    "chat_username"     TEXT      NOT NULL DEFAULT '',
    "member_id"         INTEGER   NOT NULL,
    "member_first_name" TEXT      NOT NULL DEFAULT '',
    "member_last_name"  TEXT      NOT NULL DEFAULT '',
    "member_username"   TEXT      NOT NULL DEFAULT '',
    "member_is_bot"     BOOLEAN   NOT NULL DEFAULT FALSE,
    "status"            TEXT      NOT NULL
);

CREATE INDEX "updates__chat_id__idx"
    ON "updates" ("chat_id");
