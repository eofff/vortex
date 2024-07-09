CREATE TABLE AlgorithmStatus (
	id bigserial PRIMARY KEY,
	client_id bigint,
	VWAP boolean,
	TWAP boolean,
	HFT boolean
);