package json

import (
	"bufio"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"
)

type Config struct {
	InputPath  string
	OutputPath string
}

func Run(conf Config) {
	open, err := os.Open(conf.InputPath)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer open.Close()

	output, err := os.Create(conf.OutputPath)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer output.Close()

	re := regexp.MustCompile(`\{.*}`)
	scanner := bufio.NewScanner(open)
	writer := bufio.NewWriter(output)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			matched := re.FindString(line)
			_, err := writer.WriteString(matched + "\n")
			if err != nil {
				log.Fatal().Msgf("write failed: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Msgf("scan failed: %v", err)
	}
}
