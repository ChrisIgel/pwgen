package cmd

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
  defaultDictURL = "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"
)

var (
  dictFileName string
  separator string
  passwordCount uint32
  wordsPerPassword uint16
  minWordLength uint8
  maxWordLength uint8
)

var rootCmd = &cobra.Command{
  Use: "pwgen",
  Short: "pwgen is a password generator for memorable passwords",
  Long: `pwgen is a password generator for memorable passwords
If no dictionary is supplied a default dictionary is downloaded.`,
  Run: func(cmd *cobra.Command, args []string) {
    seedRng()

    if !fileExists(dictFileName) {
      fmt.Printf("Downloading default dict: %v\n", defaultDictURL)
      err := downloadFile(dictFileName, defaultDictURL)
      if err != nil {
        fmt.Printf("Download failed: %v\n", err)
        os.Exit(1)
      }
      fmt.Println("Download finished!")
      fmt.Println()
    }
  
    content, err := ioutil.ReadFile(dictFileName)
    if err != nil {
      fmt.Printf("Error reading file %v: %v\n", dictFileName, err)
      os.Exit(1)
    }

    var words []string
    if strings.Contains(string(content), "\r\n") {
      words = strings.Split(string(content), "\r\n")
    } else {
      words = strings.Split(string(content), "\n")
    }

    words = filter(words, func(word string) bool {
      return len(word) >= int(minWordLength) && len(word) <= int(maxWordLength)
    })

    if len(words) < int(wordsPerPassword) {
      fmt.Println("Not enough words in the dictionary for each password!")
      os.Exit(1)
    }
    
    for i := uint32(0); i < passwordCount; i++ {
      var pw string
      tmp := make([]string, len(words))
      copy(tmp, words)

      for j := uint16(0); j < wordsPerPassword; j++ {
        word := popRandom(&tmp)
        pw += word + separator
      }

      pw = strings.TrimSuffix(pw, separator)
    
      fmt.Println(pw)
    }
  },
}

// Execute runs the main pwgen command
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.Flags().Uint8Var(&minWordLength, "min-length", 3, "the minimum length of each word")
  rootCmd.Flags().Uint8Var(&maxWordLength, "max-length", 10, "the maximum length of each word")
  rootCmd.Flags().Uint16VarP(&wordsPerPassword, "words", "w", 5, "the number of words to use per password")
  rootCmd.Flags().Uint32VarP(&passwordCount, "count", "c", 10, "the number of passwords to generate")
  rootCmd.Flags().StringVarP(&dictFileName, "dict", "d", "dict.txt", "the file to read/save the dicionary from/to")
  rootCmd.Flags().StringVarP(&separator, "separator", "s", "-", "the separator to use between each word")
}

func filter(vs []string, f func(string) bool) []string {
  vsf := make([]string, 0)
  for _, v := range vs {
      if f(v) {
          vsf = append(vsf, v)
      }
  }
  return vsf
}

func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
    return false
  }
  return !info.IsDir()
}

func downloadFile(filename, url string) error {
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  out, err := os.Create(filename)
  if err != nil {
    return err
  }
  defer out.Close()

  _, err = io.Copy(out, resp.Body)
  return err
}

func popRandom(words *[]string) string {
  pos := rand.Intn(len(*words))
  word := (*words)[pos]

  (*words)[pos] = (*words)[len(*words)-1]
  *words = (*words)[:len(*words)-1]

  return word
}

func seedRng() {
  var b [8]byte
  _, err := cryptoRand.Read(b[:])
  if err != nil {
    panic("cannot seed math/rand package with cryptographically secure random number generator")
  }
  rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}