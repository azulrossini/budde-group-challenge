-- +goose Up
INSERT INTO bills (id, description, total_cents) VALUES (1, 'Team dinner', 12050);
INSERT INTO shares (bill_id, person_name, percentage_bp) VALUES
  (1, 'Alice', 5000), (1, 'Bob', 3000), (1, 'Carol', 2000);
SELECT setval('bills_id_seq', (SELECT max(id) FROM bills));

-- +goose Down
DELETE FROM shares WHERE bill_id = 1;
DELETE FROM bills WHERE id = 1;
