package repository

import "testing"

func TestNewPostgresDB(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "success",
			cfg: Config{
				Username: "postgres",
				Host:     "localhost",
				Port:     "5436",
				DBName:   "postgres",
				SSLMode:  "disable",
				Password: "qwerty",
			},
			wantErr: false,
		},
		{
			name: "error open postgres",
			cfg: Config{
				Username: "postgres",
				Host:     "localhost",
				Port:     "5436234",
				DBName:   "testdb",
				SSLMode:  "disable",
				Password: "postgres",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewPostgresDB(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostgresDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if db != nil {
				defer db.Close()
			}
		})
	}
}
