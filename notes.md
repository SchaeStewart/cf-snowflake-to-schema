# Notes

## required fields

## not required fields

### transactions

not required:

- source_transaction_id
- committee_street_1 - remove from schema
- committee_street_2 - remove from schema
- committee_city - remove from schema
- committee_state - remove from schema

required:

- contributor_id
- transaction_type
- committee_name
- canon_committee_sboe_id - not in snowflake
- transaction_category - maps to transaction_type?
- date_occurred
- amount
- form_of_payment
- purpose
- candidate_referendum_name
- declaration
- original_committee_sboe_id - used for expenditures
- original_account_id - not used

accounts:

- name
- city
- state
- zip_code
- profession
- employer_name
  // remove below from schema
- is_donor
- is_vendor
- is_person
- is_organization

committees

comm_id - not used
sboe_id
current_status - in repr, but not used
committee_name
committee_street_1
committee_street_2
committee_city
committee_state
committee_full_zip
candidate_full_name - not in snowflake
candidate_first_last_name - not in snowflake
treasurer_first_name
treasurer_middle_name
treasurer_last_name
treasurer_email
asst_treasurer_first_name
asst_treasurer_middle_name
asst_treasurer_last_name
asst_treasurer_email
treasurer_city
treasurer_state
treasurer_full_zip
treasuer_id
asst_treasurer_id
party
office
juris
