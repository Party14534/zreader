package pkg

import (
	"path"
	"strings"
)

func FilterStringsWithExtension(slice []string, extension string) (result []string) {
    for _, str := range slice {
        if(strings.Compare(path.Ext(str), extension) == 0) {
            result = append(result, str)
        } 
    }

    return result
}

