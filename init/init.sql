CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE team (
                      team_name TEXT PRIMARY KEY
);

CREATE TABLE "user" (
                        user_id   TEXT PRIMARY KEY,
                        username  TEXT NOT NULL,
                        team_name TEXT NOT NULL REFERENCES team(team_name) ON DELETE CASCADE,
                        is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE pull_request (
                              pull_request_id   TEXT PRIMARY KEY,
                              pull_request_name TEXT NOT NULL,
                              author_id         TEXT NOT NULL REFERENCES "user"(user_id),
                              status            pr_status NOT NULL DEFAULT 'OPEN',
                              assigned_reviewers TEXT[] NOT NULL DEFAULT '{}',
                              created_at        TIMESTAMPTZ DEFAULT NOW(),
                              merged_at         TIMESTAMPTZ
);


CREATE INDEX idx_user_team_name ON "user"(team_name);
CREATE INDEX idx_pr_author_id ON pull_request(author_id);
CREATE INDEX idx_pr_status ON pull_request(status);
CREATE INDEX idx_pr_reviewers ON pull_request USING GIN (assigned_reviewers);