-- Seed Organization
INSERT INTO organizations (id, uuid, name, code, slug, status, settings)
VALUES (1, '11111111-1111-1111-1111-111111111111', 'Demo Healthcare Org', 'DEMO', 'demo-healthcare', 'ACTIVE', '{}')
ON CONFLICT (code) DO NOTHING;

-- Seed Subscription
INSERT INTO subscriptions (organization_id, plan, max_branches, max_counters, max_staff, max_monthly_tickets, status)
VALUES (1, 'ENTERPRISE', 10, 50, 100, 100000, 'ACTIVE')
ON CONFLICT DO NOTHING;

-- Seed Users (Password: password123)
-- BCrypt for password123: $2a$10$3R39hcaVfcgotFc3widFyehagyNsRFC13kfaKkOCY3OyRInq462Dy
INSERT INTO users (id, uuid, organization_id, email, password_hash, full_name, role, status)
VALUES 
(0, '00000000-0000-0000-0000-000000000000', 1, 'superadmin@system.com', '$2a$10$3R39hcaVfcgotFc3widFyehagyNsRFC13kfaKkOCY3OyRInq462Dy', 'Super Admin Master', 'SUPER_ADMIN', 'ACTIVE'),
(1, '22222222-2222-2222-2222-222222222222', 1, 'owner@healthcare.com', '$2a$10$3R39hcaVfcgotFc3widFyehagyNsRFC13kfaKkOCY3OyRInq462Dy', 'Dr. Sarah Connor', 'OWNER', 'ACTIVE'),
(2, '33333333-3333-3333-3333-333333333333', 1, 'manager@healthcare.com', '$2a$10$3R39hcaVfcgotFc3widFyehagyNsRFC13kfaKkOCY3OyRInq462Dy', 'Alex Mercer', 'MANAGER', 'ACTIVE'),
(3, '44444444-4444-4444-4444-444444444444', 1, 'staff@healthcare.com', '$2a$10$3R39hcaVfcgotFc3widFyehagyNsRFC13kfaKkOCY3OyRInq462Dy', 'John Doe (Staff)', 'STAFF', 'ACTIVE'),
(4, '55555555-5555-5555-5555-555555555555', 1, 'receptionist@healthcare.com', '$2a$10$3R39hcaVfcgotFc3widFyehagyNsRFC13kfaKkOCY3OyRInq462Dy', 'Emily Watson', 'RECEPTIONIST', 'ACTIVE')
ON CONFLICT (email) DO NOTHING;

-- Seed Branch
INSERT INTO branches (id, uuid, organization_id, name, code, address, phone, status)
VALUES 
(1, '66666666-6666-6666-6666-666666666666', 1, 'Surabaya Central Branch', 'SUB01', 'Jl. Pemuda No. 45, Surabaya', '+62 31 555-0199', 'ACTIVE')
ON CONFLICT DO NOTHING;

-- Seed Services
INSERT INTO services (id, uuid, organization_id, branch_id, name, code, prefix, avg_service_time_sec, priority_weight, status)
VALUES 
(1, '77777777-7777-7777-7777-777777777777', 1, 1, 'General Consultation', 'GEN_CONSULT', 'A', 480, 1, 'ACTIVE'),
(2, '88888888-8888-8888-8888-888888888888', 1, 1, 'Payment & Cashier', 'PAYMENT', 'B', 300, 2, 'ACTIVE'),
(3, '99999999-9999-9999-9999-999999999999', 1, 1, 'Pharmacy & Medicines', 'PHARMACY', 'C', 420, 1, 'ACTIVE')
ON CONFLICT DO NOTHING;

-- Seed Counters
INSERT INTO counters (id, uuid, organization_id, branch_id, counter_number, name, status)
VALUES 
(1, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, 1, '01', 'Counter 01 — General', 'OPEN'),
(2, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 1, 1, '02', 'Counter 02 — Cashier', 'OPEN'),
(3, 'cccccccc-cccc-cccc-cccc-cccccccccccc', 1, 1, '03', 'Counter 03 — Pharmacy', 'CLOSED')
ON CONFLICT DO NOTHING;

-- Counter Services Pivot
INSERT INTO counter_services (counter_id, service_id)
VALUES (1, 1), (1, 2), (2, 2), (3, 3)
ON CONFLICT DO NOTHING;
