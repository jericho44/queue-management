-- Add Kiosk & Thermal Ticket Printer settings to branches table
ALTER TABLE branches
ADD COLUMN IF NOT EXISTS kiosk_enabled BOOLEAN NOT NULL DEFAULT true,
ADD COLUMN IF NOT EXISTS kiosk_mode VARCHAR(50) NOT NULL DEFAULT 'DUAL', -- DUAL, PAPERLESS, PHYSICAL
ADD COLUMN IF NOT EXISTS paper_size VARCHAR(20) NOT NULL DEFAULT '58mm', -- 58mm, 80mm
ADD COLUMN IF NOT EXISTS receipt_header TEXT NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS receipt_footer TEXT NOT NULL DEFAULT 'Terima kasih atas kunjungan Anda. Harap menunggu hingga nomor Anda dipanggil.',
ADD COLUMN IF NOT EXISTS auto_print BOOLEAN NOT NULL DEFAULT false;
