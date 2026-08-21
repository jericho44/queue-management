-- Revert Kiosk & Thermal Ticket Printer settings from branches table
ALTER TABLE branches
DROP COLUMN IF EXISTS auto_print,
DROP COLUMN IF EXISTS receipt_footer,
DROP COLUMN IF EXISTS receipt_header,
DROP COLUMN IF EXISTS paper_size,
DROP COLUMN IF EXISTS kiosk_mode,
DROP COLUMN IF EXISTS kiosk_enabled;
