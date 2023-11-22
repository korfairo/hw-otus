package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

type DomainStat map[string]int

var ErrGetUsers = errors.New("get users error")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, errors.Wrap(ErrGetUsers, err.Error())
	}
	return countDomains(u, domain)
}

type users [100000]User

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func getUsers(r io.Reader) (result users, err error) {
	var user User
	var i int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err = easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	suffix := "." + domain
	for _, user := range u {
		if strings.HasSuffix(user.Email, suffix) {
			secondLevelDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[secondLevelDomain]++
		}
	}
	return result, nil
}
