package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"github.com/jdxj/fare"
)

func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:                        "fare",
		Aliases:                    nil,
		SuggestFor:                 nil,
		Short:                      "fare用于计算运费",
		GroupID:                    "",
		Long:                       "",
		Example:                    "",
		ValidArgs:                  nil,
		ValidArgsFunction:          nil,
		Args:                       nil,
		ArgAliases:                 nil,
		BashCompletionFunction:     "",
		Deprecated:                 "",
		Annotations:                nil,
		Version:                    "",
		PersistentPreRun:           nil,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        nil,
		RunE:                       nil,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
		CompletionOptions:          cobra.CompletionOptions{},
		TraverseChildren:           false,
		Hidden:                     false,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
	}

	// file flag
	var (
		fileFlag  = "file"
		fileUsage = "待计算的csv文件, 可以指定目录/文件, 可以指定多次. 注意不要有相同的文件名!"
		files     []string
	)
	root.Flags().StringSliceVarP(&files, fileFlag, "f", []string{"."}, fileUsage)

	// output flag
	var (
		outputFlag  = "output"
		outputUsage = "输出目录. 不指定时输出到源文件目录."
		output      string
	)
	root.Flags().StringVarP(&output, outputFlag, "o", "", outputUsage)

	// level flag
	var (
		levelFlag  = "level"
		levelUsage = "打印调试日志, 非开发人员无需使用. debug:-1 info:0 warn:1 error:2"
		level      int
	)
	root.Flags().IntVarP(&level, levelFlag, "l", 2, levelUsage)

	root.Run = func(cmd *cobra.Command, args []string) {
		fare.SetLoggerLevel(zapcore.Level(level))

		csvFiles, err := getCSVPath(files)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		for _, v := range csvFiles {
			if err := handleCSV(output, v); err != nil {
				cmd.PrintErrln(err)
			}
		}

		cmd.Println("处理完成")
	}
	return root
}

func getCSVPath(files []string) ([]string, error) {
	var csvPath []string
	for _, root := range files {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if errors.Is(err, syscall.ENOENT) {
					return fmt.Errorf("%s: 没有这个文件或目录", path)
				}
				return err
			}
			if d.IsDir() {
				return nil
			}

			if filepath.Ext(d.Name()) != ".csv" {
				return nil
			}

			csvPath = append(csvPath, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if len(csvPath) <= 0 {
		return nil, fmt.Errorf("未发现csv文件")
	}
	return csvPath, nil
}

func handleCSV(output, src string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	var (
		dst    string
		suffix = "_费用.csv"
	)
	if output == "" {
		dst = strings.TrimSuffix(src, ".csv") + suffix
	} else {
		dst = filepath.Join(output,
			strings.TrimSuffix(filepath.Base(src), ".csv")+suffix)
	}

	target, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer target.Close()

	err = fare.CalcFare(target, file)
	if err != nil {
		return err
	}
	return target.Sync()
}
