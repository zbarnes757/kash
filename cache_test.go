package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		TTL             time.Duration
		cleanupInterval time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Cache
	}{
		{
			name: "Should generate a new Cache",
			args: args{
				TTL:             1 * time.Minute,
				cleanupInterval: 1 * time.Minute,
			},
			want: &Cache{
				entries:         map[string]entry{},
				TTL:             1 * time.Minute,
				CleanupInterval: 1 * time.Minute,
			},
		},
		{
			name: "Should generate a new Cache with TTL and cleanupInterval set to -1",
			args: args{
				TTL:             -1,
				cleanupInterval: -1,
			},
			want: &Cache{
				entries:         map[string]entry{},
				TTL:             -1,
				CleanupInterval: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.TTL, tt.args.cleanupInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Put(t *testing.T) {
	type fields struct {
		entries map[string]entry
		TTL     time.Duration
	}

	type args struct {
		key   string
		value interface{}
	}

	TTL := 1 * time.Minute
	expiryTime := time.Now().Add(TTL).Unix()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   entry
	}{
		{
			name: "Should add new entry",
			fields: fields{
				entries: map[string]entry{},
				TTL:     TTL,
			},
			args: args{
				key:   "key",
				value: "value",
			},
			want: entry{
				value:      "value",
				expiryTime: expiryTime,
			},
		},
		{
			name: "Should replace an existing value for a key",
			fields: fields{
				entries: map[string]entry{
					"key": entry{
						value: "value1",
					},
				},
				TTL: TTL,
			},
			args: args{
				key:   "key",
				value: "value2",
			},
			want: entry{
				value:      "value2",
				expiryTime: expiryTime,
			},
		},
		{
			name: "Should set expiryTime of the entry to -1 to match TTL",
			fields: fields{
				entries: map[string]entry{},
				TTL:     -1,
			},
			args: args{
				key:   "key",
				value: "value",
			},
			want: entry{
				value:      "value",
				expiryTime: -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries: tt.fields.entries,
				TTL:     tt.fields.TTL,
			}
			c.Put(tt.args.key, tt.args.value)

			got := c.entries[tt.args.key]
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCache_Put(b *testing.B) {
	type fields struct {
		entries map[string]entry
		TTL     time.Duration
	}

	type args struct {
		key   string
		value interface{}
	}

	TTL := 1 * time.Minute
	expiryTime := time.Now().Add(TTL).Unix()

	test := struct {
		name   string
		fields fields
		args   args
		want   entry
	}{
		name: "Should add new entry",
		fields: fields{
			entries: map[string]entry{},
			TTL:     TTL,
		},
		args: args{
			key:   "key",
			value: "value",
		},
		want: entry{
			value:      "value",
			expiryTime: expiryTime,
		},
	}

	c := &Cache{
		entries: test.fields.entries,
		TTL:     test.fields.TTL,
	}

	for n := 0; n < b.N; n++ {
		c.Put(test.args.key, test.args.value)
	}
}

func TestCache_Get(t *testing.T) {
	type fields struct {
		entries map[string]entry
		TTL     time.Duration
	}

	type args struct {
		key string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		{
			name: "Should return an existing entry",
			fields: fields{
				entries: map[string]entry{
					"key": entry{
						value:      "value",
						expiryTime: time.Now().Add(1 * time.Second).Unix(),
					},
				},
			},
			args: args{
				key: "key",
			},
			want:  "value",
			want1: true,
		},
		{
			name: "Should return nil and false for an non-existant entry",
			fields: fields{
				entries: map[string]entry{
					"key": entry{
						value: "value",
					},
				},
			},
			args: args{
				key: "not-a-key",
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Should lazy delete an existing entry if expired and cache has TTL",
			fields: fields{
				entries: map[string]entry{
					"key": entry{
						value:      "value",
						expiryTime: time.Now().Add(-1 * time.Minute).Unix(),
					},
				},
				TTL: 1 * time.Minute,
			},
			args: args{
				key: "key",
			},
			want:  nil,
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries: tt.fields.entries,
				TTL:     tt.fields.TTL,
			}
			got, got1 := c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func BenchmarkCache_Get(b *testing.B) {
	type fields struct {
		entries map[string]entry
		TTL     time.Duration
	}

	type args struct {
		key string
	}

	TTL := 1 * time.Minute
	expiryTime := time.Now().Add(TTL).Unix()

	test := struct {
		name   string
		fields fields
		args   args
		want   entry
	}{
		name: "Should add new entry",
		fields: fields{
			entries: map[string]entry{
				"key": entry{
					value:      "value",
					expiryTime: expiryTime,
				},
			},
			TTL: TTL,
		},
		args: args{
			key: "key",
		},
	}

	c := &Cache{
		entries: test.fields.entries,
		TTL:     test.fields.TTL,
	}

	for n := 0; n < b.N; n++ {
		c.Get(test.args.key)
	}
}

func TestCache_Delete(t *testing.T) {
	type fields struct {
		entries map[string]entry
		TTL     time.Duration
	}

	type args struct {
		key string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should delete an entry",
			fields: fields{
				entries: map[string]entry{
					"key": entry{
						value: "value",
					},
				},
			},
			args: args{
				key: "key",
			},
		},
		{
			name: "Should idemptontly delete an entry",
			fields: fields{
				entries: map[string]entry{},
			},
			args: args{
				key: "key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries: tt.fields.entries,
				TTL:     tt.fields.TTL,
			}
			c.Delete(tt.args.key)

			if len(c.entries) != 0 {
				t.Errorf("Expected c.entries to have length of 0, got %v", len(c.entries))
			}
		})
	}
}

func Test_entry_isExpired(t *testing.T) {
	type fields struct {
		value      interface{}
		expiryTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Should return true if expiry time has passed",
			fields: fields{
				value:      "value",
				expiryTime: time.Now().Add(-1 * time.Second).Unix(),
			},
			want: true,
		},
		{
			name: "Should return false if expiry time has not passed",
			fields: fields{
				value:      "value",
				expiryTime: time.Now().Add(1 * time.Second).Unix(),
			},
			want: false,
		},
		{
			name: "Should return false if expiry time is -1",
			fields: fields{
				value:      "value",
				expiryTime: -1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				value:      tt.fields.value,
				expiryTime: tt.fields.expiryTime,
			}
			if got := e.isExpired(); got != tt.want {
				t.Errorf("entry.isExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: find better way to test than by sleeping
func TestCache_processCleanupInterval(t *testing.T) {
	type fields struct {
		entries         map[string]entry
		TTL             time.Duration
		CleanupInterval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Should cleanup expired entries",
			fields: fields{
				entries: map[string]entry{
					"key": entry{
						value:      "value",
						expiryTime: time.Now().Add(-1 * time.Minute).Unix(),
					},
				},
				TTL:             1 * time.Minute,
				CleanupInterval: 1 * time.Nanosecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries:         tt.fields.entries,
				TTL:             tt.fields.TTL,
				CleanupInterval: tt.fields.CleanupInterval,
			}

			go c.processCleanupInterval()

			time.Sleep(10 * time.Millisecond)

			if len(c.entries) != 0 {
				t.Errorf("Expected length of Cache.entries to be 0, got %v", len(c.entries))
			}
		})
	}
}
