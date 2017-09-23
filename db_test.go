package tickets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db DB = DB{}

func cleanEnv() {
	os.Setenv("MONGOHQ_URL", "")
}

func TestGetDBUrl(test *testing.T) {

	assert.Equal(test, "localhost", DBUrl())

	os.Setenv("MONGOHQ_URL", "ticketstest")

	assert.Equal(test, "ticketstest", DBUrl())

	cleanEnv()
}

func TestDBName(test *testing.T) {
	assert.Equal(test, "tickets", db.Name())
}

func BenchmarkDial(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		db.Dial()
	}
}
