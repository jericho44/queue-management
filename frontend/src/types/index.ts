export interface Meta {
  pagination?: {
    current_page: number;
    last_page: number;
    per_page: number;
    total: number;
    from: number;
    to: number;
    has_next: boolean;
    has_previous: boolean;
  };
  request_id: string;
  timestamp: string;
}

export interface APIResponse<T> {
  success: boolean;
  message: string;
  data: T;
  errors?: any;
  meta: Meta;
}

export interface User {
  id: number;
  uuid: string;
  organization_id?: number;
  org_uuid?: string;
  org_name?: string;
  email: string;
  full_name: string;
  role: 'SUPER_ADMIN' | 'OWNER' | 'MANAGER' | 'STAFF' | 'RECEPTIONIST';
  phone?: string;
  status: string;
  created_at: string;
}

export interface Organization {
  id: number;
  uuid: string;
  name: string;
  code: string;
  slug: string;
  status: string;
  settings?: any;
  created_at: string;
}

export interface Branch {
  id: number;
  uuid: string;
  organization_id: number;
  name: string;
  code: string;
  address?: string;
  phone?: string;
  status: string;
  kiosk_enabled?: boolean;
  kiosk_mode?: 'DUAL' | 'PAPERLESS' | 'PHYSICAL';
  paper_size?: '58mm' | '80mm';
  receipt_header?: string;
  receipt_footer?: string;
  auto_print?: boolean;
}

export interface Service {
  id: number;
  uuid: string;
  organization_id: number;
  branch_id: number;
  branch_name?: string;
  name: string;
  code: string;
  prefix: string;
  avg_service_time_sec: number;
  priority_weight: number;
  status: string;
}

export interface Counter {
  id: number;
  uuid: string;
  organization_id: number;
  branch_id: number;
  branch_name?: string;
  counter_number: string;
  name: string;
  status: 'CLOSED' | 'OPEN' | 'BUSY' | 'PAUSED';
  current_staff_id?: number;
  staff_name?: string;
  services?: Service[];
}

export interface QueueTicket {
  id: number;
  uuid: string;
  organization_id: number;
  branch_id: number;
  branch_name?: string;
  service_id: number;
  service_name?: string;
  service_prefix?: string;
  customer_id?: number;
  customer_name?: string;
  counter_id?: number;
  counter_number?: string;
  staff_id?: number;
  staff_name?: string;
  ticket_number: string;
  sequence_number: number;
  queue_date: string;
  priority: 'NORMAL' | 'PRIORITY' | 'EMERGENCY';
  status: 'WAITING' | 'CALLED' | 'SERVING' | 'COMPLETED' | 'SKIPPED' | 'CANCELLED' | 'NO_SHOW' | 'TRANSFERRED';
  public_token: string;
  called_at?: string;
  serving_started_at?: string;
  completed_at?: string;
  estimated_wait_seconds: number;
  people_ahead?: number;
  created_at: string;
}

export interface DashboardStats {
  total_tickets_today: number;
  completed_today: number;
  waiting_count: number;
  serving_count: number;
  avg_wait_time_sec: number;
  avg_service_time_sec: number;
  no_show_count: number;
  cancelled_count: number;
  active_counters: number;
  hourly_distribution: { hour: number; count: number }[];
  service_distribution: { service_id: number; service_name: string; count: number }[];
}

export interface AuditLog {
  id: number;
  uuid: string;
  organization_id?: number;
  branch_id?: number;
  user_id?: number;
  user_name?: string;
  action: string;
  entity_type: string;
  entity_id?: number;
  old_values?: string;
  new_values?: string;
  ip_address?: string;
  created_at: string;
}

export interface UsageMeter {
  id: number;
  organization_id: number;
  billing_period: string;
  ticket_count: number;
  last_ticket_at?: string;
}

export interface InvoiceItem {
  id: number;
  invoice_id: number;
  description: string;
  item_type: 'TICKET_USAGE' | 'BASE_FEE' | 'ADDON';
  quantity: number;
  unit_price: number;
  total_price: number;
  created_at: string;
}

export interface Payment {
  id: number;
  uuid: string;
  payment_number: string;
  invoice_id: number;
  organization_id: number;
  amount: number;
  payment_method?: string;
  payment_channel?: string;
  status: 'PENDING' | 'SETTLEMENT' | 'DENIED' | 'EXPIRED' | 'CANCELLED';
  snap_token?: string;
  snap_redirect_url?: string;
  paid_at?: string;
  created_at: string;
}

export interface Invoice {
  id: number;
  uuid: string;
  invoice_number: string;
  organization_id: number;
  org_name?: string;
  billing_period: string;
  subtotal: number;
  tax_amount: number;
  total_amount: number;
  status: 'UNPAID' | 'PAID' | 'OVERDUE' | 'CANCELLED';
  issued_at: string;
  due_date: string;
  items?: InvoiceItem[];
  latest_payment?: Payment;
}

export interface SuperadminBillingStats {
  total_revenue: number;
  pending_revenue: number;
  paid_invoices_count: number;
  unpaid_invoices_count: number;
  overdue_invoices_count: number;
  current_month_tickets: number;
}
