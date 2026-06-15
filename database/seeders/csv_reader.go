package seeders

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Record = map[string]string

// seedDir adalah lokasi file CSV master data, relatif ke project root.
const seedDir = "database/seeders/file"

// ReadCSV membaca file dari folder database/seeders/file relatif ke project root.
// Baris diawali '#' diabaikan (komentar).
func ReadCSV(filename string) ([]Record, error) {
	path := filepath.Join(seedDir, filename)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv %q: %w", path, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	reader.Comment = '#'

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read headers from %q: %w", filename, err)
	}
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}

	var records []Record
	lineNum := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%q line %d: %w", filename, lineNum, err)
		}
		lineNum++

		rec := make(Record, len(headers))
		for i, h := range headers {
			if i < len(row) {
				rec[h] = strings.TrimSpace(row[i])
			}
		}
		records = append(records, rec)
	}
	return records, nil
}

// GetString mengambil string dengan fallback default.
func GetString(r Record, key, defaultVal string) string {
	if v, ok := r[key]; ok && v != "" {
		return v
	}
	return defaultVal
}

// GetStringPtr mengembalikan *string (nil jika kosong).
func GetStringPtr(r Record, key string) *string {
	if v, ok := r[key]; ok && v != "" {
		return &v
	}
	return nil
}

// GetUUID mem-parse uuid dari CSV; uuid.Nil jika kosong/invalid (pemanggil bisa
// fallback ke generate). Dipakai agar identitas seed deterministik lintas rebuild.
func GetUUID(r Record, key string) uuid.UUID {
	v, ok := r[key]
	if !ok || v == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// GetBool ("true"/"1"/"yes" = true).
func GetBool(r Record, key string, defaultVal bool) bool {
	v, ok := r[key]
	if !ok || v == "" {
		return defaultVal
	}
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "true" || v == "1" || v == "yes"
}

// GetInt64 mengambil int64 dengan fallback default.
func GetInt64(r Record, key string, defaultVal int64) int64 {
	v, ok := r[key]
	if !ok || v == "" {
		return defaultVal
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetInt16 mengambil int16 dengan fallback default.
func GetInt16(r Record, key string, defaultVal int16) int16 {
	return int16(GetInt64(r, key, int64(defaultVal)))
}

// GetInt16Ptr mengembalikan *int16 (nil jika kosong).
func GetInt16Ptr(r Record, key string) *int16 {
	v, ok := r[key]
	if !ok || v == "" {
		return nil
	}
	n, err := strconv.ParseInt(v, 10, 16)
	if err != nil {
		return nil
	}
	x := int16(n)
	return &x
}
