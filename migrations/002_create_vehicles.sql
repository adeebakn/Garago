CREATE TABLE vehicles(
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    reg_no VARCHAR(20) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_vehicle_customer FOREIGN KEY (customer_id) REFERENCES customers (id)

)