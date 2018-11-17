package cache

import (
	"reflect"
	"testing"
	"time"
)

// TODO: figure out how to test for go processes as well
func TestNew(t *testing.T) {
	type args struct {
		defaultTTL      time.Duration
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
				defaultTTL:      1 * time.Minute,
				cleanupInterval: 1 * time.Minute,
			},
			want: &Cache{
				DefaultTTL:      1 * time.Minute,
				CleanupInterval: 1 * time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.defaultTTL, tt.args.cleanupInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Put(t *testing.T) {
	type fields struct {
		entries    []entry
		DefaultTTL time.Duration
	}

	type args struct {
		key   string
		value interface{}
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Should add new entry",
			fields: fields{},
			args: args{
				key:   "key",
				value: "value",
			},
		},
		{
			name: "Should replace an existing value for a key",
			fields: fields{
				entries: []entry{
					entry{
						key:   "key",
						value: "value1",
					},
				},
			},
			args: args{
				key:   "key",
				value: "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries:    tt.fields.entries,
				DefaultTTL: tt.fields.DefaultTTL,
			}
			c.Put(tt.args.key, tt.args.value)

			v, _ := c.Get(tt.args.key)
			if v != tt.args.value {
				t.Error("Returned value did not match inserted value", v)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	type fields struct {
		entries    []entry
		DefaultTTL time.Duration
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
				entries: []entry{
					entry{
						key:   "key",
						value: "value",
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
				entries: []entry{
					entry{
						key:   "key",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries:    tt.fields.entries,
				DefaultTTL: tt.fields.DefaultTTL,
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

func TestCache_Delete(t *testing.T) {
	type fields struct {
		entries    []entry
		DefaultTTL time.Duration
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
				entries: []entry{
					entry{
						key:   "key",
						value: "value",
					},
				},
			},
			args: args{
				key: "key",
			},
		},
		{
			name:   "Should idemptontly delete an entry",
			fields: fields{},
			args: args{
				key: "key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				entries:    tt.fields.entries,
				DefaultTTL: tt.fields.DefaultTTL,
			}
			c.Delete(tt.args.key)

			if _, found := c.Get(tt.args.key); found {
				t.Error("Cache did not delete desired key")
			}
		})
	}
}

func Test_entry_IsExpired(t *testing.T) {
	type fields struct {
		key        string
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
				key:        "key",
				value:      "value",
				expiryTime: time.Now().Add(-1 * time.Second).Unix(),
			},
			want: true,
		},
		{
			name: "Should return false if expiry time has not passed",
			fields: fields{
				key:        "key",
				value:      "value",
				expiryTime: time.Now().Add(1 * time.Second).Unix(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:        tt.fields.key,
				value:      tt.fields.value,
				expiryTime: tt.fields.expiryTime,
			}
			if got := e.IsExpired(); got != tt.want {
				t.Errorf("entry.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
