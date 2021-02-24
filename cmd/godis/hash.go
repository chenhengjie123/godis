/*
Copyright © 2021 lihuacai

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "hash key operation",
}

var (
	keyfmt *prettyjson.Formatter
)

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.AddCommand(hGetAllCmd)
	rootCmd.AddCommand(hGetAllCmdShort)

	hashCmd.AddCommand(hGetCmd)
	rootCmd.AddCommand(hGetCmdShort)

	hashCmd.AddCommand(hashCopyCmd)
	rootCmd.AddCommand(hashCopyCmdShort)

	hashCmd.AddCommand(hmsetCmd)
	rootCmd.AddCommand(hmsetCmdShort)

	hashCmd.AddCommand(hashMvCmd)
	rootCmd.AddCommand(hashMvCmdShort)

	keyfmt = prettyjson.NewFormatter()
	keyfmt.Newline = " " // Replace newline with space to avoid condensed output.
	keyfmt.Indent = 0
}

var hGetAllCmdShort = hGetAllCmd
var hGetAllCmd = &cobra.Command{
	Use:     "hgetall [key]",
	Aliases: []string{"hga"},
	Short:   "hash key hgetall",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result := hGetAll(args[0])
		_, _ = colorableOut.Write(map2Json(result))
		fmt.Fprintln(outWriter)
	},
}
var hGetCmdShort = hGetCmd
var hGetCmd = &cobra.Command{
	Use:     "hget [key] [field]",
	Aliases: []string{"hg"},
	Short:   "hash key hget",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result := hGet(args[0], args[1])
		b := []byte(result)
		_, _ = colorableOut.Write(str2Json(b))
		fmt.Fprintln(outWriter)
	},
}

var hashCopyCmdShort = hashCopyCmd
var hashCopyCmd = &cobra.Command{
	Use:     "hcopy [old_key] [new_key]",
	Aliases: []string{"hcp"},
	Short:   "copy a hash key",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result := hGetAll(args[0])
		for k, v := range result {
			rdb.HSet(ctx, args[1], k, v)
		}
		newResult := hGetAll(args[1])
		fmt.Println("hgetall " + args[1])
		_, _ = colorableOut.Write(map2Json(newResult))
		fmt.Fprintln(outWriter)
	},
}

var hashMvCmdShort = hashMvCmd
var hashMvCmd = &cobra.Command{
	Use:     "hrename [old_key] [new_key]",
	Aliases: []string{"hmv"},
	Short:   "rename a hash key",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result := hGetAll(args[0])
		for k, v := range result {
			rdb.HSet(ctx, args[1], k, v)
		}
		deleteCmd(args[0])
		newResult := hGetAll(args[1])
		fmt.Println("hgetall " + args[1])
		_, _ = colorableOut.Write(map2Json(newResult))
		fmt.Fprintln(outWriter)
	},
}
var hmsetCmdShort = hmsetCmd
var hmsetCmd = &cobra.Command{
	Use:     "hmset [key] [jsonValue]",
	Aliases: []string{"hms"},
	Short:   "add a hash key, auto unpack jsonValue",
	Args:    cobra.ExactArgs(2),
	Example: `godis hash hmset test_key '{"a":1, "b": "b"}'`,
	Run: func(cmd *cobra.Command, args []string) {
		mapValue, err := JsonToMap(args[1])
		if err != nil {
			errorExit("Unmarshal with error: %+v\n", err)
		}
		for k, v := range mapValue {
			rdb.HSet(ctx, args[0], k, v)
		}
		newResult := hGetAll(args[0])
		fmt.Println("hmset success, hash key is "+args[0], ", value is ")
		_, _ = colorableOut.Write(map2Json(newResult))
		fmt.Fprintln(outWriter)
	},
}

func deleteCmd(key string) int64 {
	r, _ := rdb.Del(ctx, key).Result()
	return r
}

func hGetAll(key string) map[string]string {
	r, _ := rdb.HGetAll(ctx, key).Result()
	return r
}

func hGet(key, field string) string {
	result, _ := rdb.HGet(ctx, key, field).Result()
	return result
}

func map2Json(key map[string]string) []byte {
	jsonStr, _ := json.Marshal(key)
	b, _ := prettyjson.Format(jsonStr)
	return b
}

func str2Json(key []byte) []byte {
	if b, err := prettyjson.Format(key); err == nil {
		return b
	}
	return key
}

// Convert json string to map
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
