-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    serviceName TEXT,
	price INT,
	userID UUID,
	startDate DATE,
	endDate DATE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subscriptions;
-- +goose StatementEnd
