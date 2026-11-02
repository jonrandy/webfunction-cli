<?php

declare(strict_types=1);

namespace Reservepay;

use WebFunction\Page;
use WebFunction\Package;
use WebFunction\Pipeline;
use WebFunction\BadRequestError;
use WebFunction\UnexpectedStatusCodeError;
use WebFunction\JsonParseError;

// IMPORTANT - composer require webfunction-protocol/webfunction

/**
 * Generated client for the Reservepay Merchant API webfunction package.
 * Source: https://api.reservepay.com/merchants
 *
 * @phpstan-type PaymentMethodAttributes = array{name: string, display_name: string, flow: string, category: string, platforms: 'WEB'|'IOS'|'ANDROID'|'TERMINAL'}
 * @psalm-type PaymentMethodAttributes = array{name: string, display_name: string, flow: string, category: string, platforms: 'WEB'|'IOS'|'ANDROID'|'TERMINAL'}
 * @phpstan-type InstallationAttributes = array{installation_id: string, enabled: bool, name: string, platform: 'WEB'|'IOS'|'ANDROID'|'TERMINAL', identifier?: string|null, payment_methods: list<PaymentMethodAttributes>, email_capture: 'disabled'|'optional'|'required', mobile_capture: 'disabled'|'optional'|'required', address_capture: 'disabled'|'optional'|'required', created_at: int|float}
 * @psalm-type InstallationAttributes = array{installation_id: string, enabled: bool, name: string, platform: 'WEB'|'IOS'|'ANDROID'|'TERMINAL', identifier?: string|null, payment_methods: list<PaymentMethodAttributes>, email_capture: 'disabled'|'optional'|'required', mobile_capture: 'disabled'|'optional'|'required', address_capture: 'disabled'|'optional'|'required', created_at: int|float}
 * @phpstan-type EntryAttributes = array{entry_id: string, transaction_id: string, direction: 'CREDIT'|'DEBIT', type: string, amount: int|float, currency: string, status: 'ON_HOLD'|'EFFECTIVE'|'CANCELED', invoice_id?: string|null, statement_id?: string|null, payout_id?: string|null, topup_id?: string|null, dispute_id?: string|null, parent_id?: string|null, parent_type?: 'PAYMENT'|null, created_at: int|float, effective_at?: int|float|null, canceled_at?: int|float|null}
 * @psalm-type EntryAttributes = array{entry_id: string, transaction_id: string, direction: 'CREDIT'|'DEBIT', type: string, amount: int|float, currency: string, status: 'ON_HOLD'|'EFFECTIVE'|'CANCELED', invoice_id?: string|null, statement_id?: string|null, payout_id?: string|null, topup_id?: string|null, dispute_id?: string|null, parent_id?: string|null, parent_type?: 'PAYMENT'|null, created_at: int|float, effective_at?: int|float|null, canceled_at?: int|float|null}
 * @phpstan-type BankAccountAttributes = array{bank_account_id: string, bank_code: string, number: string, name: string, country_code: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool, owner: 'CONTACT'|'MERCHANT', primary: bool}
 * @psalm-type BankAccountAttributes = array{bank_account_id: string, bank_code: string, number: string, name: string, country_code: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool, owner: 'CONTACT'|'MERCHANT', primary: bool}
 * @phpstan-type CardAttributes = array{card_id: string, token: string, usage: 'SINGLE'|'MULTIPLE', used: bool, first_digits: string, last_digits: string, length: int|float, expiration_date: string, cardholder_name: string, scheme: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool}
 * @psalm-type CardAttributes = array{card_id: string, token: string, usage: 'SINGLE'|'MULTIPLE', used: bool, first_digits: string, last_digits: string, length: int|float, expiration_date: string, cardholder_name: string, scheme: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool}
 * @phpstan-type ContactAttributes = array{contact_id: string, external_id?: string|null, email?: string|null, mobile?: string|null, name?: string|null, favorite: bool, created_at: int|float, primary_bank_account?: BankAccountAttributes|null, default_card?: CardAttributes|null, bank_accounts: list<BankAccountAttributes>, cards: list<CardAttributes>}
 * @psalm-type ContactAttributes = array{contact_id: string, external_id?: string|null, email?: string|null, mobile?: string|null, name?: string|null, favorite: bool, created_at: int|float, primary_bank_account?: BankAccountAttributes|null, default_card?: CardAttributes|null, bank_accounts: list<BankAccountAttributes>, cards: list<CardAttributes>}
 * @phpstan-type PaymentLinkAttributes = array{payment_link_id: string, batch_id?: string|null, payment_link_url: string, status: 'active'|'disabled'|'expired'|'paid', amount: int|float, currency: string, interest_bearer?: 'customer'|'merchant'|null, description: string, external_id?: string|null, expires_at?: int|float|null, single_use: bool, capture: bool, send_confirmation_email: bool, installation: InstallationAttributes, created_at: int|float, deletable: bool}
 * @psalm-type PaymentLinkAttributes = array{payment_link_id: string, batch_id?: string|null, payment_link_url: string, status: 'active'|'disabled'|'expired'|'paid', amount: int|float, currency: string, interest_bearer?: 'customer'|'merchant'|null, description: string, external_id?: string|null, expires_at?: int|float|null, single_use: bool, capture: bool, send_confirmation_email: bool, installation: InstallationAttributes, created_at: int|float, deletable: bool}
 * @phpstan-type FeeAttributes = array{fee_id: string, kind: 'bank_account_validation'|'communication'|'dispute'|'installment'|'metadata'|'payment'|'payout'|'plan'|'refund'|'virtual_account', parent_id: string, billing_code: string, unit_amount: int|float, percentage_amount: int|float, discount_applied: true|false, statement_id: string}
 * @psalm-type FeeAttributes = array{fee_id: string, kind: 'bank_account_validation'|'communication'|'dispute'|'installment'|'metadata'|'payment'|'payout'|'plan'|'refund'|'virtual_account', parent_id: string, billing_code: string, unit_amount: int|float, percentage_amount: int|float, discount_applied: true|false, statement_id: string}
 * @phpstan-type PaymentAttributes = array{payment_id: string, external_id?: string|null, status: 'PENDING'|'AUTHORIZED'|'SUCCESSFUL'|'FAILED'|'REVERSED'|'VOIDED'|'CANCELLED'|'EXPIRED'|'REFUNDED'|'PARTIALLY_REFUNDED'|'DISPUTED', amount: int|float, currency: string, payment_method: string, payment_method_category: string, payment_method_display_name: string, card_first_digits?: string|null, card_last_digits?: string|null, card_masked_number?: string|null, card_network?: string|null, cardholder_name?: string|null, created_at: int|float, authorized_at?: int|float|null, captured_at?: int|float|null, voided_at?: int|float|null, cancelled_at?: int|float|null, refundable_at?: int|float|null, expires_at?: int|float|null, capture: 'AUTOMATIC'|'MANUAL', customer_id?: string|null, customer_external_id?: string|null, customer_email?: string|null, customer_mobile?: string|null, customer_name?: string|null, fees: list<FeeAttributes>, ref1?: string|null, ref2?: string|null, ref3?: string|null, reversible: bool, capturable: bool, voidable: bool, cancellable: bool, refundable: bool, fraud_finding?: 'risky'|'safe'|null, failure_code?: string|null, failure_message?: string|null, terminal_id?: string|null, payment_link_id?: string|null, payment_session_id?: string|null, installment_term?: int|float|null, interest_bearer?: 'customer'|'merchant'|null, installment_interest_rate?: string|null, installment_interest_amount?: int|float|null, installment_bank_name?: string|null, installment_bank_display_name?: string|null}
 * @psalm-type PaymentAttributes = array{payment_id: string, external_id?: string|null, status: 'PENDING'|'AUTHORIZED'|'SUCCESSFUL'|'FAILED'|'REVERSED'|'VOIDED'|'CANCELLED'|'EXPIRED'|'REFUNDED'|'PARTIALLY_REFUNDED'|'DISPUTED', amount: int|float, currency: string, payment_method: string, payment_method_category: string, payment_method_display_name: string, card_first_digits?: string|null, card_last_digits?: string|null, card_masked_number?: string|null, card_network?: string|null, cardholder_name?: string|null, created_at: int|float, authorized_at?: int|float|null, captured_at?: int|float|null, voided_at?: int|float|null, cancelled_at?: int|float|null, refundable_at?: int|float|null, expires_at?: int|float|null, capture: 'AUTOMATIC'|'MANUAL', customer_id?: string|null, customer_external_id?: string|null, customer_email?: string|null, customer_mobile?: string|null, customer_name?: string|null, fees: list<FeeAttributes>, ref1?: string|null, ref2?: string|null, ref3?: string|null, reversible: bool, capturable: bool, voidable: bool, cancellable: bool, refundable: bool, fraud_finding?: 'risky'|'safe'|null, failure_code?: string|null, failure_message?: string|null, terminal_id?: string|null, payment_link_id?: string|null, payment_session_id?: string|null, installment_term?: int|float|null, interest_bearer?: 'customer'|'merchant'|null, installment_interest_rate?: string|null, installment_interest_amount?: int|float|null, installment_bank_name?: string|null, installment_bank_display_name?: string|null}
 * @phpstan-type PayoutAttributes = array{payout_id: string, external_id?: string|null, status: 'SUCCESSFUL'|'SENT'|'PENDING'|'FAILED'|'CANCELLED'|'REQUESTED', amount: int|float, currency: string, created_at: int|float, paid_at?: int|float|null, sent_at?: int|float|null, bank_account_id?: string|null, bank_code?: string|null, bank_account_number?: string|null, bank_account_name?: string|null, bank_country_code?: string|null, payout_method?: string|null, mobile_number?: string|null, national_id?: string|null, fees: list<FeeAttributes>, ref1?: string|null, needs_authorization: true|false, requested_by_user_id?: string|null, requested_by_user_email?: string|null, authorized_by_user_id?: string|null, authorized_by_user_email?: string|null, rejected_by_user_id?: string|null, rejected_by_user_email?: string|null}
 * @psalm-type PayoutAttributes = array{payout_id: string, external_id?: string|null, status: 'SUCCESSFUL'|'SENT'|'PENDING'|'FAILED'|'CANCELLED'|'REQUESTED', amount: int|float, currency: string, created_at: int|float, paid_at?: int|float|null, sent_at?: int|float|null, bank_account_id?: string|null, bank_code?: string|null, bank_account_number?: string|null, bank_account_name?: string|null, bank_country_code?: string|null, payout_method?: string|null, mobile_number?: string|null, national_id?: string|null, fees: list<FeeAttributes>, ref1?: string|null, needs_authorization: true|false, requested_by_user_id?: string|null, requested_by_user_email?: string|null, authorized_by_user_id?: string|null, authorized_by_user_email?: string|null, rejected_by_user_id?: string|null, rejected_by_user_email?: string|null}
 * @phpstan-type TerminalAttributes = array{terminal_id: string, connected: true|false, external_id?: string|null, installation_id?: string|null, name: string, branch_name?: string|null, enabled: true|false, kiosk_mode: true|false, last_verified_at?: int|float|null, connected_at?: int|float|null, delegated_terminal_id?: string|null, parent_terminal_id?: string|null, challenge: string}
 * @psalm-type TerminalAttributes = array{terminal_id: string, connected: true|false, external_id?: string|null, installation_id?: string|null, name: string, branch_name?: string|null, enabled: true|false, kiosk_mode: true|false, last_verified_at?: int|float|null, connected_at?: int|float|null, delegated_terminal_id?: string|null, parent_terminal_id?: string|null, challenge: string}
 * @phpstan-type TopupAttributes = array{topup_id: string, status: 'PENDING'|'SUCCESSFUL', type: 'MANUAL'|'VIRTUAL_ACCOUNT', amount: int|float, currency: string, transaction_id?: string|null, sender_account_name?: string|null, sender_account_number?: string|null, sender_bank_code?: string|null, virtual_account_number?: string|null, created_at: int|float}
 * @psalm-type TopupAttributes = array{topup_id: string, status: 'PENDING'|'SUCCESSFUL', type: 'MANUAL'|'VIRTUAL_ACCOUNT', amount: int|float, currency: string, transaction_id?: string|null, sender_account_name?: string|null, sender_account_number?: string|null, sender_bank_code?: string|null, virtual_account_number?: string|null, created_at: int|float}
 * @phpstan-type DisputeEvidenceAttributes = array{evidence_id: string, stage: 'CHALLENGE'|'ESCALATION', comment?: string|null, filename?: string|null, content_type?: string|null, byte_size?: int|float|null, submitted: true|false, upload_url?: string|null, download_url?: string|null, created_at: int|float}
 * @psalm-type DisputeEvidenceAttributes = array{evidence_id: string, stage: 'CHALLENGE'|'ESCALATION', comment?: string|null, filename?: string|null, content_type?: string|null, byte_size?: int|float|null, submitted: true|false, upload_url?: string|null, download_url?: string|null, created_at: int|float}
 * @phpstan-type ActivityLogAttributes = array{created_at: int|float, activity_type: 'COMMUNICATION_CHANNEL_ACTIVATED'|'DISPUTE_CHALLENGED'|'DISPUTE_CONCEDED'|'DISPUTE_ESCALATED'|'INVITATION_ACCEPTED'|'PAYMENT_CANCELLED'|'PAYMENT_CAPTURED'|'PAYMENT_REVERSED'|'PAYMENT_VOIDED'|'PAYOUT_AUTHORIZED'|'PAYOUT_CREATED'|'PAYOUT_REJECTED'|'PAYOUT_REQUESTED'|'PLAN_SWITCHED'|'REFUND_CREATED'|'USER_INVITED'|'USER_REMOVED', actor_type: 'USER'|'API_KEY'|'SYSTEM', actor_email?: string|null, actor_key_id?: string|null, parent_type?: 'PAYMENT'|null, parent_id?: string|null, object_id?: string|null, object_type?: 'PAYMENT'|'REFUND'|'DISPUTE'|'PAYOUT'|null, user_email?: string|null, new_plan_id?: string|null, new_plan_name?: string|null, previous_plan_id?: string|null, previous_plan_name?: string|null}
 * @psalm-type ActivityLogAttributes = array{created_at: int|float, activity_type: 'COMMUNICATION_CHANNEL_ACTIVATED'|'DISPUTE_CHALLENGED'|'DISPUTE_CONCEDED'|'DISPUTE_ESCALATED'|'INVITATION_ACCEPTED'|'PAYMENT_CANCELLED'|'PAYMENT_CAPTURED'|'PAYMENT_REVERSED'|'PAYMENT_VOIDED'|'PAYOUT_AUTHORIZED'|'PAYOUT_CREATED'|'PAYOUT_REJECTED'|'PAYOUT_REQUESTED'|'PLAN_SWITCHED'|'REFUND_CREATED'|'USER_INVITED'|'USER_REMOVED', actor_type: 'USER'|'API_KEY'|'SYSTEM', actor_email?: string|null, actor_key_id?: string|null, parent_type?: 'PAYMENT'|null, parent_id?: string|null, object_id?: string|null, object_type?: 'PAYMENT'|'REFUND'|'DISPUTE'|'PAYOUT'|null, user_email?: string|null, new_plan_id?: string|null, new_plan_name?: string|null, previous_plan_id?: string|null, previous_plan_name?: string|null}
 * @phpstan-type ApiRequestLogAttributes = array{created_at: int|float, status: int|float, endpoint: string, ip_address?: string|null, source: 'API'|'DASHBOARD', user_email?: string|null, user_agent?: string|null}
 * @psalm-type ApiRequestLogAttributes = array{created_at: int|float, status: int|float, endpoint: string, ip_address?: string|null, source: 'API'|'DASHBOARD', user_email?: string|null, user_agent?: string|null}
 * @phpstan-type BatchAttributes = array{batch_id: string, type: 'PAYMENT_LINK'|'PAYOUT', status: 'CREATED'|'PREPARING'|'PREPARED'|'PROCESSING'|'COMPLETED', output_url?: string|null, preview_url?: string|null, error_message?: string|null, created_at: int|float}
 * @psalm-type BatchAttributes = array{batch_id: string, type: 'PAYMENT_LINK'|'PAYOUT', status: 'CREATED'|'PREPARING'|'PREPARED'|'PROCESSING'|'COMPLETED', output_url?: string|null, preview_url?: string|null, error_message?: string|null, created_at: int|float}
 * @phpstan-type CommunicationChannelAttributes = array{communication_channel_id: string, delivery: 'email'|'sms', email?: string|null, sender_name?: string|null, active: true|false, status: 'active'|'pending_setup'|'inactive'}
 * @psalm-type CommunicationChannelAttributes = array{communication_channel_id: string, delivery: 'email'|'sms', email?: string|null, sender_name?: string|null, active: true|false, status: 'active'|'pending_setup'|'inactive'}
 * @phpstan-type DisputeAttributes = array{dispute_id: string, payment_id: string, contact_id?: string|null, message?: string|null, status: 'OPEN'|'CHALLENGED'|'REOPENED'|'ESCALATED'|'WON'|'LOST'|'EXPIRED', challenge_received_at?: int|float|null, escalation_received_at?: int|float|null, won_at?: int|float|null, lost_at?: int|float|null, amount: int|float, currency: string, reason: string, suggested_evidence: 'cancellation_policy'|'cancellation_request_records'|'common_evidence'|'customer_communication'|'device_id_ip_data'|'marketing_materials'|'payment_logs'|'payment_processing_records'|'receipt'|'refund_issuance'|'refund_policy'|'refund_records'|'return_policy'|'shipping_documentation'|'supplier_certification'|'technical_support_logs'|'transaction_history'|'transaction_logs', deadline_at?: int|float|null, created_at: int|float, fees: list<FeeAttributes>, conceded: true|false, evidences: list<DisputeEvidenceAttributes>}
 * @psalm-type DisputeAttributes = array{dispute_id: string, payment_id: string, contact_id?: string|null, message?: string|null, status: 'OPEN'|'CHALLENGED'|'REOPENED'|'ESCALATED'|'WON'|'LOST'|'EXPIRED', challenge_received_at?: int|float|null, escalation_received_at?: int|float|null, won_at?: int|float|null, lost_at?: int|float|null, amount: int|float, currency: string, reason: string, suggested_evidence: 'cancellation_policy'|'cancellation_request_records'|'common_evidence'|'customer_communication'|'device_id_ip_data'|'marketing_materials'|'payment_logs'|'payment_processing_records'|'receipt'|'refund_issuance'|'refund_policy'|'refund_records'|'return_policy'|'shipping_documentation'|'supplier_certification'|'technical_support_logs'|'transaction_history'|'transaction_logs', deadline_at?: int|float|null, created_at: int|float, fees: list<FeeAttributes>, conceded: true|false, evidences: list<DisputeEvidenceAttributes>}
 * @phpstan-type NotificationTopicAttributes = array{id: string, name: 'api_key'|'billing'|'dispute'|'fraud'|'incident'|'installation_key'|'payout'|'plan'|'planned_maintenance'|'team'|'webhook', display_name: string, category?: 'api_key'|'billing'|'dispute'|'fraud'|'incident'|'installation_key'|'payout'|'plan'|'planned_maintenance'|'team'|'webhook'|null}
 * @psalm-type NotificationTopicAttributes = array{id: string, name: 'api_key'|'billing'|'dispute'|'fraud'|'incident'|'installation_key'|'payout'|'plan'|'planned_maintenance'|'team'|'webhook', display_name: string, category?: 'api_key'|'billing'|'dispute'|'fraud'|'incident'|'installation_key'|'payout'|'plan'|'planned_maintenance'|'team'|'webhook'|null}
 * @phpstan-type EmailRecipientAttributes = array{email: string, language: 'en'|'th', notification_topics: list<NotificationTopicAttributes>}
 * @psalm-type EmailRecipientAttributes = array{email: string, language: 'en'|'th', notification_topics: list<NotificationTopicAttributes>}
 * @phpstan-type InvoiceAttributes = array{invoice_id: string, number: string, status: 'due'|'settled', amount: int|float, net_amount: int|float, vat_amount: int|float, currency: string, issued_on: string, pdf_download_url?: string|null, receipt_id?: string|null, receipt_number?: string|null, receipt_pdf_download_url?: string|null, receipt_xml_download_url?: string|null, receipt_created_at?: int|float|null, created_at: int|float}
 * @psalm-type InvoiceAttributes = array{invoice_id: string, number: string, status: 'due'|'settled', amount: int|float, net_amount: int|float, vat_amount: int|float, currency: string, issued_on: string, pdf_download_url?: string|null, receipt_id?: string|null, receipt_number?: string|null, receipt_pdf_download_url?: string|null, receipt_xml_download_url?: string|null, receipt_created_at?: int|float|null, created_at: int|float}
 * @phpstan-type KeyAttributes = array{key_id: string, key: string, name: string, revoked: bool, created_at: int|float, permissions: 'authorize_payout'|'manage_contacts'|'manage_developer_settings'|'manage_disputes'|'manage_payouts'|'manage_plan'|'manage_settings'|'manage_team'|'manage_terminals'|'manage_transactions'|'view_contacts'|'view_disputes'|'view_payouts'|'view_plan'|'view_settings'|'view_statements'|'view_team'|'view_transactions'}
 * @psalm-type KeyAttributes = array{key_id: string, key: string, name: string, revoked: bool, created_at: int|float, permissions: 'authorize_payout'|'manage_contacts'|'manage_developer_settings'|'manage_disputes'|'manage_payouts'|'manage_plan'|'manage_settings'|'manage_team'|'manage_terminals'|'manage_transactions'|'view_contacts'|'view_disputes'|'view_payouts'|'view_plan'|'view_settings'|'view_statements'|'view_team'|'view_transactions'}
 * @phpstan-type MetadataAttributes = array{key: string, metadata?: mixed|null, bytesize: int|float, version: int|float, updated_at: int|float, created_at: int|float, expires_at?: int|float|null}
 * @psalm-type MetadataAttributes = array{key: string, metadata?: mixed|null, bytesize: int|float, version: int|float, updated_at: int|float, created_at: int|float, expires_at?: int|float|null}
 * @phpstan-type MonthlyBillingArchiveAttributes = array{id: string, url: string}
 * @psalm-type MonthlyBillingArchiveAttributes = array{id: string, url: string}
 * @phpstan-type EventDeliveryAttributes = array{event_id: string, event_delivery_id: string, status_code: int|float, failure_code?: string|null, url?: string|null, occurrences: int|float, created_at: int|float, payload: array<string, mixed>}
 * @psalm-type EventDeliveryAttributes = array{event_id: string, event_delivery_id: string, status_code: int|float, failure_code?: string|null, url?: string|null, occurrences: int|float, created_at: int|float, payload: array<string, mixed>}
 * @phpstan-type PaymentActivityAttributes = array{created_at: int|float, activity: 'dispute_lost'|'dispute_opened'|'dispute_reopened'|'dispute_won'|'payment_authorized'|'payment_cancelled'|'payment_checked_out'|'payment_completed'|'payment_created'|'payment_redirect_url_accessed'|'payment_reversed'|'payment_selected'|'payment_voided'|'refund_completed'|'refund_created'|'refund_failed'|'refund_initiated', path?: string|null, ip_address?: string|null, severity: 'INFO'|'WARN'|'ERROR', request?: array<string, mixed>|null, response?: mixed|null, event_deliveries?: list<EventDeliveryAttributes>|null, payment_status?: string|null}
 * @psalm-type PaymentActivityAttributes = array{created_at: int|float, activity: 'dispute_lost'|'dispute_opened'|'dispute_reopened'|'dispute_won'|'payment_authorized'|'payment_cancelled'|'payment_checked_out'|'payment_completed'|'payment_created'|'payment_redirect_url_accessed'|'payment_reversed'|'payment_selected'|'payment_voided'|'refund_completed'|'refund_created'|'refund_failed'|'refund_initiated', path?: string|null, ip_address?: string|null, severity: 'INFO'|'WARN'|'ERROR', request?: array<string, mixed>|null, response?: mixed|null, event_deliveries?: list<EventDeliveryAttributes>|null, payment_status?: string|null}
 * @phpstan-type PlannedMaintenanceAttributes = array{id: string, service_name: string, affected_services: 'api'|'vault'|'relay'|'wallet'|'merchant-dashboard'|'website'|'developer', status: 'pending'|'in_progress'|'completed', start_at: int|float, end_at: int|float, expected_impact: 'no_impact'|'intermittent_service'|'temporary_unavailable', expected_impact_description: string, purpose: string, action_required: 'no_action'|'pause_transactions'}
 * @psalm-type PlannedMaintenanceAttributes = array{id: string, service_name: string, affected_services: 'api'|'vault'|'relay'|'wallet'|'merchant-dashboard'|'website'|'developer', status: 'pending'|'in_progress'|'completed', start_at: int|float, end_at: int|float, expected_impact: 'no_impact'|'intermittent_service'|'temporary_unavailable', expected_impact_description: string, purpose: string, action_required: 'no_action'|'pause_transactions'}
 * @phpstan-type PlanSkuAttributes = array{code: string, category: 'DISPUTE'|'INSTALLMENT'|'PLAN'|'PAYOUT'|'REFUND'|'METADATA'|'VIRTUAL_ACCOUNT'|'BANK_ACCOUNT_VALIDATION'|'COMMUNICATION'|'PAYMENT', unit_price?: int|float|null, percentage_fee?: int|float|null, fee_calculation_method: 'SUM'|'HIGHER_OF'|'LOWER_OF', discount_unit_price?: int|float|null, discount_percentage_fee?: int|float|null, discount_minimum_billable_amount?: int|float|null, discount_maximum_billable_amount?: int|float|null, discount_maximum_monthly_uses?: int|float|null, discount_shared_monthly_uses: true|false}
 * @psalm-type PlanSkuAttributes = array{code: string, category: 'DISPUTE'|'INSTALLMENT'|'PLAN'|'PAYOUT'|'REFUND'|'METADATA'|'VIRTUAL_ACCOUNT'|'BANK_ACCOUNT_VALIDATION'|'COMMUNICATION'|'PAYMENT', unit_price?: int|float|null, percentage_fee?: int|float|null, fee_calculation_method: 'SUM'|'HIGHER_OF'|'LOWER_OF', discount_unit_price?: int|float|null, discount_percentage_fee?: int|float|null, discount_minimum_billable_amount?: int|float|null, discount_maximum_billable_amount?: int|float|null, discount_maximum_monthly_uses?: int|float|null, discount_shared_monthly_uses: true|false}
 * @phpstan-type PlanAttributes = array{id?: string|null, tier_code: 'pay_as_you_go'|'startup'|'professional'|'enterprise', tier_name: string, features: list<string>, revision_time: int|float, selectable: bool, rank: int|float, skus: list<PlanSkuAttributes>}
 * @psalm-type PlanAttributes = array{id?: string|null, tier_code: 'pay_as_you_go'|'startup'|'professional'|'enterprise', tier_name: string, features: list<string>, revision_time: int|float, selectable: bool, rank: int|float, skus: list<PlanSkuAttributes>}
 * @phpstan-type RefundAttributes = array{refund_id: string, external_id?: string|null, payment_id: string, status: string, amount: int|float, currency: string, created_at: int|float, refunded_at?: int|float|null, full: true|false}
 * @psalm-type RefundAttributes = array{refund_id: string, external_id?: string|null, payment_id: string, status: string, amount: int|float, currency: string, created_at: int|float, refunded_at?: int|float|null, full: true|false}
 * @phpstan-type RoleAttributes = array{id: string, name: string, description: string, custom: true|false, permissions: 'authorize_payout'|'view_team'|'view_settings'|'view_contacts'|'view_transactions'|'view_statements'|'view_payouts'|'view_plan'|'view_disputes'|'manage_team'|'manage_settings'|'manage_contacts'|'manage_transactions'|'manage_payouts'|'manage_plan'|'manage_developer_settings'|'manage_terminals'|'manage_disputes'}
 * @psalm-type RoleAttributes = array{id: string, name: string, description: string, custom: true|false, permissions: 'authorize_payout'|'view_team'|'view_settings'|'view_contacts'|'view_transactions'|'view_statements'|'view_payouts'|'view_plan'|'view_disputes'|'manage_team'|'manage_settings'|'manage_contacts'|'manage_transactions'|'manage_payouts'|'manage_plan'|'manage_developer_settings'|'manage_terminals'|'manage_disputes'}
 * @phpstan-type StatementAttributes = array{date: string, payment_amount: int|float, payment_count: int|float, refund_amount: int|float, refund_count: int|float, dispute_amount?: int|float|null, dispute_count?: int|float|null, payout_amount: int|float, payout_count: int|float, fees: int|float, net_transaction_volume: int|float, invoice_amount: int|float, invoice_status: 'unbilled'|'due'|'settled', invoice_pdf_download_url?: string|null, invoice_id?: string|null, change_in_balance: int|float, csv_download_url?: string|null, pdf_download_url?: string|null}
 * @psalm-type StatementAttributes = array{date: string, payment_amount: int|float, payment_count: int|float, refund_amount: int|float, refund_count: int|float, dispute_amount?: int|float|null, dispute_count?: int|float|null, payout_amount: int|float, payout_count: int|float, fees: int|float, net_transaction_volume: int|float, invoice_amount: int|float, invoice_status: 'unbilled'|'due'|'settled', invoice_pdf_download_url?: string|null, invoice_id?: string|null, change_in_balance: int|float, csv_download_url?: string|null, pdf_download_url?: string|null}
 * @phpstan-type InvitationAttributes = array{invitation_id: string, email: string, role_id: string, inviter_id: string}
 * @psalm-type InvitationAttributes = array{invitation_id: string, email: string, role_id: string, inviter_id: string}
 * @phpstan-type MerchantMembershipAttributes = array{user_id: string, name: string, email: string, email_confirmed: true|false, mfa_enabled: true|false, role_id: string, preferred_language: string}
 * @psalm-type MerchantMembershipAttributes = array{user_id: string, name: string, email: string, email_confirmed: true|false, mfa_enabled: true|false, role_id: string, preferred_language: string}
 * @phpstan-type VirtualAccountAttributes = array{virtual_account_id: string, state: 'pending'|'opening'|'active'|'closing'|'failed'|'closed', number?: string|null, bank_name: string, display_name: string, currency: string, created_at: int|float}
 * @psalm-type VirtualAccountAttributes = array{virtual_account_id: string, state: 'pending'|'opening'|'active'|'closing'|'failed'|'closed', number?: string|null, bank_name: string, display_name: string, currency: string, created_at: int|float}
 * @phpstan-type WebhookAttributes = array{webhook_id: string, name: string, enabled: bool, api_version: string, destination_url: string, protocol: 'HTTP'|'HTTPS', signature_type: 'SHA256'|'HMAC_SHA256'|'ED25519', verification_key: string, events: 'PAYMENT_COMPLETED'|'PAYMENT_VOIDED'|'PAYMENT_REVERSED'|'PAYMENT_CANCELLED'|'PAYMENT_AUTHORIZED'|'PAYMENT_EXPIRED'|'PAYOUT_COMPLETED'|'TOPUP_SETTLED', created_at: int|float}
 * @psalm-type WebhookAttributes = array{webhook_id: string, name: string, enabled: bool, api_version: string, destination_url: string, protocol: 'HTTP'|'HTTPS', signature_type: 'SHA256'|'HMAC_SHA256'|'ED25519', verification_key: string, events: 'PAYMENT_COMPLETED'|'PAYMENT_VOIDED'|'PAYMENT_REVERSED'|'PAYMENT_CANCELLED'|'PAYMENT_AUTHORIZED'|'PAYMENT_EXPIRED'|'PAYOUT_COMPLETED'|'TOPUP_SETTLED', created_at: int|float}
 * @phpstan-type MerchantConfigurationAttributes = array{funding_currency: 'THB', time_zone: string, time_offset: string, minimum_automated_payout_amount: int|float, maximum_automated_payout_amount: int|float, high_value_payout_threshold_amount: int|float, payouts_blocked: true|false}
 * @psalm-type MerchantConfigurationAttributes = array{funding_currency: 'THB', time_zone: string, time_offset: string, minimum_automated_payout_amount: int|float, maximum_automated_payout_amount: int|float, high_value_payout_threshold_amount: int|float, payouts_blocked: true|false}
 * @phpstan-type CurrentPlanAttributes = array{id?: string|null, tier_code: 'pay_as_you_go'|'startup'|'professional'|'enterprise', tier_name: string, features: list<string>, revision_time: int|float, selectable: bool, rank: int|float, skus: list<PlanSkuAttributes>, pending_plan_id?: string|null}
 * @psalm-type CurrentPlanAttributes = array{id?: string|null, tier_code: 'pay_as_you_go'|'startup'|'professional'|'enterprise', tier_name: string, features: list<string>, revision_time: int|float, selectable: bool, rank: int|float, skus: list<PlanSkuAttributes>, pending_plan_id?: string|null}
 * @phpstan-type MerchantSettingAttributes = array{preferred_communication_language: 'en'|'th', payout_contact_authorization_required: true|false, payout_allow_self_authorization: true|false, payout_skip_authorization_for_api_key: true|false}
 * @psalm-type MerchantSettingAttributes = array{preferred_communication_language: 'en'|'th', payout_contact_authorization_required: true|false, payout_allow_self_authorization: true|false, payout_skip_authorization_for_api_key: true|false}
 * @phpstan-type KeyPersonAttributes = array{key_person_id: string, full_name: string, first_name: string, middle_name?: string|null, last_name: string, first_name_local_language?: string|null, middle_name_local_language?: string|null, last_name_local_language?: string|null, date_of_birth: string, nationality?: string|null, type: 'UBO'|'AUTHORIZED_SIGNATORY', percentage_of_ownership?: int|float|null, authorized_director?: bool|null, national_id?: string|null, national_id_expiry_date?: string|null, national_id_without_expiry_date: true|false, passport_number?: string|null, passport_expiry_date?: string|null, address_street_1?: string|null, address_street_2?: string|null, postal_code?: string|null, city?: string|null, primary_administrative_division?: string|null, country?: string|null, current_address_street_1?: string|null, current_address_street_2?: string|null, current_address_postal_code?: string|null, current_address_city?: string|null, current_address_primary_administrative_division?: string|null, current_address_country?: string|null, use_residence_registration_as_current_address: true|false}
 * @psalm-type KeyPersonAttributes = array{key_person_id: string, full_name: string, first_name: string, middle_name?: string|null, last_name: string, first_name_local_language?: string|null, middle_name_local_language?: string|null, last_name_local_language?: string|null, date_of_birth: string, nationality?: string|null, type: 'UBO'|'AUTHORIZED_SIGNATORY', percentage_of_ownership?: int|float|null, authorized_director?: bool|null, national_id?: string|null, national_id_expiry_date?: string|null, national_id_without_expiry_date: true|false, passport_number?: string|null, passport_expiry_date?: string|null, address_street_1?: string|null, address_street_2?: string|null, postal_code?: string|null, city?: string|null, primary_administrative_division?: string|null, country?: string|null, current_address_street_1?: string|null, current_address_street_2?: string|null, current_address_postal_code?: string|null, current_address_city?: string|null, current_address_primary_administrative_division?: string|null, current_address_country?: string|null, use_residence_registration_as_current_address: true|false}
 * @phpstan-type UploadAttributes = array{upload_id: string, filename: string, file_size: int|float, file_type: string, url: string}
 * @psalm-type UploadAttributes = array{upload_id: string, filename: string, file_size: int|float, file_type: string, url: string}
 * @phpstan-type RequiredDocumentAttributes = array{label: string, allow_multiple: bool, uploaded_documents: list<UploadAttributes>}
 * @psalm-type RequiredDocumentAttributes = array{label: string, allow_multiple: bool, uploaded_documents: list<UploadAttributes>}
 * @phpstan-type SalesChannelAttributes = array{sales_channel_id: string, platform: 'WEB'|'IOS'|'ANDROID'|'TERMINAL', url: string, demo_url?: string|null, instructions?: string|null}
 * @psalm-type SalesChannelAttributes = array{sales_channel_id: string, platform: 'WEB'|'IOS'|'ANDROID'|'TERMINAL', url: string, demo_url?: string|null, instructions?: string|null}
 * @phpstan-type OnboardingAttributes = array{onboarding_id: string, test_merchant_id: string, live_merchant_id?: string|null, status: 'DRAFT'|'PENDING'|'APPROVED'|'REJECTED', readiness: array<string, mixed>, entity_type: 'UNKNOWN'|'BUSINESS'|'INDIVIDUAL', merchant_display_name: string, legal_name?: string|null, category_code: string, product_types: 'DIGITAL_GOODS'|'PHYSICAL_GOODS'|'SERVICES', business?: array<string, mixed>|null, individual?: array<string, mixed>|null, key_people: list<KeyPersonAttributes>, primary_bank_account?: array<string, mixed>|null, auto_payout_frequency: 'NEVER'|'DAILY'|'WEEKLY', expected_monthly_transaction_amount: int|float, expected_maximum_transaction_amount: int|float, required_documents: list<RequiredDocumentAttributes>, sales_channels: list<SalesChannelAttributes>, preferred_communication_language: 'en'|'th', requested_changes_reason?: string|null, rejection_reason?: string|null}
 * @psalm-type OnboardingAttributes = array{onboarding_id: string, test_merchant_id: string, live_merchant_id?: string|null, status: 'DRAFT'|'PENDING'|'APPROVED'|'REJECTED', readiness: array<string, mixed>, entity_type: 'UNKNOWN'|'BUSINESS'|'INDIVIDUAL', merchant_display_name: string, legal_name?: string|null, category_code: string, product_types: 'DIGITAL_GOODS'|'PHYSICAL_GOODS'|'SERVICES', business?: array<string, mixed>|null, individual?: array<string, mixed>|null, key_people: list<KeyPersonAttributes>, primary_bank_account?: array<string, mixed>|null, auto_payout_frequency: 'NEVER'|'DAILY'|'WEEKLY', expected_monthly_transaction_amount: int|float, expected_maximum_transaction_amount: int|float, required_documents: list<RequiredDocumentAttributes>, sales_channels: list<SalesChannelAttributes>, preferred_communication_language: 'en'|'th', requested_changes_reason?: string|null, rejection_reason?: string|null}
 */
final class Client
{
    private \WebFunction\Client $client;

    public function __construct(?string $bearerAuth = null, ?string $version = null)
    {
        $this->client = \WebFunction\Client::fromPackageEndpoint('https://api.reservepay.com/merchants', bearerAuth: $bearerAuth, version: $version);
    }

    /**
     * Escape hatch for calling an endpoint this client doesn't have a typed method for.
     */
    public function call(string $name, array $args = []): mixed
    {
        return $this->client->call($name, $args);
    }

    public function getPackage(): ?Package
    {
        return $this->client->getPackage();
    }

    public function setBearerAuth(?string $bearerAuth): void
    {
        $this->client->setBearerAuth($bearerAuth);
    }

    public function setVersion(?string $version): void
    {
        $this->client->setVersion($version);
    }

    public function setPipeline(?Pipeline $pipeline): void
    {
        $this->client->setPipeline($pipeline);
    }

    /**
     * ## Abort a payment on a terminal
     *
     * This endpoint cancels any payment request that was sent to a terminal via start-payment-on-terminal.
     * This only works if the payment was not yet started.
     *
     * It returns true if successful.
     * @param array{terminal_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The terminal could not be found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function abortPaymentOnTerminal(array $args): bool
    {
        return $this->client->call('abort-payment-on-terminal', $args);
    }

    /**
     * ## Add contact
     *
     * Adds a new contact for the merchant. Either `contact_email` or `contact_mobile` must be provided.
     * The service will check for existing contacts with the same ID, email, or mobile number
     * and raise a conflict error if any duplicates are found.
     * @param array{external_id?: string, contact_email?: string, contact_mobile?: string, contact_name?: string, favorite?: bool} $args
     * @return string
     * @throws BadRequestError CONFLICT - A contact with the same ID, email, or mobile number already exists.
     * Each contact must have unique values for these fields.
     * @throws BadRequestError VALIDATION_ERROR - Validation failed. Either `contact_email` or `contact_mobile` must be provided.
     * If `contact_email` is provided, it must be a valid email format.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addContact(array $args = []): string
    {
        return $this->client->call('add-contact', $args);
    }

    /**
     * ## Add a key person
     *
     * Add a key person to your business merchant account. Key people are individuals
     * who have significant control or influence over your business, such as:
     *
     * - Ultimate Beneficial Owners (UBOs) who own 25% or more of the business
     * - Authorized Signatories who can legally represent and make decisions for
     *   the business
     *
     * This endpoint is only available for business merchants, not
     * individual merchants.
     *
     * You'll need to provide basic information about the person as well as their
     * identification documents:
     *
     * - For local nationals (nationality matches merchant's country): National ID
     *   details are required
     * - For foreign nationals: Passport details are required
     *
     * The endpoint returns a unique identifier for the key person that can be used to
     * reference them in other API calls.
     * @param array{first_name: string, middle_name?: string, last_name: string, first_name_local_language?: string, middle_name_local_language?: string, last_name_local_language?: string, date_of_birth: string, nationality: string, type: 'UBO'|'AUTHORIZED_SIGNATORY', percentage_of_ownership?: int|float, authorized_director?: bool, national_id?: string, national_id_expiry_date?: string, national_id_without_expiry_date?: bool, passport_number?: string, passport_expiry_date?: string, address_street_1?: string, address_street_2?: string, postal_code?: string, city?: string, primary_administrative_division?: string, country?: string, current_address_street_1?: string, current_address_street_2?: string, current_address_postal_code?: string, current_address_city?: string, current_address_primary_administrative_division?: string, current_address_country?: string, use_residence_registration_as_current_address?: bool} $args
     * @return string
     * @throws BadRequestError NOT_BUSINESS - This error occurs when attempting to add a key person to an individual merchant
     * account. Key persons can only be added to business merchant accounts. If you
     * need to upgrade to a business account, please contact support.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addKeyPerson(array $args): string
    {
        return $this->client->call('add-key-person', $args);
    }

    /**
     * ## Add a payment link
     *
     * Create a new payment link that can be shared with customers to collect payments.
     * A payment link is a unique URL that redirects customers to a hosted checkout page.
     *
     * If no `installation_*` attributes are provided, the merchant's default link installation will be used.
     * If any `installation_*` attributes (like `installation_email_capture`) are provided,
     * a new one-off installation will be created for this payment link,
     * based on the merchant's default link installation attributes and the provided attributes.
     * @param array{amount: int|float, currency: string, interest_bearer?: 'customer'|'merchant', description: string, external_id?: string, expires_at?: int|float, single_use?: true|false, capture?: true|false, send_confirmation_email?: true|false, installation_name?: string, installation_payment_methods?: string, installation_email_capture?: 'disabled'|'optional'|'required', installation_mobile_capture?: 'disabled'|'optional'|'required', installation_address_capture?: 'disabled'|'optional'|'required'} $args
     * @return array{payment_link_id: string, batch_id?: string|null, payment_link_url: string, status: 'active'|'disabled'|'expired'|'paid', amount: int|float, currency: string, interest_bearer?: 'customer'|'merchant'|null, description: string, external_id?: string|null, expires_at?: int|float|null, single_use: bool, capture: bool, send_confirmation_email: bool, installation: InstallationAttributes, created_at: int|float, deletable: bool}
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addPaymentLink(array $args): array
    {
        return $this->client->call('add-payment-link', $args);
    }

    /**
     * ## Add a Sales Channel
     *
     * Add a new sales channel for the merchant.
     * @param array{platform: 'WEB'|'IOS'|'ANDROID'|'TERMINAL', url?: string, demo_url?: string, instructions?: string} $args
     * @return string
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addSalesChannel(array $args): string
    {
        return $this->client->call('add-sales-channel', $args);
    }

    /**
     * ## Add a new terminal
     *
     * Creates a new terminal under the authenticated merchant account.
     * @param array{name: string, external_id?: string, installation_id?: string, branch_name?: string, kiosk_mode?: bool, delegated_terminal_id?: string} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - The installation_id or delegated_terminal_id does not exist in this merchant account.
     * @throws BadRequestError INVALID_ARGUMENT - The delegated terminal already has a delegation configured, it is already used by another terminal, or kiosk_mode cannot be combined with delegated_terminal_id.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addTerminal(array $args): string
    {
        return $this->client->call('add-terminal', $args);
    }

    /**
     * ## Add to metadata set
     *
     * Adds an item to a metadata set. A metadata set is an array of unique values within your metadata. This endpoint ensures
     * that the metadata at the given path is treated as a set:
     * - If the path is empty or does not exist, it is initialized as an empty array.
     * - If the value already exists in the set, it is not added again.
     * - The metadata at the path must be an array (or empty).
     *
     * #### Example
     * If current metadata is `{"tags": ["a"]}`, adding `"b"` to path `"tags"` results in `{"tags": ["a", "b"]}`.
     * Adding `"a"` again results in no change.
     *
     * #### Expiration (TTL)
     * Adding to a set also updates the `expires_at` timestamp for the entire metadata record.
     * @param array{key: string, path?: string, value: mixed, version?: int|float, ttl?: int|float} $args
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addToMetadataSet(array $args): bool
    {
        return $this->client->call('add-to-metadata-set', $args);
    }

    /**
     * ## Add a topup
     *
     * Add funds to your test merchant balance. This simulates an incoming bank
     * transfer settling into your account, letting you build up a balance so you can
     * exercise other endpoints such as sending payouts.
     *
     * Optionally supply a `virtual_account_number` to simulate the transfer
     * arriving into one of your virtual accounts.
     *
     * This endpoint is only available to test merchant accounts. Live merchants are
     * funded through real bank transfers and cannot top up their balance through the
     * API.
     *
     * The topup settles immediately, so the funds are available in your balance as
     * soon as this call returns.
     * @param array{amount: string, virtual_account_number?: string, transaction_id?: string, sender_account_name?: string, sender_account_number?: string, sender_bank_code?: string} $args
     * @return string
     * @throws BadRequestError NOT_ALLOWED - Topups can only be added to test merchant accounts. This error occurs when
     * the authenticated merchant is a live account.
     * @throws BadRequestError NOT_FOUND - This error occurs when a `virtual_account_number` is supplied but no
     * virtual account with that number exists in your merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function addTopup(array $args): string
    {
        return $this->client->call('add-topup', $args);
    }

    /**
     * ## Assign a role to a user in the merchant account
     *
     * This endpoint allows you to assign an existing role to a user within your merchant account.
     * @param array{user_id: string, role_id?: string} $args
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function assignRole(array $args): bool
    {
        return $this->client->call('assign-role', $args);
    }

    /**
     * ## List available API versions
     *
     * Get a list of all API versions available that are currently enabled for your merchant account.
     *
     * The endpoint returns an array of API versions that you can use in other endpoints.
     * @return list<string>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function availableApiVersions(): array
    {
        return $this->client->call('available-api-versions');
    }

    /**
     * ## Balance history
     *
     * Returns a paginated list of **ledger entries** for the merchant's default balance, newest first.
     * Each row includes the owning **transaction** (`transaction_id`) alongside the entry fields.
     * There is **no time window**: all entries for the balance are included, ordered by recency.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the balance history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return BalanceHistoryPage
     * @throws BadRequestError NOT_FOUND - No balance exists for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function balanceHistory(array $args = []): BalanceHistoryPage
    {
        $raw = $this->client->call('balance-history', $args);
        /** @var Page $raw */
        return new BalanceHistoryPage($raw);
    }

    /**
     * ## Cancel payment
     *
     * Cancels a payment that should not go through. The customer is not charged for a cancelled payment.
     *
     * Payments in a `PENDING` status, and payments in a `SUCCESSFUL` status that have not been
     * settled yet, can be cancelled. Once a payment is cancelled, its status changes to `CANCELLED`.
     * @param array{payment_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment ID.
     * @throws BadRequestError NOT_CANCELLABLE - Payment is not eligible for cancellation. This can happen if the payment is
     * not in a `PENDING` or `SUCCESSFUL` status, has already been settled, or its
     * payment method does not support cancellations.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function cancelPayment(array $args): bool
    {
        return $this->client->call('cancel-payment', $args);
    }

    /**
     * ## Capture payment
     *
     * Captures a previously authorized payment. Capture is the second step of a two-step payment flow
     * (authorize and capture). This endpoint should be called when you are ready to complete the
     * transaction and receive funds, such as when shipping goods or providing a service.
     *
     * Only payments that are in an `AUTHORIZED` status can be captured. Once a payment is captured,
     * its status changes to `SUCCESSFUL` (or `FAILED` if the capture fails).
     * @param array{payment_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment ID.
     * @throws BadRequestError NOT_CAPTURABLE - Payment is not eligible for capture. This can happen if the payment
     * is not in an authorized state, the capture window has passed, or the
     * integration does not support captures.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function capturePayment(array $args): bool
    {
        return $this->client->call('capture-payment', $args);
    }

    /**
     * ## Challenge dispute
     *
     * Confirms the submission of evidence and formally challenges the dispute.
     *
     * A non-refundable challenge fee is charged to the merchant's balance via the daily statement upon
     * submission. This fee is waived if the original payment was protected by 3D Secure. Failure to
     * challenge the dispute within the deadline results in an automatic loss in favor of the customer.
     * @param array{dispute_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The specified dispute was not found.
     * @throws BadRequestError INVALID_TRANSITION - The dispute is not in a state that can be challenged (e.g., already won or lost).
     * @throws BadRequestError MISSING_EVIDENCE - No evidence has been uploaded for this dispute action.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function challengeDispute(array $args): bool
    {
        return $this->client->call('challenge-dispute', $args);
    }

    /**
     * ## Clone a payment link
     *
     * Create a new payment link based on an existing one.
     * The new link will have the same amount, currency, description, and other settings.
     * @param array{payment_link_id: string, external_id?: string} $args
     * @return array{payment_link_id: string, batch_id?: string|null, payment_link_url: string, status: 'active'|'disabled'|'expired'|'paid', amount: int|float, currency: string, interest_bearer?: 'customer'|'merchant'|null, description: string, external_id?: string|null, expires_at?: int|float|null, single_use: bool, capture: bool, send_confirmation_email: bool, installation: InstallationAttributes, created_at: int|float, deletable: bool}
     * @throws BadRequestError not_found - The original payment link could not be found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function clonePaymentLink(array $args): array
    {
        return $this->client->call('clone-payment-link', $args);
    }

    /**
     * ## Concede dispute
     *
     * Concedes the dispute and accepts the loss. This action is irreversible. The disputed funds are
     * permanently deducted from the merchant's balance.
     * @param array{dispute_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The specified dispute was not found.
     * @throws BadRequestError INVALID_TRANSITION - The dispute is not in a state that can be conceded (e.g., already won or lost).
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function concedeDispute(array $args): bool
    {
        return $this->client->call('concede-dispute', $args);
    }

    /**
     * ## Confirm branding upload
     *
     * Confirms the upload of a branding asset by providing the `upload_id`. This effectively attaches
     * the uploaded file to the merchant branding.
     * @param array{upload_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no upload with the provided ID exists in your
     * merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function confirmBrandingUpload(array $args): bool
    {
        return $this->client->call('confirm-branding-upload', $args);
    }

    /**
     * ## Confirm dispute evidence upload
     *
     * Confirms the upload of an evidence file for a dispute by providing the `upload_id`.
     * @param array{dispute_id: string, upload_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The specified dispute or upload was not found.
     * @throws BadRequestError not_uploadable - The dispute is not in a state that allows evidence upload.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function confirmDisputeEvidenceUpload(array $args): bool
    {
        return $this->client->call('confirm-dispute-evidence-upload', $args);
    }

    /**
     * ## Confirm Onboarding Upload
     *
     * Confirms the upload of an onboarding document by providing the `upload_id`.
     *
     * This effectively attaches the uploaded file to the relevant merchant, business, individual or key person.
     * @param array{upload_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no upload with the provided ID exists in your
     * merchant account.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function confirmOnboardingUpload(array $args): bool
    {
        return $this->client->call('confirm-onboarding-upload', $args);
    }

    /**
     * ## Disable an SDK installation
     *
     * Temporarily disable an SDK installation. This will prevent the installation from
     * being able to make API calls like tokenizing cards or creating payments.
     *
     * Common reasons to disable an installation:
     *
     * - You're decommissioning an old version of your application
     * - You want to temporarily stop accepting payments on a specific platform
     * - You suspect unauthorized or malicious use from a specific installation
     * - You're performing maintenance or updates on your application
     *
     * You can re-enable a disabled installation at any time using the
     * `enable-installation` endpoint.
     *
     * Note that disabling an installation only affects its ability to make API calls -
     * it does not delete the installation or any of its settings. The installation ID
     * remains valid and can be used again once re-enabled.
     * @param array{installation_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no installation with the provided ID exists in your
     * merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function disableInstallation(array $args): bool
    {
        return $this->client->call('disable-installation', $args);
    }

    /**
     * ## Disable a payment link
     *
     * Disables a payment link so it can no longer be used to collect payments.
     * @param array{payment_link_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment link with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment link ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function disablePaymentLink(array $args): bool
    {
        return $this->client->call('disable-payment-link', $args);
    }

    /**
     * ## Disable a terminal
     *
     * Temporarily disable a terminal. This does not delete the terminal.
     * @param array{terminal_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No terminal exists with the provided ID in this merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function disableTerminal(array $args): bool
    {
        return $this->client->call('disable-terminal', $args);
    }

    /**
     * ## Disconnect webhook
     *
     * Disconnects an existing webhook endpoint. Once disconnected, the endpoint will no longer receive
     * any event notifications.
     * @param array{webhook_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when the webhook endpoint with the provided ID could not be found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function disconnectWebhook(array $args): bool
    {
        return $this->client->call('disconnect-webhook', $args);
    }

    /**
     * ## Enable an SDK installation
     *
     * Re-enable a previously disabled SDK installation. Once enabled, the installation
     * can resume making API calls like tokenizing cards or creating payments.
     *
     * Common reasons to re-enable an installation:
     *
     * - You've completed maintenance or updates on your application
     * - You're ready to resume accepting payments on a specific platform
     * - You've resolved any security concerns about unauthorized usage
     * - You're bringing a previously decommissioned application back online
     *
     * Note that enabling an installation takes effect immediately and restores the
     * exact same configuration and settings from before it was disabled. This also
     * does not require any changes to your integration code and can be reversed at any
     * time by using the `disable-installation` endpoint.
     *
     * For security best practices, we recommend only enabling installations that are
     * actively being used.
     * @param array{installation_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no installation with the provided ID exists in your
     * merchant account. Double check that you're using the correct installation ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function enableInstallation(array $args): bool
    {
        return $this->client->call('enable-installation', $args);
    }

    /**
     * ## Enable a payment link
     *
     * Enables a previously disabled payment link so it can be used to collect payments again.
     * @param array{payment_link_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment link with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment link ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function enablePaymentLink(array $args): bool
    {
        return $this->client->call('enable-payment-link', $args);
    }

    /**
     * ## Enable a terminal
     *
     * Re-enable a previously disabled terminal.
     * @param array{terminal_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No terminal exists with the provided ID in this merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function enableTerminal(array $args): bool
    {
        return $this->client->call('enable-terminal', $args);
    }

    /**
     * ## Escalate dispute
     *
     * Requests a final review of the dispute by the bank.
     *
     * An escalation hold is applied to the merchant's balance to help cover the potential escalation fee.
     * This hold is released once the dispute is resolved, regardless of the outcome. If the dispute is lost, a separate
     * escalation fee is charged to the merchant's balance via the daily statement. The hold and
     * fee are independent and their amounts may differ; refer to the dispute escalation
     * documentation for specific fee schedules by region and provider.
     * @param array{dispute_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The specified dispute was not found.
     * @throws BadRequestError INVALID_TRANSITION - The dispute is not in a state that can be escalated.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function escalateDispute(array $args): bool
    {
        return $this->client->call('escalate-dispute', $args);
    }

    /**
     * ## Filter contacts
     *
     * Returns a paginated list of contacts for the authenticated merchant with filtering capabilities.
     *
     * Filter by name, email, mobile (substring match) and creation date. Free-text `terms` also matches
     * contact ID and external ID. Results are ordered by creation date descending.
     * @param array{per_page?: int|float, page?: int|float, name?: string, email?: string, mobile?: string, terms?: string, created_at?: string} $args
     * @return FilterContactsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function filterContacts(array $args = []): FilterContactsPage
    {
        $raw = $this->client->call('filter-contacts', $args);
        /** @var Page $raw */
        return new FilterContactsPage($raw);
    }

    /**
     * ## Filter payment links
     *
     * Returns a paginated list of payment links for the authenticated merchant with advanced filtering capabilities.
     *
     * This endpoint provides comprehensive search functionality using database-driven filtering,
     * allowing you to filter links by status, usage, batch, amount ranges, and date ranges.
     *
     * The response includes pagination controls with `previous` and `next` cursors that can be used
     * to navigate through the payment link history efficiently.
     *
     * ### Features:
     * - **Full-text search**: Search across payment link descriptions and UIDs
     * - **Advanced filtering**: Filter by status (active, disabled, expired, paid), usage (single, multiple), and batch ID
     * - **Range support**: Support for amount ranges and date ranges
     * @param array{per_page?: int|float, page?: int|float, terms?: string, status?: 'ACTIVE'|'DISABLED'|'EXPIRED'|'PAID', usage?: 'SINGLE'|'MULTIPLE', batch_id?: string, amount?: string, created_at?: string, expires_at?: string} $args
     * @return FilterPaymentLinksPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function filterPaymentLinks(array $args = []): FilterPaymentLinksPage
    {
        $raw = $this->client->call('filter-payment-links', $args);
        /** @var Page $raw */
        return new FilterPaymentLinksPage($raw);
    }

    /**
     * ## Filter payments
     *
     * Returns a paginated list of payments for the authenticated merchant with advanced filtering capabilities.
     *
     * This endpoint provides comprehensive search functionality using the SearchIndex infrastructure,
     * allowing you to filter payments by various criteria including status, payment method, currency,
     * amount ranges, terminal, and date ranges.
     *
     * The response includes pagination controls with `previous` and `next` cursors that can be used
     * to navigate through the payment history efficiently.
     *
     * ### Features:
     * - **Full-text search**: Search across payment content
     * - **Advanced filtering**: Filter by status, payment method, currency, amount, contact, terminal, dates, and fraud finding
     * - **Contact filtering**: Filter by contact internal ID or external ID
     * - **Date range support**: Support for specific dates, datetimes, date ranges, datetime ranges, and keywords
     * - **Pagination**: Efficient cursor-based pagination
     * - **Merchant-scoped**: Only returns payments for the authenticated merchant
     * @param array{per_page?: int|float, page?: int|float, status?: 'PENDING'|'SUCCESSFUL'|'AUTHORIZED'|'VOIDED'|'REVERSED'|'CANCELLED'|'FAILED'|'EXPIRED'|'DISPUTED'|'REFUNDED', created_at?: string, authorized_at?: string, captured_at?: string, payment_method?: string, currency?: string, amount?: string, contact_id?: string, terminal_id?: string, terms?: string, fraud_finding?: ''|'safe'|'risky'} $args
     * @return FilterPaymentsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function filterPayments(array $args = []): FilterPaymentsPage
    {
        $raw = $this->client->call('filter-payments', $args);
        /** @var Page $raw */
        return new FilterPaymentsPage($raw);
    }

    /**
     * ## Filter Payouts
     *
     * Returns a paginated list of payouts for the authenticated merchant with advanced filtering capabilities.
     *
     * This endpoint provides comprehensive search functionality using the SearchIndex infrastructure,
     * allowing you to filter payouts by various criteria including status, payout method, currency,
     * amount ranges, and date ranges.
     *
     * The response includes pagination controls with `previous` and `next` cursors that can be used
     * to navigate through the payout history efficiently.
     *
     * ### Features:
     * - **Full-text search**: Search across payout content
     * - **Advanced filtering**: Filter by status, payout method, currency, amount, and dates
     * - **Date range support**: Support for specific dates, datetimes, date ranges, datetime ranges, and keywords
     * - **Pagination**: Efficient cursor-based pagination
     * - **Merchant-scoped**: Only returns payouts for the authenticated merchant
     * @param array{per_page?: int|float, page?: int|float, status?: 'SUCCESSFUL'|'SENT'|'PENDING'|'FAILED'|'CANCELLED'|'REQUESTED', created_at?: string, paid_at?: string, sent_at?: string, payout_method?: string, currency?: string, amount?: string, bank_code?: string, contact_id?: string, terms?: string} $args
     * @return FilterPayoutsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function filterPayouts(array $args = []): FilterPayoutsPage
    {
        $raw = $this->client->call('filter-payouts', $args);
        /** @var Page $raw */
        return new FilterPayoutsPage($raw);
    }

    /**
     * ## Filter Terminals
     *
     * Returns a paginated list of terminals for the authenticated merchant with search capabilities.
     *
     * This endpoint provides full-text search functionality using the SearchIndex infrastructure,
     * allowing you to search terminals by name, external ID, and branch name.
     * @param array{per_page?: int|float, page?: int|float, terms?: string} $args
     * @return FilterTerminalsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function filterTerminals(array $args = []): FilterTerminalsPage
    {
        $raw = $this->client->call('filter-terminals', $args);
        /** @var Page $raw */
        return new FilterTerminalsPage($raw);
    }

    /**
     * ## Filter topups
     *
     * Returns a paginated list of topups for the authenticated merchant with advanced filtering capabilities.
     *
     * This endpoint provides comprehensive search functionality using the SearchIndex infrastructure,
     * allowing you to filter topups by funding type (manual or virtual account), virtual account number,
     * amount ranges, and date ranges.
     *
     * The response includes pagination controls with `previous` and `next` cursors that can be used
     * to navigate through the topup history efficiently.
     *
     * ### Features:
     * - **Funding type**: Filter by manual transfers or virtual account topups
     * - **Virtual account**: Filter by the virtual account number that received the funds
     * - **Amount ranges**: Filter by an exact amount or a range
     * - **Date range support**: Support for specific dates, datetimes, date ranges, datetime ranges, and keywords
     * - **Pagination**: Efficient cursor-based pagination
     * - **Merchant-scoped**: Only returns topups for the authenticated merchant
     * @param array{per_page?: int|float, page?: int|float, type?: 'MANUAL'|'VIRTUAL_ACCOUNT', virtual_account_number?: string, amount?: string, created_at?: string, terms?: string} $args
     * @return FilterTopupsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function filterTopups(array $args = []): FilterTopupsPage
    {
        $raw = $this->client->call('filter-topups', $args);
        /** @var Page $raw */
        return new FilterTopupsPage($raw);
    }

    /**
     * ## Find bank account
     *
     * Retrieves details about a specific bank account. The `bank_account_id` can be either the internal
     * UID (starting with 'bka_') or your own external ID.
     *
     * This endpoint returns both active and soft-deleted bank accounts from the merchant's directory.
     * @param array{bank_account_id: string} $args
     * @return array{bank_account_id: string, bank_code: string, number: string, name: string, country_code: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool, owner: 'CONTACT'|'MERCHANT', primary: bool}
     * @throws BadRequestError NOT_FOUND - No bank account found with the provided bank_account_id in this merchant's directory.
     * The bank account may not exist or may not belong to the merchant's directory.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findBankAccount(array $args): array
    {
        return $this->client->call('find-bank-account', $args);
    }

    /**
     * ## Find batch
     *
     * Retrieve a batch by ID, including its current status, any prepare or processing error,
     * the preview download URL when prepare has completed, and the result download URL when
     * processing has completed.
     * @param array{batch_id: string} $args
     * @return array{batch_id: string, type: 'PAYMENT_LINK'|'PAYOUT', status: 'CREATED'|'PREPARING'|'PREPARED'|'PROCESSING'|'COMPLETED', output_url?: string|null, preview_url?: string|null, error_message?: string|null, created_at: int|float}
     * @throws BadRequestError NOT_FOUND - The batch was not found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findBatch(array $args): array
    {
        return $this->client->call('find-batch', $args);
    }

    /**
     * ## Find card
     *
     * Find a card by its ID or token. Either card_id or token must be provided.
     *
     * The card_id can be either the internal UID (starting with 'crd_') or your own external ID.
     *
     * The token can be either the token UID (starting with 'tkn_') or the raw card number associated with the token.
     *
     * This endpoint returns both active and soft-deleted cards from the merchant's directory.
     * @param array{card_id?: string, token?: string} $args
     * @return array{card_id: string, token: string, usage: 'SINGLE'|'MULTIPLE', used: bool, first_digits: string, last_digits: string, length: int|float, expiration_date: string, cardholder_name: string, scheme: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool}
     * @throws BadRequestError NOT_FOUND - No card found with the provided card_id or token in this merchant's directory.
     * The card may not exist or may not belong to the merchant's directory.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findCard(array $args = []): array
    {
        return $this->client->call('find-card', $args);
    }

    /**
     * ## Find a communication delivery
     *
     * Retrieves a communication delivery queued by send-payment-link, including its latest status.
     * @param array{communication_delivery_id: string} $args
     * @return array{communication_delivery_id: string, communication_channel_id: string, payment_link_id: string, delivery: 'email'|'sms', recipient: string, contact_id?: string|null, status: 'queued'|'processing'|'accepted'|'failed', created_at: int|float, accepted_at?: int|float|null}
     * @throws BadRequestError NOT_FOUND - The communication delivery does not exist for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findCommunicationDelivery(array $args): array
    {
        return $this->client->call('find-communication-delivery', $args);
    }

    /**
     * ## Find dispute
     *
     * Returns the details of a specific dispute.
     * @param array{dispute_id: string} $args
     * @return array{dispute_id: string, payment_id: string, contact_id?: string|null, message?: string|null, status: 'OPEN'|'CHALLENGED'|'REOPENED'|'ESCALATED'|'WON'|'LOST'|'EXPIRED', challenge_received_at?: int|float|null, escalation_received_at?: int|float|null, won_at?: int|float|null, lost_at?: int|float|null, amount: int|float, currency: string, reason: string, suggested_evidence: 'cancellation_policy'|'cancellation_request_records'|'common_evidence'|'customer_communication'|'device_id_ip_data'|'marketing_materials'|'payment_logs'|'payment_processing_records'|'receipt'|'refund_issuance'|'refund_policy'|'refund_records'|'return_policy'|'shipping_documentation'|'supplier_certification'|'technical_support_logs'|'transaction_history'|'transaction_logs', deadline_at?: int|float|null, created_at: int|float, fees: list<FeeAttributes>, conceded: true|false, evidences: list<DisputeEvidenceAttributes>}
     * @throws BadRequestError NOT_FOUND - The specified dispute was not found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findDispute(array $args): array
    {
        return $this->client->call('find-dispute', $args);
    }

    /**
     * ## Find an SDK installation
     *
     * Retrieve details about a specific SDK installation in your merchant account.
     * This endpoint returns information about the installation's:
     *
     * - Name and platform (web, iOS, Android or Terminal)
     * - Unique identifier (domain, bundle ID or package name)
     * - Enabled payment methods
     * - Current status (enabled or disabled)
     *
     * The installation ID required by this endpoint is the same ID that was returned
     * when you created the installation using the add-installation endpoint.
     * @param array{installation_id: string} $args
     * @return array{installation_id: string, enabled: bool, name: string, platform: 'WEB'|'IOS'|'ANDROID'|'TERMINAL', identifier?: string|null, payment_methods: list<PaymentMethodAttributes>, email_capture: 'disabled'|'optional'|'required', mobile_capture: 'disabled'|'optional'|'required', address_capture: 'disabled'|'optional'|'required', created_at: int|float}
     * @throws BadRequestError NOT_FOUND - This error occurs when no installation with the provided ID exists in your
     * merchant account. Double check that you're using the correct installation ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findInstallation(array $args): array
    {
        return $this->client->call('find-installation', $args);
    }

    /**
     * ## Find an invoice
     *
     * Retrieve details about a specific invoice in your merchant account.
     * @param array{invoice_id: string} $args
     * @return array{invoice_id: string, number: string, status: 'due'|'settled', amount: int|float, net_amount: int|float, vat_amount: int|float, currency: string, issued_on: string, pdf_download_url?: string|null, receipt_id?: string|null, receipt_number?: string|null, receipt_pdf_download_url?: string|null, receipt_xml_download_url?: string|null, receipt_created_at?: int|float|null, created_at: int|float}
     * @throws BadRequestError NOT_FOUND - This error occurs when no invoice with the provided ID exists in your
     * merchant account. Double check that you're using the correct invoice ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findInvoice(array $args): array
    {
        return $this->client->call('find-invoice', $args);
    }

    /**
     * ## Find an API key
     *
     * Retrieve details about a specific API key in your merchant account. This
     * endpoint returns information about the key's:
     *
     * - Name and description
     * - Current status (enabled or disabled)
     *
     * The key ID required by this endpoint is the same ID that was returned when you created
     * the key using the `generate-key` endpoint.
     * @param array{key_id: string} $args
     * @return array{key_id: string, key: string, name: string, revoked: bool, created_at: int|float, permissions: 'authorize_payout'|'manage_contacts'|'manage_developer_settings'|'manage_disputes'|'manage_payouts'|'manage_plan'|'manage_settings'|'manage_team'|'manage_terminals'|'manage_transactions'|'view_contacts'|'view_disputes'|'view_payouts'|'view_plan'|'view_settings'|'view_statements'|'view_team'|'view_transactions'}
     * @throws BadRequestError NOT_FOUND - This error occurs when no API key with the provided ID exists in your merchant
     * account. Double check that you're using the correct key ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findKey(array $args): array
    {
        return $this->client->call('find-key', $args);
    }

    /**
     * ## Find a payment
     *
     * Retrieve details about a specific payment by its unique identifier or external ID.
     * @param array{payment_id: string} $args
     * @return array{payment_id: string, external_id?: string|null, status: 'PENDING'|'AUTHORIZED'|'SUCCESSFUL'|'FAILED'|'REVERSED'|'VOIDED'|'CANCELLED'|'EXPIRED'|'REFUNDED'|'PARTIALLY_REFUNDED'|'DISPUTED', amount: int|float, currency: string, payment_method: string, payment_method_category: string, payment_method_display_name: string, card_first_digits?: string|null, card_last_digits?: string|null, card_masked_number?: string|null, card_network?: string|null, cardholder_name?: string|null, created_at: int|float, authorized_at?: int|float|null, captured_at?: int|float|null, voided_at?: int|float|null, cancelled_at?: int|float|null, refundable_at?: int|float|null, expires_at?: int|float|null, capture: 'AUTOMATIC'|'MANUAL', customer_id?: string|null, customer_external_id?: string|null, customer_email?: string|null, customer_mobile?: string|null, customer_name?: string|null, fees: list<FeeAttributes>, ref1?: string|null, ref2?: string|null, ref3?: string|null, reversible: bool, capturable: bool, voidable: bool, cancellable: bool, refundable: bool, fraud_finding?: 'risky'|'safe'|null, failure_code?: string|null, failure_message?: string|null, terminal_id?: string|null, payment_link_id?: string|null, payment_session_id?: string|null, installment_term?: int|float|null, interest_bearer?: 'customer'|'merchant'|null, installment_interest_rate?: string|null, installment_interest_amount?: int|float|null, installment_bank_name?: string|null, installment_bank_display_name?: string|null}
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findPayment(array $args): array
    {
        return $this->client->call('find-payment', $args);
    }

    /**
     * ## Find a payment link
     *
     * Retrieve details about a specific payment link in your merchant account.
     * @param array{payment_link_id: string} $args
     * @return array{payment_link_id: string, batch_id?: string|null, payment_link_url: string, status: 'active'|'disabled'|'expired'|'paid', amount: int|float, currency: string, interest_bearer?: 'customer'|'merchant'|null, description: string, external_id?: string|null, expires_at?: int|float|null, single_use: bool, capture: bool, send_confirmation_email: bool, installation: InstallationAttributes, created_at: int|float, deletable: bool}
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment link with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment link ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findPaymentLink(array $args): array
    {
        return $this->client->call('find-payment-link', $args);
    }

    /**
     * ## Find a payout
     *
     * Retrieves a specific payout by its unique identifier or external ID.
     * @param array{payout_id: string} $args
     * @return array{payout_id: string, external_id?: string|null, status: 'SUCCESSFUL'|'SENT'|'PENDING'|'FAILED'|'CANCELLED'|'REQUESTED', amount: int|float, currency: string, created_at: int|float, paid_at?: int|float|null, sent_at?: int|float|null, bank_account_id?: string|null, bank_code?: string|null, bank_account_number?: string|null, bank_account_name?: string|null, bank_country_code?: string|null, payout_method?: string|null, mobile_number?: string|null, national_id?: string|null, fees: list<FeeAttributes>, ref1?: string|null, needs_authorization: true|false, requested_by_user_id?: string|null, requested_by_user_email?: string|null, authorized_by_user_id?: string|null, authorized_by_user_email?: string|null, rejected_by_user_id?: string|null, rejected_by_user_email?: string|null}
     * @throws BadRequestError NOT_FOUND - This error occurs when no payout with the provided ID exists in your
     * merchant account. Double check that you're using the correct payout ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findPayout(array $args): array
    {
        return $this->client->call('find-payout', $args);
    }

    /**
     * ## Find a refund
     *
     * Retrieves a specific refund by its unique identifier or external ID.
     * @param array{refund_id: string} $args
     * @return array{refund_id: string, external_id?: string|null, payment_id: string, status: string, amount: int|float, currency: string, created_at: int|float, refunded_at?: int|float|null, full: true|false}
     * @throws BadRequestError NOT_FOUND - This error occurs when no refund with the provided ID exists in your
     * merchant account. Double check that you're using the correct refund ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findRefund(array $args): array
    {
        return $this->client->call('find-refund', $args);
    }

    /**
     * ## Find a terminal
     *
     * Find a terminal by ID and return its details.
     * @param array{terminal_id: string} $args
     * @return array{terminal_id: string, connected: true|false, external_id?: string|null, installation_id?: string|null, name: string, branch_name?: string|null, enabled: true|false, kiosk_mode: true|false, last_verified_at?: int|float|null, connected_at?: int|float|null, delegated_terminal_id?: string|null, parent_terminal_id?: string|null, challenge: string}
     * @throws BadRequestError NOT_FOUND - No terminal exists with the provided ID in this merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findTerminal(array $args): array
    {
        return $this->client->call('find-terminal', $args);
    }

    /**
     * ## Find a topup
     *
     * Retrieve details about a specific balance top-up in your merchant account.
     * @param array{topup_id: string} $args
     * @return array{topup_id: string, status: 'PENDING'|'SUCCESSFUL', type: 'MANUAL'|'VIRTUAL_ACCOUNT', amount: int|float, currency: string, transaction_id?: string|null, sender_account_name?: string|null, sender_account_number?: string|null, sender_bank_code?: string|null, virtual_account_number?: string|null, created_at: int|float}
     * @throws BadRequestError NOT_FOUND - This error occurs when no topup with the provided ID exists in your
     * merchant account. Double check that you're using the correct topup ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findTopup(array $args): array
    {
        return $this->client->call('find-topup', $args);
    }

    /**
     * ## Find a virtual account
     *
     * Retrieves a specific virtual account by its unique identifier.
     * @param array{virtual_account_id: string} $args
     * @return array{virtual_account_id: string, state: 'pending'|'opening'|'active'|'closing'|'failed'|'closed', number?: string|null, bank_name: string, display_name: string, currency: string, created_at: int|float}
     * @throws BadRequestError NOT_FOUND - This error occurs when no virtual account with the provided ID exists in your
     * merchant account. Double check that you're using the correct virtual account ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findVirtualAccount(array $args): array
    {
        return $this->client->call('find-virtual-account', $args);
    }

    /**
     * ## Find a webhook
     *
     * Retrieves details about a specific webhook in merchant account.
     * @param array{webhook_id: string} $args
     * @return array{webhook_id: string, name: string, enabled: bool, api_version: string, destination_url: string, protocol: 'HTTP'|'HTTPS', signature_type: 'SHA256'|'HMAC_SHA256'|'ED25519', verification_key: string, events: 'PAYMENT_COMPLETED'|'PAYMENT_VOIDED'|'PAYMENT_REVERSED'|'PAYMENT_CANCELLED'|'PAYMENT_AUTHORIZED'|'PAYMENT_EXPIRED'|'PAYOUT_COMPLETED'|'TOPUP_SETTLED', created_at: int|float}
     * @throws BadRequestError NOT_FOUND - This error occurs when no webhook with the provided ID exists in your
     * merchant account. Double check that you're using the correct webhook ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findWebhook(array $args): array
    {
        return $this->client->call('find-webhook', $args);
    }

    /**
     * ## Forget card
     *
     * Forget (soft delete) a card from a contact. The card will be marked as deleted but the record will be preserved for audit purposes.
     *
     * Cards include a `deletable` attribute that indicates whether the card can be deleted. The frontend should check this attribute before attempting to delete a card. If a card cannot be deleted, set another card as the default using the `set-default-card` endpoint, which will update the `deletable` attribute accordingly.
     *
     * The endpoint returns `true` when the card is successfully forgotten.
     * @param array{contact_id: string, card_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id, or the card with the provided card_id does not exist or does not belong to the contact.
     * @throws BadRequestError INVALID_OPERATION - Cannot delete primary card when contact has multiple cards. Set another card as primary first.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function forgetCard(array $args): bool
    {
        return $this->client->call('forget-card', $args);
    }

    /**
     * ## Generate branding upload URL
     *
     * Generates a pre-signed URL for uploading branding assets (logos, sidebar images) to the merchant's
     * storage. Returns an object containing the URL and the object key (which includes the timestamped filename).
     * @param array{filename: string, type: 'light_logo'|'dark_logo'|'sidebar_image', md5_hash: string, byte_size: int|float} $args
     * @return array{url: string, upload_id: string, headers: array<string, mixed>}
     * @throws BadRequestError NOT_FOUND - This error occurs when the merchant branding is not found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function generateBrandingUploadUrl(array $args): array
    {
        return $this->client->call('generate-branding-upload-url', $args);
    }

    /**
     * ## Generate dispute evidence upload URL
     *
     * Generates a pre-signed URL for uploading evidence for a dispute. Returns an object containing
     * the URL and the object key.
     * @param array{dispute_id: string, filename: string, md5_hash: string, byte_size: int|float, comment?: string} $args
     * @return array{url: string, upload_id: string, headers: array<string, mixed>}
     * @throws BadRequestError NOT_FOUND - The specified dispute was not found.
     * @throws BadRequestError INVALID_STATE - The dispute is not in a state that allows evidence upload.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function generateDisputeEvidenceUploadUrl(array $args): array
    {
        return $this->client->call('generate-dispute-evidence-upload-url', $args);
    }

    /**
     * ## Generate Onboarding Upload URL
     * Generates a pre-signed URL for uploading onboarding documents to the merchant's storage.
     * Returns an object containing the URL and the object key (which includes the timestamped filename).
     * @param array{filename: string, type: string, md5_hash: string, byte_size: int|float} $args
     * @return array{url: string, upload_id: string, headers: array<string, mixed>}
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function generateOnboardingUploadUrl(array $args): array
    {
        return $this->client->call('generate-onboarding-upload-url', $args);
    }

    /**
     * ## Generate a temporary terminal challenge
     *
     * Generate a unique 8-digit challenge for a terminal that is not connected and is enabled.
     * The challenge expires 5 minutes after it is issued. Calling this again before
     * expiry returns the same challenge.
     * @param array{terminal_id: string} $args
     * @return array{temporary_challenge: string, expires_at: int|float}
     * @throws BadRequestError NOT_FOUND - No terminal exists with the provided ID in this merchant account.
     * @throws BadRequestError NOT_ENABLED - The terminal exists but is disabled.
     * @throws BadRequestError ALREADY_CONNECTED - The terminal is already connected to a device.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function generateTemporaryTerminalChallenge(array $args): array
    {
        return $this->client->call('generate-temporary-terminal-challenge', $args);
    }

    /**
     * ## Get merchant balance
     *
     * Returns the current balance information for the authenticated merchant.
     *
     * The response includes the current balance amount and currency.
     * @return array{amount: int|float, currency: string}
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getBalance(): array
    {
        return $this->client->call('get-balance');
    }

    /**
     * ## Initiate the payment flow
     *
     * Start a new payment flow using a payment session from our SDK. This endpoint
     * creates a payment and returns a unique payment ID that you can use to track
     * its status.
     *
     * The payment flow typically involves:
     *
     * 1. Your customer selecting a payment method in the SDK
     * 2. You initiating the payment with the session ID from step 1
     * 3. The customer completing the payment via their selected payment method
     * 4. Our system notifying you of the final payment status via webhook, or you
     *    polling the status of the payment
     *
     * The payment ID returned by this endpoint can be used to track the payment status
     * and details.
     * @param array{external_id?: string, amount: string, currency: string, payment_session_id: string, card_id?: string, capture: bool, return_url: string, customer_id?: string, customer_email?: string, customer_mobile?: string, customer_name?: string, save_card?: string, interest_bearer?: 'customer'|'merchant'} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - The payment session ID provided doesn't match any active session in
     * our system, or the card_id doesn't match any card belonging to the
     * customer. This usually means the session has expired or was
     * already used, or the card has been deleted. The customer should
     * select their payment method again in the SDK to generate a new
     * session.
     * @throws BadRequestError CONFLICT - The payment session amount and currency do not match the amount and
     * currency provided in the request. The payment session was created
     * with specific values that must match when initiating the payment flow.
     * Ensure the amount and currency in your request match the values used
     * when creating the payment session. The payment external_id must also
     * be unique within the merchant account.
     * @throws BadRequestError NOT_SUPPORTED - The payment method selected does not support the requested operation.
     * This occurs when capture is set to false but the payment method's
     * integrations do not support authorization without capture. Choose a
     * different payment method or set capture to true.
     * @throws BadRequestError CUSTOMER_REQUIRED - A customer identifier is required when the payment session uses a saved
     * card token. Provide customer_id, customer_email, or customer_mobile to
     * identify the card's owner.
     * @throws BadRequestError CUSTOMER_MISMATCH - The customer identifier provided does not match the saved card on the
     * payment session. When the session was created with a saved card token,
     * the customer_id or customer_email must belong to the card's owner.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function initiatePaymentFlow(array $args): string
    {
        return $this->client->call('initiate-payment-flow', $args);
    }

    /**
     * ## Link bank account to a contact
     *
     * Link a bank account to a contact. This creates a new bank account record associated with the contact.
     *
     * - The bank account is added to the contact's list of bank accounts
     * - If the contact has no bank account, this new account will automatically become the primary one
     * - If the `primary` parameter is set to `true`, this account will be set as the primary bank account (replacing any existing primary account)
     * - If the contact already has a primary bank account and `primary` is not set, this new account will be added as a secondary account
     *
     * The endpoint returns a unique identifier for the bank account that you can use to reference it in other API calls.
     * @param array{contact_id: string, bank_code: string, number: string, name: string, country_code: string, primary?: bool, external_id?: string} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id.
     * @throws BadRequestError COUNTRY_MISMATCH - The bank account's country code does not match your payment provider's country.
     * @throws BadRequestError INVALID_BANK_CODE - The bank code is not valid for Thailand. Bank codes must be 3-digit codes assigned by the Bank of Thailand.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function linkBankAccount(array $args): string
    {
        return $this->client->call('link-bank-account', $args);
    }

    /**
     * ## List activity logs
     *
     * Returns a paginated list of merchant-scoped activity for supported transaction,
     * billing, and team events. Each row includes actor information, the parent resource,
     * and optional plan or refund identifiers.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the activity history.
     *
     * Use the search arguments to narrow results by activity type, actor, object,
     * parent resource, or creation date.
     * @param array{per_page?: int|float, after?: string, before?: string, activity_type?: 'communication_channel_activated'|'dispute'|'dispute_challenged'|'dispute_conceded'|'dispute_escalated'|'invitation_accepted'|'payment_cancelled'|'payment_captured'|'payment_reversed'|'payment_voided'|'payout'|'payout_authorized'|'payout_created'|'payout_rejected'|'payout_requested'|'plan_switched'|'refund_created'|'user_invited'|'user_removed', actor_email?: string, actor_key_id?: string, object_id?: string, parent_id?: string, created_at?: string, terms?: string} $args
     * @return ListActivityLogsPage
     * @throws BadRequestError NOT_FOUND - The actor_email or actor_key_id does not belong to this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listActivityLogs(array $args = []): ListActivityLogsPage
    {
        $raw = $this->client->call('list-activity-logs', $args);
        /** @var Page $raw */
        return new ListActivityLogsPage($raw);
    }

    /**
     * ## List API request logs
     *
     * Returns a paginated list of HTTP requests recorded for this live merchant.
     * Each row includes the endpoint, HTTP status, client IP address, request source,
     * dashboard user email when applicable, and client or browser identification.
     * @param array{per_page?: int|float, after?: string, before?: string, status?: 200|400|500, source?: 'API'|'DASHBOARD', endpoint?: string, user_email?: string, ip_address?: string, created_at?: string, terms?: string} $args
     * @return ListApiRequestLogsPage
     * @throws BadRequestError FORBIDDEN - API request logs are only available for live merchants.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listApiRequestLogs(array $args = []): ListApiRequestLogsPage
    {
        $raw = $this->client->call('list-api-request-logs', $args);
        /** @var Page $raw */
        return new ListApiRequestLogsPage($raw);
    }

    /**
     * ## List batches
     *
     * Returns a paginated list of batches for the authenticated merchant, newest first.
     *
     * Optionally filter by `type` to return only batches of that kind (for example `PAYMENT_LINK`).
     * @param array{type?: 'PAYMENT_LINK'|'PAYOUT', per_page?: int|float, after?: string, before?: string} $args
     * @return ListBatchesPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listBatches(array $args = []): ListBatchesPage
    {
        $raw = $this->client->call('list-batches', $args);
        /** @var Page $raw */
        return new ListBatchesPage($raw);
    }

    /**
     * ## List communication channels
     *
     * Returns the communication channels configured for this merchant.
     *
     * This endpoint requires authentication.
     * @return list<CommunicationChannelAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listCommunicationChannels(): array
    {
        return $this->client->call('list-communication-channels');
    }

    /**
     * ## List contacts
     *
     * Returns a paginated list of all contacts for the authenticated merchant.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListContactsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listContacts(array $args = []): ListContactsPage
    {
        $raw = $this->client->call('list-contacts', $args);
        /** @var Page $raw */
        return new ListContactsPage($raw);
    }

    /**
     * ## List disputes
     *
     * Returns a paginated list of all disputes for the authenticated merchant.
     * Responds with a cursor-paginated page object containing `page`, `next`, and `previous` cursors.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListDisputesPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listDisputes(array $args = []): ListDisputesPage
    {
        $raw = $this->client->call('list-disputes', $args);
        /** @var Page $raw */
        return new ListDisputesPage($raw);
    }

    /**
     * ## List disputes by payment
     *
     * Returns a paginated list of disputes for a specific payment.
     *
     * The `payment_id` is the payment UID (for example `pay_...`).
     *
     * The response uses cursor-based pagination with `previous` and `next` cursors.
     * @param array{payment_id: string, per_page?: int|float, after?: string, before?: string} $args
     * @return ListDisputesByPaymentPage
     * @throws BadRequestError NOT_FOUND - No payment found with the provided payment_id for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listDisputesByPayment(array $args): ListDisputesByPaymentPage
    {
        $raw = $this->client->call('list-disputes-by-payment', $args);
        /** @var Page $raw */
        return new ListDisputesByPaymentPage($raw);
    }

    /**
     * ## List email subscribers
     *
     * Retrieve all email subscription objects for this merchant. Each object contains the email address,
     * language, and all subscribed notification topics for that email.
     *
     * This endpoint requires authentication and will return an array of email subscription objects,
     * where each object represents one email address with all its notification topic subscriptions.
     * @return list<EmailRecipientAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listEmailSubscribers(): array
    {
        return $this->client->call('list-email-subscribers');
    }

    /**
     * ## List all SDK installations
     *
     * Retrieve a list of all SDK installations in your merchant account. This endpoint returns an array of
     * installation objects, each containing details about:
     *
     * - Name and platform (web, iOS, Android or Terminal)
     * - Unique identifier (domain, bundle ID or package name)
     * - Enabled payment methods
     * - Current status (enabled or disabled)
     *
     * See `find-installation` for the complete list of response attributes.
     *
     * Common use cases for this endpoint:
     *
     * - Auditing which installations are currently active
     * - Checking the configuration of multiple installations
     * - Finding installations that need to be updated or disabled
     * - Getting an overview of your integration across different platforms
     * @return list<InstallationAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listInstallations(): array
    {
        return $this->client->call('list-installations');
    }

    /**
     * ## List invoices
     *
     * Returns a paginated list of invoices for the authenticated merchant.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the invoice history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListInvoicesPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listInvoices(array $args = []): ListInvoicesPage
    {
        $raw = $this->client->call('list-invoices', $args);
        /** @var Page $raw */
        return new ListInvoicesPage($raw);
    }

    /**
     * ## List all API keys
     *
     * Retrieve a list of all API keys associated with your merchant account. This endpoint returns an array of
     * key objects containing details about each key, including:
     *
     * - Key name and ID
     * - Creation time
     * - Current status (enabled or disabled)
     *
     * See `generate-key` or `find-key` for the complete list of response attributes.
     *
     * Common use cases for this endpoint:
     *
     * - Auditing active API keys
     * - Identifying unused or outdated keys that should be rotated
     *
     * Note: For security reasons, this endpoint only returns metadata about your keys - it never returns the
     * actual key values. Make sure to securely store your key values when they are first generated, as they
     * cannot be retrieved later.
     * @return list<KeyAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listKeys(): array
    {
        return $this->client->call('list-keys');
    }

    /**
     * ## List metadata
     *
     * Returns a paginated list of all metadata records stored for your merchant account.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListMetadataPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listMetadata(array $args = []): ListMetadataPage
    {
        $raw = $this->client->call('list-metadata', $args);
        /** @var Page $raw */
        return new ListMetadataPage($raw);
    }

    /**
     * ## List monthly billing archives
     *
     * Returns a paginated list of monthly billing archives for the authenticated merchant.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the payment history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListMonthlyBillingArchivesPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listMonthlyBillingArchives(array $args = []): ListMonthlyBillingArchivesPage
    {
        $raw = $this->client->call('list-monthly-billing-archives', $args);
        /** @var Page $raw */
        return new ListMonthlyBillingArchivesPage($raw);
    }

    /**
     * ## List payment activities
     *
     * Retrieves all activities, events, and payment session information for a given payment.
     * @param array{payment_id: string} $args
     * @return list<PaymentActivityAttributes>
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPaymentActivities(array $args): array
    {
        return $this->client->call('list-payment-activities', $args);
    }

    /**
     * ## List payments for a payment link
     *
     * Returns a paginated list of payments associated with a specific payment link.
     * @param array{payment_link_id: string, per_page?: int|float, after?: string, before?: string} $args
     * @return ListPaymentLinkPaymentsPage
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment link with the provided ID exists in your
     * merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPaymentLinkPayments(array $args): ListPaymentLinkPaymentsPage
    {
        $raw = $this->client->call('list-payment-link-payments', $args);
        /** @var Page $raw */
        return new ListPaymentLinkPaymentsPage($raw);
    }

    /**
     * ## List payment links
     *
     * Returns a paginated list of payment links for the authenticated merchant.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the payment link history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListPaymentLinksPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPaymentLinks(array $args = []): ListPaymentLinksPage
    {
        $raw = $this->client->call('list-payment-links', $args);
        /** @var Page $raw */
        return new ListPaymentLinksPage($raw);
    }

    /**
     * ## List payments
     *
     * Returns a paginated list of payments for the authenticated merchant.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the payment history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListPaymentsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPayments(array $args = []): ListPaymentsPage
    {
        $raw = $this->client->call('list-payments', $args);
        /** @var Page $raw */
        return new ListPaymentsPage($raw);
    }

    /**
     * ## List payments by contact
     *
     * Returns a paginated list of payments for a specific contact belonging to the authenticated merchant.
     *
     * The `contact_id` can be either the internal UID (starting with `ctc_`) or the external ID
     * provided when the contact was created.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the payment history for the contact.
     * @param array{contact_id: string, per_page?: int|float, after?: string, before?: string} $args
     * @return ListPaymentsByContactPage
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id. The contact may not exist or may belong to a different merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPaymentsByContact(array $args): ListPaymentsByContactPage
    {
        $raw = $this->client->call('list-payments-by-contact', $args);
        /** @var Page $raw */
        return new ListPaymentsByContactPage($raw);
    }

    /**
     * ## List payments by terminal
     *
     * Returns a paginated list of payments for a specific terminal, scoped by each payment's
     * `created_on` date in the merchant's time zone.
     *
     * When `created_at` is omitted, results are limited to payments from the last three months
     * through today (same window as list-payments-by-contact).
     *
     * The `terminal_id` is the terminal UID (for example `trm_...`).
     *
     * The response uses cursor-based pagination with `previous` and `next` cursors.
     * @param array{terminal_id: string, created_at?: string, per_page?: int|float, after?: string, before?: string} $args
     * @return ListPaymentsByTerminalPage
     * @throws BadRequestError NOT_FOUND - No terminal found with the provided terminal_id for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPaymentsByTerminal(array $args): ListPaymentsByTerminalPage
    {
        $raw = $this->client->call('list-payments-by-terminal', $args);
        /** @var Page $raw */
        return new ListPaymentsByTerminalPage($raw);
    }

    /**
     * ## List payout activities
     *
     * Retrieves all activities, events, and payout session information for a given payout.
     * @param array{payout_id: string} $args
     * @return list<array{created_at: int|float, activity: 'payout_authorized'|'payout_completed'|'payout_created'|'payout_rejected'|'payout_requested', path?: string|null, ip_address?: string|null, severity: 'INFO'|'WARN'|'ERROR', request?: array<string, mixed>|null, response?: mixed|null, event_deliveries?: list<EventDeliveryAttributes>|null, payout_status?: string|null}>
     * @throws BadRequestError NOT_FOUND - This error occurs when no payout with the provided ID exists in your
     * merchant account. Double check that you're using the correct payout ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPayoutActivities(array $args): array
    {
        return $this->client->call('list-payout-activities', $args);
    }

    /**
     * ## List payouts
     *
     * Returns a paginated list of payouts for the authenticated merchant.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the payout history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListPayoutsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPayouts(array $args = []): ListPayoutsPage
    {
        $raw = $this->client->call('list-payouts', $args);
        /** @var Page $raw */
        return new ListPayoutsPage($raw);
    }

    /**
     * ## List payouts by contact
     *
     * Returns a paginated list of payouts for a specific contact belonging to the authenticated merchant.
     *
     * The `contact_id` can be either the internal UID (starting with `ctc_`) or the external ID
     * provided when the contact was created.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the payout history for the contact.
     * @param array{contact_id: string, per_page?: int|float, after?: string, before?: string} $args
     * @return ListPayoutsByContactPage
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id. The contact may not exist or may belong to a different merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPayoutsByContact(array $args): ListPayoutsByContactPage
    {
        $raw = $this->client->call('list-payouts-by-contact', $args);
        /** @var Page $raw */
        return new ListPayoutsByContactPage($raw);
    }

    /**
     * ## List all permissions
     *
     * Retrieves a list of all permissions available for roles within the merchant account.
     * These permissions can be assigned to custom roles created by the merchant.
     * @return list<string>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPermissions(): array
    {
        return $this->client->call('list-permissions');
    }

    /**
     * ## List planned maintenances
     *
     * Retrieve recent and upcoming planned maintenance windows for our services.
     * @param array{language?: 'en'|'th'} $args
     * @return list<PlannedMaintenanceAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPlannedMaintenances(array $args = []): array
    {
        return $this->client->call('list-planned-maintenances', $args);
    }

    /**
     * ## List all current plans
     * Retrieve a list of all current plans available for the merchant.
     * Plans represent different levels of service or pricing plans currently available.
     * The `selectable` flag indicates if a plan can be subscribed to via self-service.
     * @return list<PlanAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPlans(): array
    {
        return $this->client->call('list-plans');
    }

    /**
     * ## List refunds by payment
     *
     * Returns a paginated list of refunds for a specific payment belonging to the authenticated merchant.
     *
     * The payment_id can be either the internal UID (starting with 'pay_') or the external ID
     * provided when the payment was created.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the refund history for the payment.
     * @param array{payment_id: string, per_page?: int|float, after?: string, before?: string} $args
     * @return ListRefundsByPaymentPage
     * @throws BadRequestError NOT_FOUND - No payment found with the provided payment_id.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listRefundsByPayment(array $args): ListRefundsByPaymentPage
    {
        $raw = $this->client->call('list-refunds-by-payment', $args);
        /** @var Page $raw */
        return new ListRefundsByPaymentPage($raw);
    }

    /**
     * ## List all roles
     *
     * Retrieve a list of all roles available in the merchant account.
     * This endpoint returns the roles that can be assigned to users within the merchant account.
     * This includes default roles defined by the system as well as any custom roles created by the merchant.
     * @return list<RoleAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listRoles(): array
    {
        return $this->client->call('list-roles');
    }

    /**
     * ## List statements
     *
     * Returns a paginated list of statements for the authenticated merchant grouped by month.
     *
     * The response includes month-based pagination where each page contains all statements
     * for a specific month.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through months.
     * @param array{month?: string} $args
     * @return ListStatementsPage
     * @throws BadRequestError invalid_month - Invalid month format. Expected YYYY-MM format (e.g., '2024-12')
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listStatements(array $args = []): ListStatementsPage
    {
        $raw = $this->client->call('list-statements', $args);
        /** @var Page $raw */
        return new ListStatementsPage($raw);
    }

    /**
     * ## List terminals
     *
     * Paginated list of terminals for the merchant.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListTerminalsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listTerminals(array $args = []): ListTerminalsPage
    {
        $raw = $this->client->call('list-terminals', $args);
        /** @var Page $raw */
        return new ListTerminalsPage($raw);
    }

    /**
     * ## List topups
     *
     * Returns a paginated list of balance top-ups for the authenticated merchant.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the topup history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListTopupsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listTopups(array $args = []): ListTopupsPage
    {
        $raw = $this->client->call('list-topups', $args);
        /** @var Page $raw */
        return new ListTopupsPage($raw);
    }

    /**
     * ## List all user invitations
     *
     * Retrieve a list of all email addresses that have been invited to join your merchant account.
     * This endpoints only returns invitations that have not yet been accepted by the invited people.
     * @return list<InvitationAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listUserInvitations(): array
    {
        return $this->client->call('list-user-invitations');
    }

    /**
     * ## List all users
     *
     * Retrieve a list of all users who are members of your merchant account.
     * This endpoint returns an array of user objects containing details about each user.
     *
     * Common use cases for this endpoint:
     *
     * - Auditing user access to your merchant account
     * - Identifying users who can perform actions on behalf of the merchant
     * - Managing user permissions and roles within the merchant account
     * @return list<MerchantMembershipAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listUsers(): array
    {
        return $this->client->call('list-users');
    }

    /**
     * ## List virtual accounts
     *
     * Returns a paginated list of the authenticated merchant's active virtual accounts.
     *
     * The response includes cursor-based pagination with `previous` and `next` cursors
     * that can be used to navigate through the virtual account history.
     * @param array{per_page?: int|float, after?: string, before?: string} $args
     * @return ListVirtualAccountsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listVirtualAccounts(array $args = []): ListVirtualAccountsPage
    {
        $raw = $this->client->call('list-virtual-accounts', $args);
        /** @var Page $raw */
        return new ListVirtualAccountsPage($raw);
    }

    /**
     * ## List webhooks
     *
     * Returns an array of all webhook endpoints for the merchant.
     *
     * See `find-webhook` for the complete list of response attributes.
     * @return list<WebhookAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listWebhooks(): array
    {
        return $this->client->call('list-webhooks');
    }

    /**
     * ## Pause a webhook endpoint
     *
     * Temporarily disable a webhook connection.
     * @param array{webhook_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when the webhook endpoint with the provided ID could not be found.
     * Double check that you're using the correct webhook ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function pauseWebhook(array $args): bool
    {
        return $this->client->call('pause-webhook', $args);
    }

    /**
     * ## Reconfigure SDK installation
     *
     * Modify the configuration of an existing SDK installation in your merchant account. This endpoint allows you to:
     *
     * - Change the installation name
     * - Update which payment methods are enabled
     * - Configure email and mobile capture settings
     *
     * Common use cases for this endpoint:
     *
     * - Rebranding your application with a new name
     * - Adding or removing payment methods as your business needs change
     * - Standardizing installation names across platforms
     * - Adjusting payment method configuration after testing
     * - Updating customer data collection requirements
     *
     * The changes take effect immediately. Any subsequent API calls from this installation will use the new configuration.
     *
     * To see which payment methods are available for your merchant account, use the `supported-payment-methods` endpoint.
     *
     * Note: This endpoint only updates the installation's configuration - it does not affect the installation's enabled/disabled status.
     * Use the dedicated enable/disable endpoints to change an installation's status.
     * @param array{installation_id: string, name: string, payment_methods: string, email_capture?: 'disabled'|'optional'|'required', mobile_capture?: 'disabled'|'optional'|'required', address_capture?: 'disabled'|'optional'|'required'} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The installation with the provided ID could not be found in your merchant account. Verify that you're using the correct installation ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function reconfigureInstallation(array $args): bool
    {
        return $this->client->call('reconfigure-installation', $args);
    }

    /**
     * ## Reconnect a webhook endpoint
     *
     * Update an existing webhook endpoint connection details.
     * @param array{webhook_id: string, name: string, protocol: 'HTTP'|'HTTPS', signature_type?: 'SHA256'|'HMAC_SHA256'|'ED25519', api_version: string, destination_url: string, events: 'PAYMENT_COMPLETED'|'PAYMENT_VOIDED'|'PAYMENT_REVERSED'|'PAYMENT_CANCELLED'|'PAYMENT_AUTHORIZED'|'PAYMENT_EXPIRED'|'PAYOUT_COMPLETED'|'TOPUP_SETTLED'} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when the webhook endpoint with the provided ID could not be found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function reconnectWebhook(array $args): bool
    {
        return $this->client->call('reconnect-webhook', $args);
    }

    /**
     * ## Refund a payment
     *
     * Creates a refund for a captured payment that has been billed.
     * Supports partial refunds by specifying an amount less than the
     * remaining refundable amount. Defaults to full refund if amount
     * is not specified.
     * @param array{payment_id: string, amount?: int|float, external_id?: string} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - No payment with the provided ID exists in your merchant account.
     * @throws BadRequestError NOT_REFUNDABLE - Payment is not eligible for refund. This can happen if the payment
     * has not been billed yet, is already fully refunded, or the integration
     * does not support refunds.
     * @throws BadRequestError AMOUNT_EXCEEDS_LIMIT - The refund amount exceeds the remaining refundable amount.
     * @throws BadRequestError CONFLICT - A refund with the provided external ID already exists in your merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function refundPayment(array $args): string
    {
        return $this->client->call('refund-payment', $args);
    }

    /**
     * ## Remove branding document
     *
     * Removes a branding asset (logo or sidebar image) that has been uploaded and confirmed. This is
     * useful if you need to replace an asset or if an asset was uploaded in error.
     * @param array{upload_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no upload with the provided ID exists in your merchant branding.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeBrandingDocument(array $args): bool
    {
        return $this->client->call('remove-branding-document', $args);
    }

    /**
     * ## Remove contact
     *
     * Removes a contact by its ID. The `contact_id` can be either the internal UID (starting with `ctc_`)
     * or the external ID provided when the contact was created.
     * @param array{contact_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id.
     * The contact may not exist or may belong to a different merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeContact(array $args): bool
    {
        return $this->client->call('remove-contact', $args);
    }

    /**
     * ## Remove from metadata set
     *
     * Removes an item from a metadata set. A metadata set is an array of unique values. This endpoint
     * removes a specific value from the set at the given path.
     * - If the value is not found in the set, an error is returned.
     * - The metadata at the path must be an array.
     *
     * #### Example
     * If current metadata is `{"tags": ["a", "b"]}`, removing `"a"` from path `"tags"` results in `{"tags": ["b"]}`.
     * Attempting to remove `"c"` results in an error.
     * @param array{key: string, path?: string, value: mixed, version?: int|float} $args
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeFromMetadataSet(array $args): bool
    {
        return $this->client->call('remove-from-metadata-set', $args);
    }

    /**
     * ## Remove in metadata
     *
     * Removes a specific nested `path` from a metadata object. This is useful for cleaning up specific
     * parts of a complex metadata structure without overwriting the entire object.
     * @param array{key: string, path: string, version?: int|float} $args
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeInMetadata(array $args): bool
    {
        return $this->client->call('remove-in-metadata', $args);
    }

    /**
     * ## Remove an SDK installation
     *
     * Permanently delete an SDK installation from your merchant account. This will:
     *
     * - Remove the installation configuration
     * - Invalidate any API calls from this installation
     *
     * Common use cases for this endpoint:
     *
     * - Cleaning up old or unused installations
     * - Decommissioning applications that have been replaced
     *
     * Warning: This action cannot be undone. Once an installation is removed, you'll need to create a new
     * installation if you want to restore access. Consider using the `disable-installation` endpoint instead
     * if you only need to temporarily prevent API access.
     *
     * The endpoint returns true if the installation was successfully removed.
     * @param array{installation_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The installation with the provided ID could not be found in your merchant account. Verify that you're using the correct installation ID.
     * @throws BadRequestError CANNOT_BE_DELETED - This installation cannot be removed as it is used by Terminals or Payment Links.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeInstallation(array $args): bool
    {
        return $this->client->call('remove-installation', $args);
    }

    /**
     * ## Remove a key person
     *
     * Permanently remove a key person from your business merchant account during the onboarding phase. This
     * endpoint allows you to remove a key person's details from your business profile while you are still
     * setting up your account.
     *
     * You may need this endpoint if you made a mistake while entering key person details during onboarding,
     * or if there were last-minute changes to your business's ownership structure before completing the
     * registration process.
     *
     * Note that this endpoint only works for merchant accounts that are registered as businesses and only during
     * the initial onboarding phase. Individual merchant accounts cannot use this endpoint, and you cannot
     * remove key persons after completing onboarding.
     *
     * This action cannot be undone. Once a key person is removed, you'll need to add them as a
     * new key person if you want to restore their details.
     *
     * The endpoint returns true if the key person was successfully removed.
     * @param array{key_person_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The key person with the provided ID could not be found in your merchant account. Verify that you're using the correct key person ID.
     * @throws BadRequestError NOT_BUSINESS - This merchant account is not registered as a business. Only business merchant accounts can manage key persons.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeKeyPerson(array $args): bool
    {
        return $this->client->call('remove-key-person', $args);
    }

    /**
     * ## Remove an Onboarding Document
     *
     * This endpoint allows you to remove an onboarding document that has been uploaded
     * during the onboarding process. This is useful if you need to replace a document or
     * if a document was uploaded in error.
     * @param array{upload_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no upload with the provided ID exists in your merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeOnboardingDocument(array $args): bool
    {
        return $this->client->call('remove-onboarding-document', $args);
    }

    /**
     * ## Remove a payment link
     *
     * Soft-delete a payment link. The link is hidden from listings and lookups but the
     * underlying record and any associated payments are preserved for historical and
     * audit purposes.
     * @param array{payment_link_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No payment link with the provided ID exists in your merchant account.
     * @throws BadRequestError NOT_REMOVABLE - The payment link has one or more associated payments and cannot be removed.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removePaymentLink(array $args): bool
    {
        return $this->client->call('remove-payment-link', $args);
    }

    /**
     * ## Remove a sales channel
     *
     * Remove a sales channel for the merchant.
     * @param array{sales_channel_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no sales channel with the provided ID exists in your
     * merchant account.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeSalesChannel(array $args): bool
    {
        return $this->client->call('remove-sales-channel', $args);
    }

    /**
     * ## Remove a terminal
     *
     * Remove a terminal from the merchant account by performing a soft delete.
     * @param array{terminal_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No terminal exists with the provided ID in this merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function removeTerminal(array $args): bool
    {
        return $this->client->call('remove-terminal', $args);
    }

    /**
     * ## Request a payment
     *
     * Start a card payment directly from your server without requiring customer
     * interaction. This is a merchant-initiated transaction (MIT) endpoint that
     * combines payment method selection and payment initiation into a single call.
     *
     * This endpoint supports card payments only. You can pay using either a saved
     * card ID or a card token from the `tokenize-card` endpoint.
     *
     * Since this is a merchant-initiated transaction, the 3D Secure (3DS)
     * authentication flow is skipped.
     *
     * Your merchant account must have the MIT flag enabled for the card payment
     * method to use this endpoint.
     * @param array{card_id?: string, token?: string, amount: string, currency: string, installation_id: string, capture: bool, external_id?: string, customer_id?: string, customer_email?: string, customer_mobile?: string, customer_name?: string, save_card?: string} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - The installation, card, or token could not be found. Double check that
     * you are using the correct IDs and that the card belongs to the specified
     * customer.
     * @throws BadRequestError NOT_ALLOWED - Merchant initiated transactions are not allowed for this payment method.
     * Contact support to enable the MIT flag on your card payment method.
     * @throws BadRequestError NOT_ENABLED - This installation is currently disabled.
     * @throws BadRequestError TOKEN_ALREADY_USED - This token has already been used. Request a new token from tokenize-card.
     * @throws BadRequestError NOT_SUPPORTED - The payment method does not support the requested operation.
     * This occurs when capture is set to false but the payment method's
     * integrations do not support authorization without capture.
     * @throws BadRequestError CONFLICT - A payment with the same external_id already exists in this merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function requestPayment(array $args): string
    {
        return $this->client->call('request-payment', $args);
    }

    /**
     * ## Reset a terminal
     *
     * Reset the device binding and revoke the current API key for the given terminal.
     * @param array{terminal_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No terminal exists with the provided ID in this merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function resetTerminal(array $args): bool
    {
        return $this->client->call('reset-terminal', $args);
    }

    /**
     * ## Resume payment on terminal
     *
     * Sends an in-progress payment to a child (delegated) terminal for display.
     * The parent terminal must have a `delegated_terminal_id` configured.
     *
     * The child terminal will receive a real-time event and navigate to the QR code screen.
     * @param array{terminal_id: string, payment_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The terminal or payment could not be found.
     * @throws BadRequestError NO_DELEGATION - The terminal has no delegated terminal configured.
     * @throws BadRequestError NO_SESSION - No payment session found for this payment.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function resumePaymentOnTerminal(array $args): bool
    {
        return $this->client->call('resume-payment-on-terminal', $args);
    }

    /**
     * ## Resume a webhook endpoint
     *
     * Enable temporarily disabled webhook connection.
     * @param array{webhook_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when the webhook endpoint with the provided ID could not be found.
     * Double check that you're using the correct webhook ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function resumeWebhook(array $args): bool
    {
        return $this->client->call('resume-webhook', $args);
    }

    /**
     * ## Retrieve branding
     *
     * Returns the current visual branding settings for your merchant account. This includes colors,
     * typography, button styling, display preferences, and uploaded assets (logos, sidebar images) for
     * your payment interface.
     *
     * The response includes all branding configuration options that can be customized through the
     * `update-branding` endpoint, as well as URLs to any uploaded branding assets.
     *
     * Common use cases for this endpoint include:
     *
     * - Checking your current branding configuration
     * - Previewing how your payment interface will look to customers
     * - Verifying that uploaded logos and images are properly configured
     * - Getting branding settings for integration with your own systems
     *
     * This endpoint requires authentication and will only return branding details for the authenticated merchant.
     * @return array{primary_color: string, secondary_color: string, sidebar_color: string, sidebar_background_type: 'color'|'image', tile_sidebar_image: bool, hide_logo: bool, heading_font_family: string, heading_font_variant: 'Thin'|'ExtraLight'|'Light'|'Regular'|'Medium'|'SemiBold'|'Bold'|'ExtraBold'|'Black', body_font_family: string, body_font_variant: 'Thin'|'ExtraLight'|'Light'|'Regular'|'Medium'|'SemiBold'|'Bold'|'ExtraBold'|'Black', use_iso_code: bool, hide_non_significant_zeros: bool, show_amount_in_button: bool, symbol_placement: 'use_default'|'before_amount'|'after_amount', default_button_text: string, button_shape: 'rounded_none'|'rounded_full'|'rounded_xs'|'rounded_sm'|'rounded_md'|'rounded_lg'|'rounded_xl', light_logo_url?: string|null, light_logo_id?: string|null, light_logo_filename?: string|null, light_logo_file_type?: string|null, dark_logo_url?: string|null, dark_logo_id?: string|null, dark_logo_filename?: string|null, dark_logo_file_type?: string|null, sidebar_image_url?: string|null, sidebar_image_id?: string|null, sidebar_image_filename?: string|null, sidebar_image_file_type?: string|null}
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveBranding(): array
    {
        return $this->client->call('retrieve-branding');
    }

    /**
     * ## Retrieve branding fonts
     *
     * Returns a list of available fonts that can be used for branding your payment interface. This endpoint
     * returns a comprehensive list of fonts that are supported for customizing the typography of your
     * payment forms.
     *
     * The response includes:
     *
     * - Font family names (e.g., 'Noto Sans', 'Roboto')
     * - Font categories (sans-serif, serif, monospace, display, handwriting)
     * - Available font weight variants for each font family
     *
     * Common use cases for this endpoint include:
     *
     * - Populating font selection dropdowns in your branding interface
     * - Validating font choices before updating branding settings
     * - Displaying available font options to users
     * - Building font preview functionality
     *
     * This endpoint requires authentication and returns the same font list for all merchants, as fonts are
     * globally available across the platform.
     *
     * **Note:** Font variants are filtered to exclude italic styles and are converted to readable weight names
     * (e.g., 'Bold', 'Medium', 'Light').
     * @return list<array{family: string, category: 'sans-serif'|'serif'|'monospace'|'display'|'handwriting', variants: 'Thin'|'ExtraLight'|'Light'|'Regular'|'Medium'|'SemiBold'|'Bold'|'ExtraBold'|'Black'}>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveBrandingFont(): array
    {
        return $this->client->call('retrieve-branding-font');
    }

    /**
     * ## Retrieve checkout
     *
     * Retrieves the details of an existing checkout session. You can use this to check the current
     * status of a checkout, see which payment methods were used, and get customer information.
     * @param array{checkout_id: string} $args
     * @return array{checkout_id: string, amount: string, currency: string, interest_bearer?: 'customer'|'merchant'|null, installation_id: string, customer_id?: string|null, customer_email?: string|null, customer_mobile?: string|null, customer_name?: string|null, redirect_url: string, payment_id?: string|null}
     * @throws BadRequestError NOT_FOUND - This error occurs when no checkout with the provided ID exists in your
     * merchant account. Double check that you're using the correct checkout ID.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveCheckout(array $args): array
    {
        return $this->client->call('retrieve-checkout', $args);
    }

    /**
     * ## Retrieve merchant configuration
     *
     * Retrieves the read-only system configuration set up for your merchant account by Reservepay.
     *
     * To view or modify self-service parameters that you control directly, see `retrieve-merchant-settings` and `update-merchant-settings`.
     * @return MerchantConfigurationAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveConfiguration(): array
    {
        return $this->client->call('retrieve-configuration');
    }

    /**
     * ## Retrieve contact
     *
     * Retrieves a contact by its ID. The `contact_id` can be either the internal UID (starting with `ctc_`)
     * or the external ID provided when the contact was created.
     * @param array{contact_id: string} $args
     * @return array{contact_id: string, external_id?: string|null, email?: string|null, mobile?: string|null, name?: string|null, favorite: bool, created_at: int|float, primary_bank_account?: BankAccountAttributes|null, default_card?: CardAttributes|null, bank_accounts: list<BankAccountAttributes>, cards: list<CardAttributes>}
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id.
     * The contact may not exist or may belong to a different merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveContact(array $args): array
    {
        return $this->client->call('retrieve-contact', $args);
    }

    /**
     * ## Retrieve current plan
     *
     * Retrieves the current billing plan for the merchant account.
     * This endpoint returns the plan details including the plan ID, name, tier, and any associated metadata.
     * It is used to check the current subscription status and plan features available to the merchant.
     *
     * The response includes information about the plan's tier, which indicates the level of service provided.
     * It also reports `pending_plan_id`, the plan a scheduled change will promote to, or `null` when no change is scheduled.
     * @return CurrentPlanAttributes
     * @throws BadRequestError NO_PLAN_SELECTED - The merchant is not currently subscribed to any plan. A plan must be subscribed to via the `select-plan` endpoint first.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveCurrentPlan(): array
    {
        return $this->client->call('retrieve-current-plan');
    }

    /**
     * ## Retrieve dispute performance
     *
     * Returns dispute performance data for the authenticated merchant based on the specified date range.
     *
     * The endpoint supports various input formats:
     *
     * - Specific date: `YYYY-MM-DD` (e.g., "2024-01-15")
     * - Date range: `YYYY-MM-DD..YYYY-MM-DD` (e.g., "2024-01-01..2024-01-31")
     *
     * The response includes aggregated metrics such as total amount and count of disputes.
     *
     * ### Timezone Handling
     *
     * The service uses the merchant's configured timezone for all date boundaries and grouping.
     * Disputes are grouped based on the merchant's local date/time.
     * @param array{period: string, include_won?: bool, include_lost?: bool, include_in_progress?: bool} $args
     * @return array{start_date: string, end_date: string, amount: int|float, count: int|float, average?: int|float|null, change?: int|float|null, granularity: 'hourly'|'daily'|'weekly'|'monthly', transactions: list<array<string, mixed>>}
     * @throws BadRequestError INVALID_DATE_FORMAT - Invalid date format. Please use YYYY-MM-DD or YYYY-MM-DD..YYYY-MM-DD.
     * @throws BadRequestError INVALID_DATE_RANGE - End date cannot be before start date.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveDisputePerformance(array $args): array
    {
        return $this->client->call('retrieve-dispute-performance', $args);
    }

    /**
     * ## Retrieve email subscriber
     *
     * Retrieve all notification topic subscriptions for a specific email address within this merchant.
     * Returns a single object with the email, language, and all subscribed notification topics.
     *
     * This endpoint requires authentication and will return subscription details for the specified email.
     * @param array{email: string} $args
     * @return EmailRecipientAttributes
     * @throws BadRequestError NOT_FOUND - No subscriber found for the specified email address.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveEmailSubscriber(array $args): array
    {
        return $this->client->call('retrieve-email-subscriber', $args);
    }

    /**
     * ## Retrieve merchant settings
     *
     * Retrieves current self-service settings for your merchant account. These are values you control directly and can modify using `update-merchant-settings`.
     *
     * To view read-only account settings managed by Reservepay, use `retrieve-configuration`.
     * @return MerchantSettingAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveMerchantSettings(): array
    {
        return $this->client->call('retrieve-merchant-settings');
    }

    /**
     * ## Retrieve metadata
     *
     * Retrieves metadata associated with a specific `key`. You can also specify an optional dot-separated
     * `path` to retrieve only a subset of the metadata.
     *
     * Metadata is a flexible way to store additional information for your account, and can be used for
     * storing configuration, preferences, or any other data you need.
     * @param array{key: string, path?: string, data_only?: bool} $args
     * @return MetadataAttributes|mixed
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveMetadata(array $args): mixed
    {
        return $this->client->call('retrieve-metadata', $args);
    }

    /**
     * ## Retrieve monthly billing archive
     *
     * Retrieves a monthly billing archive by its ID (i.e., month and year in `YYYY-MM` format).
     * @param array{monthly_billing_archive_id: string} $args
     * @return array{id: string, url: string}
     * @throws BadRequestError NOT_FOUND - No archive found with the provided monthly_billing_archive_id.
     * This could indicate that there was no activity for that month, or that the archive has not yet been created.
     * @throws BadRequestError INVALID - The provided monthly_billing_archive_id could not be parsed as a valid month (expected format: 'YYYY-MM').
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveMonthlyBillingArchive(array $args): array
    {
        return $this->client->call('retrieve-monthly-billing-archive', $args);
    }

    /**
     * ## Retrieve notification topics
     *
     * Retrieves all available notification topics that merchants can subscribe to. Topics are used to organize
     * different types of notifications such as incidents, payments, developer updates, etc.
     *
     * The response is an array of notification topic objects containing the topic ID, name, and display name.
     *
     * This endpoint requires authentication and returns all available topics that the merchant can potentially
     * subscribe to for email notifications.
     * @return list<NotificationTopicAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveNotificationTopics(): array
    {
        return $this->client->call('retrieve-notification-topics');
    }

    /**
     * ## Retrieve your onboarding details
     *
     * Get comprehensive details about your onboarding. This includes your onboarding ID, display name, legal name,
     * entity type (business or individual), business category, enabled product types, and auto-payout frequency settings.
     * Your primary bank account information is also returned.
     *
     * For business accounts, this endpoint returns additional information about the business entity and its key people.
     * For individual accounts, it includes the individual's details instead.
     *
     * Common use cases for this endpoint include checking your onboarding status, verifying your business or individual details,
     * and ensuring your primary bank account is set up correctly for payouts.
     *
     * This endpoint requires authentication and will only return details for the authenticated merchant.
     * @return OnboardingAttributes
     * @throws BadRequestError NO_ONBOARDING - This error occurs when the merchant has no onboarding available.
     * Ensure that the merchant has started the onboarding process before attempting to retrieve it.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveOnboardingDetails(): array
    {
        return $this->client->call('retrieve-onboarding-details');
    }

    /**
     * ## Retrieve primary bank account
     *
     * Get the merchant's current primary bank account used for payouts. This is the account configured via `set-primary-bank-account` and used for:
     *
     * - Receiving settlement funds from successful payments
     * - Automated payouts based on your payout schedule
     * - Manual payout requests
     * @return array{bank_account_id: string, bank_code: string, number: string, name: string, country_code: string, created_at: int|float, deleted_at?: int|float|null, deletable: bool, owner: 'CONTACT'|'MERCHANT', primary: bool}
     * @throws BadRequestError NOT_FOUND - No primary bank account has been configured for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrievePrimaryBankAccount(): array
    {
        return $this->client->call('retrieve-primary-bank-account');
    }

    /**
     * ## Retrieve refund performance
     *
     * Returns refund performance data for the authenticated merchant based on the specified date range.
     *
     * The endpoint supports various input formats:
     *
     * - Specific date: `YYYY-MM-DD` (e.g., "2024-01-15")
     * - Date range: `YYYY-MM-DD..YYYY-MM-DD` (e.g., "2024-01-01..2024-01-31")
     *
     * The response includes aggregated metrics such as total refunded, transaction count, and average transaction value.
     *
     * - **`exclude_voided`**: The default is `false`, so voided transactions appear alongside your refunds. When `true`, only refund transactions appear.
     *
     * ### Timezone Handling
     *
     * The service uses the merchant's configured timezone for all date boundaries and grouping.
     * Refunds are grouped by `refunded_at` in the merchant's local date/time.
     * For example, a refund at 2025-10-09 22:51:28 UTC would appear in 2025-10-10 data for a merchant in timezone "Asia/Bangkok" (GMT+7).
     *
     * ### Granularity Logic
     *
     * The response granularity is automatically determined based on the date range:
     *
     * - **0 days** (single date): Hourly breakdown
     * - **1-83 days**: Daily breakdown
     * - **84-167 days**: Weekly breakdown
     * - **168+ days**: Monthly breakdown
     *
     * This affects how the `transactions` array is structured and the `granularity` field in the response.
     * @param array{period: string, exclude_voided?: bool} $args
     * @return array{start_date: string, end_date: string, amount: int|float, count: int|float, average?: int|float|null, change?: int|float|null, granularity: 'hourly'|'daily'|'weekly'|'monthly', transactions: list<array<string, mixed>>}
     * @throws BadRequestError INVALID_DATE_FORMAT - This error occurs when the date_input parameter is not in a supported format.
     * Please use one of the supported formats: YYYY-MM-DD, YYYY-MM-DD..YYYY-MM-DD, or keywords.
     * @throws BadRequestError INVALID_DATE_RANGE - This error occurs when the end date is before the start date in a date range.
     * Please ensure the end date is after or equal to the start date.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveRefundPerformance(array $args): array
    {
        return $this->client->call('retrieve-refund-performance', $args);
    }

    /**
     * ## Retrieve sales performance
     *
     * Returns sales performance data for the authenticated merchant based on the specified date range.
     *
     * The endpoint supports various input formats:
     *
     * - Specific date: `YYYY-MM-DD` (e.g., "2024-01-15")
     * - Date range: `YYYY-MM-DD..YYYY-MM-DD` (e.g., "2024-01-01..2024-01-31")
     *
     * The response includes aggregated metrics such as total sales, transaction count, and average transaction value.
     *
     * - **`exclude_voided`**: The default is `true` — only non-void transactions are counted. When `false`, voided transactions are included in those totals.
     *
     * ### Timezone Handling
     *
     * The service uses the merchant's configured timezone for all date boundaries and grouping.
     * This ensures that payments are grouped according to the merchant's local date/time.
     * For example, a payment at 2025-10-09 22:51:28 UTC would appear in 2025-10-10 data for a merchant in timezone "Asia/Bangkok" (GMT+7).
     *
     * ### Granularity Logic
     *
     * The response granularity is automatically determined based on the date range:
     *
     * - **0 days** (single date): Hourly breakdown
     * - **1-83 days**: Daily breakdown
     * - **84-167 days**: Weekly breakdown
     * - **168+ days**: Monthly breakdown
     *
     * This affects how the `transactions` array is structured and the `granularity` field in the response.
     * @param array{period: string, exclude_voided?: bool, terminal_id?: string} $args
     * @return array{start_date: string, end_date: string, amount: int|float, count: int|float, average?: int|float|null, change?: int|float|null, granularity: 'hourly'|'daily'|'weekly'|'monthly', transactions: list<array<string, mixed>>}
     * @throws BadRequestError INVALID_DATE_FORMAT - This error occurs when the date_input parameter is not in a supported format.
     * Please use one of the supported formats: YYYY-MM-DD, YYYY-MM-DD..YYYY-MM-DD, or keywords.
     * @throws BadRequestError INVALID_DATE_RANGE - This error occurs when the end date is before the start date in a date range.
     * Please ensure the end date is after or equal to the start date.
     * @throws BadRequestError NOT_FOUND - The terminal specified by `terminal_id` could not be found for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function retrieveSalesPerformance(array $args): array
    {
        return $this->client->call('retrieve-sales-performance', $args);
    }

    /**
     * ## Reverse payment
     *
     * Reverses a previously authorized but not yet captured payment. Reversing a payment releases the
     * temporary hold on the customer's funds. This is useful when a transaction is cancelled before it
     * is captured, or when an authorization is no longer needed.
     *
     * Only payments that are in an `AUTHORIZED` status can be reversed. Once a payment is successfully
     * reversed, its status changes to `REVERSED`.
     * @param array{payment_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment ID.
     * @throws BadRequestError NOT_REVERSIBLE - Payment is not eligible for reversal. This can happen if the payment
     * is not in an authorized state, the reversal window has passed, or the
     * integration does not support reversals.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function reversePayment(array $args): bool
    {
        return $this->client->call('reverse-payment', $args);
    }

    /**
     * ## Revoke an API key
     *
     * Permanently revoke an API key to prevent it from being used for authentication. Once revoked, any
     * requests made with this key will be rejected. This is useful when:
     *
     * - An API key has been compromised
     * - You're rotating old API keys as part of security best practices
     * - An employee with access to the key leaves your organization
     *
     * The revocation takes effect immediately and cannot be undone. If you need to restore access,
     * you'll need to create a new API key.
     *
     * The endpoint returns true if the key was successfully revoked.
     * @param array{key_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The API key with the provided ID could not be found in your merchant account. Verify that you're using the correct key_id.
     * @throws BadRequestError ALREADY_REVOKED - The API key has already been revoked. No action is needed.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function revokeKey(array $args): bool
    {
        return $this->client->call('revoke-key', $args);
    }

    /**
     * ## Save card
     *
     * Converts a single-use token into a reusable card that can be charged multiple times. This is useful for:
     *
     * - Implementing "save card for later" functionality
     * - Setting up subscription billing
     * - Enabling one-click payments
     *
     * The token must be created with multi-use capability (`usage="multi"`). Single-use tokens cannot be
     * saved as cards. See our [SDK documentation](/api/sdk) for more details on card tokenization.
     *
     * You can optionally provide your own `external_id` to help track the card in your system. This ID
     * can be used instead of the card ID in other API calls.
     *
     * This card will be associated with the specified contact. If `default` is set to `true`,
     * the card will be set as the contact's default card.
     *
     * The endpoint returns the unique card ID that you can use for future charges.
     *
     * Note: A token can only be saved as a card once. Attempting to save the same token multiple times
     * will result in an error.
     * @param array{external_id?: string, token: string, contact_id: string, contact_email?: string, contact_mobile?: string, contact_name?: string, default?: bool} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - The token with the provided ID could not be found, or the contact with the provided contact_id could not be found. Verify that you're using the correct IDs.
     * @throws BadRequestError NOT_REUSABLE - The token is single-use and cannot be saved as a card. Create a new multi-use token instead.
     * @throws BadRequestError ALREADY_SAVED - This token has already been converted to a saved card. Use the existing card ID instead.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function saveCard(array $args): string
    {
        return $this->client->call('save-card', $args);
    }

    /**
     * ## Select plan
     *
     * Selects a billing plan for the merchant account. This endpoint allows merchants to choose a plan
     * that suits their business needs. The plan must be one of the plans currently available for the provider.
     *
     * If you upgrade, your new plan starts right away and you’ll only pay the difference for the rest of
     * the month. If you downgrade, you’ll keep your current features until the end of the month, and
     * your new plan will begin on the 1st.
     * @param array{plan_id: string} $args
     * @return bool
     * @throws BadRequestError INVALID_PLAN - The specified plan ID was not found.
     * @throws BadRequestError RETIRED_PLAN - The specified plan ID is retired and no longer available for selection.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function selectPlan(array $args): bool
    {
        return $this->client->call('select-plan', $args);
    }

    /**
     * ## Send a payment link
     *
     * Queues a payment link for delivery through an active email or SMS communication channel.
     * The recipient is validated and normalized according to the selected channel.
     * @param array{payment_link_id: string, communication_channel_id: string, recipient: string} $args
     * @return array{communication_delivery_id: string, communication_channel_id: string, payment_link_id: string, delivery: 'email'|'sms', recipient: string, contact_id?: string|null, status: 'queued'|'processing'|'accepted'|'failed', created_at: int|float, accepted_at?: int|float|null}
     * @throws BadRequestError NOT_FOUND - The payment link or communication channel does not exist for this merchant.
     * @throws BadRequestError BAD_REQUEST - The payment link, channel, or recipient cannot be used for this delivery.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function sendPaymentLink(array $args): array
    {
        return $this->client->call('send-payment-link', $args);
    }

    /**
     * ## Send a payout
     *
     * Send a payout to the primary bank account.
     * @param array{external_id?: string, amount: string} $args
     * @return string
     * @throws BadRequestError CONFLICT - A payout with the same external_id already exists in this merchant account.
     * @throws BadRequestError NOT_FOUND - The bank account does not exist or the payout method is not configured for this merchant.
     * @throws BadRequestError INSUFFICIENT_BALANCE - The merchant does not have sufficient balance to complete this payout.
     * The requested amount exceeds the available balance in the merchant's account.
     * @throws BadRequestError INVALID_AMOUNT - The payout amount exceeds the merchant or payout method limits.
     * @throws BadRequestError PAYOUTS_BLOCKED - Payouts are currently blocked for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function sendPayout(array $args): string
    {
        return $this->client->call('send-payout', $args);
    }

    /**
     * ## Send a payout to a contact
     *
     * Sends a payout to an existing contact, or creates the contact in the same request. Provide `recipient_id`
     * to use an existing contact. Without it, the contact is matched by `recipient_email` or `recipient_mobile`,
     * or created when none matches.
     *
     * - **Bank account**: Omit the bank fields to use the contact's primary bank account,
     *   or include `bank_account_id` to specify a particular account. To create a new bank account
     *   on the contact, omit `bank_account_id` and provide `bank_account_number`, `bank_account_name`,
     *   `bank_code`, and `country_code` instead. The two approaches are mutually exclusive.
     *   A newly created contact has no bank account, so provide the bank details.
     * - **Mobile number (PromptPay)**: Set `payout_method` to `mobile_number` and provide `mobile_number`.
     *   Accepts local format (e.g. `0812223333`) or E.164 (e.g. `+66812223333`).
     * - **National ID (PromptPay)**: Set `payout_method` to `national_id` and provide `national_id` (13 digits).
     *
     * When using a bank account, it must belong to the contact and be active (not soft-deleted).
     * @param array{recipient_id?: string, payout_method?: string, bank_account_id?: string, bank_account_number?: string, bank_account_name?: string, bank_code?: string, country_code?: string, mobile_number?: string, national_id?: string, external_id?: string, amount: string, recipient_mobile?: string, recipient_name?: string, recipient_email?: string, update_contact?: true|false, send_confirmation_email?: true|false} $args
     * @return string
     * @throws BadRequestError CONFLICT - A payout with the same external_id already exists in this merchant account.
     * @throws BadRequestError NOT_FOUND - No contact found with the provided recipient_id, the bank account does not exist or does not belong to the contact, or the payout method is not configured for this merchant.
     * @throws BadRequestError INVALID_OPERATION - The contact does not have a primary bank account. Please specify a bank_account_id.
     * @throws BadRequestError INSUFFICIENT_BALANCE - The merchant does not have sufficient balance to complete this payout.
     * The requested amount exceeds the available balance in the merchant's account.
     * @throws BadRequestError INVALID_MOBILE_NUMBER - The mobile number must be a valid local number (e.g. 0812223333) or E.164 number (e.g. +66812223333).
     * @throws BadRequestError INVALID_NATIONAL_ID - The national ID must be exactly 13 digits.
     * @throws BadRequestError INVALID_AMOUNT - The payout amount exceeds the merchant or payout method limits.
     * @throws BadRequestError PAYOUTS_BLOCKED - Payouts are currently blocked for this merchant.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function sendPayoutToContact(array $args): string
    {
        return $this->client->call('send-payout-to-contact', $args);
    }

    /**
     * ## Set default card for a contact
     *
     * Set an existing card as the default card for a contact. The card must already be associated with the contact.
     *
     * If the contact already has a default card set up, using this endpoint will:
     *
     * - Replace the existing default card with the specified one
     *
     * The endpoint returns `true` when the default card is successfully set.
     * @param array{contact_id: string, card_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The specified contact was not found, or the card does not exist for this contact.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function setDefaultCard(array $args): bool
    {
        return $this->client->call('set-default-card', $args);
    }

    /**
     * ## Set entity type
     *
     * Specify whether you are registering as a business or an individual merchant. This is a critical first
     * step in the onboarding process that determines:
     *
     * - What information you'll need to provide
     * - Which features and products will be available to you
     * - Your verification requirements
     *
     * For business merchants:
     *
     * - You'll need to provide business registration details
     * - You can add multiple key people (owners, directors)
     *
     * For individual merchants:
     *
     * - Only personal identification is required
     * - Simpler verification process
     *
     * Important: This setting can only be configured once during initial onboarding and cannot be changed
     * later. Make sure you select the correct entity type for your situation.
     *
     * The endpoint returns true if the entity type was successfully set.
     * @param array{entity_type: 'BUSINESS'|'INDIVIDUAL'} $args
     * @return bool
     * @throws BadRequestError ALREADY_SET - The entity type has already been set for this merchant account and cannot be changed. You must create a new merchant account if you need a different entity type.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function setEntityType(array $args): bool
    {
        return $this->client->call('set-entity-type', $args);
    }

    /**
     * ## Set merchant classification
     *
     * Specify your business category and the types of products/services you offer. This information helps us:
     *
     * - Configure the right payment features for your business
     * - Apply appropriate transaction limits and fees
     * - Meet regulatory requirements for your business type
     *
     * You'll need to provide:
     *
     * 1. A Merchant Category Code (MCC) - This is a 4-digit code that classifies your primary business activity
     * 2. Product types - The categories of products or services you sell
     *
     * You can optionally declare the transaction volume you expect to process. We use those
     * figures to monitor your account for unexpected activity. Both amounts are expressed in
     * the smallest currency unit (subunits) of your funding currency, must be a positive
     * integer greater than 0, and can be set or changed at any point while onboarding is in
     * draft. Omitted arguments are left unchanged.
     *
     * Important: Your business category and product types can only be set once during initial
     * onboarding and cannot be changed later. Make sure you select the correct classification
     * that best represents your business.
     *
     * The endpoint returns true if the classification was successfully set.
     * @param array{category_code: string, product_types: 'DIGITAL_GOODS'|'PHYSICAL_GOODS'|'SERVICES', expected_monthly_transaction_amount?: int|float, expected_maximum_transaction_amount?: int|float} $args
     * @return bool
     * @throws BadRequestError ALREADY_SET - The merchant classification has already been set for this account and cannot be
     * changed. You must create a new merchant account if you need a
     * different classification.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function setMerchantClassification(array $args): bool
    {
        return $this->client->call('set-merchant-classification', $args);
    }

    /**
     * ## Set primary bank account
     *
     * Configure the bank account where you want to receive payouts from your payment transactions. This bank account will be used for:
     *
     * - Receiving settlement funds from successful payments
     * - Automated payouts based on your payout schedule
     * - Manual payout requests
     *
     * If you already have a primary bank account set up, using this endpoint will:
     *
     * - Replace your existing primary bank account with the new one
     * - Keep the old bank account record for audit purposes
     * - Immediately start using the new account for future payouts
     *
     * Important: The bank account must be in the same country as your payment provider. For example, if your provider is in Thailand, you must provide a Thai bank account.
     *
     * The endpoint returns a unique identifier for the bank account that you can use to reference it in other API calls.
     * @param array{bank_code: string, number: string, name: string, country_code: string} $args
     * @return string
     * @throws BadRequestError COUNTRY_MISMATCH - The bank account's country code does not match your payment provider's country. You must provide a bank account in the same country as your provider to receive payouts.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function setPrimaryBankAccount(array $args): string
    {
        return $this->client->call('set-primary-bank-account', $args);
    }

    /**
     * ## Set primary bank account for a contact
     *
     * Set an existing bank account as the primary bank account for a contact. The bank account must already be associated with the contact.
     *
     * If the contact already has a primary bank account set up, using this endpoint will:
     *
     * - Replace the existing primary bank account with the specified one
     *
     * The endpoint returns `true` when the primary bank account is successfully set.
     * @param array{contact_id: string, bank_account_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id, or the bank account with the provided bank_account_id does not exist or does not belong to the contact.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function setPrimaryContactBankAccount(array $args): bool
    {
        return $this->client->call('set-primary-contact-bank-account', $args);
    }

    /**
     * ## Start checkout
     *
     * Initiates a new checkout session. A checkout session represents a customer's intent to pay and
     * provides a secure URL where they can complete their transaction.
     *
     * This endpoint is typically used in web integrations where you want to redirect your customer to a
     * hosted payment page.
     * @param array{installation_id: string, amount: string, currency: string, interest_bearer?: 'customer'|'merchant', return_url: string, customer_id?: string, customer_email?: string, customer_mobile?: string, customer_name?: string} $args
     * @return array{checkout_id: string, amount: string, currency: string, interest_bearer?: 'customer'|'merchant'|null, installation_id: string, customer_id?: string|null, customer_email?: string|null, customer_mobile?: string|null, customer_name?: string|null, redirect_url: string, payment_id?: string|null}
     * @throws BadRequestError NOT_FOUND - The installation could not be found with the provided ID.
     * @throws BadRequestError NOT_ENABLED - This installation is currently disabled.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function startCheckout(array $args): array
    {
        return $this->client->call('start-checkout', $args);
    }

    /**
     * ## Start a payment on a terminal
     *
     * This endpoint initiates a payment flow on a specific terminal. It selects a payment method and notifies the terminal device to start the payment process.
     *
     * When called, it will notify the terminal device to start the payment process.
     *
     * It returns true if the notification was sent.
     * @param array{external_id?: string, terminal_id: string, amount: int|float, currency: string, payment_method: string, capture?: bool} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - The terminal could not be found.
     * @throws BadRequestError INVALID_ARGUMENT - The payment method is not allowed or the amount is invalid.
     * @throws BadRequestError NO_INSTALLATION - The terminal is not properly configured with an installation.
     * @throws BadRequestError NOT_ENABLED - This installation is currently disabled.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function startPaymentOnTerminal(array $args): bool
    {
        return $this->client->call('start-payment-on-terminal', $args);
    }

    /**
     * ## Subscribe to topics
     *
     * Subscribes an email address to one or more notification topics for this merchant.
     * The email will receive notifications when events related to those topics occur.
     *
     * This endpoint supports both create and update operations:
     *
     * - **Create**: When no subscriptions exist for the email+merchant combination.
     * - **Update**: When subscriptions exist, it will replace existing subscriptions for that email+merchant with the new ones.
     *
     * You can subscribe to multiple topics at once by providing an array of topic IDs.
     * The language parameter determines the language of the notifications sent to this email.
     * @param array{email: string, language: 'en'|'th', notification_topic_ids: array} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - One or more notification topics with the specified IDs could not be found.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function subscribeToTopics(array $args): bool
    {
        return $this->client->call('subscribe-to-topics', $args);
    }

    /**
     * ## Supported payment methods
     *
     * Get a list of all payment methods that are currently enabled for your merchant account.
     *
     * The endpoint returns an array of payment methods that you can use with other API endpoints. Each payment method
     * lists the `platforms` it is available on, so you only offer the payment methods that fit the platform the payment happens on.
     * @return list<array{name: string, display_name: string, flow: string, category: string, platforms: 'WEB'|'IOS'|'ANDROID'|'TERMINAL'}>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function supportedPaymentMethods(): array
    {
        return $this->client->call('supported-payment-methods');
    }

    /**
     * ## Sync object metadata
     *
     * Metadata allows you to store and manage arbitrary JSON data associated with your merchant account.
     *
     * ### Keys and Objects
     *
     * Metadata is organized into objects, each identified by a unique `key`. The `key` is a string of your choice that serves as a unique identifier for the metadata record.
     *
     * ### Nested Paths
     *
     * You can use dot-separated paths (e.g., `settings.theme` or `user.profile.id`) to update specific parts of a metadata object without overwriting the entire thing.
     *
     * - If the path does not exist, it will be created.
     * - If no path is provided, the entire metadata object for the given key is replaced.
     *
     * ### Expiration (TTL)
     *
     * Each metadata record has an expiration date.
     *
     * - By default, metadata expires after 2 years from its last update.
     * - You can specify a custom `ttl` (Time To Live) in seconds when syncing.
     * - Metadata can be deleted by syncing it with a `ttl` of `0`.
     * - Expired metadata can be reactivated by syncing it again.
     *
     * ### Optimistic Locking
     *
     * To prevent concurrent updates from overwriting each other, you can use the `version` argument. If provided, the update will only succeed if it matches the current `version` (lock_version) of the record.
     *
     * ### Limits
     *
     * - Individual metadata records are limited to 4 KB.
     * - For test accounts, the total size of all metadata records is limited to 10 MB.
     * @param array{key: string, path?: string, metadata: mixed, version?: int|float, ttl?: int|float} $args
     * @return string
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function syncMetadata(array $args): string
    {
        return $this->client->call('sync-metadata', $args);
    }

    /**
     * ## Turn off auto payout
     *
     * Disable automatic payouts to your bank account. When auto payout is turned off:
     *
     * - Funds from your transactions will accumulate in your merchant balance
     * - You'll need to manually request payouts when you want to withdraw funds
     * - You can still view your available balance and payout history
     * - You can turn auto payout back on at any time
     *
     * This is useful if you want more control over when you receive your funds, for example:
     *
     * - To accumulate a larger balance before withdrawing
     * - To manually time your payouts with your business needs
     * - To temporarily hold funds during account changes
     *
     * The endpoint returns true if auto payout was successfully disabled.
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function turnOffAutoPayout(): bool
    {
        return $this->client->call('turn-off-auto-payout');
    }

    /**
     * ## Turn on auto payout
     *
     * Enable automatic transfers of your available balance to your bank account on a regular schedule. When auto payout is turned on:
     *
     * - Your available balance will automatically be transferred to your primary bank account
     * - Payouts will occur according to your chosen frequency (daily or weekly)
     * - You don't need to manually request payouts
     * - You can still view your payout history and track upcoming payouts
     * - You can turn auto payout off at any time if you want to switch back to manual payouts
     *
     * This is useful if you want a hands-off approach to managing your funds, receive your money as soon as it's available.
     *
     * Note that you must have a primary bank account set up before enabling auto payout. Use the [set-primary-bank-account](/api/merchants/set-primary-bank-account) endpoint first if you haven't configured one yet.
     * @param array{frequency: 'DAILY'|'WEEKLY'} $args
     * @return bool
     * @throws BadRequestError NO_PRIMARY_BANK_ACCOUNT - This error occurs when the merchant does not have a primary bank account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function turnOnAutoPayout(array $args): bool
    {
        return $this->client->call('turn-on-auto-payout', $args);
    }

    /**
     * ## Unlink bank account from a contact
     *
     * Unlink (soft delete) a bank account from a contact. The bank account will be marked as deleted but the record will be preserved for audit purposes.
     *
     * Bank accounts include a `deletable` attribute that indicates whether the bank account can be deleted. The frontend should check this attribute before attempting to delete a bank account. If a bank account cannot be deleted, set another bank account as the primary using the `set-primary-contact-bank-account` endpoint, which will update the `deletable` attribute accordingly.
     *
     * The endpoint returns `true` when the bank account is successfully unlinked.
     * @param array{contact_id: string, bank_account_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id, or the bank account with the provided bank_account_id does not exist or does not belong to the contact.
     * @throws BadRequestError INVALID_OPERATION - Cannot delete primary bank account when contact has multiple accounts. Set another account as primary first.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function unlinkBankAccount(array $args): bool
    {
        return $this->client->call('unlink-bank-account', $args);
    }

    /**
     * ## Unsubscribe from all topics
     *
     * Unsubscribe a specific email address from all notification topics within this merchant.
     * This will remove all existing subscriptions for the specified email address.
     *
     * This endpoint requires authentication and will return true if the unsubscribe operation was successful.
     * @param array{email: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No subscriptions found for the specified email address.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function unsubscribeFromAllTopics(array $args): bool
    {
        return $this->client->call('unsubscribe-from-all-topics', $args);
    }

    /**
     * ## Update a bank account
     *
     * Update the details of an existing bank account linked to a contact. Only the fields you provide will be updated — omitted fields will remain unchanged.
     *
     * The endpoint returns the updated bank account object.
     * @param array{contact_id: string, bank_account_id: string, name?: string, primary?: bool, external_id?: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id, or the bank account does not exist or does not belong to the contact.
     * @throws BadRequestError INVALID_PRIMARY - Cannot unset primary on the current primary bank account. Set another bank account as primary instead.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateBankAccount(array $args): bool
    {
        return $this->client->call('update-bank-account', $args);
    }

    /**
     * ## Update branding
     *
     * Update the visual branding settings for your merchant account. This includes colors, typography,
     * button styling, and display preferences for your payment interface.
     *
     * All parameters are optional - you can update individual settings without affecting others. The
     * endpoint returns true if the branding was successfully updated.
     *
     * **Note:** For uploading logo and sidebar images, use the `generate-branding-upload-url` and
     * `confirm-branding-upload` endpoints.
     * @param array{primary_color?: string, secondary_color?: string, sidebar_color?: string, sidebar_background_type?: 'color'|'image', tile_sidebar_image?: bool, hide_logo?: bool, heading_font_family?: string, heading_font_variant?: 'Thin'|'ExtraLight'|'Light'|'Regular'|'Medium'|'SemiBold'|'Bold'|'ExtraBold'|'Black', body_font_family?: string, body_font_variant?: 'Thin'|'ExtraLight'|'Light'|'Regular'|'Medium'|'SemiBold'|'Bold'|'ExtraBold'|'Black', use_iso_code?: bool, hide_non_significant_zeros?: bool, show_amount_in_button?: bool, symbol_placement?: 'use_default'|'before_amount'|'after_amount', default_button_text?: string, button_shape?: 'rounded_none'|'rounded_full'|'rounded_xs'|'rounded_sm'|'rounded_md'|'rounded_lg'|'rounded_xl'} $args
     * @return bool
     * @throws BadRequestError INVALID_RECORD - This error occurs when the record contains invalid parameters. Common cases
     * include invalid hex color formats, unsupported font families, or invalid
     * font variants for the specified font family.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateBranding(array $args = []): bool
    {
        return $this->client->call('update-branding', $args);
    }

    /**
     * ## Update business details
     *
     * Updates the business details for a merchant during the onboarding process. This endpoint is only available
     * while your merchant account is in the onboarding phase. Once onboarding is complete, these details cannot
     * be modified through the API.
     *
     * This endpoint only works for merchants registered as a business entity (entity_type = "business"). Individual
     * merchants cannot use this endpoint.
     *
     * You can either:
     *
     * - Provide all required fields for a complete update
     * - Set partial_update=true to update only specific fields
     *
     * All dates should be provided in ISO8601 format (YYYY-MM-DD).
     * @param array{partial_update?: bool, legal_name: string, registration_number: string, registration_date?: string, address_street_1: string, address_street_2?: string, postal_code: string, city: string, country: string, primary_administrative_division: string, operation_address_street_1?: string, operation_address_street_2?: string, operation_address_postal_code?: string, operation_address_city?: string, operation_address_country?: string, operation_address_primary_administrative_division?: string, use_business_registration_as_operation_address?: bool, branch: string, legal_structure: 'PRIVATE_LIMITED_COMPANY'|'PUBLIC_LIMITED_COMPANY'|'LIMITED_PARTNERSHIP'|'REGISTERED_ORDINARY_PARTNERSHIP'|'ORDINARY_PARTNERSHIP'|'GOVERNMENT_AGENCY'|'FOUNDATION'|'ASSOCIATION'|'COOPERATIVE'|'CHAMBER_OF_COMMERCE'|'SCHOOL_OR_UNIVERSITY'|'PLACE_OF_WORSHIP'|'FOREIGN_COMPANY', kym_contact_person: string, kym_contact_phone: string} $args
     * @return bool
     * @throws BadRequestError NOT_BUSINESS - This error occurs when attempting to update business details for a merchant that is not registered as a business entity. Business details can only be updated for business merchants, not individual merchants.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateBusinessDetails(array $args): bool
    {
        return $this->client->call('update-business-details', $args);
    }

    /**
     * ## Update contact
     *
     * Updates an existing contact. The `contact_id` can be either the internal UID
     * (starting with `ctc_`) or the external ID.
     * The service will check for conflicts with other contacts if the new values already exist.
     * @param array{contact_id: string, contact_email?: string, contact_mobile?: string, contact_name?: string, favorite?: bool} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - No contact found with the provided contact_id.
     * The contact may not exist or may belong to a different merchant.
     * @throws BadRequestError CONFLICT - A contact with the same email or mobile number already exists.
     * Each contact must have unique values for these fields.
     * @throws BadRequestError VALIDATION_ERROR - Validation failed. At least one of contact_email or contact_mobile must be provided.
     * If contact_email is provided, it must be a valid email format.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateContact(array $args): bool
    {
        return $this->client->call('update-contact', $args);
    }

    /**
     * ## Update individual details
     *
     * Updates personal and contact information for an individual merchant account during the onboarding process.
     * This endpoint is only available:
     *
     * - During the initial onboarding phase
     * - For individual merchants (entity_type = "individual")
     * - For local nationals only (foreign nationals cannot register as individual merchants)
     *
     * Once onboarding is complete, these details cannot be modified through the API.
     *
     * You can update information in two ways:
     *
     * - Provide all required fields for a complete update
     * - Set partial_update=true to update only specific fields
     *
     * For identification, you must provide:
     *
     * - Your national ID number
     * - National ID expiry date
     * - For Thai nationals: The 12-digit laser code from the back of your ID card
     *
     * All dates should be provided in ISO8601 format (YYYY-MM-DD). The merchant must be at least 18 years old.
     *
     * Note: If you are a business merchant or your onboarding is complete, this endpoint will not be available.
     * Foreign nationals must register as a business entity instead of an individual merchant.
     * @param array{partial_update?: bool, first_name: string, middle_name?: string, last_name: string, first_name_local_language?: string, middle_name_local_language?: string, last_name_local_language?: string, nationality: string, national_id: string, national_id_expiry_date: string, national_id_without_expiry_date?: bool, thai_national_id_laser_code?: string, phone: string, date_of_birth: string, address_street_1: string, address_street_2?: string, postal_code: string, city: string, country: string, primary_administrative_division: string, operation_address_street_1?: string, operation_address_street_2?: string, operation_address_postal_code?: string, operation_address_city?: string, operation_address_country?: string, operation_address_primary_administrative_division?: string, current_address_street_1?: string, current_address_street_2?: string, current_address_postal_code?: string, current_address_city?: string, current_address_country?: string, current_address_primary_administrative_division?: string, use_residence_registration_as_current_address?: bool, use_residence_registration_as_operation_address?: bool} $args
     * @return bool
     * @throws BadRequestError NOT_INDIVIDUAL - This error occurs when attempting to update individual details for a business merchant account.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateIndividualDetails(array $args): bool
    {
        return $this->client->call('update-individual-details', $args);
    }

    /**
     * ## Update an API key
     *
     * Update the name of an existing API key for your merchant account. This allows you to change the descriptive
     * name that helps identify the key's purpose.
     *
     * The key's actual value and functionality remain unchanged - only the display name is modified.
     * @param array{key_id: string, name: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no active API key with the provided ID exists in your merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateKey(array $args): bool
    {
        return $this->client->call('update-key', $args);
    }

    /**
     * ## Update a key person
     *
     * Update details of a key person in your business merchant account. Key people are individuals
     * who have significant control or influence over your business, such as:
     *
     * - Ultimate Beneficial Owners (UBOs) who own 25% or more of the business
     * - Authorized Signatories who can legally represent and make decisions for
     *   the business
     *
     * This endpoint is only available for business merchants, not
     * individual merchants.
     * @param array{key_person_id: string, first_name: string, middle_name?: string, last_name: string, first_name_local_language?: string, middle_name_local_language?: string, last_name_local_language?: string, date_of_birth: string, nationality: string, type: 'UBO'|'AUTHORIZED_SIGNATORY', percentage_of_ownership?: int|float, authorized_director?: bool, national_id?: string, national_id_expiry_date?: string, national_id_without_expiry_date?: bool, passport_number?: string, passport_expiry_date?: string, address_street_1?: string, address_street_2?: string, postal_code?: string, city?: string, primary_administrative_division?: string, country?: string, current_address_street_1?: string, current_address_street_2?: string, current_address_postal_code?: string, current_address_city?: string, current_address_primary_administrative_division?: string, current_address_country?: string, use_residence_registration_as_current_address?: bool} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no key person with the provided ID exists in your
     * merchant account.
     * @throws BadRequestError NOT_EDITABLE - This error occurs when the onboarding status of the merchant is no longer
     * in draft.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateKeyPerson(array $args): bool
    {
        return $this->client->call('update-key-person', $args);
    }

    /**
     * ## Update merchant settings
     *
     * Updates self-service settings for your merchant account, such as communication language and payout authorization requirements.
     *
     * To view system properties set up for you by Reservepay, use `retrieve-configuration`.
     * @param array{preferred_communication_language?: 'en'|'th', payout_contact_authorization_required?: bool, payout_allow_self_authorization?: bool, payout_skip_authorization_for_api_key?: bool} $args
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateMerchantSettings(array $args = []): bool
    {
        return $this->client->call('update-merchant-settings', $args);
    }

    /**
     * ## Update a sales channel
     *
     * Update a sales channel for the merchant. Note that the platform cannot be changed after creation.
     * @param array{sales_channel_id: string, url?: string, demo_url?: string, instructions?: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no sales channel with the provided ID exists in your
     * merchant account.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateSalesChannel(array $args): bool
    {
        return $this->client->call('update-sales-channel', $args);
    }

    /**
     * ## Update terminal
     *
     * Updates an existing terminal. Returns the `terminal_id`.
     * @param array{terminal_id: string, name?: string, external_id?: string, installation_id?: string, branch_name?: string, kiosk_mode?: bool, delegated_terminal_id?: string} $args
     * @return string
     * @throws BadRequestError NOT_FOUND - Terminal or installation not found.
     * @throws BadRequestError INVALID_ARGUMENT - Terminal cannot delegate to itself, the delegated terminal already has a delegation configured or is already used by another terminal, or kiosk_mode cannot be combined with delegated_terminal_id.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function updateTerminal(array $args): string
    {
        return $this->client->call('update-terminal', $args);
    }

    /**
     * ## Validate bank account
     *
     * Check with the bank that a bank account exists and that its holder name matches the name
     * you expect, before adding the bank account to a contact.
     *
     * The account holder name from the bank is compared against `name` and `local_name`.
     * Titles such as `Mr.` or `นาย` are ignored.
     * `match` is `true` when `score` is greater than or equal to `sensitivity`.
     *
     * A fee is charged per call. Contact support for pricing.
     * @param array{bank_code: string, account_number: string, name?: string, local_name?: string, sensitivity: int|float} $args
     * @return array{available: bool, match: bool, score: int|float, similarity: int|float, distance?: int|float|null, normalized_name?: string|null}
     * @throws BadRequestError NOT_SUPPORTED - Bank account validation is not available for your account.
     * @throws BadRequestError LOOKUP_FAILED - The bank account could not be looked up with the bank. Please try again later.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function validateBankAccount(array $args): array
    {
        return $this->client->call('validate-bank-account', $args);
    }

    /**
     * ## Void payment
     *
     * Voids a previously captured payment that has not yet been settled. Voids are typically used to cancel
     * a transaction before funds are officially transferred. Unlike refunds, which occur after settlement,
     * a void happens before the transaction is finalized on the customer's statement.
     *
     * Only payments that are in a `SUCCESSFUL` (captured) status and have not yet been settled can be voided.
     * Once a payment is successfully voided, its status changes to `VOIDED`.
     * @param array{payment_id: string} $args
     * @return bool
     * @throws BadRequestError NOT_FOUND - This error occurs when no payment with the provided ID exists in your
     * merchant account. Double check that you're using the correct payment ID.
     * @throws BadRequestError NOT_VOIDABLE - Payment is not eligible for void. This can happen if the payment is
     * not in a captured state, has already settled into a statement, the
     * void window has passed, or the integration does not support voids.
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function voidPayment(array $args): bool
    {
        return $this->client->call('void-payment', $args);
    }

}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type EntryAttributes from Client
 * @psalm-import-type EntryAttributes from Client
 */
final class BalanceHistoryPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<EntryAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type ContactAttributes from Client
 * @psalm-import-type ContactAttributes from Client
 */
final class FilterContactsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<ContactAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentLinkAttributes from Client
 * @psalm-import-type PaymentLinkAttributes from Client
 */
final class FilterPaymentLinksPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentLinkAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentAttributes from Client
 * @psalm-import-type PaymentAttributes from Client
 */
final class FilterPaymentsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PayoutAttributes from Client
 * @psalm-import-type PayoutAttributes from Client
 */
final class FilterPayoutsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PayoutAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type TerminalAttributes from Client
 * @psalm-import-type TerminalAttributes from Client
 */
final class FilterTerminalsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<TerminalAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type TopupAttributes from Client
 * @psalm-import-type TopupAttributes from Client
 */
final class FilterTopupsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<TopupAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type ActivityLogAttributes from Client
 * @psalm-import-type ActivityLogAttributes from Client
 */
final class ListActivityLogsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<ActivityLogAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type ApiRequestLogAttributes from Client
 * @psalm-import-type ApiRequestLogAttributes from Client
 */
final class ListApiRequestLogsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<ApiRequestLogAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type BatchAttributes from Client
 * @psalm-import-type BatchAttributes from Client
 */
final class ListBatchesPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<BatchAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type ContactAttributes from Client
 * @psalm-import-type ContactAttributes from Client
 */
final class ListContactsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<ContactAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type DisputeAttributes from Client
 * @psalm-import-type DisputeAttributes from Client
 */
final class ListDisputesPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<DisputeAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type DisputeAttributes from Client
 * @psalm-import-type DisputeAttributes from Client
 */
final class ListDisputesByPaymentPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<DisputeAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type InvoiceAttributes from Client
 * @psalm-import-type InvoiceAttributes from Client
 */
final class ListInvoicesPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<InvoiceAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type MetadataAttributes from Client
 * @psalm-import-type MetadataAttributes from Client
 */
final class ListMetadataPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<MetadataAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type MonthlyBillingArchiveAttributes from Client
 * @psalm-import-type MonthlyBillingArchiveAttributes from Client
 */
final class ListMonthlyBillingArchivesPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<MonthlyBillingArchiveAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentAttributes from Client
 * @psalm-import-type PaymentAttributes from Client
 */
final class ListPaymentLinkPaymentsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentLinkAttributes from Client
 * @psalm-import-type PaymentLinkAttributes from Client
 */
final class ListPaymentLinksPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentLinkAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentAttributes from Client
 * @psalm-import-type PaymentAttributes from Client
 */
final class ListPaymentsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentAttributes from Client
 * @psalm-import-type PaymentAttributes from Client
 */
final class ListPaymentsByContactPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PaymentAttributes from Client
 * @psalm-import-type PaymentAttributes from Client
 */
final class ListPaymentsByTerminalPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PaymentAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PayoutAttributes from Client
 * @psalm-import-type PayoutAttributes from Client
 */
final class ListPayoutsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PayoutAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type PayoutAttributes from Client
 * @psalm-import-type PayoutAttributes from Client
 */
final class ListPayoutsByContactPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<PayoutAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type RefundAttributes from Client
 * @psalm-import-type RefundAttributes from Client
 */
final class ListRefundsByPaymentPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<RefundAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type StatementAttributes from Client
 * @psalm-import-type StatementAttributes from Client
 */
final class ListStatementsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<StatementAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type TerminalAttributes from Client
 * @psalm-import-type TerminalAttributes from Client
 */
final class ListTerminalsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<TerminalAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type TopupAttributes from Client
 * @psalm-import-type TopupAttributes from Client
 */
final class ListTopupsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<TopupAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 *
 * @phpstan-import-type VirtualAccountAttributes from Client
 * @psalm-import-type VirtualAccountAttributes from Client
 */
final class ListVirtualAccountsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<VirtualAccountAttributes>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

