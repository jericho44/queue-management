-- 1. Billing Plans (SaaS Pricing Rates)
CREATE TABLE IF NOT EXISTS billing_plans (
    id BIGSERIAL PRIMARY KEY,
    plan_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    rate_per_ticket DECIMAL(12, 2) NOT NULL DEFAULT 500.00, -- IDR per ticket
    base_monthly_fee DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert Default Plan
INSERT INTO billing_plans (plan_code, name, rate_per_ticket, base_monthly_fee)
VALUES ('POSTPAID_STANDARD', 'Postpaid Standard Plan', 500.00, 0.00)
ON CONFLICT (plan_code) DO NOTHING;

-- 2. Billing Usage Meters (Monthly Ticket Accumulation per Tenant)
CREATE TABLE IF NOT EXISTS billing_usage_meters (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    billing_period VARCHAR(7) NOT NULL, -- Format: YYYY-MM
    ticket_count INT NOT NULL DEFAULT 0,
    last_ticket_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_org_billing_period UNIQUE (organization_id, billing_period)
);

CREATE INDEX IF NOT EXISTS idx_usage_meters_org_period ON billing_usage_meters(organization_id, billing_period);

-- 3. Invoices (Master Header Table)
CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(100) NOT NULL UNIQUE,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    billing_period VARCHAR(7) NOT NULL, -- YYYY-MM
    subtotal DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    tax_amount DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'UNPAID', -- UNPAID, PAID, OVERDUE, CANCELLED
    issued_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_invoices_uuid ON invoices(uuid);
CREATE INDEX IF NOT EXISTS idx_invoices_org ON invoices(organization_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);

-- 4. Invoice Items (Rincian Satuan Item/Layanan pada Invoice)
CREATE TABLE IF NOT EXISTS invoice_items (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    item_type VARCHAR(50) NOT NULL DEFAULT 'TICKET_USAGE', -- TICKET_USAGE, BASE_FEE, ADDON
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    total_price DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice ON invoice_items(invoice_id);

-- 5. Payments (Transaksi & Pembayaran Midtrans Gateway)
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    payment_number VARCHAR(100) NOT NULL UNIQUE,
    invoice_id BIGINT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    amount DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    payment_method VARCHAR(50), -- Qris, BankTransfer, CreditCard
    payment_channel VARCHAR(50), -- Gopay, BCA, Mandiri, Visa
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, SETTLEMENT, DENIED, EXPIRED, CANCELLED
    snap_token VARCHAR(255),
    snap_redirect_url TEXT,
    paid_at TIMESTAMPTZ,
    raw_gateway_response JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payments_uuid ON payments(uuid);
CREATE INDEX IF NOT EXISTS idx_payments_invoice ON payments(invoice_id);
CREATE INDEX IF NOT EXISTS idx_payments_org ON payments(organization_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

-- 6. Billing Webhook Logs (Midtrans Audit Trail)
CREATE TABLE IF NOT EXISTS billing_webhook_logs (
    id BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(100) NOT NULL,
    transaction_status VARCHAR(50) NOT NULL,
    raw_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_logs_order ON billing_webhook_logs(order_id);
