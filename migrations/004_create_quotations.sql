CREATE TABLE quotations (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    vehicle_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Draft',
    total_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_quotation_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id),

    CONSTRAINT fk_quotation_vehicle
        FOREIGN KEY (vehicle_id)
        REFERENCES vehicles(id)
);