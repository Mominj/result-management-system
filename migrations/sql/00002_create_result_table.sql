-- +goose Up
-- SQL in this section is executed when the migration is applied.


CREATE TABLE result
(
    id SERIAL PRIMARY KEY,
    year VARCHAR(20) NOT NULL,
    exam_type VARCHAR(20) NOT NULL,
    groups VARCHAR(20) NOT NULL,
    roll_number VARCHAR(20) NOT NULL,
    name VARCHAR(20) NOT NULL,
    bangla BIGINT,
    english  BIGINT,
    ict BIGINT,
    physics BIGINT,
    chemistry BIGINT,
    biology BIGINT,
    higher_mathematics BIGINT,
    agriculture_education BIGINT,
    geography BIGINT,
    psychology BIGINT,
    statistics BIGINT,
    accounting BIGINT,
    economics BIGINT,
    business_organization_and_management BIGINT,
    finance_banking_and_insurance BIGINT,
    history BIGINT,
    islamic_history_and_culture BIGINT,
    sociology BIGINT,
    logic BIGINT
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE result;