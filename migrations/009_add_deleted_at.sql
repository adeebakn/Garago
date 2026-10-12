-- Customers
ALTER TABLE customers
ADD COLUMN deleted_at TIMESTAMP NULL;

-- Vehicles
ALTER TABLE vehicles
ADD COLUMN deleted_at TIMESTAMP NULL;

-- Services
ALTER TABLE services
ADD COLUMN deleted_at TIMESTAMP NULL;