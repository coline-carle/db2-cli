meta:
  id: wdb6
  file-extension: db2
  endian: le

seq:
  - id: header
    type: header

  - id: fields
    type: field_format
    repeat: expr
    repeat-expr: header.field_count

  - id: records
    size: header.record_size
    repeat: expr
    repeat-expr: header.record_count

  - id: string_table
    size: header.string_table_size
    if: header.has_offset_map != true

  - id: ids
    type: u4
    repeat: expr
    repeat-expr: header.record_count
    if: header.has_non_inlne_id

  - id: common_data_table
    type: common_data_table
    if: header.has_common_data_table
    size: header.common_data_table_size

types:
  common_data_table:
    seq:
      - id: num_columns_in_table
        type: u4

      - id: common_data_table
        type: common_data_table_entry
        repeat: expr
        repeat-expr: num_columns_in_table

  common_data_table_entry:
    seq:
      - id: count
        type: u4

      - id: type
        type: u1

      - id: common_data_map_entry
        type: common_data_map_entry
        repeat: expr
        repeat-expr: count

  common_data_map_entry:
    seq:
      - id: id
        type: u4

      - id: type
        type: u4

  field_format:
    seq:
      - id: size
        type: u2

      - id: position
        type: u2

    instances:
      byte_size:
        value: (32 - size) / 8

  header:
    seq:
      - id: magic
        contents: WDB6

      - id: record_count
        type: u4

      - id: field_count
        type: u4

      - id: record_size
        type: u4

      - id: string_table_size
        type: u4

      - id: table_hash
        type: u4

      - id: layout_hash
        type: u4

      - id: min_id
        type: u4

      - id: max_id
        type: u4

      - id: locale
        type: u4

      - id: copy_table_size
        type: u4

      - id: flags
        type: u2

      - id: id_index
        type: u2

      - id: total_field_count
        type: u4

      - id: common_data_table_size
        type: u4

    instances:
      has_offset_map:
        value: (flags & 0x01) == 0x01

      has_secondary_key:
        value: (flags & 0x02) == 0x02

      has_non_inlne_id:
        value: (flags & 0x04) == 0x04

      has_common_data_table:
        value: common_data_table_size > 0

      offset_map_pos:
        value: string_table_size
        if: has_offset_map
