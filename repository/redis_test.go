package repository

import "testing"

func TestNewRedisDB(t *testing.T) {
	tests := []struct {
		name    string
		cfg     RedisConfig
		wantErr bool
	}{
		{
			name: "success",
			cfg: RedisConfig{
				Addr:     "localhost:6379",
				Password: "qwerty",
				DB:       0,
			},
			wantErr: false,
		},
		{
			name: "error open redis",
			cfg: RedisConfig{
				Addr:     "localhost:63379",
				Password: "dasdasvdcc vd",
				DB:       0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewRedisDB(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if db != nil {
				defer db.Close()
			}
		})
	}

}
