package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test the function on different structures and other types.
type (
	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	MinValidator struct {
		Min100      int   `validate:"min:100"`
		SliceMin100 []int `validate:"min:100"`
	}

	MaxValidator struct {
		Max1000      int   `validate:"max:1000"`
		SliceMax1000 []int `validate:"max:1000"`
	}

	IntegerInValidator struct {
		Digit  int   `validate:"in:0,1,2,3,4,5,6,7,8,9"`
		Digits []int `validate:"in:0,1,2,3,4,5,6,7,8,9"`
	}

	IntegerComplex struct {
		Age []int `validate:"min:18|max:70"`
	}

	LenValidator struct {
		Phone  string   `validate:"len:11"`
		Phones []string `validate:"len:11"`
	}

	RegexpValidator struct {
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Emails []string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	StringInValidator struct {
		Season  string   `validate:"in:winter,spring,summer,autumn"`
		Seasons []string `validate:"in:winter,spring,summer,autumn"`
	}

	User struct {
		ID      string `json:"id" validate:"len:36"`
		Name    string
		Age     int             `validate:"min:18|max:50"`
		Email   string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role    UserRole        `validate:"in:admin,stuff"`
		Phones  []string        `validate:"len:11"`
		Address Address         `validate:"nested"`
		meta    json.RawMessage //nolint:unused
	}

	UserRole string

	Address struct {
		City   string `validate:"in:New-York,London,Moscow"`
		Street string
		House  int `validate:"min:0|max:1000"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		wantErr error
	}{
		{
			name: "no validate tags",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			wantErr: nil,
		},
		{
			name:    "not a struct",
			in:      "test",
			wantErr: ErrPassedArgumentIsNotStruct,
		},
		{
			name: "integer, min, positive test",
			in: MinValidator{
				Min100:      100,
				SliceMin100: []int{101, 102, 103, 100_000_000},
			},
			wantErr: nil,
		},
		{
			name: "integer, min, negative test",
			in: MinValidator{
				Min100:      99,
				SliceMin100: []int{101, 102, 10, 100_000_000},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Min100",
					Err:   ErrMinRoolViolated,
				},
				ValidationError{
					Field: "SliceMin100",
					Err:   ErrMinRoolViolated,
				},
			},
		},
		{
			name: "integer, max, positive test",
			in: MaxValidator{
				Max1000:      1000,
				SliceMax1000: []int{1, 2, 3, 4, 5, 999},
			},
			wantErr: nil,
		},
		{
			name: "integer, max, negative test",
			in: MaxValidator{
				Max1000:      3000,
				SliceMax1000: []int{1, 2, 3, 4, 5, 1001},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Max1000",
					Err:   ErrMaxRoolViolated,
				},
				ValidationError{
					Field: "SliceMax1000",
					Err:   ErrMaxRoolViolated,
				},
			},
		},
		{
			name: "integer, in, positive test",
			in: IntegerInValidator{
				Digit:  6,
				Digits: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			wantErr: nil,
		},
		{
			name: "integer, in, negative test",
			in: IntegerInValidator{
				Digit:  10,
				Digits: []int{0, 1, 2, 3, 4, 25, 6, 7, 8, 9},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Digit",
					Err:   ErrInRoolViolated,
				},
				ValidationError{
					Field: "Digits",
					Err:   ErrInRoolViolated,
				},
			},
		},
		{
			name: "integer complex test, min and max, positive test",
			in: IntegerComplex{
				Age: []int{18, 19, 20, 21, 22, 40, 69, 70},
			},
			wantErr: nil,
		},
		{
			name: "integer complex test, min and max, negative test",
			in: IntegerComplex{
				Age: []int{17, 18, 19, 20, 21, 22, 40, 69, 70, 71, 100},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrMinRoolViolated,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrMaxRoolViolated,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrMaxRoolViolated,
				},
			},
		},
		{
			name: "string, len, positive test",
			in: LenValidator{
				Phone:  "79001112233",
				Phones: []string{"79001112233", "79112223344", "79223334455"},
			},
			wantErr: nil,
		},
		{
			name: "string, len, negative test",
			in: LenValidator{
				Phone:  "790011122334455",
				Phones: []string{"79001112233", "791122233445566", "79223334455"},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Phone",
					Err:   ErrLenRoolViolated,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrLenRoolViolated,
				},
			},
		},
		{
			name: "string, regexp, positive test",
			in: RegexpValidator{
				Email:  "test@mail.ru",
				Emails: []string{"johndoe@icloud.com", "janedoe@outlook.com", "peterjohnson1976@gmail.com"},
			},
			wantErr: nil,
		},
		{
			name: "string, regexp, negative test",
			in: RegexpValidator{
				Email:  "test@@mail.ru",
				Emails: []string{"johndoe@icloud..com", "janedoe@outlook.com", "peterjohnson1976@gmail.com"},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrRegexpRoolViolated,
				},
				ValidationError{
					Field: "Emails",
					Err:   ErrRegexpRoolViolated,
				},
			},
		},
		{
			name: "string, in, positive test",
			in: StringInValidator{
				Season:  "winter",
				Seasons: []string{"winter", "spring", "summer", "autumn"},
			},
			wantErr: nil,
		},
		{
			name: "string, in, negative test",
			in: StringInValidator{
				Season:  "monday",
				Seasons: []string{"winter", "spring", "sammer", "autumn"},
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Season",
					Err:   ErrInRoolViolated,
				},
				ValidationError{
					Field: "Seasons",
					Err:   ErrInRoolViolated,
				},
			},
		},
		{
			name: "complex test with nested struct, positive",
			in: User{
				ID:     "012345678901234567890123456789012345",
				Name:   "Vasily",
				Age:    19,
				Email:  "vasilystepanov@gmail.com",
				Role:   "stuff",
				Phones: []string{"79001112233"},
				Address: Address{
					City:   "Moscow",
					Street: "Pushkin street",
					House:  5,
				},
				meta: nil,
			},
			wantErr: nil,
		},
		{
			name: "complex test with nested struct, negative",
			in: User{
				ID:     "012345678901234567890123456789012345",
				Name:   "Vasily",
				Age:    19,
				Email:  "vasilystepanov@@gmail.com",
				Role:   "stuff",
				Phones: []string{"79001112233"},
				Address: Address{
					City:   "Moskva",
					Street: "Pushkin street",
					House:  -1,
				},
				meta: nil,
			},
			wantErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrRegexpRoolViolated,
				},
				ValidationError{
					Field: "City",
					Err:   ErrInRoolViolated,
				},
				ValidationError{
					Field: "House",
					Err:   ErrMinRoolViolated,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			gotErr := Validate(tt.in)

			require.Equal(t, tt.wantErr, gotErr, fmt.Sprintf("test case: %s", tt.name))
		})
	}
}
