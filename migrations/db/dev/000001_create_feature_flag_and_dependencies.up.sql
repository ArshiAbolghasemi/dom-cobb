CREATE TABLE feature_flags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    "name" VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX idx_feature_flags_name
    ON feature_flags (name);

CREATE INDEX idx_feature_flags_deleted_at
    ON feature_flags (deleted_at);


CREATE TABLE flag_dependencies (
    flag_id BIGINT NOT NULL,
    depends_on_flag_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (flag_id, depends_on_flag_id),

    CONSTRAINT fk_flag_dependencies_flag_id
        FOREIGN KEY (flag_id)
        REFERENCES feature_flags (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_flag_dependencies_depends_on_flag_id
        FOREIGN KEY (depends_on_flag_id)
        REFERENCES feature_flags (id)
        ON DELETE CASCADE
);
